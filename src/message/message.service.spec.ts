import { Test, TestingModule } from '@nestjs/testing';
import { MessageService } from './message.service';
import { S3Service } from '../s3/s3.service';
import { getRepositoryToken } from '@nestjs/typeorm';
import { Message } from './entities/message.entity';
import { HttpException } from '@nestjs/common';
import { CreateMessageDto } from './dto/create-message.dto';
import { UpdateMessageDto } from './dto/update-message.dto';

describe('MessageService', () => {
  let service: MessageService;

  const mockMessageRepository = {
    find: jest.fn(),
    findOne: jest.fn(),
    save: jest.fn(),
    create: jest.fn(),
    delete: jest.fn(),
    update: jest.fn(),
  };

  const mockS3Service = {
    uploadFile: jest.fn(),
  };

  beforeEach(async () => {
    jest.clearAllMocks();

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

  // =========================
  // CREATE
  // =========================
  describe('create', () => {
    it('should create message without file', async () => {
      const dto = {
        content: 'hello',
        name: 'test',
      } as CreateMessageDto;

      mockMessageRepository.create.mockReturnValue({
        ...dto,
        image: null,
      });

      mockMessageRepository.save.mockResolvedValue({
        ...dto,
        image: null,
      });

      const result = await service.create(dto, null);

      expect(mockMessageRepository.create).toHaveBeenCalledWith({
        ...dto,
        image: null,
      });

      expect(result).toEqual({
        ...dto,
        image: null,
      });
    });

    it('should create message with file', async () => {
      const dto = {
        content: 'hello',
        name: 'test',
      } as CreateMessageDto;
      const file = { originalname: 'test.png' } as Express.Multer.File;

      mockS3Service.uploadFile.mockResolvedValue({
        url: 'http://image.com',
      });

      mockMessageRepository.create.mockReturnValue({
        ...dto,
        image: 'http://image.com',
      });

      mockMessageRepository.save.mockResolvedValue({
        ...dto,
        image: 'http://image.com',
      });

      const result = await service.create(dto, file);

      expect(mockS3Service.uploadFile).toHaveBeenCalledWith('photos', file);

      expect(result.image).toBe('http://image.com');
    });
  });

  // =========================
  // FIND ALL
  // =========================
  describe('findAll', () => {
    it('should return all messages', async () => {
      const messages = [{ id: 1 }];

      mockMessageRepository.find.mockResolvedValue(messages);

      const result = await service.findAll();

      expect(result).toEqual(messages);
    });
  });

  // =========================
  // FIND ONE
  // =========================
  describe('findOne', () => {
    it('should throw if id not provided', async () => {
      await expect(service.findOne(null)).rejects.toThrow(HttpException);
    });

    it('should return message', async () => {
      const message = { id: 1 };

      mockMessageRepository.findOne
        .mockResolvedValueOnce(message) // verifyMessage
        .mockResolvedValueOnce(message); // findOne final

      const result = await service.findOne(1);

      expect(result).toEqual(message);
    });

    it('should throw if message not found', async () => {
      mockMessageRepository.findOne.mockResolvedValue(null);

      await expect(service.findOne(1)).rejects.toThrow('Message not found');
    });
  });

  // =========================
  // UPDATE
  // =========================
  describe('update', () => {
    it('should throw if id not provided', async () => {
      await expect(
        service.update(null, {} as UpdateMessageDto),
      ).rejects.toThrow(HttpException);
    });

    it('should update message', async () => {
      const message = { id: 1 };

      mockMessageRepository.findOne.mockResolvedValue(message);
      mockMessageRepository.update.mockResolvedValue({});

      await service.update(1, {
        content: 'updated',
        name: 'test',
      } as UpdateMessageDto);

      expect(mockMessageRepository.update).toHaveBeenCalledWith(1, {
        content: 'updated',
        name: 'test',
      });
    });
  });

  // =========================
  // REMOVE
  // =========================
  describe('remove', () => {
    it('should throw if id not provided', async () => {
      await expect(service.remove(null)).rejects.toThrow(HttpException);
    });

    it('should delete message', async () => {
      const message = { id: 1 };

      mockMessageRepository.findOne.mockResolvedValue(message);
      mockMessageRepository.delete.mockResolvedValue({});

      await service.remove(1);

      expect(mockMessageRepository.delete).toHaveBeenCalledWith(1);
    });
  });

  // =========================
  // VERIFY MESSAGE (extra)
  // =========================
  describe('verifyMessage', () => {
    it('should pass if message exists', async () => {
      mockMessageRepository.findOne.mockResolvedValue({ id: 1 });

      await expect(service.verifyMessage(1)).resolves.toBeUndefined();
    });

    it('should throw if message does not exist', async () => {
      mockMessageRepository.findOne.mockResolvedValue(null);

      await expect(service.verifyMessage(1)).rejects.toThrow(
        'Message not found',
      );
    });
  });
});
