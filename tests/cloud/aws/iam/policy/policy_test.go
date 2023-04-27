package policy_test

import (
	"context"

	"github.com/astronetes/sdk-go/cloud/aws/iam"
	"github.com/astronetes/sdk-go/tests"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS - IAM - Policy", func() {
	Describe("Working with IAM Policy", func() {
		Context("with default AWS config", func() {
			config, err := config.LoadDefaultConfig(context.TODO())
			Expect(err).To(BeNil())

			client := awsiam.NewFromConfig(config)
			ctx := context.TODO()

			policyName := tests.RandString("policy-test")
			policyDescription := "Policy test"
			policyArn := ""

			It("should create the Policy", func() {
				request := iam.CreatePolicyRequest{
					Name:        policyName,
					Description: policyDescription,
					Document: iam.PolicyDocument{
						Version: "2012-10-17",
						Statement: []iam.PolicyStatement{
							{
								Effect:   "Allow",
								Action:   []string{"*"},
								Resource: aws.String("*"),
							},
						},
					},
				}
				policy, err := iam.CreatePolicy(ctx, client, request)
				Expect(err).To(BeNil())
				Expect(policy).NotTo(BeNil())
				Expect(policy.Arn).NotTo(BeNil())
				policyArn = *policy.Arn
			})

			It("should fail if the Policy already exists", func() {
				request := iam.CreatePolicyRequest{
					Name: policyName,
				}
				policy, err := iam.CreatePolicy(ctx, client, request)
				Expect(err).NotTo(BeNil())
				Expect(policy).To(BeNil())
			})

			It("should get the Policy already created", func() {
				policy, err := iam.GetPolicy(ctx, client, policyArn)
				Expect(err).To(BeNil())
				Expect(policy).NotTo(BeNil())
				Expect(policy.PolicyName).To(Equal(aws.String(policyName)))
				Expect(policy.Description).To(Equal(aws.String(policyDescription)))
			})

			It("should delete the Policy already created", func() {
				err := iam.DeletePolicy(ctx, client, policyArn)
				Expect(err).To(BeNil())
			})

			It("should not fail if delete the Policy that is already deleted", func() {
				err := iam.DeletePolicy(ctx, client, policyArn)
				Expect(err).To(BeNil())
			})

		})
	})
})
