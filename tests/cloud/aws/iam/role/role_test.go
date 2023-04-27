package role_test

import (
	"context"

	"github.com/astronetes/sdk-go/cloud/aws/iam"
	"github.com/astronetes/sdk-go/cloud/aws/sts"
	"github.com/astronetes/sdk-go/tests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	awssts "github.com/aws/aws-sdk-go-v2/service/sts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS - IAM - Role", func() {
	Describe("Working with IAM Role", func() {
		Context("with default AWS config", func() {
			config, err := config.LoadDefaultConfig(context.TODO())
			Expect(err).To(BeNil())

			client := awsiam.NewFromConfig(config)
			stsClient := awssts.NewFromConfig(config)
			ctx := context.TODO()

			accountArn, err := sts.GetAccountID(ctx, stsClient)
			Expect(err).To(BeNil())
			Expect(accountArn).NotTo(BeNil())

			roleName := tests.RandString("role-test")
			assumeRolePolicyDocument := iam.PolicyDocument{
				Version: "2012-10-17",
				Statement: []iam.PolicyStatement{
					{
						Effect: "Allow",
						Principal: map[string]interface{}{
							"AWS": accountArn,
						},
						Action: []string{"sts:AssumeRole"},
					},
				},
			}

			It("should create the Role", func() {
				request := iam.CreateRoleRequest{
					Name:                     roleName,
					AssumeRolePolicyDocument: assumeRolePolicyDocument,
				}
				err := iam.CreateRole(ctx, client, request)
				Expect(err).To(BeNil())
			})

			It("should fail if the Role already exists", func() {
				request := iam.CreateRoleRequest{
					Name:                     roleName,
					AssumeRolePolicyDocument: assumeRolePolicyDocument,
				}
				err := iam.CreateRole(ctx, client, request)
				Expect(err).NotTo(BeNil())
			})

			It("should get the Role already created", func() {
				role, err := iam.GetRole(ctx, client, roleName)
				Expect(err).To(BeNil())
				Expect(role).NotTo(BeNil())
				Expect(role.RoleName).To(Equal(aws.String(roleName)))
			})

			It("should delete the Role already created", func() {
				err := iam.DeleteRole(ctx, client, roleName)
				Expect(err).To(BeNil())
			})

			It("should not fail if delete the Role that is already deleted", func() {
				err := iam.DeleteRole(ctx, client, roleName)
				Expect(err).To(BeNil())
			})

		})
	})
})
