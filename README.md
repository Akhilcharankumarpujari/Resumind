# 🎯 Resumind — AI-Powered Resume Analyzer & ATS Optimizer

[![React](https://img.shields.io/badge/React-19.0-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![React Router](https://img.shields.io/badge/React_Router-v7-CA4245?logo=reactrouter&logoColor=white)](https://reactrouter.com/)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![Go Fiber](https://img.shields.io/badge/Fiber-v2.52-00ACD7?logo=gofiber&logoColor=white)](https://gofiber.io/)
[![GORM](https://img.shields.io/badge/GORM-v1.31-008080?logo=go&logoColor=white)](https://gorm.io/)
[![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)](https://www.mysql.com/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?logo=tailwindcss&logoColor=white)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Resumind** is a full-stack, AI-powered resume analysis and applicant tracking system (ATS) optimization platform. Upload any resume PDF alongside a job description to receive instant, actionable feedback across ATS compatibility, content strength, tone, structure, and skill alignment—driven by ultra-fast LLM inference via Groq AI.

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
- ⚡ **Production-Grade Architecture**: Designed for modern decoupled deployments (Vercel Frontend + Render Backend + Aiven MySQL).

---

## 🏗️ Architecture Overview

```
                      ┌───────────────────────────────┐
                      │    Client / Web Browser       │
                      └──────────────┬────────────────┘
                                     │
                 ┌───────────────────┴───────────────────┐
                 │                                       │
                 ▼                                       ▼
     ┌──────────────────────┐                ┌──────────────────────┐
     │  Vercel Frontend     │                │   Render Backend     │
     │  (React Router v7)   │─── API Calls ──►   (Go / Fiber v2)    │
     └──────────────────────┘                └──────────┬───────────┘
                                                        │
                         ┌──────────────────────────────┼──────────────────────────────┐
                         ▼                              ▼                              ▼
              ┌─────────────────────┐        ┌─────────────────────┐        ┌─────────────────────┐
              │     Groq AI API     │        │    Aiven MySQL DB   │        │ Local Uploads Dir   │
              │  (Llama-3 / GPT)    │        │  (TLS / GORM ORM)   │        │   (PDF Storage)     │
              └─────────────────────┘        └─────────────────────┘        └─────────────────────┘
```

---

## 🛠️ Tech Stack

### Frontend
- **Framework**: React 19 + React Router v7 (SPA / Server-capable)
- **Build Tool**: Vite 6
- **Styling**: Tailwind CSS v4 + `tw-animate-css`
- **State Management**: Zustand
- **File Upload**: `react-dropzone`
- **Icons**: Lucide React

### Backend
- **Language**: Go (Golang)
- **Web Framework**: Fiber v2 (FastHTTP engine)
- **ORM & Database**: GORM + `go-sql-driver/mysql` (MySQL 8.0 / Managed Aiven)
- **Authentication**: JWT (`golang-jwt/jwt/v5`) + Google Identity API (`google.golang.org/api/idtoken`)
- **PDF Extraction**: `github.com/ledongthuc/pdf`
- **AI Engine**: Groq Client (`github.com/sashabaranov/go-openai`)

---

## 🚀 Getting Started Locally

### Prerequisites
- **Node.js**: v20 or higher
- **Go**: 1.22 or higher
- **MySQL**: Local MySQL server (port 3306/3307) or Docker
- **Groq API Key**: Get a free key at [console.groq.com](https://console.groq.com/keys)
- **Google OAuth Client ID**: Get one at [console.cloud.google.com](https://console.cloud.google.com/apis/credentials)

---

### 1. Clone the Repository

```bash
git clone https://github.com/Akhilcharankumarpujari/Resumind.git
cd Resumind
```

---

### 2. Configure Environment Variables

#### Backend (`backend/.env`)
Create a `.env` file in the `backend/` directory:

```ini
# MySQL Connection (Local Docker or MySQL)
MYSQL_DSN=root:akhil123@tcp(127.0.0.1:3307)/resumind?parseTime=true

# Google OAuth 2.0 Credentials
GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Groq AI API Key
GROQ_API_KEY=gsk_your_groq_api_key_here
GROQ_MODEL=openai/gpt-oss-120b

# JWT Authentication Secret (Minimum 32 characters)
JWT_SECRET=your-32-character-secret-key-goes-here

# Server Settings
PORT=8081
UPLOADS_DIR=./uploads
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

#### Frontend (`.env`)
Create a `.env` file in the project root:

```ini
VITE_GOOGLE_CLIENT_ID=your-google-client-id.apps.googleusercontent.com
VITE_API_URL=http://localhost:8081
```

---

### 3. Run Database & Backend

#### Option A: Using Local Docker Compose (Recommended)
```bash
docker-compose up -d db
cd backend
go run .
```

#### Option B: Manual Setup
Make sure your local MySQL instance is running and create the `resumind` database:
```sql
CREATE DATABASE IF NOT EXISTS resumind;
```
Then run the Go backend:
```bash
cd backend
go run .
```
The backend will run on `http://localhost:8081` and automatically migrate database schemas.

---

### 4. Run Frontend

In a new terminal window at the project root:

```bash
npm install
npm run dev
```

Open `http://localhost:5173` in your browser.

---

## 🐳 Running Full Stack with Docker Compose

To start the Database, Backend, Frontend, and Nginx reverse proxy together:

```bash
docker-compose up --build
```

Access the application at `http://localhost`.

---

## 🌐 Production Deployment Guide

Resumind is configured for a multi-cloud production architecture:

### 1. Backend → [Render](https://render.com)
- **Environment**: Go runtime
- **Root Directory**: `backend`
- **Build Command**: `go build -o server .`
- **Start Command**: `./server`
- **Environment Variables**:
  - `MYSQL_DSN`: `user:password@tcp(host:22026)/defaultdb?parseTime=true&tls=skip-verify`
  - `GROQ_API_KEY`: `gsk_...`
  - `GOOGLE_CLIENT_ID`: `...apps.googleusercontent.com`
  - `JWT_SECRET`: Generate a secure 32+ character random string
  - `ALLOWED_ORIGINS`: `https://your-app.vercel.app`

### 2. Frontend → [Vercel](https://vercel.com)
- **Framework Preset**: Vite / React Router
- **Build Command**: `npm run build`
- **Output Directory**: `build/client`
- **Environment Variables**:
  - `VITE_API_URL`: `https://resumind-backend.onrender.com`
  - `VITE_GOOGLE_CLIENT_ID`: `...apps.googleusercontent.com`

### 3. Database → [Aiven MySQL](https://aiven.io)
- Ensure your Aiven MySQL instance enables SSL (`tls=skip-verify` in `MYSQL_DSN`).
- `backend/db.go` automatically normalizes `tls=true` and `tls=skip-verify` configurations for managed cloud CAs.

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

Run backend unit and integration tests (DSN parsing, network validation, TLS configuration):

```bash
cd backend
go test -v ./...
```

Run TypeScript typechecks on the frontend:

```bash
npm run typecheck
```

---

## 📋 Environment Variables Summary

| Variable | Scope | Description | Required |
| :--- | :--- | :--- | :---: |
| `MYSQL_DSN` | Backend | MySQL connection string (`user:pass@tcp(host:port)/dbname?parseTime=true`) | Yes |
| `GROQ_API_KEY` | Backend | Groq Cloud AI API key | Yes |
| `GROQ_MODEL` | Backend | AI Model (default: `openai/gpt-oss-120b`) | No |
| `GOOGLE_CLIENT_ID` | Backend/Frontend | Google OAuth 2.0 Client ID | Yes |
| `GOOGLE_CLIENT_SECRET` | Backend | Google OAuth 2.0 Client Secret | Yes |
| `JWT_SECRET` | Backend | Secret string for signing session JWT tokens | Yes |
| `VITE_API_URL` | Frontend | Full URL of Go backend API | Yes (Prod) |
| `ALLOWED_ORIGINS` | Backend | CORS allowed origin origins list (comma-separated) | Yes (Prod) |

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request
