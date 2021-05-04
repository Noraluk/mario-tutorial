package screenService

import (
	"log"
	"server/common"
	"server/constants"
	marioEntity "server/mario/entities"

	bgEntity "server/background/entities"
	bgService "server/background/services"
	screenEntity "server/screen/entities"
)

type Screen interface {
	GetScreen(camera *screenEntity.Camera, mario *marioEntity.Mario) *screenEntity.Screen
}

type screen struct {
	bgService bgService.Background
}

func New(bgService bgService.Background) Screen {
	return &screen{
		bgService: bgService,
	}
}

func (s *screen) GetScreen(camera *screenEntity.Camera, mario *marioEntity.Mario) *screenEntity.Screen {
	level, err := s.bgService.GetBackground()
	if err != nil {
		log.Fatal(err)
	}

	if mario.Position.X >= constants.HALF_SCREEN && mario.Position.X < 3392-constants.HALF_SCREEN {
		camera.Position.X = mario.Position.X - constants.HALF_SCREEN
	}

	cameraStart := camera.Position.X
	cameraEnd := int(camera.Position.X+constants.TILE_SILE) + camera.Size.Width

	for i, bg := range level.Backgrounds {
		newRanges := []bgEntity.Ranges{}
		for _, val := range bg.Ranges {
			x1 := val.X1 * constants.TILE_SILE
			x2 := val.X2 * constants.TILE_SILE

			if (cameraStart > float64(x1) && cameraStart > float64(x2)) || (x1 > cameraEnd && x2 > cameraEnd) {
				continue
			}
			newRange := bgEntity.Ranges{X1: x1 / constants.TILE_SILE, X2: x2 / constants.TILE_SILE, Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: int(constants.TILE_SILE), Height: int(constants.TILE_SILE)}}
			if cameraStart > float64(x1) && float64(x2) > cameraStart {
				newRange = bgEntity.Ranges{X1: int(cameraStart / constants.TILE_SILE), X2: x2 / int(constants.TILE_SILE), Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: (val.X2 - int(cameraStart)) / constants.TILE_SILE, Height: int(constants.TILE_SILE)}}
			} else if cameraEnd > x1 && x2 > cameraEnd {
				newRange = bgEntity.Ranges{X1: x1 / constants.TILE_SILE, X2: cameraEnd / constants.TILE_SILE, Y1: val.Y1, Y2: val.Y2, TileSize: common.Size{Width: (cameraEnd - val.X1) / constants.TILE_SILE, Height: int(constants.TILE_SILE)}}
			}
			newRanges = append(newRanges, newRange)
		}
		level.Backgrounds[i].Ranges = newRanges
	}

	return screenEntity.NewScreen(level.Backgrounds, camera, mario)
}
