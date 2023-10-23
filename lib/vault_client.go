package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	Config    Target
	apiClient *vault.Client
}

type Target struct {
	Hostname   string
	Token      string
	KV_version int
}

func (c *VaultClient) init() error {
	vaultConfig := vault.DefaultConfig()
	err := vaultConfig.ConfigureTLS(
		&vault.TLSConfig{
			CACert:        "",
			CACertBytes:   []byte{},
			CAPath:        "",
			ClientCert:    "",
			ClientKey:     "",
			TLSServerName: c.Config.Hostname,
			Insecure:      true,
		},
	)
	if err != nil {
		log.Fatalf("%v", err)
	}
	vaultConfig.Address = c.Config.Hostname
	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}
	vaultClient.SetToken(c.Config.Token)
	c.apiClient = vaultClient

	return nil
}
func (c VaultClient) WriteAllSecrets(basepath string, secrets map[string]string, stripNFirstElementsFromPath int) error {
	c.init()

	for k, v := range secrets {
		k = strings.Join(strings.Split(k, "/")[stripNFirstElementsFromPath:], "/")
		fmt.Printf("Migrating %s\n", k)
		var writePath string
		if c.Config.KV_version != 2 {
			writePath = fmt.Sprintf("%s/%s", basepath, k)
		} else {
			writePath = fmt.Sprintf("%s/data/%s", basepath, k)
		}
		var value map[string]interface{}
		json.Unmarshal([]byte(v), &value)
		data := map[string]interface{}{
			"data": value,
		}
		_, err := c.apiClient.Logical().Write(writePath, data)
		if err != nil {
			log.Fatalf("%s", err)
		}
	}

	return nil
}
func (c VaultClient) GetAllSecrets(basepath string, key string) map[string]string {
	c.init()

	if c.Config.KV_version == 2 {
		return *c.recurse(basepath, key, map[string]string{})
	}
	return *c.recurse(basepath, key, map[string]string{})
}

func (c VaultClient) recurse(mountpoint, key string, secretsMap map[string]string) *map[string]string {
	var path string
	if c.Config.KV_version != 2 {
		path = mountpoint
	} else {
		path = fmt.Sprintf("%s/metadata", mountpoint)
	}
	if key != "" {
		path = fmt.Sprintf("%s/%s", path, key)
	}
	var secrets *vault.Secret
	secrets, err := c.apiClient.Logical().List(path)
	if err != nil {
		log.Fatalf("%v", err)
	}
	keys := secrets.Data["keys"].([]interface{})
	for _, v := range keys {
		k := v.(string)
		if k[len(k)-1:] == "/" {
			secretsMap = *c.recurse(mountpoint, fmt.Sprintf("%s%s", key, k), secretsMap)
		} else {
			var secretPath string
			fullPath := fmt.Sprintf("%s%s", key, k)
			fmt.Printf("found secret at %s\n", fullPath)
			if c.Config.KV_version == 2 {
				secretPath = fmt.Sprintf("%s/data/%s", strings.Split(mountpoint, "/")[0], fullPath)
			} else {
				secretPath = fmt.Sprintf("%s/%s", strings.Split(mountpoint, "/")[0], fullPath)
			}
			secret, err := c.apiClient.Logical().Read(secretPath)
			if err != nil {
				log.Fatalf("%v", err)
			}
			var jsonString []byte
			if c.Config.KV_version == 2 {
				jsonString, err = json.Marshal(secret.Data["data"])
			} else {
				jsonString, err = json.Marshal(secret.Data)
			}
			if err != nil {
				log.Fatalf("%v", err)
			}
			secretsMap[fullPath] = string(jsonString)
		}
	}

	return &secretsMap
}
