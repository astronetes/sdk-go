package loggroup_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/astronetes/sdk-go/cloud/aws/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awscloudwatchlogs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

var _ = Describe("AWS - Cloudwatchlogs - LogGroup", func() {
	Describe("Working with LogGroup", func() {
		Context("with default AWS config", func() {
			config, err := config.LoadDefaultConfig(context.TODO())
			Expect(err).To(BeNil())

			client := awscloudwatchlogs.NewFromConfig(config)
			ctx := context.TODO()

			logGroupName := "test"
			retentionDays := int32(3)
			newRetentionDays := int32(7)

			It("should create the LogGroup", func() {
				request := cloudwatchlogs.CreateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: retentionDays,
				}
				err := cloudwatchlogs.CreateLogGroup(ctx, client, request)
				Expect(err).To(BeNil())
			})

			It("should fail if the LogGroup already exists", func() {
				request := cloudwatchlogs.CreateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: retentionDays,
				}
				err := cloudwatchlogs.CreateLogGroup(ctx, client, request)
				Expect(err).NotTo(BeNil())
			})

			It("should not fail if the LogGroup already exists", func() {
				request := cloudwatchlogs.CreateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: retentionDays,
				}
				err := cloudwatchlogs.CreateOrUpdateLogGroup(ctx, client, request)
				Expect(err).To(BeNil())
			})

			It("should get the LogGroup already created", func() {
				request := cloudwatchlogs.GetLogGroupRequest{
					Name: logGroupName,
				}
				result, err := cloudwatchlogs.GetLogGroup(ctx, client, request)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.LogGroup()).NotTo(BeNil())
				Expect(result.LogGroup().RetentionInDays).To(Equal(aws.Int32(retentionDays)))
			})

			It("should update the LogGroup", func() {
				request := cloudwatchlogs.UpdateLogGroupRequest{
					Name:          logGroupName,
					RetentionDays: newRetentionDays,
				}
				err := cloudwatchlogs.UpdateLogGroup(ctx, client, request)
				Expect(err).To(BeNil())

				getRequest := cloudwatchlogs.GetLogGroupRequest{
					Name: logGroupName,
				}
				result, err := cloudwatchlogs.GetLogGroup(ctx, client, getRequest)
				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.LogGroup()).NotTo(BeNil())
				Expect(result.LogGroup().RetentionInDays).To(Equal(aws.Int32(newRetentionDays)))
			})

			It("should delete the LogGroup already created", func() {
				request := cloudwatchlogs.DeleteLogGroupRequest{
					Name: logGroupName,
				}
				err := cloudwatchlogs.DeleteLogGroup(ctx, client, request)
				Expect(err).To(BeNil())
			})

			It("should not fail if delete the LogGroup that is already deleted", func() {
				request := cloudwatchlogs.DeleteLogGroupRequest{
					Name: logGroupName,
				}
				err := cloudwatchlogs.DeleteLogGroup(ctx, client, request)
				Expect(err).To(BeNil())
			})
		})
	})
})
