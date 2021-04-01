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
import { BlocksService } from 'src/blocks/blocks.service';
import { BlockDto } from 'src/blocks/blocks.model';

@WebSocketGateway()
export class EventsGateway
  implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer() server: Server;

  private logger: Logger = new Logger('MessageGateway');
  private blockService = new BlocksService();

  @SubscribeMessage('start')
  public start(client: Socket, payload: any): void {
    const blocks = ['ground', 'sky'];
    blocks.forEach((block) => {
      client.emit('onStart', this.blockService.getBlock(block));
    });
  }

  public afterInit(server: Server): void {
    return this.logger.log('Init');
  }

  public handleDisconnect(client: Socket): void {
    return this.logger.log(`Client disconnected: ${client.id}`);
  }

  public handleConnection(client: Socket): void {
    return this.logger.log(`Client connected: ${client.id}`);
  }
}
