package cloudwatchlogs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func CreateLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req CreateLogGroupRequest) CreateLogGroupResponse {
	request := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(req.Name),
	}
	response, err := client.CreateLogGroup(ctx, request)
	if err != nil {
		return CreateLogGroupResponse{
			request:  req,
			response: response,
			error:    err,
		}
	}

	retentionRequest := &cloudwatchlogs.PutRetentionPolicyInput{
		LogGroupName:    request.LogGroupName,
		RetentionInDays: aws.Int32(req.RetentionDays),
	}
	retentionResponse, err := client.PutRetentionPolicy(ctx, retentionRequest)
	return CreateLogGroupResponse{
		request:           req,
		response:          response,
		retentionResponse: retentionResponse,
		error:             err,
	}
}

type CreateLogGroupRequest struct {
	Name          string
	RetentionDays int32
}

type CreateLogGroupResponse struct {
	request           CreateLogGroupRequest
	response          *cloudwatchlogs.CreateLogGroupOutput
	retentionResponse *cloudwatchlogs.PutRetentionPolicyOutput
	error
}

func (r CreateLogGroupResponse) Error() error {
	return r.error
}

func (r CreateLogGroupResponse) RetentionDays() int32 {
	return r.request.RetentionDays
}
