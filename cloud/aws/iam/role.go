package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func CreateOrUpdateRole(ctx context.Context, client *iam.Client, req CreateRoleRequest) error {
	if _, err := GetRole(ctx, client, req.Name); err != nil {
		if !IsNotFoundError(err) {
			return err
		} else {
			return CreateRole(ctx, client, req)
		}
	}

	return nil
}

func CreateRole(ctx context.Context, client *iam.Client, req CreateRoleRequest) error {
	policyDocument, err := getPolicyDocumentAsJson(req.AssumeRolePolicyDocument)
	if err != nil {
		return err
	}

	request := &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(string(policyDocument)),
		RoleName:                 aws.String(req.Name),
	}
	_, err = client.CreateRole(ctx, request)
	return err
}

func GetRole(ctx context.Context, client *iam.Client, name string) (*types.Role, error) {
	request := &iam.GetRoleInput{
		RoleName: aws.String(name),
	}
	output, err := client.GetRole(ctx, request)
	if err != nil {
		return nil, err
	}

	return output.Role, nil
}

func DeleteRole(ctx context.Context, client *iam.Client, name string) error {
	_, err := GetRole(ctx, client, name)
	if err != nil {
		if IsNotFoundError(err) {
			return nil
		} else {
			return err
		}
	}

	if err := DetachAllPoliciesFromRole(ctx, client, name); err != nil {
		return err
	}

	request := &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	}
	_, err = client.DeleteRole(ctx, request)
	if err != nil && !IsNotFoundError(err) {
		return err
	}

	return nil
}

func DetachAllPoliciesFromRole(ctx context.Context, client *iam.Client, roleName string) error {
	awsRoleName := aws.String(roleName)

	result, err := client.ListAttachedRolePolicies(context.TODO(), &iam.ListAttachedRolePoliciesInput{
		RoleName: awsRoleName,
	})
	if err != nil {
		return err
	}

	for _, policy := range result.AttachedPolicies {

		_, err := client.DetachRolePolicy(context.TODO(), &iam.DetachRolePolicyInput{
			PolicyArn: policy.PolicyArn,
			RoleName:  awsRoleName,
		})
		if err != nil {
			return err

		}
	}

	return nil
}
