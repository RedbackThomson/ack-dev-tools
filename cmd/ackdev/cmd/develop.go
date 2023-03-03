// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	ackcmd "github.com/aws-controllers-k8s/code-generator/cmd/ack-generate/command"
	ackconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	ackwizard "github.com/aws-controllers-k8s/dev-tools/pkg/generate"
)

var (
	defaultCacheDir string
	sdkDir          string

	optGenAllRepos         bool
	optCacheDir            string
	optGeneratorConfigPath string
	optAWSSDKGoVersion     string
	optIgnoreServices      string
)

const (
	DefaultAPIVersion = "v1alpha1"
)

func init() {
	hd, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("unable to determine $HOME: %s\n", err)
		os.Exit(1)
	}
	defaultCacheDir = filepath.Join(hd, ".cache", appName)

	listCmd.AddCommand(listDependenciesCmd)
	listCmd.AddCommand(listRepositoriesCmd)
	listCmd.AddCommand(getConfigCmd)

	developCmd.PersistentFlags().BoolVarP(&optGenAllRepos, "all", "a", false, "all repositories")
	developCmd.PersistentFlags().StringVar(
		&optCacheDir, "cache-dir", defaultCacheDir, "Path to directory to store cached files (including clone'd aws-sdk-go repo)",
	)
	developCmd.PersistentFlags().StringVar(
		&optGeneratorConfigPath, "generator-config-path", "", "Path to file containing instructions for code generation to use",
	)
	developCmd.PersistentFlags().StringVar(
		&optAWSSDKGoVersion, "aws-sdk-go-version", "", "Version of github.com/aws/aws-sdk-go used to generate apis and controllers files",
	)
	developCmd.PersistentFlags().StringVar(
		&optIgnoreServices, "ignore", "", "List of service model names to ignore",
	)
}

var developCmd = &cobra.Command{
	Use:     "develop",
	Aliases: []string{"dev"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Opens the developer CLI for a controller",
	RunE:    developController,
}

func developController(cmd *cobra.Command, args []string) (err error) {
	if !optGenAllRepos && len(args) == 0 {
		return errors.New("requires the name of a single service")
	}

	if optGeneratorConfigPath == "" {
		return errors.New("flag --generator-config-path is required")
	}

	ctx, cancel := acksdk.ContextWithSigterm(context.Background())
	defer cancel()

	if len(args) > 0 {
		return generateSingleController(ctx, args[0])
	} else {
		modelNames, err := getSDKModelNames()
		if err != nil {
			return err
		}

		ignore := strings.Split(optIgnoreServices, ",")
		ignoreMap := make(map[string]struct{})
		for _, ig := range ignore {
			ignoreMap[ig] = struct{}{}
		}

		for _, modelName := range modelNames {
			_, ignored := ignoreMap[modelName]
			if ignored {
				continue
			}

			fmt.Printf("Generate generator.yaml for %s\n", modelName)
			if err := generateSingleController(ctx, modelName); err != nil {
				return err
			}
		}
	}

	return nil
}

func generateSingleController(ctx context.Context, svcAlias string) error {
	cfg := ackgenerate.DefaultConfig
	modelName := strings.ToLower(cfg.ModelName)
	if modelName == "" {
		modelName = svcAlias
	}

	sdkDir, err := acksdk.EnsureRepo(ctx, optCacheDir, false, optAWSSDKGoVersion, filepath.Dir(optGeneratorConfigPath))
	if err != nil {
		return err
	}

	sdkHelper := acksdk.NewHelper(sdkDir, cfg)
	sdkAPI, err := sdkHelper.API(modelName)
	if err != nil {
		modelName, err = ackcmd.FallBackFindServiceID(sdkDir, svcAlias)
		if err != nil {
			return err
		}
		// Retry using path found by querying service ID
		sdkAPI, err = sdkHelper.API(modelName)
		if err != nil {
			fmt.Printf("service %s not found\n", svcAlias)
			return nil
		}
	}

	m, err := ackmodel.New(
		sdkAPI, svcAlias, DefaultAPIVersion, cfg,
	)
	if err != nil {
		return err
	}

	cfg, err = ackconfig.New(optGeneratorConfigPath, cfg)
	if err != nil {
		return err
	}

	initialState, err := ackwizard.InitialState(&cfg, m, svcAlias, modelName, DefaultAPIVersion)
	if err != nil {
		return err
	}

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}
	p := tea.NewProgram(initialState, tea.WithAltScreen())
	res, err := p.Run()
	if err != nil {
		return err
	}

	resState, ok := res.(ackwizard.Wizard)
	if !ok {
		return errors.New("unable to parse final state as Wizard")
	}

	if err = writeGenerator(resState.Config(), svcAlias); err != nil {
		return err
	}

	return nil
}

func writeGenerator(config *ackconfig.Config, svcAlias string) error {
	y, err := yaml.Marshal(*config)
	if err != nil {
		return err
	}

	return os.WriteFile(optGeneratorConfigPath, y, 0644)
}

func getSDKModelNames() ([]string, error) {
	modelNames := make([]string, 0)

	modelsDir := path.Join(sdkDir, "models", "apis")
	serviceDirs, _ := ioutil.ReadDir(modelsDir)
	for _, serviceDir := range serviceDirs {
		if !serviceDir.IsDir() {
			continue
		}

		modelNames = append(modelNames, serviceDir.Name())

		// versions, _ := ioutil.ReadDir(sdkDir)
		// versionNames := make([]string, 0)
		// for _, vs := range versions {
		// 	if !vs.IsDir() {
		// 		continue
		// 	}
		// 	versionNames = append(versionNames, vs.Name())
		// }
		// sort.Strings(versionNames)

		// latestVersion := versionNames[len(versionNames) - 1]

		// apiFile, err := os.Open(fmt.Sprintf("%s/%s/%s/api-2.json", sdkDir, serviceDir.Name(), latestVersion))
		// if err != nil {
		// 	return nil, err
		// }

		// scanner := bufio.NewScanner(apiFile)

		// // Ignore the first 3 lines
		// scanner.Scan()
		// scanner.Scan()
		// scanner.Scan()

		// scanner.Scan()
	}

	return modelNames, nil
}
