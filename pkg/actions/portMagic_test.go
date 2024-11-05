package actions

import (
	"os"
	"testing"

	"github.com/6oof/mole/pkg/consts"
	"github.com/stretchr/testify/assert"
)

func TestSaveReservedPorts(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	p := ports{1000}

	err := saveReservedPorts(p)
	assert.Nil(t, err, "ports should save")
}

func TestGetReservedAndUsedPorts(t *testing.T) {
	consts.Testing = true

	tmp := os.TempDir()
	consts.BasePath = tmp
	defer os.RemoveAll(tmp)

	p := ports{12345}

	saveReservedPorts(p)

	prts, err := getReservedAndUsedPorts()
	assert.Nil(t, err, "should be able to get ports")
	assert.Contains(t, prts, 12345, "ports we saved should exits")
}
