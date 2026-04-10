import { Controller, Get } from '@nestjs/common';
import { ApiTags, ApiOperation, ApiResponse } from '@nestjs/swagger';

@ApiTags('Health')
@Controller()
export class AppController {
 @Get()
 @ApiOperation({ summary: 'Health check de la API' })
 @ApiResponse({ status: 200, description: 'API funcionando correctamente' })
 getHealth() {
 return {
 status: 'ok',
 message: 'API centralizada de usuarios',
 version: '1.0.0',
 timestamp: new Date().toISOString(),
 };
 }
}

