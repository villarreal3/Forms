import { NestFactory } from '@nestjs/core';
import { ValidationPipe, Logger } from '@nestjs/common';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import helmet from 'helmet';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import { Request, Response } from 'express';

async function bootstrap() {
 const logger = new Logger('Bootstrap');
 const app = await NestFactory.create(AppModule);

 const configService = app.get(ConfigService);
 const isDevelopment = configService.get('NODE_ENV') === 'development';

 // Configuración de seguridad con Helmet (relajada para desarrollo)
 app.use(
 helmet({
 contentSecurityPolicy: isDevelopment ? false : undefined,
 crossOriginEmbedderPolicy: isDevelopment ? false : undefined,
 }),
 );

 // CORS configurado desde variables de entorno
 const allowedOrigins = configService.get<string>('ALLOWED_ORIGINS')?.split(',') || [];
 app.enableCors({
 origin: (origin, callback) => {
 // Permitir peticiones sin origin (Swagger, Postman, curl, server-to-server)
 if (!origin) {
 callback(null, true);
 return;
 }
 
 // Verificar si el origin está permitido
 if (allowedOrigins.some(allowed => origin.startsWith(allowed.trim()))) {
 callback(null, true);
 } else if (isDevelopment) {
 // En desarrollo, permitir todo
 callback(null, true);
 } else {
 callback(new Error('Not allowed by CORS'));
 }
 },
 credentials: true,
 methods: ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS'],
 allowedHeaders: ['Content-Type', 'Authorization', 'X-Requested-With'],
 });

 // Prefijo global de la API
 const apiPrefix = configService.get<string>('API_PREFIX') || 'api/v1';
 app.setGlobalPrefix(apiPrefix);

 // Validación global de DTOs
 app.useGlobalPipes(
 new ValidationPipe({
 whitelist: true,
 forbidNonWhitelisted: true,
 transform: true,
 transformOptions: {
 enableImplicitConversion: true,
 },
 }),
 );

 // Configuración de Swagger
 const config = new DocumentBuilder()
 .setTitle('API centralizada de usuarios')
 .setDescription(
 'Sistema institucional unificado para gestión, autenticación y autorización de usuarios',
 )
 .setVersion('1.0')
 .addBearerAuth(
 {
 type: 'http',
 scheme: 'bearer',
 bearerFormat: 'JWT',
 name: 'JWT',
 description: 'Ingrese su token JWT',
 in: 'header',
 },
 'JWT-auth',
 )
 .addTag('Auth', 'Endpoints de autenticación')
 .addTag('Users', 'Gestión de usuarios')
 .addTag('Roles', 'Gestión de roles')
 .addTag('Departments', 'Gestión de departamentos')
 .addTag('Audit', 'Auditoría del sistema')
 .build();

 const document = SwaggerModule.createDocument(app, config);
 SwaggerModule.setup('docs', app, document, {
 customSiteTitle: 'Users API - Documentación',
 customCss: '.swagger-ui .topbar { display: none }',
 });

 // Exponer la especificación en formato JSON para el gateway
 app.getHttpAdapter().get('/docs-json', (req: Request, res: Response) => {
 res.json(document);
 });

 const port = configService.get<number>('PORT') || 3000;
 await app.listen(port, '0.0.0.0');

 logger.log(`🚀 API de usuarios ejecutándose en: http://localhost:${port}/${apiPrefix}`);
 logger.log(`📚 Documentación Swagger disponible en: http://localhost:${port}/docs`);
}

bootstrap();

