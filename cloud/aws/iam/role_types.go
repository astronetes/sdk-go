package iam

type CreateRoleRequest struct {
	Name                     string
	AssumeRolePolicyDocument PolicyDocument
}
