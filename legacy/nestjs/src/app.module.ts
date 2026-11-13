import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { S3Service } from './s3/s3.service';
import { ConfigModule } from '@nestjs/config';
import { MessageModule } from './message/message.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { DataSource } from 'typeorm';
import { Message } from './message/entities/message.entity';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
    TypeOrmModule.forRoot({
      type: 'postgres',
      host: process.env.DB_HOST!,
      port: parseInt(process.env.DB_PORT!, 10),
      username: process.env.DB_USER!,
      password: process.env.DB_PASS!,
      database: process.env.DB_NAME!,
      autoLoadEntities: true,
      synchronize: true,
      entities: [Message],
    }),
    MessageModule,
  ],
  controllers: [AppController],
  providers: [AppService, S3Service],
})
export class AppModule {
  constructor(private dataSource: DataSource) {}
}
