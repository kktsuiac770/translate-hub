# TranslateHub

TranslateHub is a web service that helps manage translation projects using Google's Gemini AI model for initial translations with support for human review and corrections.

## Features

- Project Management
  - Create translation projects with source and target languages
  - List and manage multiple translation projects
- Task Management
  - Upload text files for translation
  - Automatic initial translation using Gemini AI
  - Track task status and progress
- Translation Workflow
  - Submit translation changes
  - Review and approve/reject changes
  - Maintain translation history

## API Endpoints

### Projects
- `POST /projects` - Create a new translation project
- `GET /projects` - List all projects
- `GET /projects/{id}` - Get a specific project

### Tasks
- `POST /projects/{id}/tasks` - Create a new translation task
- `GET /projects/{id}/tasks` - List tasks for a project
- `GET /tasks` - List all tasks
- `GET /tasks/{id}` - Get a specific task

### Translation Changes
- `POST /changes` - Submit a translation change
- `POST /review` - Review (approve/reject) a translation change

## Setup

1. Set environment variables:
```sh
GOOGLE_GENAI_API_KEY=your_api_key
```

2. Start the server:
```sh
go run main.go
```

The server will start on port 8080.

## Example Usage

Create a new project:
```sh
curl -X POST http://localhost:8080/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "Japanese Translation","source_lang": "en","target_lang": "jp"}'
```

Upload a file for translation:
```sh
curl -X POST http://localhost:8080/projects/1/tasks \
  -F "file=@dialogue.txt" \
  -F "name=Chapter 1" \
  -F "creator=alice"
```

## Tech Stack

- Go (backend)
- PostgreSQL (database)
- Google Gemini AI (machine translation)


## TODO

### Authentication & Authorization
- [ ] User authentication system (JWT/OAuth2)
- [ ] Role-based access control (Admin/Manager/Translator)
- [ ] Session management
- [ ] Password reset flow

### User Experience
- [ ] Real-time notifications system
  - [ ] WebSocket integration
  - [ ] Email notifications
  - [ ] In-app notifications
- [ ] Translation UI improvements
  - [ ] Rich text editor
  - [ ] Translation memory
  - [ ] Side-by-side comparison view
- [ ] Changes tracking system
  - [ ] Version history
  - [ ] Diff viewer
  - [ ] Comment system

### Monitoring & Operations
- [ ] Logging system
  - [ ] Structured logging
  - [ ] Log aggregation (ELK/Loki)
- [ ] Metrics collection
  - [ ] Prometheus integration
  - [ ] Grafana dashboards
  - [ ] Performance monitoring
- [ ] Error tracking
  - [ ] Error reporting service integration
  - [ ] Alert system

### Deployment & Infrastructure
- [ ] Containerization
  - [ ] Docker setup
  - [ ] Docker Compose for local development
- [ ] CI/CD pipeline
  - [ ] Automated testing
  - [ ] Build automation
  - [ ] Deployment automation
- [ ] Infrastructure as Code
  - [ ] Kubernetes manifests
  - [ ] Terraform configurations
- [ ] Environment management
  - [ ] Development
  - [ ] Staging
  - [ ] Production

### Additional Features
- [ ] API rate limiting
- [ ] Data backup and recovery
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Translation quality metrics
- [ ] Batch processing for large files
- [ ] Export/Import functionality
- [ ] Team collaboration features
- [ ] Project templates
- [ ] Custom workflow definitions