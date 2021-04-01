import { Injectable } from '@nestjs/common';
import { TilesService } from 'src/tiles/tiles.service';
import { BlockBoundary, BlockDto } from './blocks.model';

@Injectable()
export class BlocksService {
  tilesService: TilesService = new TilesService();

  ground = { x1: 0, y1: 10, x2: 30, y2: 13 };
  sky = { x1: 0, y1: 0, x2: 30, y2: 10 };
  boundary = {
    ground: this.ground,
    sky: this.sky,
  };

  public getBlock(key: string): BlockDto {
    return {
      tile: this.tilesService.getTile(key),
      boundary: this.boundary[key],
    };
  }
}
