package actions

import (
	"os"
	"path"
	"testing"

	"github.com/zulubit/mole/pkg/consts"
	"github.com/zulubit/mole/pkg/enums"
	"github.com/stretchr/testify/assert"
)

func TestLinkUnlinkServices(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	np := Project{
		Name: "test-project",
	}
	addProject(np)

	sfp := path.Join(tmp, "projects", np.Name, "mole_services")
	os.MkdirAll(sfp, 0755)
	os.WriteFile(path.Join(sfp, "asdf.service"), []byte("asdf"), 0755)
	os.WriteFile(path.Join(sfp, "aaa.container"), []byte("asdf"), 0755)

	_, err := os.ReadFile(path.Join(sfp, "asdf.service"))
	assert.Nil(t, err, "test was setup correctly")

	err = LinkServices(np.Name, enums.Static)
	assert.ErrorContains(t, err, "invalid service type static", "services were linked")

	err = LinkServices(np.Name, enums.Podman)
	assert.Nil(t, err, "services were linked")
	assert.FileExists(t, path.Join(tmp, ".config", "containers", "systemd", "mole-test-project-asdf.service"), "service files were named correctly")
	assert.FileExists(t, path.Join(tmp, ".config", "containers", "systemd", "mole-test-project-aaa.container"), "service files were named correctly")

	err = LinkServices(np.Name, enums.Systemd)
	assert.Nil(t, err, "services were linked")
	assert.FileExists(t, path.Join(tmp, ".config", "systemd", "user", "mole-test-project-asdf.service"), "service files were named correctly")
	assert.FileExists(t, path.Join(tmp, ".config", "systemd", "user", "mole-test-project-aaa.container"), "service files were named correctly")

	err = UnlinkServices(np.Name)
	assert.Nil(t, err, "services were linked")
	assert.NoFileExists(t, path.Join(tmp, ".config", "systemd", "user", "mole-test-project-asdf.service"), "service files were removed")
	assert.NoFileExists(t, path.Join(tmp, ".config", "systemd", "user", "mole-test-project-aaa.container"), "service files were removed")
	assert.NoFileExists(t, path.Join(tmp, ".config", "containers", "systemd", "mole-test-project-asdf.service"), "service files were removed")
	assert.NoFileExists(t, path.Join(tmp, ".config", "containers", "systemd", "mole-test-project-aaa.container"), "service files were removed")

}
