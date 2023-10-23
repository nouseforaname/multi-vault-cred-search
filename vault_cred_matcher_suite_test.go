package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestVaultCredMatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VaultCredMatcher Suite")
}
