package iam

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
)

func IsNotFoundError(err error) bool {
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *types.NoSuchEntityException:
				return true
			default:
				return false
			}
		}
	}

	return false
}
