package internal

import (
	"context"
	"fmt"

	"github.com/astronetes/sdk-go/log"
	"helm.sh/helm/v3/pkg/action"
)

func Uninstall(_ context.Context, action *action.Uninstall,
	release string,
) error {
	res, err := action.Run(release)
	if err != nil {
		return fmt.Errorf("error installing helm chart: '%w", err)
	}

	log.Log.V(1).Info("installation completed successfully",
		"release", res.Release.Name, "info", res.Info,
	)

	return nil
}
