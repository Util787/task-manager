# Task Manager API

RESTful API for task management built with Go and Gin framework.

---

## Quick Start

### Prerequisites
- [Go 1.24.3+](https://golang.org/doc/install)

### 1. Clone Repository 📂
```bash
git clone https://github.com/Util787/task-manager.git
cd task-manager
```

### 2. Configure `.env` ⚙️
Create a `.env` file and configure according to your environment (`prod`, `dev`, or `local`):

```env
ENV=local
HTTP_PORT=8080
HTTP_READ_HEADER_TIMEOUT=5s
HTTP_WRITE_TIMEOUT=10s
HTTP_READ_TIMEOUT=10s 
```

### 3. Run the Application ▶️
```bash
go run cmd/main.go
```

## API Documentation
Swagger UI: `http://localhost:8080/swagger/index.html`

