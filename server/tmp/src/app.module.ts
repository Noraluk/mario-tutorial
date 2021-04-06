import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { EventsGateway } from './events/events.gateway';
import { MarioService } from './mario/mario.service';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    EventsGateway,
  ],
  controllers: [],
  providers: [MarioService],
})
export class AppModule {}
