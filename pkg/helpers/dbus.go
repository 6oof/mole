package helpers

import (
	"context"

	"github.com/coreos/go-systemd/v22/dbus"
)

// we should always close the connection when no longer using it
func ContactDbus() (*dbus.Conn, error) {
	ctx := context.Background()

	conn, err := dbus.NewUserConnectionContext(ctx)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
