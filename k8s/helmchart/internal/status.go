package internal

import (
	"context"
	"fmt"

	"github.com/astronetes/sdk-go/log"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
)

func IsCompleted(_ context.Context, action *action.Status,
	name string,
) (bool, error) {
	res, err := action.Run(name)
	if err != nil {
		return false, fmt.Errorf("error installing helm chart: '%w", err)
	}

	log.Log.V(1).Info("installation completed successfully",
		"release", res.Name, "status", res.Info.Status,
	)

	if res.Info.Status == release.StatusDeployed {
		return true, nil
	}
	if res.Info.Status == release.StatusPendingInstall {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status '%v'", res.Info.Status)
}

func GetStatus(_ context.Context, action *action.Status, name string) (release.Status, error) {
	res, err := action.Run(name)
	if err != nil {
		return release.StatusUnknown, fmt.Errorf("error installing helm chart: '%w", err)
	}

	log.Log.V(1).Info("installation completed successfully",
		"release", res.Name, "status", res.Info.Status,
	)

	return res.Info.Status, nil
}
