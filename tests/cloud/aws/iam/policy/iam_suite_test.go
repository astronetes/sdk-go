package policy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAwsIamPolicy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS - IAM - Policy Suite")
}
