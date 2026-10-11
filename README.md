# 🎯 Resumind — AI-Powered Resume Analyzer & ATS Optimizer

[![React](https://img.shields.io/badge/React-19.0-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![React Router](https://img.shields.io/badge/React_Router-v7-CA4245?logo=reactrouter&logoColor=white)](https://reactrouter.com/)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Go Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?logo=gofiber&logoColor=white)](https://gofiber.io/)
[![GORM](https://img.shields.io/badge/GORM-v1.31-008080?logo=go&logoColor=white)](https://gorm.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)](https://docker.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Resumind** is a full-stack, AI-powered resume analysis and applicant tracking system (ATS) optimization platform structured as a clean production monorepo. Upload any resume PDF alongside a target job description to receive instant, actionable feedback across ATS compatibility, content strength, tone, structure, and skill alignment—powered by ultra-fast LLM inference via Groq AI.

---

## 📁 Repository Structure

```
resumind/
├── frontend/                     # React 19 + React Router v7 SPA
│   ├── app/                      # Application routes, components, and state
│   ├── constants/                # UI constants and styling presets
│   ├── public/                   # Static assets, icons, and branding
│   ├── types/                    # TypeScript interfaces and type declarations
│   ├── package.json              # Frontend dependencies and scripts
│   ├── package-lock.json         # Lockfile
│   ├── vite.config.ts            # Vite 6 config with /api proxy
│   ├── react-router.config.ts    # React Router SPA configuration
│   ├── tsconfig.json             # TypeScript configuration
│   ├── Dockerfile                # Frontend multi-stage container build
│   └── .env.example              # Frontend environment template
├── backend/                      # High-performance Go (Fiber v2) REST API
│   ├── *.go                      # API handlers, models, AI, DB, and auth logic
│   ├── go.mod                    # Go module definition
│   ├── go.sum                    # Go module checksums
│   ├── Dockerfile                # Backend multi-stage container build
│   └── .env.example              # Backend environment template
├── nginx/                        # Reverse proxy configuration
│   └── nginx.conf                # Unified routing for SPA and /api/
├── .github/                      # CI/CD Workflows
│   └── workflows/
│       └── ci.yml                # Monorepo automated test & build check
├── docker-compose.yml            # Multi-container orchestration (DB, Backend, Frontend, Proxy)
├── .dockerignore                 # Docker build exclusion rules
├── .gitignore                    # Monorepo Git ignore rules
├── .env.example                  # Root environment template for Docker Compose
└── README.md                     # Documentation
```

---

## ✨ Features

- 📄 **Instant PDF Text Extraction & AI Analysis**: Upload your PDF resume and target job description for automated extraction and LLM analysis.
- 📊 **Detailed Category Breakdown**:
  - **Overall Score**: Comprehensive match rating (0–100).
  - **ATS Readability**: Keyword optimization, formatting cleanliness, and parser safety.
  - **Tone & Style**: Professional voice, action verb usage, and impact phrasing.
  - **Content Quality**: Quantifiable achievements, clarity, and relevance.
  - **Structure & Formatting**: Section organization and logical layout.
  - **Skill Alignment**: Hard vs. soft skills missing from your target role.
- 💡 **Actionable Improvement Tips**: Specific suggestions categorized by priority (`critical`, `warning`, `suggestion`).
- 🔐 **Google OAuth 2.0 & Session Security**: Secure single-sign-on (SSO) with JWT authentication supported by both HTTP-Only cookies and Bearer tokens.
- 📁 **Resume History & Management**: View past resume scores, inspect detailed feedback reports, download uploaded PDFs, or delete records.
- 🐳 **Production-Ready Docker Orchestration**: Complete local and cloud deployment setup with Postgres, Go Backend, React SPA, and Nginx reverse proxy.

---

## 🛠️ Tech Stack

### Frontend (`frontend/`)
- **Framework**: React 19 + React Router v7 (SPA mode)
- **Build Tool**: Vite 6
- **Styling**: Tailwind CSS v4 + `tw-animate-css`
- **State Management**: Zustand
- **File Upload**: `react-dropzone`

### Backend (`backend/`)
- **Language**: Go (Golang 1.22+)
- **Web Framework**: Fiber v2 (FastHTTP engine)
- **ORM & Database**: GORM + `gorm.io/driver/postgres` (PostgreSQL 16+)
- **Authentication**: JWT (`golang-jwt/jwt/v5`) + Google Identity API
- **PDF Extraction**: `github.com/ledongthuc/pdf`
- **AI Inference**: Groq Cloud API (`openai/gpt-oss-120b` or Llama 3)

---

## 🚀 Local Development

### 1. Prerequisites
- **Node.js**: v20 or higher
- **Go**: v1.22 or higher
- **Docker & Docker Compose**: Optional, for containerized run
- **Groq API Key**: [console.groq.com/keys](https://console.groq.com/keys)
- **Google OAuth Client ID**: [console.cloud.google.com/apis/credentials](https://console.cloud.google.com/apis/credentials)

---

### 2. Running Locally (Separated Services)

#### Step A: Start Database & Backend
```bash
# 1. Start a local PostgreSQL instance (or use Docker):
docker run --name resumind-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=resumind -p 5432:5432 -d postgres:16

# 2. Configure backend environment:
cd backend
cp .env.example .env
# Edit .env with your GROQ_API_KEY and GOOGLE credentials

# 3. Start the Go backend:
go run .
# Backend starts listening on http://localhost:8082
```

#### Step B: Start Frontend
```bash
# 1. In a separate terminal:
cd frontend
cp .env.example .env

# 2. Install dependencies & run development server:
npm ci
npm run dev
# Frontend starts on http://localhost:5173 with Vite proxying /api to http://localhost:8082
```

---

## 🐳 Running Full Stack with Docker Compose

To start PostgreSQL, Backend, Frontend, and Nginx proxy with a single command:

```bash
# 1. Create root .env file from template:
cp .env.example .env
# Fill in POSTGRES_PASSWORD, GROQ_API_KEY, and Google OAuth credentials

# 2. Build and launch all services:
docker compose up --build -d

# 3. View running containers:
docker compose ps
```

The stack will be available at:
- **Application (Proxy)**: `http://localhost`
- **API Health Check**: `http://localhost/health`
- **PostgreSQL**: `127.0.0.1:5432`

To stop the containers:
```bash
docker compose down
```

---

## 🌐 AWS EC2 Production Deployment

Resumind runs in production on AWS EC2 (`resumind-server`) at **`http://52.23.111.29.nip.io/`**.

> **Note**: A separate project (ReachInbox) runs on EC2 instance `3.89.202.15` under security group `reachinbox-sg`. Do not modify or restart that instance.

### Safe Deployment Procedure to `resumind-server`

#### 1. Build the Production Frontend
On your local machine or CI:
```bash
cd frontend
npm ci
npm run build
# Generates static assets in frontend/build/client/
```

#### 2. Deploy Frontend Assets to EC2
```bash
# Backup existing web directory on EC2:
ssh -i resumind-key.pem ec2-user@52.23.111.29 "sudo cp -r /var/www/resumind/dist /var/www/resumind/dist.bak"

# Upload new build:
scp -i resumind-key.pem -r frontend/build/client ec2-user@52.23.111.29:/tmp/dist

# Replace and set correct permissions:
ssh -i resumind-key.pem ec2-user@52.23.111.29 "sudo rm -rf /var/www/resumind/dist && \
  sudo cp -r /tmp/dist /var/www/resumind/dist && \
  sudo chown -R nginx:nginx /var/www/resumind/dist && \
  sudo find /var/www/resumind/dist -type d -exec chmod 755 {} + && \
  sudo find /var/www/resumind/dist -type f -exec chmod 644 {} + && \
  sudo systemctl reload nginx"
```

#### 3. Update & Restart Backend (if Go source changed)
```bash
# Upload backend updates and restart systemd service:
ssh -i resumind-key.pem ec2-user@52.23.111.29 "sudo systemctl restart resumind && sudo systemctl status resumind"
```

#### 4. Verify Live Production
```bash
curl -i http://52.23.111.29.nip.io/health
curl -i http://52.23.111.29.nip.io/
```

---

## 📡 API Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :---: |
| `GET` | `/health` | Service health status check | ❌ |
| `POST` | `/api/auth/google` | Authenticate via Google OAuth ID token | ❌ |
| `GET` | `/api/auth/me` | Retrieve authenticated user profile | ✅ |
| `POST` | `/api/auth/logout` | Terminate session & clear cookies | ❌ |
| `POST` | `/api/analyze` | Upload resume PDF + job details for AI feedback | ✅ |
| `GET` | `/api/resumes` | Fetch list of all analyzed resumes for user | ✅ |
| `GET` | `/api/resume/:id` | Fetch specific resume analysis details | ✅ |
| `GET` | `/api/resume/:id/download` | Download uploaded resume PDF file | ✅ |
| `DELETE` | `/api/resume/:id` | Delete resume record and associated file | ✅ |

---

## 🧪 Testing

### Backend Unit & Integration Tests
```bash
cd backend
go test -v ./...
```

### Frontend Typechecking & Verification
```bash
cd frontend
npm run typecheck
npm run build
```

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
