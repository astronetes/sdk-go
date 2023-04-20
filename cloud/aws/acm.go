package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/astronetes/sdk-go/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	acmTypes "github.com/aws/aws-sdk-go-v2/service/acm/types"
)

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

func (c *acmClient) RequestCertificate(ctx context.Context, domain string) (string, error) {
	logger := log.FromContext(ctx)
	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		logger.V(log.Error).Info("invalid domain, It should have at least two blocks")
		return "", fmt.Errorf("invalid domain, It should have at least two blocks")
	}
	validationDomain := strings.Join(domainParts[len(domainParts)-2:], ".")
	request := &acm.RequestCertificateInput{
		DomainName: aws.String(domain),
		DomainValidationOptions: []acmTypes.DomainValidationOption{
			{
				DomainName:       aws.String(domain),
				ValidationDomain: aws.String(validationDomain),
			},
		},
		ValidationMethod: acmTypes.ValidationMethodDns,
	}
	response, err := c.client.RequestCertificate(ctx, request)
	if err != nil {
		logger.V(log.Error).Info("unexpected error requesting the certificate: '%w'", err)
		return "", err
	}
	return *response.CertificateArn, nil
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
