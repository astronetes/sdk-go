package helmchart

import (
	"context"
	"fmt"
	"os"

	"github.com/astronetes/sdk-go/k8s/helmchart/internal"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/util"
)

type Client interface {
	Install(ctx context.Context, spec Spec, fn func(install *action.Install)) error
}

type client struct {
	cfg *action.Configuration
}

type ClientOption func(*ClientBuilder)

func WithRESTClientGetter(restClientGetter genericclioptions.RESTClientGetter) ClientOption {
	return func(c *ClientBuilder) {
		c.restClientGetter = restClientGetter
		c.kubeClient = &kube.Client{
			Factory: util.NewFactory(restClientGetter),
			Log: func(s string, i ...interface{}) {
				fmt.Printf(s, i)
			},
		}
	}
}

func WithRegistryClient(client *registry.Client) ClientOption {
	return func(c *ClientBuilder) {
		c.registryClient = client
	}
}

func WithReleaseStorage(releaseStorage *storage.Storage) ClientOption {
	return func(c *ClientBuilder) {
		c.releaseStorage = releaseStorage
	}
}

type ClientBuilder struct {
	kubeClient       *kube.Client
	releaseStorage   *storage.Storage
	restClientGetter genericclioptions.RESTClientGetter
	registryClient   *registry.Client
}

func NewClient(opts ...ClientOption) (Client, error) {
	var (
		defaultRESTClientGetter = genericclioptions.NewConfigFlags(false)
		defaultReleaseStorage   = storage.Init(driver.NewMemory())
		kubeClient              = &kube.Client{
			Factory: util.NewFactory(defaultRESTClientGetter),
			Log: func(s string, i ...interface{}) {
				fmt.Printf(s, i)
			},
		}
	)
	defaultRegistryClient, err := registry.NewClient(
		registry.ClientOptDebug(false),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing registry client: '%v'", err)
	}

	builder := &ClientBuilder{
		kubeClient:       kubeClient,
		releaseStorage:   defaultReleaseStorage,
		restClientGetter: defaultRESTClientGetter,
		registryClient:   defaultRegistryClient,
	}
	for _, opt := range opts {
		opt(builder)
	}

	cfg := &action.Configuration{
		RESTClientGetter: builder.restClientGetter,
		Releases:         builder.releaseStorage,
		KubeClient:       builder.kubeClient,
		RegistryClient:   builder.registryClient,
		Capabilities:     chartutil.DefaultCapabilities,
		Log: func(s string, i ...interface{}) {
		},
	}
	cfg.Releases.Log = func(v string, args ...interface{}) {
		fmt.Printf(v, args)
	}
	return &client{cfg: cfg}, nil
}

func (c *client) Install(ctx context.Context, spec Spec, fn func(install *action.Install)) error {
	chart, values, err := spec.chartAndValues()
	if err != nil {
		return err
	}
	action := action.NewInstall(c.cfg)
	fn(action)
	return internal.Install(ctx, action, chart, values)
}
