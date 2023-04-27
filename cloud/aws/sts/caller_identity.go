package sts

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetAccountID(ctx context.Context, client *sts.Client) (*string, error) {
	input := &sts.GetCallerIdentityInput{}
	output, err := client.GetCallerIdentity(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.Account, nil
}

func GetAccountArn(ctx context.Context, client *sts.Client) (*string, error) {
	accountId, err := GetAccountID(ctx, client)
	if err != nil {
		return nil, err
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:root", *accountId)
	return aws.String(arn), nil
}
