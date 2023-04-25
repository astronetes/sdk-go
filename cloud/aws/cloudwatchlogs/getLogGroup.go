package cloudwatchlogs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req GetLogGroupRequest) GetLogGroupResponse {
	request := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(req.Name),
	}
	response, err := client.DescribeLogGroups(ctx, request)
	if err != nil {
		return GetLogGroupResponse{
			error: err,
		}
	}

	logGroup := getLogGroupByName(response, req.Name)
	if logGroup == nil {
		return GetLogGroupResponse{
			response: response,
			error:    fmt.Errorf("LogGroup not found with name %v", req.Name),
		}
	}

	return GetLogGroupResponse{
		logGroup: logGroup,
		response: response,
	}
}

type GetLogGroupRequest struct {
	Name string
}

type GetLogGroupResponse struct {
	response *cloudwatchlogs.DescribeLogGroupsOutput
	logGroup *types.LogGroup
	error
}

func (r GetLogGroupResponse) Error() error {
	return r.error
}

func (r GetLogGroupResponse) LogGroup() *types.LogGroup {
	return r.logGroup
}

func getLogGroupByName(response *cloudwatchlogs.DescribeLogGroupsOutput, name string) *types.LogGroup {
	for _, logGroup := range response.LogGroups {
		if *logGroup.LogGroupName == name {
			return &logGroup
		}
	}

	return nil
}
