import { Injectable } from '@nestjs/common';
import { SceneDto } from './scenes.model';

@Injectable()
export class SceneService {
  scene: SceneDto = { width: 200, height: 200, x: 0, y: 0 };

  public getScene(): SceneDto {
    return this.scene;
  }
}
