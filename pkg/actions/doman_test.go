package actions

import (
	"os"
	"path"
	"testing"

	"github.com/6oof/mole/pkg/consts"
	"github.com/stretchr/testify/assert"
)

var successDomainProxy = "www.test.com {\n    redir https://test.com{uri}\n}\n\ntest.com {\n    reverse_proxy 127.0.0.1:3000\n}"
var successDomainStatic = "www.test.com {\n    redir https://test.com{uri}\n}\n\ntest.com {\n    root * /home/projects/test/\n    file_server\n}"

func TestAddDomainProxy(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	p := Project{
		Name: "test",
	}
	err := addProject(p)
	assert.Nil(t, err, "project should be added")

	err = AddDomainProxy("test", "test.com", 3000)
	assert.Nil(t, err, "domain should be added")

	d, err := os.ReadFile(path.Join(consts.BasePath, "domains", p.Name+".caddy"))
	assert.Nil(t, err, "domain should exist")
	assert.Contains(t, string(d), "test.com")
	assert.Contains(t, string(d), "127.0.0.1:3000")
	assert.Equal(t, successDomainProxy, string(d), "created .caddy content should match expected")

	// making sure the domain passed gets validated
	err = AddDomainProxy("test", "tes", 3000)
	assert.Error(t, err, "domain should validate and fail")

	// making sure the domain changes if ran again
	err = AddDomainProxy("test", "two.com", 3001)
	d, err = os.ReadFile(path.Join(consts.BasePath, "domains", p.Name+".caddy"))
	assert.Nil(t, err, "domain should be able to be changed")
	assert.Contains(t, string(d), "two.com")
	assert.Contains(t, string(d), "127.0.0.1:3001")
}

func TestAddDomainStatic(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	p := Project{
		Name: "test",
	}
	err := addProject(p)
	assert.Nil(t, err, "project should be added")

	err = AddDomainStatic("test", "test.com", "")
	assert.Nil(t, err, "domain should be added")

	d, err := os.ReadFile(path.Join(consts.BasePath, "domains", p.Name+".caddy"))
	assert.Nil(t, err, "domain should exist")
	assert.Contains(t, string(d), "test.com")
	assert.Contains(t, string(d), "root * /home/projects/test/")
	assert.Equal(t, successDomainStatic, string(d), "created .caddy content should match expected")

	// making sure the domain passed gets validated
	err = AddDomainStatic("test", "tes", "")
	assert.Error(t, err, "domain should validate and fail")

	// making sure the domain changes if ran again
	err = AddDomainStatic("test", "two.com", "inner")
	d, err = os.ReadFile(path.Join(consts.BasePath, "domains", p.Name+".caddy"))
	assert.Nil(t, err, "domain should be able to be changed")
	assert.Contains(t, string(d), "two.com")
	assert.Contains(t, string(d), "root * /home/projects/test/inner")

}

func TestSetupDomains(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	err := SetupDomains("test@email.com")
	assert.Nil(t, err, "base caddy config should be setup")

	d, err := os.ReadFile(path.Join(consts.BasePath, "caddy", "main.caddy"))
	assert.Nil(t, err, "main config should exist")
	assert.Contains(t, string(d), "test@email.com", "main config should include email")

	// test if email changes on rerun
	err = SetupDomains("test2@email.com")
	assert.Nil(t, err, "base caddy config should be setup")

	d, err = os.ReadFile(path.Join(consts.BasePath, "caddy", "main.caddy"))
	assert.Nil(t, err, "main config should exist")
	assert.Contains(t, string(d), "test2@email.com", "main config should include email")

}

func TestDeleteDomain(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	p := Project{
		Name: "test",
	}
	err := addProject(p)
	assert.Nil(t, err, "project should be added")

	err = AddDomainStatic("test", "test.com", "")
	assert.Nil(t, err, "domain should be added")

	err = DeleteProjectDomain("test")
	assert.Nil(t, err, "domain should be able to be deleted")

	_, err = os.ReadFile(path.Join(consts.BasePath, "domains", p.Name+".caddy"))
	assert.Error(t, err, "domain shouldn't exist")
}
