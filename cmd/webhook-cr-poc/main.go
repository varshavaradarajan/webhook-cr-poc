package main

import (
	"github.com/varshavaradarajan/webhook-cr-poc/dokswebhooks"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme = runtime.NewScheme()
)

func main() {
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	ll := ctrl.Log.WithName("dokswebhooks")
	ll.Info("getting config")
	config := ctrl.GetConfigOrDie()

	ll.Info("constructing cr client")
	c, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		ll.Error(err, "failed to construct cr client")
		os.Exit(1)
	}

	ll.Info("setting up webhook server")
	// default server running at port 9443, looking for tls.crt, tls.key in /tmp/k8s-webhook-server/serving-certs
	server := webhook.Server{}

	ll.Info("registering paths")
	server.Register("/validate-doks-lb-service", &webhook.Admission{Handler: &dokswebhooks.DOKSLBServiceValidator{Client: c, Log: ll}})
	ll.Info("starting webhook server")
	if err := server.StartStandalone(ctrl.SetupSignalHandler(), scheme); err != nil {
		ll.Error(err, "failed to start webhook server")
		os.Exit(1)
	}
	// expose a Debug server on 8081 to handle liveness and readiness checks
}
