package config

import (
	bgEntity "server/background/entities"
)

type Config interface {
	SetCollider(position bgEntity.Position, tileName string)
	GetCollider(position bgEntity.Position) string
}

type config struct {
	Colliders map[bgEntity.Position]string
}

func New() Config {
	colliders := make(map[bgEntity.Position]string)
	return &config{
		Colliders: colliders,
	}
}

func (c *config) SetCollider(position bgEntity.Position, tileName string) {
	c.Colliders[position] = tileName
}

func (c *config) GetCollider(position bgEntity.Position) string {
	return c.Colliders[position]
}
