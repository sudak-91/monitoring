package command

import "github.com/sudak-91/monitoring/internal/pkg/client"

type Commands interface {
	SetUUID()
	GetUUID() string
}

type Command struct {
	client *client.Client
}

func NewCommand(c *client.Client) *Command {
	var cmd Command
	cmd.client = c
	return &cmd
}
