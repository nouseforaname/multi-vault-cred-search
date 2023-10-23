package cmd_test

import (
	"github.com/nouseforaname/vault-cred-matcher/cmd"
	"github.com/nouseforaname/vault-cred-matcher/flags"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cmd", func() {

	Describe("no config.yml found", func() {
		It("gives a file not found error when there is no defaulut config.yml", func() {
			cmd, err := cmd.NewCmd()
			Expect(err).To(Not(BeNil()))
			Expect(cmd).To(BeNil())
			Expect(err).Should(MatchError(ContainSubstring("FileNotFound")))
		})
	})
	Describe("with a config.yml found", func() {
		BeforeEach(func() {
			existingConfig := "assets/config.yml"
			flags.ConfigPath = &existingConfig
		})
		It("creates a command when the path to config.yml can be found", func() {
			cmd, err := cmd.NewCmd()
			Expect(err).To(BeNil())
			Expect(cmd).To(Not(BeNil()))
		})
		It("returns the right values", func() {
			cmd, err := cmd.NewCmd()
			Expect(err).To(BeNil())
			Expect(cmd).To(Not(BeNil()))
			Expect(cmd.GetConfigValues()).To(Equal(flags.FlagValues{
				Action:         "search",
				ConfigPath:     "assets/config.yml",
				SecretListFile: "",
			}))
		})
	})
})
