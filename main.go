package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/rand"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var setupLog = ctrl.Log.WithName("setup")

func main() {
	var etcdServers string
	var apiCertFile string
	var apiKeyFile string
	var clientCAFile string
	var tokenAuthFile string
	var apiPort int

	flag.StringVar(&apiCertFile, "api-cert-file", "", "API server certificate file")
	flag.StringVar(&apiKeyFile, "api-key-file", "", "API server key file")
	flag.StringVar(&etcdServers, "etcd-servers", "http://localhost:2379/", "ETCD servers")
	flag.StringVar(&clientCAFile, "client-ca-file", "", "CA for user authentication")
	flag.StringVar(&tokenAuthFile, "token-auth-file", "", "CSV file with tokens for user authentication")
	flag.IntVar(&apiPort, "api-port", 6443, "API server port")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	rand.Seed(uint64(time.Now().UnixNano()))

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	ctx := ctrl.SetupSignalHandler()

	apiConfig := defaultAPIServerConfig()
	apiConfig.BindPort = apiPort
	apiConfig.ETCDServers = strings.Split(etcdServers, ",")
	apiConfig.APICertFile = apiCertFile
	apiConfig.APIKeyFile = apiKeyFile
	if clientCAFile != "" {
		apiConfig.ClientCACertFile = clientCAFile
	}
	if tokenAuthFile != "" {
		apiConfig.TokenAuthFile = tokenAuthFile
	}

	if err := startAPIServer(ctx, apiConfig); err != nil {
		setupLog.Error(err, "unable to start API server")
		os.Exit(1)
	}
}
