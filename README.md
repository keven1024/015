<div align="center"><a name="readme-top"></a>

# 015

015 (/ˈzɪərəʊ wʌn faɪv/, "zero-one-five") is a self-hosted temporary file sharing platform. Focused on providing one-time, temporary file and text upload, processing, and sharing services. The project name originates from [Ichigo](https://darling-in-the-franxx.fandom.com/wiki/Ichigo) from DARLING in the FRANXX.

A modern file sharing website built with Vue 3 + Nuxt 4 + Go, supporting file upload, text sharing, image compression, concurrent processing, instant transfer functionality, and more, featuring a complete sharing management and access control system.

![015 Platform Overview](/.github/image/0.png)

English | [中文](README-zh.md)

</div>

## 🌟 Features

### Core Functionality

- 🖼️ **High-Performance File Upload** - Supports large file chunked uploads with frontend file hash calculation for instant transfer
- 📱 **Responsive Design** - Modern UI based on Tailwind V4 + Reka UI, adapts to various devices
- ⚡ **Concurrent Processing** - Uses a Web Worker for frontend hash calculation and a backend queue for task processing
- 🌐 **Multi-language Support** - Supports Simplified Chinese, Traditional Chinese, English, Japanese, Korean, French, and German
- 🔗 **Share Management** - Flexible share link generation and management system

### File Processing

- 🔄 **Smart Instant Transfer** - Detects existing uploads from the file hash and size to avoid uploading duplicate data
- 📷 **Image Compression** - Supports asynchronous compression of uploaded images
- 🖼️ **File Preview** - Previews image and video content and displays type and basic information for other files
- 📊 **Upload Statistics** - Displays upload progress and file information in real time
- 🌈 **Resumable Uploads** - Resumes uploads from chunks already stored by the server

### Advanced Features

- 🎛️ **Share Control** - Supports password protection, download count limits, and expiration settings
- 🔍 **Pickup Code System** - Provides short pickup codes for easier sharing
- ⚡ **Queue Processing** - Asynchronous task processing based on Redis and Asynq
- 🗂️ **File Management** - Complete file lifecycle management
- 📷 **Image Processing** - Image compression and format conversion
- 🏷️ **Download Control** - JWT-based download token management

## 📸 Screenshots

| File Selection Upload Page | Text Input Upload Page    |
| -------------------------- | ------------------------- |
| ![](/.github/image/1.png)  | ![](/.github/image/2.png) |

| Multiple File Upload       | Upload Progress Heatmap   |
| -------------------------- | ------------------------- |
| ![](/.github/image/3.png)  | ![](/.github/image/4.png) |

| Upload Progress Bar        | Upload Success Page       |
| -------------------------- | ------------------------- |
| ![](/.github/image/5.png)  | ![](/.github/image/6.png) |

## 🚀 Quick Start

### Docker

1. Download files
   - config.example.yaml
   - docker-compose.yml

2. Rename config.example.yaml to config.yaml after configuration

3. Start
```bash
docker compose up -d
```

4. Visit `http://localhost:8080`

## 🏗️ Technical Architecture

### Frontend Tech Stack

- **Vue 3** - Progressive JavaScript framework
- **Nuxt 4** - Vue.js full-stack framework
- **TypeScript** - Complete type safety
- **Tailwind CSS** - Atomic CSS framework
- **Reka UI** - Modern component library
- **Pinia** - State management
- **TanStack Query** - Data fetching and caching
- **Nuxt File-based Routing** - Page and route management
- **Nuxt I18n / Vue I18n** - Internationalization support

### Backend Tech Stack

- **Go 1.25.5** - High-performance server-side language
- **Echo** - High-performance HTTP framework
- **Redis** - Caching and session storage
- **Asynq** - Asynchronous task queue
- **JWT** - Download access tokens
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

## 📁 Project Structure

```
015/
├── front/                    # Frontend application (Vue 3 + Nuxt 4)
│   ├── components/           # Vue and UI components
│   ├── pages/                # File-based page routes
│   ├── composables/          # Composable functions
│   ├── i18n/                 # Frontend localization resources
│   ├── assets/               # Styles and static assets
│   ├── plugins/              # Nuxt plugins
│   └── server/               # Nuxt server configuration endpoint
├── backend/                  # API service (Go + Echo)
│   ├── internal/
│   │   ├── controllers/      # HTTP controllers
│   │   ├── services/         # Business logic
│   │   └── utils/            # Backend utilities
│   └── middleware/           # HTTP middleware
├── worker/                   # Asynchronous task processing (Go + Asynq)
│   ├── internal/
│   │   ├── services/         # File, image, notification, and text services
│   │   ├── tasks/            # Queue task handlers
│   │   └── utils/            # Worker utilities
│   └── middleware/           # Worker middleware
├── pkg/                      # Shared Go Workspace modules
│   ├── geoip/                # IP geolocation
│   ├── i18n/                 # Backend localization
│   ├── mail/                 # Email templates and delivery
│   ├── models/               # Redis data models
│   ├── services/             # Shared services
│   └── utils/                # Shared utilities
├── go.work                   # Go Workspace configuration
└── pnpm-workspace.yaml       # pnpm Workspace configuration
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
- Multiple file upload support
- Upload statistics page
- Multi-language support
- Maximum upload limits
- Backend queue system and Worker file processing
- Resume upload (backend calculates uploaded parts and returns)
- Image format conversion and compression

### Planned Features 🚧

- Image OCR copy
- Document to Markdown conversion
- Text translation/summarization

## 🤝 Contributing

Welcome to submit Issues and Pull Requests to improve this project.

## 📄 License

This project is licensed under AGPLV3.

## 🔗 Related Links

- [Vue 3 Documentation](https://vuejs.org/)
- [Nuxt 4 Documentation](https://nuxt.com/)
- [Echo Framework Documentation](https://echo.labstack.com/)
- [Asynq Documentation](https://github.com/hibiken/asynq)
