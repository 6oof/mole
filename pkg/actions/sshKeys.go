package actions

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/keygen"
	"github.com/zulubit/mole/pkg/consts"

	"golang.org/x/crypto/ssh"
)

// FindOrCreateDeployKey creates a new SSH deploy key if one does not already exist,
// and returns the authorized key string representation.
func FindOrCreateDeployKey() (string, error) {
	deployKeyPath := path.Join(consts.GetBasePath(), ".ssh", "id_rsa")

	kp, err := keygen.New(deployKeyPath, keygen.WithKeyType(keygen.RSA), keygen.WithBitSize(2048), keygen.WithWrite())
	if err != nil {
		return "", fmt.Errorf("error creating SSH key pair: %v", err)
	}

	return kp.AuthorizedKey(), nil
}

// FindOrCreateActionsKey creates an SSH key pair specifically for actions if it does not already exist.
// It adds the public key to the authorized_keys file and returns the private key as a string.
func FindOrCreateActionsKey() (string, error) {
	deployKeyPath := path.Join(consts.GetBasePath(), ".ssh", "actions_rsa")

	// Generate a new key pair
	_, err := keygen.New(deployKeyPath, keygen.WithKeyType(keygen.RSA), keygen.WithBitSize(2048), keygen.WithWrite())
	if err != nil {
		return "", fmt.Errorf("error creating SSH key pair: %v", err)
	}

	// Ensure the .ssh directory exists
	err = ensureShhDirectory()
	if err != nil {
		return "", err
	}

	// Read the public key from the generated file
	publicKeyPath := deployKeyPath + ".pub"
	publicKeyContent, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading public key: %v", err)
	}

	// Add the correct public key to authorized_keys if it doesn't already exist
	_, err = checkAuthorizedExists(string(publicKeyContent))
	if err == nil {
		err = AddAuthorizedKeys(string(publicKeyContent), "Key used for Actions")
		if err != nil {
			return "", err
		}
	}

	// Read and return the private key
	privateKeyContent, err := os.ReadFile(deployKeyPath)
	if err != nil {
		return "", fmt.Errorf("error reading private key: %v", err)
	}

	return string(privateKeyContent), nil
}

// checkAuthorizedExists checks if a given public key exists in the authorized_keys file.
// Returns the existing keys if successful, or an error otherwise.
func checkAuthorizedExists(publicKey string) (string, error) {
	authorizedKeysPath := path.Join(consts.GetBasePath(), ".ssh", "authorized_keys")

	existingKeys, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		return "", fmt.Errorf("failed to read authorized_keys file: %v", err)
	}

	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey)); err != nil {
		return "", fmt.Errorf("invalid public key format: %v", err)
	}

	if strings.Contains(string(existingKeys), publicKey) {
		return "", fmt.Errorf("public key already exists in authorized_keys")
	}

	return string(existingKeys), nil
}

// ensureShhDirectory ensures the .ssh directory and the authorized_keys file exist, creating them if necessary.
func ensureShhDirectory() error {
	authorizedKeysPath := path.Join(consts.GetBasePath(), ".ssh", "authorized_keys")

	if err := os.MkdirAll(path.Join(consts.GetBasePath(), ".ssh"), 0700); err != nil {
		return fmt.Errorf("failed to create project volume directory: %w", err)
	}

	if _, err := os.Stat(authorizedKeysPath); os.IsNotExist(err) {
		if err := os.WriteFile(authorizedKeysPath, []byte{}, 0600); err != nil {
			return fmt.Errorf("failed to create authorized_keys file: %v", err)
		}
	}

	return nil
}

// AddAuthorizedKeys appends a given public key to the authorized_keys file,
// validating its format and ensuring it is not already present.
func AddAuthorizedKeys(publicKey, comment string) error {

	authorizedKeysPath := path.Join(consts.GetBasePath(), ".ssh", "authorized_keys")

	err := ensureShhDirectory()
	if err != nil {
		return err
	}

	existingKeys, err := checkAuthorizedExists(publicKey)
	if err != nil {
		return err
	}

	// Format the new entry with a comment, key, and an empty line
	newEntry := fmt.Sprintf("# %s\n%s\n\n", comment, publicKey)
	updatedKeys := existingKeys + newEntry

	// Write back the updated keys to the file
	if err := os.WriteFile(authorizedKeysPath, []byte(updatedKeys), 0600); err != nil {
		return fmt.Errorf("failed to append public key to authorized_keys: %v", err)
	}

	return nil
}
