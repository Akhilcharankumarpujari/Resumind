<div align="center">
  <br />
    <img src="public/readme/hero.webp" alt="Project Banner">
  <br />

  <div>
    <img alt="React" src="https://img.shields.io/badge/React-4c84f3?style=for-the-badge&logo=react&logoColor=white">
    <img alt="Go" src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white">
    <img alt="Fiber" src="https://img.shields.io/badge/Fiber-000000?style=for-the-badge&logo=gofiber&logoColor=white">
    <img alt="MySQL" src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white">
    <img alt="Docker" src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white">
    <img alt="Tailwind" src="https://img.shields.io/badge/TailwindCSS-38B2AC?style=for-the-badge&logo=tailwindcss&logoColor=white">
    <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white">
  </div>

  <h3 align="center">Resumind - Production-Ready AI Resume Analyzer & ATS Optimizer</h3>
  <p align="center">Migrated from Puter.js to a secure Go Fiber backend with MySQL GORM, Google Sign-In, and Groq Llama 3.3 AI.</p>
</div>

---

## 📋 Table of Contents

1. ✨ [Introduction](#introduction)
2. ⚙️ [Tech Stack](#tech-stack)
3. 🔋 [Features](#features)
4. 🤸 [Quick Start (Local Setup)](#quick-start)
5. 🐳 [Docker Compose (Single Command Launch)](#docker-compose)
6. 🚀 [Deployment to Cloud](#deployment)

---

## ✨ Introduction

**Resumind** is a production-ready web application that helps candidates optimize their resumes for Applicant Tracking Systems (ATS) and job descriptions.

This repository features a fully refactored, robust architecture:
- **Go backend** (using **Fiber**) secures API keys and handles heavy tasks like PDF parsing.
- **MySQL (GORM)** manages user profiles, database tables, and historical data logs.
- **Google OAuth 2.0** provides secure login.
- **Groq API (Llama 3.3 70B)** parses and rates resumes in under 3 seconds.

---

## ⚙️ Tech Stack

### Backend
* **Go (Golang)**: Standard compiled web runtime.
* **Fiber (v2)**: High-performance routing engine.
* **MySQL 8.0 & GORM ORM**: Automated schema migrations and data queries.
* **Google token verification**: Authenticates client-side Google sign-in tokens on the server.
* **JWT & HttpOnly Cookies**: Issues signed sessions that JS cannot access (XSS safe).
* **ledongthuc/pdf**: Full-text, multi-page PDF coordinate text extraction.
* **Groq LPU Inference Cloud**: Integrates `llama-3.3-70b-versatile` in strict JSON mode.

### Frontend
* **React Router v7**: Runs in Server-Side Rendering (SSR) mode.
* **Zustand**: Clean client session state engine.
* **Tailwind CSS v4**: Styling framework.

---

## 🔋 Features

* **Secure Authentication**: Google Sign-In with server-side validation and secure HTTP-Only JWT cookies.
* **Multi-page Resume Parsing**: Extracts plaintext characters from multi-page PDFs on the server.
* **Interactive Live Mock Report**: Allows logged-out visitors to play with an interactive preview card.
* **Personalized Dashboard**: View analysis history and manage previously uploaded documents.
* **Granular AI Evaluations**: Generates ATS score, layout, formatting, skills gap, keyword density, and tone tips.
* **Optimized Client Bundle**: Removed heavy client-side converters (`pdfjs-dist`), reducing bundle size by **1.03 MB**.

---

## 🤸 Quick Start (Local Setup)

### Prerequisites
Make sure you have the following installed locally:
- [Git](https://git-scm.com/)
- [Node.js](https://nodejs.org/en) (v18+)
- [Go (Golang)](https://go.dev/) (v1.22+)
- [MySQL Server](https://www.mysql.com/)

---

### Step 1: Clone & Install Dependencies

1. Clone the repository:
   ```bash
   git clone https://github.com/Akhilcharankumarpujari/Resumind.git
   cd Resumind
   ```

2. Install frontend dependencies:
   ```bash
   npm install
   ```

3. Install Go backend dependencies:
   ```bash
   cd backend
   go mod tidy
   cd ..
   ```

---

### Step 2: Database Setup

1. Open your MySQL client and create a new schema database:
   ```sql
   CREATE DATABASE resumind CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

---

### Step 3: Configure Environment variables

1. **Backend Environment**:
   Create a `.env` file in the `backend/` directory:
   ```env
   MYSQL_DSN=root:yourpassword@tcp(127.0.0.1:3306)/resumind?parseTime=true
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   GROQ_API_KEY=your-groq-api-key
   JWT_SECRET=your-long-random-jwt-signing-secret
   PORT=8080
   UPLOADS_DIR=./uploads
   ```

2. **Frontend Environment**:
   Create a `.env` file in the root directory:
   ```env
   VITE_GOOGLE_CLIENT_ID=your-google-client-id
   ```

---

### Step 4: Run the Application

Open two terminals:

* **Terminal 1 (Go Backend)**:
  ```bash
  cd backend
  go run .
  ```
  *(Go server will start on port `8080` and auto-migrate MySQL tables).*

* **Terminal 2 (Vite Frontend)**:
  ```bash
  npm run dev
  ```
  *(Vite server will start on port `5173` with dynamic routing to `/api` requests).*

---

## 🐳 Docker Compose (Single Command Launch)

If you have Docker installed, you can start the entire stack (Database, Nginx Proxy, Go Backend, and React Frontend) with a single command.

1. Create a `.env` file in the root directory:
   ```env
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   VITE_GOOGLE_CLIENT_ID=your-google-client-id
   GROQ_API_KEY=your-groq-api-key
   JWT_SECRET=your-random-jwt-key
   ```

2. Run the compose environment:
   ```bash
   docker-compose up -d --build
   ```

3. Open **`http://localhost`** in your browser. All requests to `/api` will be proxied automatically via Nginx.

---

## 🚀 Deployment to Cloud

For detailed instructions on cloud host deployments (using Render, Fly.io, or AWS EC2), check out the [Production Deployment Guide](backend/../.gemini/antigravity/brain/e5c90a6b-3f1b-4466-b9d1-0825c64999c7/deployment_guide.md).
