package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func CreateOrUpdateLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req CreateLogGroupRequest) error {
	getReq := GetLogGroupRequest{
		Name: req.Name,
	}
	if _, err := GetLogGroup(ctx, client, getReq); err != nil {
		if !IsNotFoundError(err) {
			return err
		} else {
			return CreateLogGroup(ctx, client, req)
		}
	}

	updateReq := UpdateLogGroupRequest{
		Name:          req.Name,
		RetentionDays: req.RetentionDays,
	}
	return UpdateLogGroup(ctx, client, updateReq)
}
