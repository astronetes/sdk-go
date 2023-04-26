package cloudwatchlogs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetLogGroup(ctx context.Context, client *cloudwatchlogs.Client, req GetLogGroupRequest) (*GetLogGroupResponse, error) {
	request := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(req.Name),
	}
	response, err := client.DescribeLogGroups(ctx, request)
	if err != nil {
		return nil, err
	}

	logGroup := getLogGroupByName(response.LogGroups, req.Name)
	if logGroup == nil {
		return nil, fmt.Errorf("LogGroup not found with name %v", req.Name)
	}

	return &GetLogGroupResponse{
		logGroup: logGroup,
	}, nil
}

type GetLogGroupRequest struct {
	Name string
}

type GetLogGroupResponse struct {
	logGroup *types.LogGroup
}

func (r GetLogGroupResponse) LogGroup() *types.LogGroup {
	return r.logGroup
}

func getLogGroupByName(logGroups []types.LogGroup, name string) *types.LogGroup {
	for _, logGroup := range logGroups {
		if *logGroup.LogGroupName == name {
			return &logGroup
		}
	}

	return nil
}
