package role_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAwsIamRole(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS - IAM - Role Suite")
}
