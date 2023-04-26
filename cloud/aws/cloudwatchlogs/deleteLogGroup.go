package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func DeleteLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req DeleteLogGroupRequest) error {
	request := &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(req.Name),
	}
	_, err := client.DeleteLogGroup(ctx, request)
	if !IsNotFoundError(err) {
		return err
	}

	return nil
}

type DeleteLogGroupRequest struct {
	Name string
}
