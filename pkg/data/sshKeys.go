package data

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/charmbracelet/keygen"

	"golang.org/x/crypto/ssh"
)

func FindOrCreateDeployKey() (string, error) {
	deployKeyPath := path.Join(consts.BasePath, ".ssh", "id_rsa")

	kp, err := keygen.New(deployKeyPath, keygen.WithKeyType(keygen.RSA), keygen.WithBitSize(2048), keygen.WithWrite())
	if err != nil {
		return "", fmt.Errorf("error creating SSH key pair: %v", err)
	}

	return kp.AuthorizedKey(), nil
}

func AddAuthorizedKeys(publicKey string) error {
	authorizedKeysPath := filepath.Join(consts.BasePath, ".ssh", "authorized_keys")

	if _, err := os.Stat(authorizedKeysPath); os.IsNotExist(err) {
		if err := os.WriteFile(authorizedKeysPath, []byte{}, 0644); err != nil {
			return fmt.Errorf("failed to create authorized_keys file: %v", err)
		}
	}

	existingKeys, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		return fmt.Errorf("failed to read authorized_keys file: %v", err)
	}

	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey)); err != nil {
		return fmt.Errorf("invalid public key format: %v", err)
	}

	if strings.Contains(string(existingKeys), publicKey) {
		return fmt.Errorf("public key already exists in authorized_keys")
	}

	updatedKeys := string(existingKeys) + publicKey + "\n"
	if err := os.WriteFile(authorizedKeysPath, []byte(updatedKeys), 0600); err != nil {
		return fmt.Errorf("failed to append public key to authorized_keys: %v", err)
	}

	return nil
}