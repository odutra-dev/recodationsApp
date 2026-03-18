import { NestFactory } from '@nestjs/core';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { AppModule } from './app.module';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);

  const config = new DocumentBuilder()
    .setTitle('Recondations API')
    .setDescription('API for Recondations application')
    .setVersion('1.0')
    .addTag('recondations')
    .build();
  const documentFactory = () => SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api', app, documentFactory, {
    swaggerOptions: {
      jsonDocumentUrl: 'swagger/json',
    },
  });

  await app.listen(process.env.PORT ?? 3000);
}
bootstrap();
