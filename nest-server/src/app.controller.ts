import { Controller, Get } from '@nestjs/common';
import { EventPattern, MessagePattern } from '@nestjs/microservices';
import { AppService } from './app.service';

@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}
  @EventPattern('exchange')
  public async mock(data: any) {
    console.log('go-messsage', data);
  }
  @Get()
  getHello(): string {
    return this.appService.getHello();
  }
}
