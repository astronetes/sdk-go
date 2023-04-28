package certificate_test

import (
	"context"
	"fmt"

	"github.com/astronetes/sdk-go/cloud/aws/acm"
	"github.com/astronetes/sdk-go/tests"
	"github.com/aws/aws-sdk-go-v2/config"
	awsacm "github.com/aws/aws-sdk-go-v2/service/acm"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AWS - ACM - Certificate", func() {
	Describe("Working with ACM Certificate", func() {
		Context("with default AWS config", func() {
			config, err := config.LoadDefaultConfig(context.TODO())
			Expect(err).To(BeNil())

			client := awsacm.NewFromConfig(config)
			ctx := context.TODO()

			domainPrefix := tests.RandString("certificate-test")
			domain := fmt.Sprintf("%s.astrokube.es", domainPrefix)
			certificateArn := ""

			It("should create the Certificate", func() {
				request := acm.RequestCertificateRequest{
					Domain: domain,
				}
				certificate, err := acm.RequestCertificate(ctx, client, request)
				Expect(err).To(BeNil())
				Expect(certificate).NotTo(BeNil())
				Expect(certificate.Arn).NotTo(BeNil())
				certificateArn = *certificate.Arn
			})

			It("should not fail if the Certificate already exists", func() {
				request := acm.RequestCertificateRequest{
					Domain: domain,
				}
				certificate, err := acm.RequestCertificate(ctx, client, request)
				Expect(err).To(BeNil())
				Expect(certificate).NotTo(BeNil())
			})

			It("should get the Certificate already created", func() {
				certificate, err := acm.GetCertificate(ctx, client, certificateArn)
				Expect(err).To(BeNil())
				Expect(certificate).NotTo(BeNil())
			})

			It("should delete the Certificate already created", func() {
				err := acm.DeleteCertificate(ctx, client, certificateArn)
				Expect(err).To(BeNil())
			})

			It("should not fail if delete the Certificate that is already deleted", func() {
				err := acm.DeleteCertificate(ctx, client, certificateArn)
				Expect(err).To(BeNil())
			})

		})
	})
})
