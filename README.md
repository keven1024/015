<div align="center"><a name="readme-top"></a>

# 015

015 (/ˈzɪərəʊ wʌn faɪv/, "zero-one-five") is a self-hosted temporary file sharing platform. Focused on providing one-time, temporary file and text upload, processing, and sharing services. The project name originates from [Ichigo](https://darling-in-the-franxx.fandom.com/wiki/Ichigo) from DARLING in the FRANXX.

A modern file sharing website built with Vue 3 + Nuxt 3 + Go, supporting file upload, text sharing, image compression, concurrent processing, instant transfer functionality, and more, featuring a complete sharing management and access control system.

![015 Platform Overview](/.github/image/0.png)

English | [中文](README-zh.md)

</div>

## 🌟 Features

### Core Functionality

🖼️ **High-Performance File Upload** - Supports large file chunked uploads with frontend file hash calculation for instant transfer  
📱 **Responsive Design** - Modern UI based on Tailwind V4 + Reka UI, adapts to various devices  
⚡ **Concurrent Processing** - Uses Web Worker for frontend hash calculation, backend queue system for task processing  
🌐 **Multi-language Support** - Complete Chinese and English internationalization support  
🔗 **Share Management** - Flexible sharing link generation and management system

### File Processing

🔄 **Smart Instant Transfer** - Frontend instant transfer detection based on file hash + file size, avoiding duplicate uploads  
📷 **Image Compression** - Automatic image compression functionality supporting multiple formats  
🖼️ **File Preview** - Supports preview of images, videos, audio, documents, and various file types  
📊 **Upload Statistics** - Real-time display of upload progress and file information  
🌈 **Resume Upload** - Supports resuming uploads after interruption

### Advanced Features

🎛️ **Share Control** - Supports password protection, download count limits, and expiration time settings  
🔍 **Pickup Code System** - Supports pickup code sharing, simplifying sharing difficulty  
⚡ **Queue Processing** - Asynchronous task processing system based on Redis + Asynq  
🗂️ **File Management** - Complete file lifecycle management  
📷 **Image Processing** - Image compression, format conversion, and other processing features  
🏷️ **Download Control** - Download token management system based on JWT

## 📸 Screenshots

File selection upload page
![](/.github/image/1.png)

Text input upload page
![](/.github/image/2.png)

File selection upload page - supports multiple file uploads
![](/.github/image/3.png)

File uploading page - similar to GitHub's file heatmap showing upload progress
![](/.github/image/4.png)

File uploading page - similar to qBittorrent's progress bar showing file upload progress
![](/.github/image/5.png)

File upload success page
![](/.github/image/6.png)

## 🏗️ Technical Architecture

### Frontend Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Nuxt 3** - Vue.js full-stack framework
- **TypeScript** - Complete type safety
- **Tailwind CSS** - Atomic CSS framework
- **Reka UI** - Modern component library
- **Pinia** - State management
- **TanStack Query** - Data fetching and caching
- **Vue Router** - Routing management
- **i18next** - Internationalization support

### Backend Tech Stack

- **Go 1.23** - High-performance server-side language
- **Echo** - High-performance HTTP framework
- **Redis** - Caching and session storage
- **Asynq** - Asynchronous task queue
- **JWT** - Authentication
- **Zap** - Structured logging

### Build System

- **Node.js** - Server-side runtime
- **pnpm** - Fast package manager
- **Husky** - Git hooks management
- **Prettier** - Code formatting
- **Lint-staged** - Staged file checking

### Storage Architecture

- **File Storage** - Local file system storage
- **Redis Cache** - Share information and file metadata caching
- **Queue System** - Asynchronous task processing queue

## 🚀 Quick Start

### Requirements

- Node.js 18+
- Go 1.23+
- Redis 6+
- pnpm 9+

### Install Dependencies

```bash
# Install root dependencies
pnpm install

# Install frontend dependencies
cd front && pnpm install

# Install backend dependencies
cd backend && go mod download

# Install Worker dependencies
cd worker && go mod download
```

### Environment Configuration

```bash
# Copy environment variables file
cp .env.example .env

# Configure necessary environment variables
REDIS_URL=redis://localhost:6379
UPLOAD_PATH=/.github/uploads
download_secret=your_download_secret
```

### Start Services

```bash
# Start all services in development mode
pnpm run dev

# Or start separately
pnpm run dev:front    # Frontend service (port 5000)
pnpm run dev:backend  # Backend service (port 1323)
pnpm run dev:worker   # Worker service
```

## 📁 Project Structure

```
015/
├── front/                 # Frontend application (Vue 3 + Nuxt 3)
│   ├── components/       # Vue components
│   │   ├── pages/           # Page routes
│   │   ├── composables/     # Composable functions
│   │   ├── i18n/           # Internationalization files
│   │   └── assets/         # Static assets
│   └── middleware/     # Middleware
├── backend/             # Backend service (Go + Echo)
│   ├── internal/       # Internal packages
│   │   ├── controllers/ # Controllers
│   │   ├── models/     # Data models
│   │   ├── services/   # Business logic
│   │   └── utils/      # Utility functions
│   └── middleware/     # Middleware
├── worker/             # Asynchronous task processing (Go + Asynq)
│   ├── internal/       # Internal packages
│   │   ├── tasks/      # Task processors
│   │   └── utils/      # Utility functions
│   └── middleware/     # Middleware
└── tmp/               # Temporary files
```

## 🔧 Development Guide

### Code Standards

- Use Prettier for code formatting
- Use Husky + lint-staged for pre-commit checking
- Follow TypeScript type safety standards

### Commit Standards

```bash
# Code formatting will run automatically before commit
git add .
git commit -m "feat: add new feature"
```

### Build and Deploy

```bash
# Build frontend
cd front && pnpm run build

# Build backend (requires Go environment)
cd backend && go build -o main .

# Build Worker
cd worker && go build -o worker .
```

## 📝 Development Roadmap

### Completed Features ✅

- Frontend hash calculation and instant transfer
- Concurrent chunked upload (using Web Worker)
- File upload/text upload and sharing
- Upload statistics page
- Multi-language support
- Maximum upload limits
- Backend queue system and Worker file processing

### Planned Features 🚧

- Resume upload (backend calculates uploaded parts and returns)
- Image format conversion and compression
- Image OCR copy
- Document to Markdown conversion
- Text translation/summarization
- Support for multiple file uploads

## 🤝 Contributing

Welcome to submit Issues and Pull Requests to improve this project.

## 📄 License

This project is licensed under AGPLV3.

## 🔗 Related Links

- [Vue 3 Documentation](https://vuejs.org/)
- [Nuxt 3 Documentation](https://nuxt.com/)
- [Echo Framework Documentation](https://echo.labstack.com/)
- [Asynq Documentation](https://github.com/hibiken/asynq)
