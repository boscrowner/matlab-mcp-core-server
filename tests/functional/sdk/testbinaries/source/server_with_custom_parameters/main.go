// Copyright 2026 The MathWorks, Inc.

package main

import (
	"context"
	"os"

	"github.com/matlab/matlab-mcp-core-server/pkg/config"
	"github.com/matlab/matlab-mcp-core-server/pkg/i18n"
	"github.com/matlab/matlab-mcp-core-server/pkg/server"
)

func main() {
	serverDefinition := server.Definition[any]{
		Name:         "custom-parameters",
		Title:        "Custom Parameters",
		Instructions: "This is the Custom Parameters test binary",

		Parameters: []server.Parameter{
			CustomParameter(),
			CustomRecordedParameter(),
		},

		DependenciesProvider: func(dependenciesProviderResources server.DependenciesProviderResources) (any, i18n.Error) {
			logger := dependenciesProviderResources.Logger()
			cfg := dependenciesProviderResources.Config()

			customParameter := CustomParameter()
			customParameterValue, err := config.Get(cfg, customParameter)
			if err != nil {
				return nil, err
			}

			logger.With(customParameter.GetID(), customParameterValue).Info("Config value from dependency provider")

			return nil, nil
		},

		ToolsProvider: func(toolsProviderResources server.ToolsProviderResources[any]) []server.Tool {
			logger := toolsProviderResources.Logger()
			cfg := toolsProviderResources.Config()

			// This is purely for example purposes.
			// You should retrieve config values in the DependenciesProvider or during tool handlers.
			customParameter := CustomParameter()
			customParameterValue, err := config.Get(cfg, customParameter)
			if err != nil {
				return nil
			}

			logger.With(customParameter.GetID(), customParameterValue).Info("Config value from tools provider")

			return nil
		},
	}
	serverInstance := server.New(serverDefinition)

	exitCode := serverInstance.StartAndWaitForCompletion(context.Background())

	os.Exit(exitCode)
}
