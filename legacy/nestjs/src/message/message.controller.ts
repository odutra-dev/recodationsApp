import {
  Controller,
  Get,
  Post,
  Body,
  Patch,
  Param,
  Delete,
  HttpCode,
  HttpStatus,
  UseInterceptors,
  UploadedFile,
  HttpException,
} from '@nestjs/common';
import { MessageService } from './message.service';
import { CreateMessageDto } from './dto/create-message.dto';
import { UpdateMessageDto } from './dto/update-message.dto';
import { FileInterceptor } from '@nestjs/platform-express';
import { ApiBody, ApiConsumes } from '@nestjs/swagger';

@Controller('message')
export class MessageController {
  constructor(private readonly messageService: MessageService) {}

  @HttpCode(HttpStatus.CREATED)
  @Post()
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    schema: {
      type: 'object',
      required: ['content', 'name'],
      properties: {
        name: {
          type: 'string',
          description: 'Name of the message sender',
          example: 'John Doe',
        },
        content: {
          type: 'string',
          description: 'Content of the message',
          example: 'Greetings friend! This is a message to your Wedding.',
        },
        file: {
          type: 'string',
          format: 'binary',
          description: 'Your image file for the message (optional)',
        },
      },
    },
  })
  @UseInterceptors(
    FileInterceptor('file', {
      fileFilter: (req, file, cb) => {
        if (!file.mimetype.startsWith('image/')) {
          return cb(
            new HttpException(
              'Only image files are allowed',
              HttpStatus.BAD_REQUEST,
            ),
            false,
          );
        }
        cb(null, true);
      },
    }),
  )
  create(
    @Body() createMessageDto: CreateMessageDto,
    @UploadedFile() file: Express.Multer.File,
  ) {
    return this.messageService.create(createMessageDto, file);
  }

  @Get()
  findAll() {
    return this.messageService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.messageService.findOne(+id);
  }

  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    schema: {
      type: 'object',
      required: ['content', 'name'],
      properties: {
        name: {
          type: 'string',
          description: 'Name of the message sender',
          example: 'John Doe',
        },
        content: {
          type: 'string',
          description: 'Content of the message',
          example: 'Greetings friend! This is a message to your Wedding.',
        },
        file: {
          type: 'string',
          format: 'binary',
          description: 'Your image file for the message (optional)',
        },
      },
    },
  })
  @UseInterceptors(
    FileInterceptor('file', {
      fileFilter: (req, file, cb) => {
        if (!file.mimetype.startsWith('image/')) {
          return cb(
            new HttpException(
              'Only image files are allowed',
              HttpStatus.BAD_REQUEST,
            ),
            false,
          );
        }
        cb(null, true);
      },
    }),
  )
  @Patch(':id')
  update(
    @Param('id') id: string,
    @Body() updateMessageDto: UpdateMessageDto,
    @UploadedFile() file: Express.Multer.File,
  ) {
    return this.messageService.update(+id, updateMessageDto, file);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.messageService.remove(+id);
  }
}
