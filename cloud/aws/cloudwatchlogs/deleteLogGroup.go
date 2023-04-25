package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func DeleteLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req DeleteLogGroupRequest) DeleteLogGroupResponse {
	request := &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(req.Name),
	}
	response, err := client.DeleteLogGroup(ctx, request)
	return DeleteLogGroupResponse{
		response: response,
		error:    err,
	}
}

type DeleteLogGroupRequest struct {
	Name string
}

type DeleteLogGroupResponse struct {
	response *cloudwatchlogs.DeleteLogGroupOutput
	error
}

func (r DeleteLogGroupResponse) Error() error {
	return r.error
}
