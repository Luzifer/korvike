package functions

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	registerFunction("mustVault", tplVaultKeyFetch)

	registerFunction("vault", func(name string, key string) string {
		v, err := tplVaultKeyFetch(name, key)
		if err != nil {
			return ""
		}
		return v
	})
}

func tplVaultKeyFetch(name string, key string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("path is not set")
	}
	if key == "" {
		return "", fmt.Errorf("key is not set")
	}

	client, err := vaultClientFromEnvOrFile()
	if err != nil {
		return "", err
	}

	secret, err := client.Logical().Read(name)
	if err != nil {
		return "", fmt.Errorf("reading secret: %s", err)
	}

	if secret != nil && secret.Data != nil {
		if val, ok := secret.Data[key].(string); ok {
			return val, nil
		}
	}

	return "", fmt.Errorf("requested value %q in key %q was not found in Vault and no default was set", key, name)
}

func vaultClientFromEnvOrFile() (*api.Client, error) {
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv(api.EnvVaultAddress),
	})
	if err != nil {
		return nil, fmt.Errorf("creating Vault client: %w", err)
	}

	switch {
	case os.Getenv(api.EnvVaultToken) != "":
		client.SetToken(os.Getenv(api.EnvVaultToken))

	case os.Getenv("VAULT_ROLE_ID") != "":
		if err = setVaultTokenFromRoleID(client); err != nil {
			return nil, fmt.Errorf("fetch VAULT_TOKEN: %w", err)
		}

	case hasTokenFile():
		if f, err := homedir.Expand("~/.vault-token"); err == nil {
			//#nosec:G304 // File is a well-known file location
			if b, err := os.ReadFile(f); err == nil {
				client.SetToken(string(b))
			}
		}

	default:
		return nil, errors.New("neither VAULT_TOKEN nor VAULT_ROLE_ID was found in ENV and no ~/.vault-token file is present")
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
		return fmt.Errorf("fetching authentication token: %w", lserr)
	}

	client.SetToken(loginSecret.Auth.ClientToken)
	return nil
}
