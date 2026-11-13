import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsString, Length } from 'class-validator';

export class CreateMessageDto {
  @IsString()
  @IsNotEmpty()
  @Length(1, 255)
  @ApiProperty({
    description: 'Name of the message sender',
    type: String,
    example: 'John Doe',
  })
  name: string;

  @IsString()
  @IsNotEmpty()
  @ApiProperty({
    description: 'Content of the message',
    type: String,
    example: 'Greetings friend! This is a message to your Wedding.',
  })
  content: string;
}
