// Heavily inspired from vcluster e2e test

package framework

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	PollTimeout          = time.Minute
	DefaultClientTimeout = 32 * time.Second // the default in client-go is 32
)

var DefaultFramework = &Framework{}

type Framework struct {
	// The context to use for testing
	Context context.Context

	// ClientSet is the kubernetes client of the current
	// host kubernetes cluster were we are testing in
	ClientSet *kubernetes.Clientset

	// CtrlClient is the kubernetes client originally supposed
	// to write controllers. It provides some convinience methods
	// e.g. create objects
	CtrlClient ctrlclient.Client

	// Scheme is the global scheme to use
	Scheme *runtime.Scheme

	// ClientTimeout value used in the clients
	ClientTimeout time.Duration
}

func CreateFramework(ctx context.Context, scheme *runtime.Scheme) error {
	restConfig, err := config.GetConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	ctrlClient, err := ctrlclient.New(restConfig, client.Options{})
	if err != nil {
		return err
	}

	// create the framework
	DefaultFramework = &Framework{
		Context:       ctx,
		ClientSet:     clientset,
		CtrlClient:    ctrlClient,
		Scheme:        scheme,
		ClientTimeout: DefaultClientTimeout,
	}

	return nil
}

func (f *Framework) Cleanup() error {
	return nil
}
