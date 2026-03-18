import { Injectable } from '@nestjs/common';
import {
  S3Client,
  PutObjectCommand,
  HeadBucketCommand,
  CreateBucketCommand,
} from '@aws-sdk/client-s3';
@Injectable()
export class S3Service {
  private s3Client: S3Client;

  constructor() {
    this.s3Client = new S3Client({
      region: process.env.AWS_REGION,
      endpoint: process.env.AWS_S3_ENDPOINT!,
      credentials: {
        accessKeyId: process.env.AWS_ACCESS_KEY_ID!,
        secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY!,
      },
      forcePathStyle: true,
    });
  }

  async ensureBucketExists(bucket: string) {
    try {
      await this.s3Client.send(new HeadBucketCommand({ Bucket: bucket }));
    } catch {
      await this.s3Client.send(new CreateBucketCommand({ Bucket: bucket }));
      console.log(`Bucket ${bucket} created successfully`);
    }
  }

  async uploadFile(
    bucketName: string,
    file: Express.Multer.File,
  ): Promise<{ message: string; url: string }> {
    const nameFile = file.originalname + '-' + Date.now();

    const command = new PutObjectCommand({
      Bucket: bucketName,
      Key: nameFile,
      Body: file.buffer,
      ContentType: file.mimetype,
    });

    try {
      await this.ensureBucketExists(bucketName);

      await this.s3Client.send(command);

      console.log(`File uploaded successfully to ${bucketName}/${nameFile}`);
    } catch (error) {
      console.error('Error uploading file:', error);
      throw error;
    }

    return {
      message: 'File uploaded successfully',
      url: `${process.env.AWS_S3_ENDPOINT}/${bucketName}/${nameFile}`,
    };
  }
}
