package acm

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
)

func DeleteCertificate(ctx context.Context, client *acm.Client, certificateARN string) DeleteCertificateResponse {
	if certificateARN == "" {
		return DeleteCertificateResponse{
			error: fmt.Errorf("missing required certificate ARN"),
		}
	}
	input := &acm.DeleteCertificateInput{
		CertificateArn: aws.String(certificateARN),
	}
	response, err := client.DeleteCertificate(ctx, input)
	return DeleteCertificateResponse{
		response: response,
		error:    err,
	}
}

type DeleteCertificateResponse struct {
	response *acm.DeleteCertificateOutput
	error
}

func (r DeleteCertificateResponse) Error() error {
	return r.error
}

func (r DeleteCertificateResponse) Response() *acm.DeleteCertificateOutput {
	return r.response
}
