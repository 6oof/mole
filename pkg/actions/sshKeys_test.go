package actions

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zulubit/mole/pkg/consts"
)

var tk string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJ+CEFyKmeys1gL/EsBhKIf58Rny7kQUZe5XbktCZl8HC+Oq497viBT3fH+1Nr3pimEVEWwQtgtKBofvwmWrUfFc0ABp+ofxrwh2uRqpClSzfFBzh5GlovliOnNEP73g5jh9DZCQBJOcNUljBi5mn1EHXe5l3np7Qae8IKy2YdHihb6RGFWkl8WPJSyT7ZzHq9FLe3O3ks3+Tr0/f1rTfoflJGeG4kfJawNaRmfLsjWTiFhiFXrtFMKtQ/yK6/RODCfe7rq4gLKFSQddKXDj2oftXhazSnXApfJmu9nMpRnVEXu6EBdFMWSeZTwp814FL15HtMvz4aTjqpFxJ6vzcX"

var tkw string = "ADAQABAAABAQDJ+CEFyKmeys1gL/EsBhKIf58Rny7kQUZe5XbktCZl8HC+Oq497viBT3fH+1Nr3pimEVEWwQtgtKBofvwmWrUfFc0ABp+ofxrwh2uRqpClSzfFBzh5GlovliOnNEP73g5jh9DZCQBJOcNUljBi5mn1EHXe5l3np7Qae8IKy2YdHihb6RGFWkl8WPJSyT7ZzHq9FLe3O3ks3+Tr0/f1rTfoflJGeG4kfJawNaRmfLsjWTiFhiFXrtFMKtQ/yK6/RODCfe7rq4gLKFSQddKXDj2oftXhazSnXApfJmu9nMpRnVEXu6EBdFMWSeZTwp814FL15HtMvz4aTjqpFxJ6vzcX"

func TestFindOrCreateDeployKey(t *testing.T) {
	consts.Testing = true

	consts.BasePath = os.TempDir()
	defer os.RemoveAll(consts.BasePath)

	key, err := FindOrCreateDeployKey()
	assert.Nil(t, err, "deploy key was created")
	assert.Contains(t, key, "ssh-rsa", "something resembling a key is returned")
}

func TestAddAuthorizedKeys(t *testing.T) {
	consts.Testing = true

	consts.BasePath = os.TempDir()
	defer os.RemoveAll(consts.BasePath)

	err := AddAuthorizedKeys(tk, "test-key")
	assert.Nil(t, err, "deploy key was created")

	f, _ := os.ReadFile(path.Join(consts.GetBasePath(), ".ssh", "authorized_keys"))
	assert.Contains(t, string(f), tk, "key should be in the file")

	err = AddAuthorizedKeys(tk, "test-key")
	assert.ErrorContains(t, err, "public key already exists in authorized_keys", "error should be correct")

	err = AddAuthorizedKeys(tkw, "test-key")
	assert.ErrorContains(t, err, "invalid public key format", "error should be correct")
}

func TestFindOrCreateActionsKey(t *testing.T) {
	consts.Testing = true

	consts.BasePath = os.TempDir()
	defer os.RemoveAll(consts.BasePath)

	privateKey, err := FindOrCreateActionsKey()
	assert.Nil(t, err, "actions key was created")
	assert.Contains(t, privateKey, "--BEGIN", "private key should be returned as string")

	authorizedKeysPath := path.Join(consts.GetBasePath(), ".ssh", "authorized_keys")
	content, err := os.ReadFile(authorizedKeysPath)
	assert.Nil(t, err, "authorized_keys file should exist")

	actionsKeyPub := path.Join(consts.GetBasePath(), ".ssh", "actions_rsa.pub")
	contentAKP, err := os.ReadFile(actionsKeyPub)
	assert.Nil(t, err, "actions_rsa.pub file should exist")

	// Verify the public key (ignoring comments in authorized_keys)
	assert.Contains(t, string(content), string(contentAKP), "actions public key should be in the authorized_keys file")
}

func TestEnsureShhDirectory(t *testing.T) {
	consts.Testing = true

	consts.BasePath = os.TempDir()
	defer os.RemoveAll(consts.BasePath)

	err := ensureShhDirectory()
	assert.Nil(t, err, "ssh directory should be ensured without errors")

	sshPath := path.Join(consts.GetBasePath(), ".ssh")
	authorizedKeysPath := path.Join(sshPath, "authorized_keys")

	_, err = os.Stat(sshPath)
	assert.Nil(t, err, ".ssh directory should exist")

	_, err = os.Stat(authorizedKeysPath)
	assert.Nil(t, err, "authorized_keys file should exist")
}

func TestCheckAuthorizedExists(t *testing.T) {
	consts.Testing = true

	consts.BasePath = os.TempDir()
	defer os.RemoveAll(consts.BasePath)

	err := AddAuthorizedKeys(tk, "test-key")
	assert.Nil(t, err, "deploy key should be added without errors")

	_, err = checkAuthorizedExists(tk)
	assert.NotNil(t, err, "key should exist in authorized_keys")

	_, err = checkAuthorizedExists(tkw)
	assert.ErrorContains(t, err, "invalid public key format", "error should be correct")
}
