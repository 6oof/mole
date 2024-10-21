package execs

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/6oof/mole/pkg/consts"
	"github.com/6oof/mole/pkg/data"
)

func FindAndEditEnv(pName string) error {

	p, err := data.FindProject(pName)
	if err != nil {
		return err
	}

	c := exec.Command("nano", ".env")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Dir = path.Join(consts.BasePath, "projects", p.Name)

	// Run the command and handle any error
	err = c.Run()
	if err != nil {
		return errors.New("Error running vi:" + err.Error())
	}

	return nil
}
