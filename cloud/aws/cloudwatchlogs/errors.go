package cloudwatchlogs

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/smithy-go"
)

func IsNotFoundError(err error) bool {
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *types.ResourceNotFoundException:
				return true
			default:
				return false
			}
		}
	}

	return false
}
