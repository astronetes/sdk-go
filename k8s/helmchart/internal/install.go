package internal

import (
	"context"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

func Install(ctx context.Context, action *action.Install,
	chart *chart.Chart, values map[string]interface{}) error {
	release, err := action.Run(chart, values)
	if err != nil {
		return err
	}
	println(release.Info.Status.String())
	return nil
}
