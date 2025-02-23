package main

import (
	"context"
	"errors"
	"strconv"
	"strings"

	apiapp "k8s.io/kubernetes/cmd/kube-apiserver/app"
)

type APIServerConfig struct {
	// Required
	BindAddr      string
	BindPort      int
	ClusterDomain string
	APICertFile   string
	APIKeyFile    string
	ETCDServers   []string

	// Optional
	// Uses a Certificate Authority file to authenticate requests
	ClientCACertFile string
	// Uses a csv file to authenticate requests with tokens
	TokenAuthFile string
	// Log level for apiserver
	LogLevel int
	// AdvertiseAddress is the IP address on which to advertise the apiserver to members of the cluster.
	AdvertiseAddress string

	// Not relevant to deployment, but required by apiserver
	ClusterCIDR string
}

func (c *APIServerConfig) Validate() error {
	if c.BindAddr == "" {
		return errors.New("BindAddr is required")
	}
	if c.BindPort == 0 {
		return errors.New("BindPort is required")
	}
	if c.APICertFile == "" {
		return errors.New("APICertFile is required")
	}
	if c.APIKeyFile == "" {
		return errors.New("APIKeyFile is required")
	}
	if len(c.ETCDServers) == 0 {
		return errors.New("ETCDServers is required")
	}
	if c.ClusterDomain == "" {
		return errors.New("ClusterDomain is required")
	}
	if c.ClusterCIDR == "" {
		return errors.New("ClusterCIDR is required")
	}
	return nil
}

func defaultAPIServerConfig() *APIServerConfig {
	return &APIServerConfig{
		BindAddr:      "0.0.0.0",
		BindPort:      6443,
		ClusterDomain: "cluster.local",
		ClusterCIDR:   "10.0.0.0/24",
		LogLevel:      4,
	}
}

func startAPIServer(ctx context.Context, cfg *APIServerConfig) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	argsMap := map[string]string{
		"authorization-mode":               "RBAC",
		"bind-address":                     cfg.BindAddr,
		"secure-port":                      strconv.Itoa(cfg.BindPort),
		"service-cluster-ip-range":         cfg.ClusterCIDR,
		"tls-cert-file":                    cfg.APICertFile,
		"tls-private-key-file":             cfg.APIKeyFile,
		"service-account-signing-key-file": cfg.APIKeyFile,
		"service-account-key-file":         cfg.APIKeyFile,
		"service-account-issuer":           "https://kubernetes.default.svc." + cfg.ClusterDomain,
		"api-audiences":                    "https://kubernetes.default.svc." + cfg.ClusterDomain,
		"etcd-servers":                     strings.Join(cfg.ETCDServers, ","),
		"v":                                strconv.Itoa(cfg.LogLevel),
		"profiling":                        "false",
		"storage-backend":                  "etcd3",
		"anonymous-auth":                   "false",
	}

	if cfg.ClientCACertFile != "" {
		argsMap["client-ca-file"] = cfg.ClientCACertFile
	}

	if cfg.TokenAuthFile != "" {
		argsMap["token-auth-file"] = cfg.TokenAuthFile
	}

	if cfg.AdvertiseAddress != "" {
		argsMap["advertise-address"] = cfg.AdvertiseAddress
	}

	command := apiapp.NewAPIServerCommand()
	apiArgs := flattenArgs(argsMap)
	command.SetArgs(apiArgs)

	return command.ExecuteContext(ctx)
}

func flattenArgs(argsMap map[string]string) []string {
	args := make([]string, 0, len(argsMap))
	for k, v := range argsMap {
		args = append(args, "--"+k+"="+v)
	}
	return args
}
