import { Injectable } from '@nestjs/common';
import { TileDto } from 'src/tiles/tiles.model';

@Injectable()
export class TilesService {
  ground: TileDto = { width: 16, height: 16, x: 0, y: 0 };
  sky: TileDto = { width: 16, height: 16, x: 3, y: 23 };

  tiles = {
    ground: this.ground,
    sky: this.sky,
  };

  public getTile(key: string): TileDto {
    return this.tiles[key];
  }
}
