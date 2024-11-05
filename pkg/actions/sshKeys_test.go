package actions

import (
	"os"
	"path"
	"testing"

	"github.com/6oof/mole/pkg/consts"
	"github.com/stretchr/testify/assert"
)

var tk string = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDJ+CEFyKmeys1gL/EsBhKIf58Rny7kQUZe5XbktCZl8HC+Oq497viBT3fH+1Nr3pimEVEWwQtgtKBofvwmWrUfFc0ABp+ofxrwh2uRqpClSzfFBzh5GlovliOnNEP73g5jh9DZCQBJOcNUljBi5mn1EHXe5l3np7Qae8IKy2YdHihb6RGFWkl8WPJSyT7ZzHq9FLe3O3ks3+Tr0/f1rTfoflJGeG4kfJawNaRmfLsjWTiFhiFXrtFMKtQ/yK6/RODCfe7rq4gLKFSQddKXDj2oftXhazSnXApfJmu9nMpRnVEXu6EBdFMWSeZTwp814FL15HtMvz4aTjqpFxJ6vzcX"

var tkw string = "ADAQABAAABAQDJ+CEFyKmeys1gL/EsBhKIf58Rny7kQUZe5XbktCZl8HC+Oq497viBT3fH+1Nr3pimEVEWwQtgtKBofvwmWrUfFc0ABp+ofxrwh2uRqpClSzfFBzh5GlovliOnNEP73g5jh9DZCQBJOcNUljBi5mn1EHXe5l3np7Qae8IKy2YdHihb6RGFWkl8WPJSyT7ZzHq9FLe3O3ks3+Tr0/f1rTfoflJGeG4kfJawNaRmfLsjWTiFhiFXrtFMKtQ/yK6/RODCfe7rq4gLKFSQddKXDj2oftXhazSnXApfJmu9nMpRnVEXu6EBdFMWSeZTwp814FL15HtMvz4aTjqpFxJ6vzcX"

func TestFindOrCreateDeployKey(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	key, err := FindOrCreateDeployKey()
	assert.Nil(t, err, "deploy key was created")
	assert.Contains(t, key, "ssh-rsa", "something resembling a key is returned")
}

func TestAddAuthorizedKeys(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	err := AddAuthorizedKeys(tk)

	assert.Nil(t, err, "deploy key was created")

	f, _ := os.ReadFile(path.Join(tmp, ".ssh", "authorized_keys"))
	assert.Contains(t, string(f), tk, "key shold be in the file")

	err = AddAuthorizedKeys(tk)
	assert.ErrorContains(t, err, "public key already exists in authorized_keys", "error should be correct")

	err = AddAuthorizedKeys(tkw)
	assert.ErrorContains(t, err, "invalid public key format", "error should be correct")
}
