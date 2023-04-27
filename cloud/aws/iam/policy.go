package iam

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func CreatePolicy(ctx context.Context, client *iam.Client, req CreatePolicyRequest) (*types.Policy, error) {
	policyDocument, err := getPolicyDocumentAsJson(req.Document)
	if err != nil {
		return nil, err
	}

	input := &iam.CreatePolicyInput{
		PolicyName:     aws.String(req.Name),
		Description:    aws.String(req.Description),
		PolicyDocument: aws.String(string(policyDocument)),
	}
	output, err := client.CreatePolicy(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.Policy, nil
}

func GetPolicy(ctx context.Context, client *iam.Client, arn string) (*types.Policy, error) {
	input := &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	}
	output, err := client.GetPolicy(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Policy, nil
}

func DeletePolicy(ctx context.Context, client *iam.Client, arn string) error {
	input := &iam.DeletePolicyInput{
		PolicyArn: aws.String(arn),
	}
	_, err := client.DeletePolicy(ctx, input)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	return nil
}

func getPolicyDocumentAsJson(policy PolicyDocument) ([]byte, error) {
	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	return policyBytes, nil
}
