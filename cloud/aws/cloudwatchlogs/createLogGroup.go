package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func CreateLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req CreateLogGroupRequest) error {
	request := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(req.Name),
	}
	if _, err := client.CreateLogGroup(ctx, request); err != nil {
		return err
	}

	retentionRequest := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    request.LogGroupName,
		RetentionInDays: aws.Int32(req.RetentionDays),
	}
	if _, err := client.PutRetentionPolicy(ctx, retentionRequest); err != nil {
		return err
	}

	return nil
}

type CreateLogGroupRequest struct {
	Name          string
	RetentionDays int32
}
