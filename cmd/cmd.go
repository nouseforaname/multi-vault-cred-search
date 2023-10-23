package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nouseforaname/vault-cred-matcher/actions"
	"github.com/nouseforaname/vault-cred-matcher/actions/migrate"
	"github.com/nouseforaname/vault-cred-matcher/actions/search"
	"github.com/nouseforaname/vault-cred-matcher/config"
	"github.com/nouseforaname/vault-cred-matcher/flags"
	"gopkg.in/yaml.v2"
)

type cmd struct {
	configValues flags.FlagValues
	Config       config.Conf
}

func (c cmd) GetConfigValues() flags.FlagValues {
	return c.configValues
}
func (c cmd) GetAction(actionType string) (actions.Action, error) {
	switch actionType {
	case "search":
		return &search.SearchVaultAction{
			Config:        c.Config,
			SourceSecrets: map[string]string{},
			TargetSecrets: map[string]string{},
		}, nil
	case "migrate":
		return &migrate.MigrateVaultAction{
			Config:        c.Config,
			SourceSecrets: map[string]string{},
			TargetSecrets: map[string]string{},
		}, nil
	}

	return nil, fmt.Errorf("ActionNotImplemented")
}
func NewCmd() (*cmd, error) {

	configValues := flags.Get()
	yamlFile, err := os.ReadFile(configValues.ConfigPath)

	if err != nil {
		log.Printf("FileNotFound: %v", configValues.ConfigPath)
		return nil, fmt.Errorf("FileNotFound: %v", configValues.ConfigPath)
	}

	conf := config.Conf{}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return &cmd{
		configValues: configValues,
		Config:       conf,
	}, nil
}
