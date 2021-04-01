import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { EventsGateway } from './events/events.gateway';
import { BlocksService } from './blocks/blocks.service';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
    }),
    EventsGateway,
  ],
  controllers: [],
  providers: [BlocksService],
})
export class AppModule {}
