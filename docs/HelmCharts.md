# Helm Charts API

The Astronetes SDK Go provides us with a powerful API to interact with Helm programatically.

## Client Configuration

The SDK creates an out-of-the-box configuration for the client. The configuration is taken from your operating system
in which the code is running.

```go
package main

import "github.com/astronetes/sdk-go/k8s/helmchart"

func main(){
	client, err := helmchart.NewClient()
	if err!=nil{
		
    }
	...
}
```

On the other hand, you can customize the client configuration programmatically with the following functions.

* WithRESTClientGetter(restClientGetter genericclioptions.RESTClientGetter) ClientOption
* WithRegistryClient(client *registry.Client) ClientOption
* WithReleaseStorage(releaseStorage *storage.Storage) ClientOption


The client is based in the creational options pattern 

```go
package main

import (
	"os"
	"github.com/astronetes/sdk-go/k8s/helmchart"
	"github.com/ghodss/yaml"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	registryClient, err := registry.NewClient(
		registry.ClientOptDebug(true),
		registry.ClientOptEnableCache(true),
		registry.ClientOptWriter(os.Stderr),
	)
	if err != nil {

	}
	client, err := helmchart.NewClient(
		helmchart.WithReleaseStorage(storage.Init(driver.NewMemory())),
		helmchart.WithRegistryClient(registryClient),
		helmchart.WithRESTClientGetter(genericclioptions.NewConfigFlags(true)),
	)
}
```

## Spec Configuration

A Spec is composed by a Chart and optionally a values map,  that is used by Helm to define the resources dynamically. The
SDK can load the packaged Chart from a wide range of file system implementations.

Let's see the following example to understand

```go
package main

import (
	
	"github.com/astronetes/sdk-go/k8s/helmchart"
)
func main() {
	spec:= helmchart.LoadPackagedChart("file:///tmp/charts/mysql-9.7.1.tgz")		
}
```

To run the above piece of code you could download the packaged chart as below:

```bash
mkdir -p /tmp/charts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm pull bitnami/mysql --version 9.7.1
mv mysql-9.7.1.tgz /tmp/charts/
```

