package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func UpdateLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req UpdateLogGroupRequest) error {
	retentionRequest := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(req.Name),
		RetentionInDays: aws.Int32(req.RetentionDays),
	}
	if _, err := client.PutRetentionPolicy(ctx, retentionRequest); err != nil {
		return err
	}

	return nil
}

type UpdateLogGroupRequest struct {
	Name          string
	RetentionDays int32
}
