package aws

import (
	"context"

	acm2 "github.com/astronetes/sdk-go/cloud/aws/acm"
	"github.com/astronetes/sdk-go/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
)

func RequestCertificate(ctx context.Context, client *acm.Client, domain string) acm2.RequestCertificateResponse {
	return acm2.RequestCertificate(ctx, client, domain)
}

type ACMClient interface {
	IsCertificateCreated(ctx context.Context, certificatedID string) (bool, error)
	IsCertificateInIssuedState(ctx context.Context, certificatedID string) (bool, error)
	RequestCertificate(ctx context.Context, domain string) (string, error)
	DeleteCertificate(ctx context.Context, certificateID string) error
	GetDomainValidationOptionsForCertificate(ctx context.Context, certificateID string) ([]acmTypes.DomainValidation, error)
}

type acmClient struct {
	client *acm.Client
}

func (c *acmClient) IsCertificateIsCreated(ctx context.Context, certificatedID string) (bool, error) {
	logger := log.FromContext(ctx)
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificatedID),
	}
	output, err := c.client.DescribeCertificate(ctx, input)
	if err != nil {
		logger.V(log.Error).Info("unexpected error obtaining the details for the certificate: '%w'", err)
		return false, err
	}

	return output != nil, nil
}

func (c *acmClient) IsCertificateInIssuedState(ctx context.Context, certificatedID string) (bool, error) {
	logger := log.FromContext(ctx)
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificatedID),
	}
	output, err := c.client.DescribeCertificate(ctx, input)
	if err != nil {
		logger.V(log.Error).Info("unexpected error obtaining the details for the certificate: '%w'", err)
		return false, err
	}

	return output.Certificate.Status == acmTypes.CertificateStatusIssued, nil
}

func (c *acmClient) GetDomainValidationOptionsForCertificate(ctx context.Context, certificateID string) ([]acmTypes.DomainValidation, error) {
	logger := log.FromContext(ctx)
	input := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateID),
	}

	response, err := c.client.DescribeCertificate(ctx, input)
	if err != nil {
		logger.V(log.Error).Info("Unable to get certificate details: '%w", err)
	}

	return response.Certificate.DomainValidationOptions, nil
}
