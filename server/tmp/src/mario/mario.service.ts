import { Injectable } from '@nestjs/common';
import { MarioDto, MarioPostion, MarioVelocity } from './mario.model';

@Injectable()
export class MarioService {
  mario: MarioDto = { width: 16, height: 16, x: 276, y: 44 };
  position: MarioPostion = { x: 0, y: 130 };
  velocity: MarioVelocity = { x: 70, y: -390 };
  deltaTime = 1 / 60;
  lastTime = 0;

  public getMario(): MarioDto {
    return this.mario;
  }

  public getPosition(deltaTime): MarioPostion {
    this.position.x += this.velocity.x * deltaTime;
    this.position.y += this.velocity.y * deltaTime;
    if (this.position.y > 130) {
      this.position.y = 130;
    }
    this.velocity.y += 15;
    return this.position;
  }
}
