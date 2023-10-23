package config

import (
	"github.com/nouseforaname/vault-cred-matcher/lib"
)

type Conf struct {
	Src    lib.Target `yaml:"src"`
	Target lib.Target `yaml:"target"`
}
