import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateMessageDto } from './dto/create-message.dto';
import { UpdateMessageDto } from './dto/update-message.dto';
import { Repository } from 'typeorm';
import { Message } from './entities/message.entity';
import { S3Service } from 'src/s3/s3.service';

@Injectable()
export class MessageService {
  constructor(
    @InjectRepository(Message)
    private messageRepository: Repository<Message>,
    private readonly s3Service: S3Service,
  ) {}

  async verifyMessage(id: number): Promise<void> {
    const message = await this.messageRepository.findOne({ where: { id } });

    if (!message) {
      throw new HttpException('Message not found', HttpStatus.NOT_FOUND);
    }
  }

  async create(
    createMessageDto: CreateMessageDto,
    file: Express.Multer.File,
  ): Promise<Message> {
    let imageUrl: string | null = null;

    if (file) {
      const upload = await this.s3Service.uploadFile('photos', file);
      imageUrl = upload.url;
    }

    const message = this.messageRepository.create({
      ...createMessageDto,
      image: imageUrl,
    });

    return this.messageRepository.save(message);
  }

  findAll(): Promise<Message[]> {
    return this.messageRepository.find();
  }

  async findOne(id: number): Promise<Message | null> {
    if (!id) {
      throw new HttpException('ID is required', HttpStatus.BAD_REQUEST);
    }

    await this.verifyMessage(id);

    return await this.messageRepository.findOne({ where: { id } });
  }

  async update(id: number, updateMessageDto: UpdateMessageDto): Promise<void> {
    if (!id) {
      throw new HttpException('ID is required', HttpStatus.BAD_REQUEST);
    }
    await this.verifyMessage(id);

    return await this.messageRepository
      .update(id, updateMessageDto)
      .then(() => undefined);
  }

  async remove(id: number): Promise<void> {
    if (!id) {
      throw new HttpException('ID is required', HttpStatus.BAD_REQUEST);
    }

    await this.verifyMessage(id);

    await this.messageRepository.delete(id);
  }
}
