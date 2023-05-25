package cliententity

import (
	"errors"

	"github.com/google/uuid"
)

type Clienter interface {
	GetUUID() (*uuid.UUID, error)
	SetUUID(*uuid.UUID)
}

type Client struct {
	uuid *uuid.UUID
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetUUID() (*uuid.UUID, error) {
	if c.uuid == nil {
		return nil, errors.New("Empty uuid")
	}
	return c.uuid, nil
}

func (c *Client) SetUUID(newUUID *uuid.UUID) {
	c.uuid = newUUID
}
