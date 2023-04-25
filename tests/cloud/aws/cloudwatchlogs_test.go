package aws_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/astronetes/sdk-go/cloud/aws/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awscloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

var _ = Describe("Cloudwatchlogs", func() {
	Describe("Working with LogGroup", func() {
		Context("with default AWS config", func() {
			config, err := config.LoadDefaultConfig(context.TODO())
			Expect(err).To(BeNil())

			client := awscloudwatchlogs.NewFromConfig(config)
			ctx := context.TODO()

			logGroupName := "test"
			retentionDays := int32(3)

			It("should create the LogGroup", func() {
				request := cloudwatchlogs.CreateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: retentionDays,
				}
				result := cloudwatchlogs.CreateLogGroup(ctx, client, request)
				Expect(result.Error()).To(BeNil())
				Expect(result.RetentionDays()).NotTo(BeNil())
				Expect(result.RetentionDays()).To(Equal(request.RetentionDays))
			})

			It("should fail if the LogGroup already exists", func() {
				request := cloudwatchlogs.CreateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: retentionDays,
				}
				result := cloudwatchlogs.CreateLogGroup(ctx, client, request)
				Expect(result.Error()).NotTo(BeNil())
			})

			It("should get the LogGroup already created", func() {
				request := cloudwatchlogs.GetLogGroupRequest{
					Name: logGroupName,
				}
				result := cloudwatchlogs.GetLogGroup(ctx, client, request)
				Expect(result.Error()).To(BeNil())
				Expect(result.LogGroup()).NotTo(BeNil())
				Expect(result.LogGroup().RetentionInDays).To(Equal(aws.Int32(retentionDays)))
			})

			It("should delete the LogGroup already created", func() {
				request := cloudwatchlogs.DeleteLogGroupRequest{
					Name: logGroupName,
				}
				result := cloudwatchlogs.DeleteLogGroup(ctx, client, request)
				Expect(result.Error()).To(BeNil())
			})
		})
	})
})
