package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
	term "github.com/nsf/termbox-go"
	"gopkg.in/yaml.v2"
)

type target struct {
	Addr  string `yaml:"addr"`
	Token string `yaml:"token"`
}
type conf struct {
	Src    target `yaml:"src"`
	Target target `yaml:"target"`
}

func main() {

	yamlFile, err := os.ReadFile("./config.yml")
	if err != nil {
		log.Fatalf("%v", err)
	}

	conf := &conf{}

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("%v", err)
	}
	secretsSrc := make(map[string]string)
	secretsTarget := make(map[string]string)

	if _, err := os.Stat("./source-data.json"); err == nil {
		fmt.Println("Source dump exists, loading export")
		bytes, err := os.ReadFile("./source-data.json")
		if err != nil {
			log.Fatalf("%v", err)
		}
		json.Unmarshal(bytes, &secretsSrc)
	} else {
		srcConfig := vault.DefaultConfig()
		err = srcConfig.ConfigureTLS(
			&vault.TLSConfig{
				Insecure:      true,
				TLSServerName: conf.Src.Addr,
			},
		)
		if err != nil {
			log.Fatalf("%v", err)
		}
		srcConfig.Address = conf.Src.Addr

		srcClient, err := vault.NewClient(srcConfig)
		if err != nil {
			log.Fatalf("unable to initialize Vault client: %v", err)
		}
		srcClient.SetToken(conf.Src.Token)

		secretsSrc["source_vault"] = srcConfig.Address

		secretsSrc = recurse("concourse", "", srcClient, secretsSrc)
		jsonString, err := json.Marshal(secretsSrc)
		os.WriteFile("./source-data.json", []byte(jsonString), 0644)

		if err != nil {
			log.Fatalf("unable to write src vault export: %v", err)
		}
	}

	if _, err := os.Stat("./target-data.json"); err == nil {
		fmt.Println("Target dump exists, loading export")
		bytesTarget, err := os.ReadFile("./target-data.json")
		if err != nil {
			log.Fatalf("%v", err)
		}
		json.Unmarshal(bytesTarget, &secretsTarget)
	} else {
		targetConfig := vault.DefaultConfig()
		err = targetConfig.ConfigureTLS(
			&vault.TLSConfig{
				Insecure:      true,
				TLSServerName: conf.Target.Addr,
			},
		)
		if err != nil {
			log.Fatalf("%v", err)
		}
		targetConfig.Address = conf.Target.Addr

		targetClient, err := vault.NewClient(targetConfig)
		if err != nil {
			log.Fatalf("unable to initialize Vault client: %v", err)
		}
		targetClient.SetToken(conf.Target.Token)

		secretsTarget["source_vault"] = targetConfig.Address

		secretsTarget = recurse("runway_concourse", "cryogenics/", targetClient, secretsTarget)
		jsonString, err := json.Marshal(secretsTarget)
		os.WriteFile("./target-data.json", []byte(jsonString), 0644)

		if err != nil {
			log.Fatalf("unable to write target vault export: %v", err)
		}
	}

	fmt.Println("You can start searching. Hit CTRL-C to stop")
	findInMapValues(secretsSrc, secretsTarget)
}
func reset() {
	term.Sync() // cosmestic purpose
}

func findInMapValues(src, target map[string]string) {
	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()
	searchString := ""
	fmt.Println("Hit enter to search for: " + searchString)
keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEnter:
				reset()
				fmt.Println("searching for " + searchString)
				srcHits := make([]string, 0)
				targetHits := make([]string, 0)
				for k, v := range src {
					if strings.Contains(v, searchString) {
						srcHits = append(srcHits, k)
					}
				}
				for k, v := range target {
					if strings.Contains(v, searchString) {
						targetHits = append(targetHits, k)
					}
				}
				fmt.Printf("%s was found in srcSecrets: \n", searchString)
				fmt.Println(strings.Join(srcHits, "\n  -"))
				fmt.Printf("%s was found in targetSecrets: \n", searchString)
				fmt.Println(strings.Join(targetHits, "\n  -"))
				searchString = ""
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyBackspace2:
				reset()
				if len(searchString) > 0 {
					searchString = searchString[:len(searchString)-1]
				}
				fmt.Println("Back enter to search for: " + searchString)

			case term.KeyBackspace:
				reset()
				if len(searchString) > 0 {
					searchString = searchString[:len(searchString)-1]
				}
				fmt.Println("Back enter to search for: " + searchString)
			default:
				// we only want to read a single character or one key pressed event
				reset()
				searchString += string(ev.Ch)
				fmt.Println("Hit enter to search for: " + searchString)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

func recurse(basepath, key string, client *vault.Client, secretsMap map[string]string) map[string]string {
	path := basepath
	if key != "" {
		path = fmt.Sprintf("%s/%s", basepath, key)
	}
	secrets, err := client.Logical().List(path)

	if err != nil {
		log.Fatalf("%v", err)
	}
	keys := secrets.Data["keys"].([]interface{})
	for _, v := range keys {
		k := v.(string)
		if k[len(k)-1:] == "/" {
			secretsMap = recurse(basepath, fmt.Sprintf("%s%s", key, k), client, secretsMap)
		} else {
			secretPath := fmt.Sprintf("%s/%s%s", basepath, key, k)
			secret, err := client.Logical().Read(secretPath)
			if err != nil {
				log.Fatalf("%v", err)
			}
			jsonString, err := json.Marshal(secret)
			if err != nil {
				log.Fatalf("%v", err)
			}
			secretsMap[secretPath] = string(jsonString)
			fmt.Printf("found secret in %s: %s\n", secretsMap["source_vault"], secretPath)

		}
	}

	return secretsMap
}
