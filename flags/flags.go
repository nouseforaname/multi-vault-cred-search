package flags

import (
	"flag"
)

var SecretListFile = flag.String("f", "", "path to file with line delimited list of secrets to search for")
var Action = flag.String("a", "search", "searching for dupes or migrating contents. valid options search(default) || migrate")
var ConfigPath = flag.String("c", "./config.yml", "searching for dupes or migrating contents. valid options search(default) || migrate")

type FlagValues struct {
	Action         string
	ConfigPath     string
	SecretListFile string
}

func Get() FlagValues {
	flag.Parse()
	return FlagValues{
		Action:         *Action,
		ConfigPath:     *ConfigPath,
		SecretListFile: *SecretListFile,
	}
}
