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
  tilesColider = [];

  public getTile(key: string): TileDto {
    return this.tiles[key];
  }

  public setTilesColider(x: number, y: number, value: string) {
    // if (this.getTilesColider(x, y)) {
    //   return;
    // }

    if (!this.tilesColider[x]) {
      this.tilesColider[x] = [];
    }

    this.tilesColider[x][y] = value;
    console.log(this.tilesColider[x][y]);
  }

  public getTilesColider(x: number, y: number): string {
    if (!this.tilesColider[x]) return;
    return this.tilesColider[x][y];
  }
}
