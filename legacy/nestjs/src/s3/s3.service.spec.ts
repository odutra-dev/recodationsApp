import { Test, TestingModule } from '@nestjs/testing';
import { S3Service } from './s3.service';

describe('S3Service', () => {
  let service: S3Service;

  const mockSend = jest.fn();

  beforeEach(async () => {
    process.env.AWS_S3_ENDPOINT = 'http://localhost:4566';

    const module: TestingModule = await Test.createTestingModule({
      providers: [S3Service],
    }).compile();

    service = module.get<S3Service>(S3Service);

    (service as any).s3Client = {
      send: mockSend,
    };

    jest.clearAllMocks();
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  describe('ensureBucketExists', () => {
    it('should not create bucket if it exists', async () => {
      mockSend.mockResolvedValueOnce({}); // HeadBucket

      await service.ensureBucketExists('bucket');

      expect(mockSend).toHaveBeenCalledTimes(1);
    });

    it('should create bucket if it does not exist', async () => {
      mockSend
        .mockRejectedValueOnce(new Error()) // HeadBucket
        .mockResolvedValueOnce({}); // CreateBucket

      await service.ensureBucketExists('bucket');

      expect(mockSend).toHaveBeenCalledTimes(2);
    });
  });

  describe('uploadFile', () => {
    it('should upload file successfully', async () => {
      mockSend
        .mockResolvedValueOnce({}) // HeadBucket
        .mockResolvedValueOnce({}); // PutObject

      const result = await service.uploadFile('bucket', {
        originalname: 'file.jpg',
        mimetype: 'image/jpeg',
        buffer: Buffer.from('test'),
      } as Express.Multer.File);

      expect(result.message).toBe('File uploaded successfully');
      expect(result.url).toContain('bucket');
      expect(mockSend).toHaveBeenCalledTimes(2);
    });

    it('should create bucket if not exists and upload', async () => {
      mockSend
        .mockRejectedValueOnce(new Error()) // HeadBucket
        .mockResolvedValueOnce({}) // CreateBucket
        .mockResolvedValueOnce({}); // PutObject

      const result = await service.uploadFile('bucket', {
        originalname: 'file.jpg',
        mimetype: 'image/jpeg',
        buffer: Buffer.from('test'),
      } as Express.Multer.File);

      expect(result.message).toBe('File uploaded successfully');
      expect(mockSend).toHaveBeenCalledTimes(3);
    });

    it('should throw error if upload fails', async () => {
      mockSend
        .mockResolvedValueOnce({}) // HeadBucket
        .mockRejectedValueOnce(new Error('Upload failed')); // PutObject

      await expect(
        service.uploadFile('bucket', {
          originalname: 'file.jpg',
          mimetype: 'image/jpeg',
          buffer: Buffer.from('test'),
        } as Express.Multer.File),
      ).rejects.toThrow('Upload failed');
    });
  });
});
