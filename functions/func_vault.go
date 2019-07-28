package functions

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/vault/api"
	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	registerFunction("vault", func(name string, v ...string) (interface{}, error) {
		if name == "" {
			return nil, fmt.Errorf("Path is not set")
		}
		if len(v) < 1 {
			return nil, fmt.Errorf("Key is not set")
		}

		client, err := vaultClientFromEnvOrFile()
		if err != nil {
			return nil, err
		}

		secret, err := client.Logical().Read(name)
		if err != nil {
			return nil, err
		}

		if secret != nil && secret.Data != nil {
			if val, ok := secret.Data[v[0]]; ok {
				return val, nil
			}
		}

		if len(v) < 2 {
			return nil, fmt.Errorf("Requested value %q in key %q was not found in Vault and no default was set", v[0], name)
		}

		return v[1], nil
	})
}

func vaultClientFromEnvOrFile() (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv(api.EnvVaultAddress),
	})
	if err != nil {
		return nil, err
	}

	switch {

	case os.Getenv(api.EnvVaultToken) != "":
		client.SetToken(os.Getenv(api.EnvVaultToken))

	case os.Getenv("VAULT_ROLE_ID") != "":
		if err = setVaultTokenFromRoleID(client); err != nil {
			return nil, fmt.Errorf("Unable to fetch VAULT_TOKEN: %s", err)
		}

	case hasTokenFile():
		if f, err := homedir.Expand("~/.vault-token"); err == nil {
			if b, err := ioutil.ReadFile(f); err == nil {
				client.SetToken(string(b))
			}
		}

	default:
		return nil, errors.New("Neither VAULT_TOKEN nor VAULT_ROLE_ID was found in ENV and no ~/.vault-token file is present")

	}

	return client, nil
}

func hasTokenFile() bool {
	if f, err := homedir.Expand("~/.vault-token"); err == nil {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}

	return false
}

func setVaultTokenFromRoleID(client *api.Client) error {
	data := map[string]interface{}{
		"role_id": os.Getenv("VAULT_ROLE_ID"),
	}

	if os.Getenv("VAULT_SECRET_ID") != "" {
		data["secret_id"] = os.Getenv("VAULT_SECRET_ID")
	}

	loginSecret, lserr := client.Logical().Write("auth/approle/login", data)
	if lserr != nil || loginSecret.Auth == nil {
		return fmt.Errorf("Unable to fetch authentication token: %s", lserr)
	}

	client.SetToken(loginSecret.Auth.ClientToken)
	return nil
}
