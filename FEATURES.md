# xcstrings-translator - New Features Documentation

## Overview
This document describes the new features added to xcstrings-translator, including user authentication, database persistence with PostgreSQL, and enhanced project management capabilities.

## New Features

### 1. User Authentication System
- **User Registration**: New users can create accounts with username, email, and password
- **User Login**: Secure authentication using JWT tokens
- **Session Management**: Automatic token storage and verification

### 2. PostgreSQL Database Integration
- **User Management**: Store user accounts and credentials (hashed passwords)
- **Project Storage**: Persist localizable.xcstrings files in the database
- **Translation Memory**: Store translated strings for reuse and reference
- **Provider Configuration**: Save translation provider settings per user

### 3. Project Management
- **Project Creation**: Create new localization projects
- **File Upload**: Store .xcstrings files in the database
- **Project Metadata**: Track project details like name, description, and source language
- **Translation History**: Maintain records of translation jobs and their outcomes

### 4. Translation Provider Configuration
- **Multiple Providers**: Support for Google, DeepL, Baidu, and OpenAI
- **Secure Storage**: Safely store API keys and provider settings
- **Default Provider**: Set default translation provider per user
- **Provider Management**: Create, update, and delete provider configurations

### 5. Enhanced Web Interface
- **Authentication UI**: Login and registration modals
- **Project Selection**: Modal for choosing existing projects or creating new ones
- **Project Creation**: Form to create new localization projects
- **Improved Upload**: Option to upload to new project or existing project

## Database Schema

### Users Table
- `id`: Primary key
- `created_at`, `updated_at`, `deleted_at`: Timestamps
- `username`: Unique username
- `email`: User's email address
- `password`: Hashed password
- `is_active`: Account status

### Projects Table
- `id`: Primary key
- `created_at`, `updated_at`, `deleted_at`: Timestamps
- `name`: Project name
- `description`: Project description
- `user_id`: Foreign key to Users table
- `file_content`: Stored xcstrings file content
- `file_name`: Original filename
- `source_language`: Project's source language
- `content_structure`: JSON field with parsed structure
- `settings`: JSON field for project settings

### Translations Table
- `id`: Primary key
- `created_at`, `updated_at`, `deleted_at`: Timestamps
- `project_id`: Foreign key to Projects table
- `key`: Original key from xcstrings file
- `source_text`: Source text to translate
- `target_text`: Translated text
- `target_language`: Language code (e.g., "zh-Hans", "ja")
- `state`: Translation state ("translated", "needs_review", etc.)
- `translation_provider`: Provider used for translation

### ProviderConfigs Table
- `id`: Primary key
- `created_at`, `updated_at`, `deleted_at`: Timestamps
- `user_id`: Foreign key to Users table
- `provider_type`: Provider type ("openai", "google", "deepl", "baidu")
- `config_data`: JSON field with provider-specific configuration
- `is_default`: Whether this is the default provider for the user

## API Endpoints

### Authentication Endpoints
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout

### Project Endpoints (require authentication)
- `GET /api/protected/projects` - Get user's projects
- `POST /api/protected/projects` - Create new project
- `GET /api/protected/projects/:id` - Get specific project
- `PUT /api/protected/projects/:id` - Update project
- `DELETE /api/protected/projects/:id` - Delete project
- `POST /api/protected/projects/:id/upload` - Upload file to project
- `POST /api/protected/projects/:id/translate` - Translate project
- `GET /api/protected/projects/:id/export` - Export project
- `GET /api/protected/projects/:id/translations` - Get project translations
- `GET /api/protected/projects/:id/missing-translations` - Get missing translations

### Provider Configuration Endpoints (require authentication)
- `GET /api/protected/providers` - Get user's provider configs
- `POST /api/protected/providers` - Create provider config
- `GET /api/protected/providers/:id` - Get specific provider config
- `PUT /api/protected/providers/:id` - Update provider config
- `DELETE /api/protected/providers/:id` - Delete provider config
- `GET /api/protected/providers/:type/default` - Get default provider config for type

## Setup Instructions

### Database Setup
1. Install PostgreSQL
2. Create a database for the application:
   ```sql
   CREATE DATABASE xcstrings;
   CREATE USER xcstrings WITH PASSWORD 'xcstrings';
   GRANT ALL PRIVILEGES ON DATABASE xcstrings TO xcstrings;
   ```
3. Set environment variables:
   - `DATABASE_URL` (optional, defaults to local setup)
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` (optional, can use DATABASE_URL instead)

### Environment Variables
- `JWT_SECRET` - Secret key for JWT token signing
- `DATABASE_URL` - PostgreSQL connection string
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Individual database connection parameters (alternative to DATABASE_URL)

### Running the Application
1. Install dependencies: `go mod tidy`
2. Build the frontend: `cd web && npm install && npm run build`
3. Run the application: `go run main.go`

## Security Considerations
- All sensitive fields (passwords, API keys) are properly sanitized before being returned to the client
- Authentication is required for project and provider management endpoints
- Passwords are hashed using bcrypt before storage
- JWT tokens provide secure session management

## Browser vs Server-Side Configuration

### Server-Side Provider Configuration
- Stored securely in the database per user
- API keys securely stored with encryption
- Configuration available across sessions
- Suitable for production use

### Browser-Side Configuration (Practice/Limited)
- Configuration stored in browser's localStorage
- Limited to 10 common languages and limited translation requests
- Useful for practice and testing without registration
- All data cleared when browser data is cleared

## Future Enhancements
- Translation caching mechanism
- Translation quality assessment
- Batch file processing
- Translation memory
- Interactive translation confirmation
- Rate limiting for free tier users
