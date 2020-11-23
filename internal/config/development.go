package config

import (
	"github.com/layer5io/meshery-adapter-library/common"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
)

var (
	// DevelopmentConfig holds the configuration for development environment
	DevelopmentConfig = configprovider.Options{
		ServerConfig:   developmentServerConfig,
		MeshSpec:       developmentMeshSpec,
		ProviderConfig: developmentProviderConfig,
		Operations:     developmentOperations,
	}

	developmentServerConfig = map[string]string{
		"name":    "linkerd-adapter",
		"port":    "10001",
		"version": "v1.0.0",
	}

	developmentMeshSpec = map[string]string{
		"name":     "linkerd",
		"status":   "none",
		"traceurl": "none",
		"version":  "none",
	}

	developmentProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "linkerd",
	}

	// Controlling the kubeconfig lifecycle with viper
	developmentKubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "kubeconfig",
	}

	// developmentOperations = getDevelopmentOperations(adapter.Operations{}) // Should be used in case of not using common operations
	developmentOperations = getOperations(common.Operations)
)
