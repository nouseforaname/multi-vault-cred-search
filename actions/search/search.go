package search

import (
	"log"

	"github.com/nouseforaname/vault-cred-matcher/config"
	"github.com/nouseforaname/vault-cred-matcher/lib"
)

type SearchVaultAction struct {
	Config            config.Conf
	SourceSecrets     map[string]string
	TargetSecrets     map[string]string
	srcVaultClient    lib.VaultClient
	targetVaultClient lib.VaultClient
}

func (a *SearchVaultAction) init() {
	a.srcVaultClient = lib.VaultClient{
		Config: a.Config.Src,
	}
	a.targetVaultClient = lib.VaultClient{
		Config: a.Config.Target,
	}
}

func (a *SearchVaultAction) Result() interface{} {
	return a.TargetSecrets
}
func (a *SearchVaultAction) Execute() {
	log.Print("Starting Search")

	a.init()

	a.SourceSecrets = a.srcVaultClient.GetAllSecrets("runway_concourse", "cryogenics/")
	a.TargetSecrets = a.targetVaultClient.GetAllSecrets("concourse", "")

}
