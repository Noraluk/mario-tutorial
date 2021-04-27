package config

import "server/common"

type Config interface {
	SetCollider(position common.Position, tileName string)
	GetCollider(position common.Position) string
	GetColliders() map[common.Position]string
}

type config struct {
	Colliders map[common.Position]string
}

func New() Config {
	colliders := make(map[common.Position]string)
	return &config{
		Colliders: colliders,
	}
}

func (c *config) SetCollider(position common.Position, tileName string) {
	c.Colliders[position] = tileName
}

func (c *config) GetCollider(position common.Position) string {
	return c.Colliders[position]
}

func (c *config) GetColliders() map[common.Position]string {
	return c.Colliders
}
