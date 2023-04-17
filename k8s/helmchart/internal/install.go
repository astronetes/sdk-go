package internal

import (
	"context"
	"fmt"

	"github.com/astronetes/sdk-go/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
)

func Install(_ context.Context, action *action.Install,
	chart *chart.Chart, values map[string]interface{},
) error {
	release, err := action.Run(chart, values)
	if err != nil {
		return fmt.Errorf("error installing helm chart: '%w", err)
	}

	log.Log.V(1).Info("release '%s' with status '%s'", release.Name, release.Info.Status)

	return nil
}
