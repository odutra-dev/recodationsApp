import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { CreateMessageDto } from './dto/create-message.dto';
import { UpdateMessageDto } from './dto/update-message.dto';
import { Repository } from 'typeorm';
import { Message } from './entities/message.entity';

@Injectable()
export class MessageService {
  constructor(
    @InjectRepository(Message)
    private messageRepository: Repository<Message>,
  ) {}
  create(createMessageDto: CreateMessageDto): Promise<Message> {
    return this.messageRepository.save(createMessageDto);
  }

  findAll(): Promise<Message[]> {
    return this.messageRepository.find();
  }

  findOne(id: number): Promise<Message | null> {
    return this.messageRepository.findOne({ where: { id } });
  }

  update(id: number, updateMessageDto: UpdateMessageDto): Promise<void> {
    return this.messageRepository
      .update(id, updateMessageDto)
      .then(() => undefined);
  }

  async remove(id: number): Promise<void> {
    await this.messageRepository.delete(id);
  }
}
