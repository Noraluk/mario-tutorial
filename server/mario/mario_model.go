package mario

import bg_entity "server/background/entities"

type Model struct {
	X        int                `json:"x"`
	Y        int                `json:"y"`
	Width    int                `json:"width"`
	Height   int                `json:"height"`
	Position bg_entity.Position `json:"position"`
}
