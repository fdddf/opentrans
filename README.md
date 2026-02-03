# OpenTrans

A modern iOS/macOS multilingual string translation management platform supporting multiple translation providers (Google, DeepL, Baidu, OpenAI, Tencent Hunyuan local model). Provides CLI command-line tools and Web UI interface. Includes complete user authentication, subscription management, App Store Connect integration, and team collaboration features.

## Features

- **Multiple Translation Providers**
  - Google Translate API
  - DeepL API
  - Baidu Translation API
  - OpenAI Compatible API (supports GPT, Claude, etc.)
  - Tencent Hunyuan Local Model

- **Web UI Management**
  - Modern interface based on Vue 3 + Vite
  - Project management
  - Translation status tracking
  - Visual workspace for translation editing

- **User & Permissions**
  - Complete user authentication system (JWT)
  - Role management (admin/user)
  - Subscription plans

- **App Store Connect Integration**
  - Sync App metadata
  - Localization information
  - Bidirectional sync support (partially complete)

- **Collaboration Features**
  - Team member management
  - Usage tracking
  - Subscription quota control

- **CLI Tools**
  - Command-line batch translation of xcstrings files

## Installation

### Prerequisites

- Go 1.24.4 or higher
- Node.js 16+ and npm
- PostgreSQL or MySQL
- (Optional) Docker

### Quick Start with Docker

```bash
docker build -t opentrans .
docker run -p 3000:3000 opentrans
```

### Manual Installation

1. Clone the repository:
```bash
git clone https://github.com/fdddf/opentrans.git
cd opentrans
```

2. Configure environment variables:
```bash
cp .env.sample .env
# Edit .env with your configuration
```

3. Install dependencies:
```bash
npm --prefix web install
```

4. Initialize database:
```bash
bash backend/scripts/db-create.sh
bash backend/scripts/db-migrate.sh
```

5. Build and run:
```bash
make -C backend binary
./backend/opentrans serve
```

## Usage

### CLI Translation

```bash
# Translate a single file
./backend/opentrans translate -i input.xcstrings -o output.xcstrings

# Configure translation providers
./backend/opentrans config set google.api_key "your-api-key"
./backend/opentrans config set deepl.api_key "your-api-key"
```

### Web UI

Start the development server:

```bash
# Backend
cd backend && go run main.go serve

# Frontend (in another terminal)
npm --prefix web run dev
```

Access the Web UI at `http://localhost:5173`

## Configuration

### Environment Variables

Key environment variables in `.env`:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=xcstrings_user
DB_PASSWORD=your_password
DB_NAME=xcstrings_translator

# JWT
JWT_SECRET=your-secret-key

# Server
SERVER_PORT=3000
BASE_URL=http://localhost:8080

# Translation Providers
GOOGLE_API_KEY=your-google-api-key
DEEPL_API_KEY=your-deepl-api-key
BAIDU_APP_ID=your-baidu-app-id
BAIDU_APP_SECRET=your-baidu-app-secret
OPENAI_API_KEY=your-openai-api-key
```

### CLI Configuration File

Edit `backend/config.yaml` for CLI-specific settings:

```yaml
global:
  source_language: "en"
  target_languages:
    - "zh-Hans"
    - "ja"
    - "ko"
  concurrency: 5

google:
  api_key: "your-google-api-key-here"
```

## Development

### Backend

```bash
# Build
make -C backend binary

# Run tests
make -C backend test

# Install to GOPATH
make -C backend install
```

### Frontend

```bash
# Development
npm --prefix web run dev

# Build for production
npm --prefix web run build

# Type check
npm --prefix web run lint
```

### Database Migrations

```bash
# Create database
bash backend/scripts/db-create.sh

# Run migrations
bash backend/scripts/db-migrate.sh

# Reset database
bash backend/scripts/db-reset.sh
```

## API Documentation

The platform provides a RESTful API with the following endpoints:

- **Authentication**: `/api/auth/*` - Register, login, user info
- **Apps**: `/api/apps/*` - CRUD operations for apps
- **Localizations**: `/api/apps/:id/localizations` - Manage app localizations
- **Apple Connect**: `/api/apple-connect/*` - App Store Connect sync
- **Translation**: `/api/translate/*` - File translation and batch operations
- **Subscriptions**: `/api/subscriptions/*` - Stripe integration
- **Users**: `/api/users/*` - User management (admin)

## Architecture

- **Backend**: Go 1.24.4 + Fiber + GORM + Uber FX
- **Frontend**: Vue 3 + TypeScript + Vite + Tailwind CSS
- **Database**: PostgreSQL (MySQL supported)
- **Authentication**: JWT
- **Payment**: Stripe

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

- GitHub Issues: https://github.com/fdddf/opentrans/issues
- Documentation: See [AGENTS.md](AGENTS.md) for detailed project documentation

## Roadmap

- [ ] Complete App Store Connect bidirectional sync
- [ ] Add translation quality scoring
- [ ] Implement real-time collaboration
- [ ] Add more translation providers
- [ ] Enhanced conflict resolution for sync operations