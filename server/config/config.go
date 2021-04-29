package config

import "server/common"

type Config interface {
	SetCollider(position common.Position, isCollide bool)
	GetCollider(position common.Position) bool
	GetColliders() map[common.Position]bool
}

type config struct {
	Colliders map[common.Position]bool
}

func New() Config {
	colliders := make(map[common.Position]bool)
	return &config{
		Colliders: colliders,
	}
}

func (c *config) SetCollider(position common.Position, tileName bool) {
	c.Colliders[position] = tileName
}

func (c *config) GetCollider(position common.Position) bool {
	return c.Colliders[position]
}

func (c *config) GetColliders() map[common.Position]bool {
	return c.Colliders
}
