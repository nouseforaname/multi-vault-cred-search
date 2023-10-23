package migrate

import (
	"log"

	"github.com/nouseforaname/vault-cred-matcher/config"
	"github.com/nouseforaname/vault-cred-matcher/lib"
)

type MigrateVaultAction struct {
	Config            config.Conf
	SourceSecrets     map[string]string
	TargetSecrets     map[string]string
	srcVaultClient    lib.VaultClient
	targetVaultClient lib.VaultClient
}

func (a *MigrateVaultAction) init() {
	a.srcVaultClient = lib.VaultClient{
		Config: a.Config.Src,
	}
	a.targetVaultClient = lib.VaultClient{
		Config: a.Config.Target,
	}
}

func (a *MigrateVaultAction) Result() interface{} {
	return a.TargetSecrets
}
func (a *MigrateVaultAction) Execute() {
	log.Print("Starting Migrate")

	a.init()
	// TODO make stripNFirstElementsFromPath configurable
	a.SourceSecrets = a.srcVaultClient.GetAllSecrets("runway_concourse", "cryogenics/")
	a.targetVaultClient.WriteAllSecrets("secret", a.SourceSecrets, 1)

}
