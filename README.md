# Recodations API

A robust RESTful API for managing messages with integrated image storage using AWS S3. Built with NestJS, PostgreSQL, and LocalStack for local S3 simulation.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Setup](#environment-setup)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Database](#database)
- [S3 Integration](#s3-integration)
- [Development](#development)

## Overview

Recodations API is a NestJS-based application designed to manage messages with support for image attachments. It provides a complete CRUD (Create, Read, Update, Delete) interface for messages and handles file uploads securely to AWS S3 (or LocalStack for local development).

The application is particularly suited for scenarios like wedding guestbook applications, event message boards, or any system requiring message management with media attachments.

## Features

- ✅ **Message Management**: Create, read, update, and delete messages
- 📸 **Image Upload**: Upload images with messages to S3 storage
- 🗄️ **PostgreSQL Database**: Persistent data storage
- 🏗️ **Docker Support**: Pre-configured Docker Compose for easy local development
- 🧪 **Comprehensive Testing**: Unit tests and E2E tests with Jest
- 📚 **Swagger Documentation**: Auto-generated API documentation
- ✔️ **Data Validation**: Request validation using class-validator
- 🔒 **Type Safety**: Full TypeScript support

## Technology Stack

- **Runtime**: Node.js
- **Framework**: NestJS 11.x
- **Language**: TypeScript 5.x
- **Database**: PostgreSQL 15+
- **ORM**: TypeORM
- **Storage**: AWS SDK S3 Client
- **Testing**: Jest, Supertest
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker, Docker Compose
- **Code Quality**: ESLint, Prettier
- **Validation**: class-validator, class-transformer

## Prerequisites

Before you begin, ensure you have installed:

- **Node.js** (v18 or higher) - [Download](https://nodejs.org/)
- **npm** (v9 or higher) - Comes with Node.js
- **Docker** and **Docker Compose** - [Download](https://www.docker.com/products/docker-desktop)

Verify installations:

```bash
node --version
npm --version
docker --version
docker-compose --version
```

## Installation

1. **Clone the repository**

```bash
git clone https://github.com/odutra-dev/recodationsApp
cd recodationsApp
```

2. **Install dependencies**

```bash
npm install
```

3. **Start Docker services** (PostgreSQL and LocalStack)

```bash
docker-compose up -d
```

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```bash
PORT=3000
AWS_REGION='us-east-1'
AWS_ACCESS_KEY_ID='test'
AWS_SECRET_ACCESS_KEY='test'
AWS_S3_BUCKET_NAME='my-bucket'
AWS_S3_ENDPOINT='http://localhost:4566'

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=recodations
```

### Environment Variables Explanation

| Variable                | Description         | Example                              |
| ----------------------- | ------------------- | ------------------------------------ |
| `PORT`                  | Server port         | `3000`                               |
| `DB_HOST`               | PostgreSQL host     | `localhost`                          |
| `DB_PORT`               | PostgreSQL port     | `5432`                               |
| `DB_USER`               | PostgreSQL username | `postgres`                           |
| `DB_PASS`               | PostgreSQL password | `postgres`                           |
| `DB_NAME`               | Database name       | `recodations`                        |
| `AWS_REGION`            | AWS region          | `us-east-1`                          |
| `AWS_ACCESS_KEY_ID`     | AWS access key      | `test` (for LocalStack)              |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key      | `test` (for LocalStack)              |
| `AWS_S3_ENDPOINT`       | S3 endpoint URL     | `http://localhost:4566` (LocalStack) |

## Running the Application

### Development Mode

```bash
npm run start:dev
```

Starts the application in watch mode with hot reloading. Any file changes will automatically restart the server.

### Standard Start

```bash
npm run start
```

Runs the application once without watch mode.

### Debug Mode

```bash
npm run start:debug
```

Starts the application with Node.js debugger enabled on port 9229.

### Production Build

```bash
npm run build
```

Compiles TypeScript to JavaScript in the `dist/` directory.

### Run Production Build

```bash
npm run start:prod
```

Runs the compiled production build from the `dist/` directory.

## API Documentation

Once the application is running, access the interactive API documentation:

- **Swagger UI**: [http://localhost:3000/api](http://localhost:3000/api)
- **OpenAPI JSON**: [http://localhost:3000/swagger/json](http://localhost:3000/swagger/json)

### Available Endpoints

#### Create Message

```http
POST /message
Content-Type: multipart/form-data

{
  "name": "John Doe",
  "content": "Greetings friend! This is a message.",
  "file": <binary-image-file> (optional)
}
```

Response (201 Created):

```json
{
  "id": 1,
  "name": "John Doe",
  "content": "Greetings friend! This is a message.",
  "image": "http://localhost:4566/photos/john-doe-1617286800000.jpg"
}
```

#### Get All Messages

```http
GET /message
```

Response (200 OK):

```json
[
  {
    "id": 1,
    "name": "John Doe",
    "content": "Greetings friend!",
    "image": "http://localhost:4566/photos/john-doe-1617286800000.jpg"
  },
  {
    "id": 2,
    "name": "Jane Smith",
    "content": "Congratulations!",
    "image": null
  }
]
```

#### Get Single Message

```http
GET /message/:id
```

Response (200 OK):

```json
{
  "id": 1,
  "name": "John Doe",
  "content": "Greetings friend!",
  "image": "http://localhost:4566/photos/john-doe-1617286800000.jpg"
}
```

#### Update Message

```http
PATCH /message/:id
Content-Type: multipart/form-data

{
  "name": "John Doe Updated",
  "content": "Updated content",
  "file": <binary-image-file> (optional)
}
```

Response (200 OK):

```json
{
  "message": "Message updated successfully"
}
```

#### Delete Message

```http
DELETE /message/:id
```

Response (200 OK):

```json
{
  "message": "Message deleted successfully"
}
```

## Project Structure

```text
recodationsApp/
├── src/
│   ├── main.ts                    # Application entry point
│   ├── app.module.ts              # Root module with configuration
│   ├── app.controller.ts          # Health check endpoint
│   ├── app.service.ts             # App-level service
│   ├── message/
│   │   ├── message.module.ts      # Message module
│   │   ├── message.controller.ts  # Message HTTP handlers
│   │   ├── message.service.ts     # Message business logic
│   │   ├── entities/
│   │   │   └── message.entity.ts  # Message database entity
│   │   ├── dto/
│   │   │   ├── create-message.dto.ts
│   │   │   └── update-message.dto.ts
│   │   └── message.service.spec.ts # Unit tests
│   └── s3/
│       ├── s3.service.ts          # S3 file upload service
│       └── s3.service.spec.ts     # S3 service tests
├── test/
│   ├── app.e2e-spec.ts           # End-to-end tests
│   └── jest-e2e.json             # E2E test configuration
├── docker-compose.yml             # Docker services configuration
├── package.json                   # Dependencies and scripts
├── tsconfig.json                  # TypeScript configuration
├── eslint.config.mjs              # ESLint configuration
└── README.md                      # This file
```

### Key Directories Explained

**`src/`** - Main application source code

- Contains all production code organized by feature modules

**`message/`** - Message feature module

- `entities/` - Database models (Message entity)
- `dto/` - Data Transfer Objects for validation
- `controller.ts` - HTTP request handlers
- `service.ts` - Business logic and database operations

**`s3/`** - AWS S3 integration service

- Handles file uploads to S3 or LocalStack
- Creates buckets automatically if they don't exist

**`test/`** - Test files

- Contains E2E tests and test configurations

## Testing

### Unit Tests

Run all unit tests:

```bash
npm run test
```

Watch mode (re-run on file changes):

```bash
npm run test:watch
```

### Coverage Report

Generate test coverage report:

```bash
npm run test:cov
```

Coverage report is generated in `coverage/` directory.

### E2E Tests

Run end-to-end tests:

```bash
npm run test:e2e
```

### Test Structure

The project includes comprehensive tests:

- **Message Service Tests** (`message.service.spec.ts`):
  - ✅ Create message with/without file
  - ✅ Retrieve all messages
  - ✅ Retrieve single message
  - ✅ Update message with/without file
  - ✅ Delete message
  - ✅ Error handling (missing ID, not found, etc.)

- **S3 Service Tests** (`s3.service.spec.ts`):
  - ✅ File upload functionality
  - ✅ Bucket creation and verification
  - ✅ Error handling

- **E2E Tests** (`app.e2e-spec.ts`):
  - ✅ Full API workflow testing
  - ✅ Integration testing with actual database

## Database

### PostgreSQL

The application uses PostgreSQL as its primary data store.

#### Database Configuration

Connection details are defined in `app.module.ts` using environment variables:

```typescript
TypeOrmModule.forRoot({
  type: 'postgres',
  host: process.env.DB_HOST, // localhost
  port: parseInt(process.env.DB_PORT), // 5432
  username: process.env.DB_USER, // postgres
  password: process.env.DB_PASS, // postgres
  database: process.env.DB_NAME, // recodations
  autoLoadEntities: true,
  synchronize: true,
  entities: [Message],
});
```

#### Auto-migration

TypeORM's `synchronize: true` automatically creates/updates database schema based on entities. This is convenient for development but **should be disabled in production** to prevent accidental data loss.

#### Message Entity

```typescript
@Entity()
export class Message {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string; // Sender name (max 255 chars)

  @Column()
  content: string; // Message content

  @Column({ type: 'text', nullable: true })
  image: string | null; // S3 image URL
}
```

### Database Access

Connect to PostgreSQL directly:

```bash
docker exec -it Local-Postgres psql -U postgres -d recodations
```

Then use SQL to explore:

```sql
SELECT * FROM message;
```

## S3 Integration

### S3 Overview

The application integrates with AWS S3 for storing message images. For local development, **LocalStack** provides a mock S3 service.

### S3Service Implementation

The `S3Service` handles:

1. **Bucket Management**
   - Automatically checks if bucket exists
   - Creates bucket if it doesn't exist
   - Logs bucket creation

2. **File Upload**
   - Generates unique filename (original name + timestamp)
   - Uploads to specified bucket
   - Returns accessible URL

3. **Error Handling**
   - Catches upload errors
   - Logs detailed error information
   - Throws error to caller

### LocalStack (Development S3)

LocalStack emulates AWS services locally. Configuration:

```yaml
localstack:
  image: localstack/localstack
  ports:
    - '4566:4566' # LocalStack endpoint
  environment:
    - SERVICES=s3 # Enable only S3
    - AWS_REGION=us-east-1
```

Access LocalStack S3:

```bash
# List buckets
aws s3 ls --endpoint-url http://localhost:4566

# List files in bucket
aws s3 ls s3://photos --endpoint-url http://localhost:4566
```

### Using AWS S3 in Production

To use real AWS S3 in production:

1. Update environment variables:

```bash
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-real-key
AWS_SECRET_ACCESS_KEY=your-real-secret
AWS_S3_ENDPOINT=https://s3.amazonaws.com  # Or your S3-compatible endpoint
```

2. Ensure AWS credentials have S3 permissions

3. Create S3 bucket in AWS Console or CLI:

```bash
aws s3 mb s3://photos --region us-east-1
```

## Development

### Code Quality

#### Format Code

```bash
npm run format
```

Formats all TypeScript code using Prettier.

#### Lint Code

```bash
npm run lint
```

Checks code for style issues and automatically fixes them.

### Building for Production

1. **Build the application**

```bash
npm run build
```

Output is in the `dist/` directory.

2. **Run production build**

```bash
npm run start:prod
```

3. **Docker considerations**
   - Build a Docker image: `docker build -t recodations:latest .`
   - Push to registry: `docker push <registry>/recodations:latest`
   - Deploy using Docker Compose or Kubernetes

### Adding New Features

#### Adding a New Endpoint

1. Create DTO in `src/feature/dto/`:

```typescript
export class CreateFeatureDto {
  @IsString()
  @IsNotEmpty()
  name: string;
}
```

2. Create Entity in `src/feature/entities/`:

```typescript
@Entity()
export class Feature {
  @PrimaryGeneratedColumn()
  id: number;

  @Column()
  name: string;
}
```

3. Create Service in `src/feature/feature.service.ts`:

```typescript
@Injectable()
export class FeatureService {
  // Business logic here
}
```

4. Create Controller in `src/feature/feature.controller.ts`:

```typescript
@Controller('feature')
export class FeatureController {
  constructor(private readonly featureService: FeatureService) {}
}
```

5. Create Module in `src/feature/feature.module.ts`:

```typescript
@Module({
  imports: [TypeOrmModule.forFeature([Feature])],
  controllers: [FeatureController],
  providers: [FeatureService],
})
export class FeatureModule {}
```

6. Add tests in `src/feature/feature.service.spec.ts`

### Debugging

#### Using VS Code Debugger

1. Start app in debug mode:

```bash
npm run start:debug
```

2. Open VS Code debug panel (Ctrl+Shift+D)

3. Run "Node" configuration (should connect to port 9229)

4. Set breakpoints in code

#### Using Console Logs

The project includes Console Ninja extension support for enhanced debugging in VS Code.

## Troubleshooting

### Common Issues

#### 1. Database Connection Error

```
Error: connect ECONNREFUSED 127.0.0.1:5432
```

**Solution:**

```bash
# Start Docker services
docker-compose up -d

# Verify PostgreSQL is running
docker ps | grep postgres
```

#### 2. S3/LocalStack Connection Error

```
Error: Network Timeout at LocalStack
```

**Solution:**

```bash
# Restart LocalStack
docker-compose restart localstack

# Verify LocalStack is running
docker logs localstack-main
```

#### 3. Port Already in Use

```
Error: listen EADDRINUSE: address already in use :::3000
```

**Solution:**

```bash
# Kill process on port 3000
lsof -ti:3000 | xargs kill -9

# Or use different port
PORT=3001 npm run start:dev
```

#### 4. Module Not Found

```
Cannot find module '@nestjs/...'
```

**Solution:**

```bash
# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

## License

UNLICENSED - See LICENSE file for details.

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Create a feature branch: `git checkout -b feature/my-feature`
2. Commit changes: `git commit -am 'Add my feature'`
3. Push to branch: `git push origin feature/my-feature`
4. Submit a Pull Request

## Support

For issues and questions, please create an issue in the repository or contact the maintainers.
