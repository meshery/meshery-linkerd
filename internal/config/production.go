package config

import (
	"github.com/layer5io/meshery-adapter-library/common"
	configprovider "github.com/layer5io/meshery-adapter-library/config/provider"
)

var (
	ProductionConfig = configprovider.Options{
		ServerConfig:   productionServerConfig,
		MeshSpec:       productionMeshSpec,
		ProviderConfig: productionProviderConfig,
		Operations:     productionOperations,
	}

	productionServerConfig = map[string]string{
		"name":    "linkerd-adapter",
		"port":    "10001",
		"version": "v1.0.0",
	}

	productionMeshSpec = map[string]string{
		"name":     "linkerd",
		"status":   "none",
		"traceurl": "none",
		"version":  "none",
	}

	productionProviderConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "yaml",
		configprovider.FileName: "linkerd",
	}

	// Controlling the kubeconfig lifecycle with viper
	productionKubeConfig = map[string]string{
		configprovider.FilePath: configRootPath,
		configprovider.FileType: "",
		configprovider.FileName: "kubeconfig",
	}

	productionOperations = getOperations(common.Operations)
)
