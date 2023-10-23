package flags_test

import (
	"github.com/nouseforaname/vault-cred-matcher/flags"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flags", func() {
	It("parses a config", func() {
		defaultAction := "search"
		defaultConfigPath := "./config.yml"
		Expect(flags.Action).To(Equal(&defaultAction))
		Expect(flags.ConfigPath).To(Equal(&defaultConfigPath))
	})

})
