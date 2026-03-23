import { Test, TestingModule } from '@nestjs/testing';
import { MessageService } from './message.service';
import { S3Service } from '../s3/s3.service';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Message } from './entities/message.entity';

describe('MessageService', () => {
  let service: MessageService;

  const mockMessageRepository = {
    find: jest.fn(),
    findOne: jest.fn(),
    save: jest.fn(),
    create: jest.fn(),
    delete: jest.fn(),
  };

  const mockS3Service = {
    upload: jest.fn(),
    delete: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        MessageService,
        {
          provide: S3Service,
          useValue: mockS3Service,
        },
        {
          provide: getRepositoryToken(Message),
          useValue: mockMessageRepository,
        },
      ],
    }).compile();

    service = module.get<MessageService>(MessageService);
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });
});
