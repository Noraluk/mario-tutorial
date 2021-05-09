package screenService

import (
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

var delay = 0

func (s *screen) GetScreen(camera *screenEntity.Camera, mario *marioEntity.Mario) *screenEntity.Screen {
	level := s.bgService.GetBackground()

	if mario.Position.X >= constants.HALF_SCREEN && mario.Position.X < 3392-constants.HALF_SCREEN {
		camera.Position.X = mario.Position.X - constants.HALF_SCREEN - 2
	}

	cameraStart := int(camera.Position.X)
	cameraEnd := int(camera.Position.X+constants.TILE_SILE) + camera.Size.Width

	newBGs := []bgEntity.Background{}
	for i, bg := range level.Backgrounds {
		if len(bg.Animation.Frames) > 0 && delay%bg.Animation.Duration == 0 {
			level.Backgrounds[i].Position = bg.Animation.Frames[bg.Animation.Current%(len(bg.Animation.Frames))]
			level.Backgrounds[i].Animation.Current++
		}

		newRanges := []bgEntity.Range{}
		for _, val := range bg.Ranges {
			x1 := val.X1 * constants.TILE_SILE
			x2 := val.X2 * constants.TILE_SILE

			if (x1 >= cameraStart && x2 >= cameraStart) && (cameraEnd >= x1 && cameraEnd >= x2) {
				newRanges = append(newRanges, bgEntity.NewRange(val.X1, val.X2, val.Y1, val.Y2))
			} else if cameraStart >= x1 && x2 >= cameraEnd {
				newRanges = append(newRanges, bgEntity.NewRange(cameraStart/int(constants.TILE_SILE), cameraEnd/int(constants.TILE_SILE), val.Y1, val.Y2))
			} else if x2 > cameraEnd && (x1 >= cameraStart && x1 <= cameraEnd) {
				newRanges = append(newRanges, bgEntity.NewRange(val.X1, cameraEnd/int(constants.TILE_SILE), val.Y1, val.Y2))
			} else if cameraStart > x1 && (x2 >= cameraStart && x2 <= cameraEnd) {
				newRanges = append(newRanges, bgEntity.NewRange(cameraStart/int(constants.TILE_SILE), val.X2, val.Y1, val.Y2))
			}
		}

		newBGs = append(newBGs, bgEntity.NewBackground(bg.Tile, bg.IsCollide, bg.Position, newRanges, bg.Animation))
	}

	delay++
	return screenEntity.NewScreen(newBGs, camera, mario)
}
