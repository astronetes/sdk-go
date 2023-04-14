package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/astronetes/sdk-go/k8s/helmchart"
	"helm.sh/helm/v3/pkg/action"
)

const (
	releaseName = "my-release"
	namespace   = "astronetes-testing"
)

func installChartByDefault() {
	pwd, _ := os.Getwd()
	baseDir := fmt.Sprintf("file://%s", pwd)
	spec := helmchart.
		LoadPackagedChart(filepath.Join(baseDir, "tmp/redis-operator-3.1.2.tgz"))
	//LoadPackagedChart(filepath.Join(baseDir, "tmp/mysql-9.7.1.tgz")).
	c, err := helmchart.NewClient()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()

	if err := c.Install(ctx, spec, func(a *action.Install) {
		a.CreateNamespace = true
		a.ReleaseName = releaseName
		a.Namespace = namespace

	}); err != nil {
		print(err.Error())
	}
}

func installChartWithVariables() {
	pwd, _ := os.Getwd()
	baseDir := fmt.Sprintf("file://%s", pwd)
	spec := helmchart.
		LoadPackagedChart(filepath.Join(baseDir, "tmp/mysql-9.7.1.tgz")).
		With("primary", map[string]interface{}{
			"podLabels": []string{"astronetes.sdk-go/version: 0.0.1"},
		})
	/**
	With("image", map[string]interface{}{
		"tag": "ivan",
	})
	*/
	c, err := helmchart.NewClient()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	if err := c.Install(ctx, spec, func(a *action.Install) {
		a.CreateNamespace = true
		a.ReleaseName = releaseName
		a.Namespace = namespace
		a.IsUpgrade = true
	}); err != nil {
		print(err.Error())
	}
}

func installChartWithVariablesAndPath() {
	pwd, _ := os.Getwd()
	baseDir := fmt.Sprintf("file://%s", pwd)
	spec := helmchart.
		LoadPackagedChart(filepath.Join(baseDir, "tmp/mysql-9.7.1.tgz")).
		WithValuesTemplate(filepath.Join(baseDir, "tmp/mysql-values.yml"))

	c, err := helmchart.NewClient()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	if err := c.Install(ctx, spec, func(a *action.Install) {
		a.CreateNamespace = true
		a.ReleaseName = releaseName
		a.Namespace = namespace
		a.IsUpgrade = true
	}); err != nil {
		print(err.Error())
	}
}

func main() {
	installChartByDefault()
	//installChartWithVariables()
	//installChartWithVariablesAndPath()
}

func chekHelmChart() {
	pwd, _ := os.Getwd()
	baseDir := filepath.Join("file://%s", pwd)

	spec := helmchart.
		LoadPackagedChart(filepath.Join(baseDir, "tmp/mysql-9.7.1.tgz"))

	c, err := helmchart.NewClient()
	if err != nil {
		os.Exit(1)
		println(err.Error())
	}
	ctx := context.Background()
	c.Install(ctx, spec, func(a *action.Install) {
		a.CreateNamespace = true
		a.ReleaseName = releaseName
		a.Namespace = namespace
	})

	/**
		spec := helmchart.LoadPackagedChart(filepath.Join(baseDir, "tmp/prometheus-nginx-exporter-0.1.0.tgz")).
			WithValuesTemplate(filepath.Join(baseDir, "k8s/helmchart/testdata/prometheus-nginx-exporter-0.1.0-values.yml")).
			With("replicaCount", 3).
			WithEntries(map[string]interface{}{
				"imageRepository": "nginx/nginx-prometheus-exporter",
				"serviceAccount": map[string]interface{}{
					"create": true,
					"name":   "my-service-account",
				},
			})

		registryClient, err := registry.NewClient(
			registry.ClientOptDebug(true),
			registry.ClientOptEnableCache(true),
			registry.ClientOptWriter(os.Stderr),
		)
		if err != nil {

		}
		client, err := helmchart.NewClient(
			helmchart.WithReleaseStorage(storage.Init(driver.)),
			helmchart.WithRegistryClient(registryClient),
			helmchart.WithRESTClientGetter(genericclioptions.NewConfigFlags(true)),
		)
		if err != nil {
			os.Exit(1)
		}

		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		if err := client.Install("astronetes-testing", "my-release", spec); err != nil {
			println(err.Error())
			os.Exit(1)
		}
	**/
}
