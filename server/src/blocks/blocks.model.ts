import { TileDto } from 'src/tiles/tiles.model';

export class BlockDto {
  boundary: BlockBoundary;
  tile: TileDto;
}

export class BlockBoundary {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
}
