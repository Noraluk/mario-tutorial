import {
  SubscribeMessage,
  WebSocketGateway,
  OnGatewayInit,
  WebSocketServer,
  OnGatewayConnection,
  OnGatewayDisconnect,
} from '@nestjs/websockets';
import { Socket } from 'socket.io';
import { Logger } from '@nestjs/common';
import { Server } from 'ws';
import { MarioService } from 'src/mario/mario.service';
import * as levels from '../../static/levels/1-1.json';
import { TilesService } from 'src/tiles/tiles.service';

@WebSocketGateway()
export class EventsGateway
  implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer() server: Server;

  private logger: Logger = new Logger('MessageGateway');
  private marioService = new MarioService();
  private tileService = new TilesService();

  private createTiles(
    x1: number,
    y1: number,
    x2: number,
    y2: number,
    value: string,
  ) {
    for (let x = x1; x < x2; x++) {
      for (let y = y1; y < y2; y++) {
        // console.log(x, y);
        this.tileService.setTilesColider(x, y, value);
      }
    }
  }

  @SubscribeMessage('start')
  public start(client: Socket, payload: any): void {
    const now = Date.now();
    setTimeout(() => {
      levels.backgrounds.forEach((bg) => {
        const tile = this.tileService.getTile(bg.tile);
        bg.ranges.forEach(([x1, x2, y1, y2]) => {
          const data = {
            tile: tile,
            boundary: {
              x1: x1,
              y1: y1,
              x2: x2,
              y2: y2,
            },
          };
          console.log(data);
          client.emit('onStart', data);
          // console.log(bg.tile, x1, y1, x2, y2);
          // this.createTiles(x1, y1, x2, y2, bg.tile);
          // console.log(this.tileService.tilesColider);
        });

        // console.log(this.tileService.tilesColider);
        // console.log('a');
        // client.emit('onStart', this.tileService.tilesColider);
      });
      // blocks.forEach((block) => {
      // client.emit('onStart', this.blockService.getBlock(block));
      // });
      // client.emit('onStart', {
      //   mario: this.marioService.getMario(),
      //   position: this.marioService.getPosition((Date.now() - now) / 1000),
      // });
    }, 10);
  }

  public afterInit(server: Server): void {
    return this.logger.log('Init');
  }

  public handleDisconnect(client: Socket): void {
    console.log(this.tileService.tilesColider);
    return this.logger.log(`Client disconnected: ${client.id}`);
  }

  public handleConnection(client: Socket): void {
    this.marioService = new MarioService();
    // console.log(levels.backgrounds);
    return this.logger.log(`Client connected: ${client.id}`);
  }
}
