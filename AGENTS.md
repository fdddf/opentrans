# Repository Guidelines

## Project Overview

**OpenTrans** 是一个现代化的 iOS/macOS 多语言字符串翻译管理平台，支持多种翻译提供商（Google、DeepL、百度、OpenAI、腾讯混元本地模型），提供 CLI 命令行工具和 Web UI 界面。平台包含完整的用户认证、订阅管理、App Store Connect 集成和团队协作功能。

### 核心功能
- **多提供商翻译**: 支持 Google Translate、DeepL、百度翻译、OpenAI 兼容 API、腾讯混元本地模型
- **Web UI 管理**: 基于 Vue 3 + Vite 的现代化界面，支持项目管理、翻译状态追踪
- **用户与权限**: 完整的用户认证系统（JWT）、角色管理（admin/user）、订阅计划
- **App Store Connect 集成**: 同步 App 元数据、本地化信息、双向同步支持（部分完成）
- **协作功能**: 团队成员管理、用量追踪、订阅额度控制
- **CLI 工具**: 支持命令行批量翻译 xcstrings 文件

---

## Project Structure & Module Organization

### Backend (`backend/`)
基于 Go 1.24.4 构建，使用 Fiber 框架作为 Web 服务器，GORM 作为 ORM，Uber FX 进行依赖注入。

```
backend/
├── cmd/                    # CLI 命令定义
│   ├── root.go            # 根命令
│   ├── translate.go       # 翻译命令
│   ├── serve.go           # 启动 Web 服务器
│   ├── generate.go        # 生成相关命令
│   ├── config.go          # 配置命令
│   ├── baidu.go           # 百度翻译配置
│   ├── deepl.go           # DeepL 配置
│   ├── google.go          # Google 翻译配置
│   ├── openai.go          # OpenAI 配置
│   ├── llama.go           # Llama 本地模型配置
│   └── migrate.go         # 数据库迁移命令
├── internal/
│   ├── auth/              # 认证模块（JWT）
│   ├── config/            # 配置管理（环境变量、配置文件）
│   ├── context/           # 请求上下文
│   ├── controllers/       # HTTP 控制器
│   │   ├── auth_controller.go
│   │   ├── app_controller.go
│   │   ├── app_localization_controller.go
│   │   ├── apple_connect_controller.go
│   │   ├── user_controller.go
│   │   ├── subscription_controller.go
│   │   ├── translation_controller.go
│   │   └── ...
│   ├── dao/               # GORM Gen 生成的 DAO 层
│   ├── database/          # 数据库模型、迁移
│   ├── email/             # 邮件发送服务
│   ├── model/             # 共享数据模型
│   ├── server/            # Fiber 服务器配置、中间件
│   ├── services/          # 业务逻辑层
│   │   ├── app_service.go
│   │   ├── app_localization_service.go
│   │   ├── appleconnect_service.go
│   │   ├── translation_service.go
│   │   ├── subscription_service.go
│   │   └── ...
│   └── translator/        # 翻译提供商实现
│       ├── translator.go  # 翻译服务核心
│       ├── google.go
│       ├── deepl.go
│       ├── baidu.go
│       ├── openai.go
│       └── llama.go
├── pkg/
│   ├── appleconnect/      # App Store Connect API 客户端
│   └── hunyuan/           # 腾讯混元语言支持
├── migrations/            # 数据库迁移文件（SQL）
├── webui/dist/            # 前端构建产物（嵌入到后端）
├── config.yaml            # CLI 配置文件
├── main.go                # 应用入口（FX 依赖注入）
└── Makefile               # 构建脚本
```

### Frontend (`web/`)
基于 Vue 3 + TypeScript + Vite 构建，使用 Tailwind CSS 进行样式设计。

```
web/
├── src/
│   ├── components/        # Vue 组件
│   │   ├── AppShell.vue   # 应用主布局
│   │   └── TranslationTest.vue
│   ├── views/             # 页面视图
│   │   ├── Login.vue
│   │   ├── Register.vue
│   │   ├── Dashboard.vue
│   │   ├── Apps.vue
│   │   ├── AppWorkspace.vue
│   │   ├── AppLocalizations.vue
│   │   ├── AppleConnectConfig.vue
│   │   ├── ProviderConfigs.vue
│   │   ├── Users.vue
│   │   ├── Subscriptions.vue
│   │   ├── Languages.vue
│   │   └── Profile.vue
│   ├── composables/       # Vue 组合式函数
│   │   └── useApi.ts      # API 调用封装
│   ├── stores/            # Pinia 状态管理
│   │   └── theme.ts
│   ├── router/            # Vue Router 配置
│   │   └── index.ts
│   ├── locales/           # 国际化文件
│   │   ├── en.json
│   │   └── zh.json
│   ├── App.vue            # 根组件
│   ├── main.ts            # 应用入口
│   └── i18n.ts            # i18n 配置
├── index.html
├── vite.config.ts         # Vite 配置（代理到后端 API）
├── tailwind.config.js     # Tailwind CSS 配置
└── package.json
```

### 其他目录
- `tests/` - CLI 测试用例和测试脚本
- `docs/` - 项目文档、计划文档
- `.github/workflows/` - CI/CD 配置
- `webui/dist/` - 独立的前端构建输出目录

---

## Build, Test, and Development Commands

### Backend
```bash
# 从项目根目录运行
make -C backend binary           # 构建 Go 二进制文件到 backend/opentrans
make -C backend ui               # 安装前端依赖并构建到 backend/webui/dist
make -C backend test             # 运行 Go 测试
make -C backend clean            # 清理构建产物
make -C backend install          # 安装到 GOPATH/bin

# 数据库迁移
bash backend/scripts/db-create.sh    # 创建数据库
bash backend/scripts/db-migrate.sh   # 运行迁移
bash backend/scripts/db-reset.sh     # 重置数据库

# CLI 使用
./backend/opentrans translate -i input.xcstrings -o output.xcstrings
./backend/opentrans serve  # 启动 Web 服务器（默认 :3000）
```

### Frontend
```bash
# 开发模式
npm --prefix web run dev          # 启动 Vite 开发服务器（默认 :5173，代理到 :3000）

# 构建
npm --prefix web run build        # 构建生产版本到 webui/dist
npm --prefix web run preview      # 预览构建结果

# 类型检查
npm --prefix web run lint         # 运行 vue-tsc 类型检查
```

### 测试
```bash
# CLI 烟雾测试
bash tests/test.sh                # 运行 CLI 测试脚本（需要 jq）

# Go 测试
cd backend && go test ./...       # 运行所有测试
cd backend && go test -v ./internal/services/  # 运行特定包测试
```

### Docker
```bash
docker build -t opentrans .  # 构建镜像
docker run -p 3000:3000 opentrans  # 运行容器
```

---

## Coding Style & Naming Conventions

### Go
- **格式化**: 使用 `gofmt` 格式化代码
- **命名约定**:
  - 本地变量: `camelCase`（如 `userID`, `appService`）
  - 导出符号: `PascalCase`（如 `AppService`, `CreateApp`）
  - 接口名: 通常以行为命名（如 `TranslationProvider`）
- **文件组织**: 每个文件主要包含一个主要类型或功能
- **错误处理**: 使用 `fmt.Errorf` 包装错误，返回给调用者
- **模块边界**:
  - `internal/` - 应用内部包，不被外部导入
  - `pkg/` - 可复用的公共工具包
  - `cmd/` - CLI 命令实现
- **依赖注入**: 使用 Uber FX 进行依赖注入（`fx.Provide`, `fx.Invoke`）

### Vue/TypeScript
- **格式化**: 2 空格缩进
- **命名约定**:
  - 组件: `PascalCase`（如 `AppShell.vue`, `Dashboard.vue`）
  - 组合式函数: `camelCase`（如 `useApi.ts`, `useAuth`）
  - 变量/函数: `camelCase`
- **类型定义**: 在 `composables/useApi.ts` 中集中定义 TypeScript 接口
- **路由**: 使用 Vue Router，路由参数通过 `props: true` 传递
- **状态管理**: 使用 Pinia（`stores/`）
- **国际化**: 使用 Vue I18n（`locales/`）

### 数据库
- **表名**: 使用复数形式（如 `apps`, `app_localizations`）
- **列名**: 使用 `camelCase`（如 `user_id`, `created_at`）
- **迁移**: 使用 SQL 迁移文件，版本号前缀（如 `000001_init_schema.up.sql`）
- **ORM**: 使用 GORM，通过 GORM Gen 生成类型安全的 DAO 层

---

## Technology Stack

### Backend
- **语言**: Go 1.24.4
- **Web 框架**: Fiber v2
- **ORM**: GORM + GORM Gen
- **数据库**: PostgreSQL（支持 MySQL）
- **依赖注入**: Uber FX
- **认证**: JWT (github.com/golang-jwt/jwt/v5)
- **配置**: Viper + 环境变量
- **CLI**: Cobra
- **迁移**: golang-migrate/migrate

### Frontend
- **框架**: Vue 3.5.26
- **语言**: TypeScript 5.4.5
- **构建工具**: Vite 5.2.8
- **样式**: Tailwind CSS 3.4.4
- **路由**: Vue Router 4.6.4
- **国际化**: Vue I18n 9.14.4
- **类型检查**: vue-tsc 2.0.21

### 翻译提供商
- Google Translate API
- DeepL API
- 百度翻译 API
- OpenAI 兼容 API（支持 GPT、Claude 等）
- 腾讯混元本地模型（通过 libyzma）

---

## Testing Guidelines

### Backend Tests
- 使用 Go 内置 `testing` 包
- 测试文件与源文件同名，添加 `_test.go` 后缀
- 测试函数以 `Test` 开头
- 运行测试: `go test -v ./...`

### Frontend Tests
- 类型检查: `vue-tsc --noEmit`
- 暂未配置单元测试框架（可考虑 Vitest）

### CLI 烟雾测试
- 测试脚本位于 `tests/test.sh`
- 示例 xcstrings 文件位于 `tests/example.xcstrings`
- 需要安装 `jq` 工具

### 集成测试
- 数据库迁移测试: `backend/test/test-migration.go`
- 建议添加端到端测试（Playwright/Cypress）

---

## Configuration & Secrets

### 环境变量
复制 `.env.sample` 到 `.env` 并配置以下变量：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=xcstrings_user
DB_PASSWORD=change_this_password
DB_NAME=xcstrings_translator
DB_SSL_MODE=disable

# JWT 配置
JWT_SECRET=your-secret-key-change-this-in-production

# 应用配置
BASE_URL=http://localhost:8080
SERVER_PORT=3000
SERVER_HOST=0.0.0.0

# 邮件配置（用户激活）
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=noreply@yourdomain.com

# Stripe 配置（订阅）
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Apple Connect 配置（可选）
APPLE_ISSUER_ID=your-issuer-id
APPLE_KEY_ID=your-key-id
APPLE_PRIVATE_KEY_PATH=/path/to/private/key.p8

# 管理员初始化
ADMIN_USERNAME=admin
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
ADMIN_CREATE_IF_NOT_EXISTS=true
```

### CLI 配置
`backend/config.yaml` 用于 CLI 命令的翻译提供商配置：

```yaml
global:
  input_file: "Localizable.xcstrings"
  output_file: "Localizable_translated.xcstrings"
  source_language: "en"
  target_languages:
    - "zh-Hans"
    - "ja"
    - "ko"
  concurrency: 5

google:
  api_key: "your-google-api-key-here"

deepl:
  api_key: "your-deepl-api-key-here"
  is_free: false

baidu:
  app_id: "your-baidu-app-id-here"
  app_secret: "your-baidu-app-secret-here"

openai:
  api_key: "your-openai-api-key-here"
  api_base_url: "https://api.openai.com"
  model: "gpt-3.5-turbo"

llama:
  model_path: "/path/to/model.gguf"
  threads: 4
```

### 安全原则
- **永不提交** `.env` 文件或包含密钥的配置
- **使用环境变量** 或密钥管理服务
- **生产环境** 必须更改默认的 JWT 密钥和管理员密码

---

## Database Schema

### 核心表结构

#### users
用户账户和订阅信息
- `id`, `username`, `email`, `password`
- `role` (admin/user)
- `is_active`, `is_activated`
- `subscription_type`, `subscription_end`
- `max_apps`, `max_translations`
- `current_usage`, `current_app_count`

#### apps
应用元数据
- `id`, `name`, `description`, `bundle_id`, `apple_id`
- `primary_locale`, `subtitle`
- `short_description`, `long_description`
- `keywords`, `support_url`, `marketing_url`, `privacy_url`
- `user_id`, `origin` (manual/synced)
- `is_ready_for_review`

#### app_localizations
应用本地化数据
- `id`, `app_id`, `language_code`
- `name`, `subtitle`, `description`
- `promotional_text`, `whats_new`
- `synced_at`, `source`, `sync_status`
- `version`, `version_state`

#### projects
xcstrings 项目
- `id`, `name`, `description`, `user_id`
- `file_content`, `file_name`
- `source_language`, `content_structure`

#### translations
翻译结果
- `id`, `project_id`, `key`
- `source_text`, `target_text`, `target_language`
- `state`, `translation_provider`

#### subscriptions
订阅记录
- `id`, `user_id`, `stripe_subscription_id`
- `status`, `plan_id`, `amount`

#### app_provider_configs
应用级别的翻译提供商配置
- `id`, `app_id`, `provider_type`
- `api_key`, `api_base_url`, `model`

#### translation_queues
异步翻译队列
- `id`, `project_id`, `app_id`
- `status`, `total_items`, `completed_items`

#### user_activities
用户活动日志
- `id`, `user_id`, `action`, `details`

### 迁移管理
- 所有迁移文件位于 `backend/migrations/`
- 使用 `golang-migrate` 工具管理
- 当前版本: `000007_add_subtitle_to_apps`

---

## API Architecture

### RESTful API 端点
后端使用 Fiber 提供 REST API，所有端点前缀为 `/api`。

#### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/me` - 获取当前用户信息

#### 应用管理
- `GET /api/apps` - 获取应用列表
- `POST /api/apps` - 创建应用
- `GET /api/apps/:id` - 获取应用详情
- `PUT /api/apps/:id` - 更新应用
- `DELETE /api/apps/:id` - 删除应用

#### 应用本地化
- `GET /api/apps/:id/localizations` - 获取本地化列表
- `POST /api/apps/:id/localizations` - 创建本地化
- `PUT /api/apps/:id/localizations/:langCode` - 更新本地化
- `DELETE /api/apps/:id/localizations/:langCode` - 删除本地化

#### Apple Connect
- `POST /api/apple-connect/sync-apps` - 同步应用列表
- `POST /api/apple-connect/sync-localizations` - 同步本地化
- `GET /api/apple-connect/config` - 获取配置
- `POST /api/apple-connect/config` - 保存配置

#### 翻译
- `POST /api/translate/file` - 上传文件翻译
- `POST /api/translate/batch` - 批量翻译
- `GET /api/translate/queue/:id` - 查询队列状态

#### 订阅
- `GET /api/subscriptions` - 获取订阅信息
- `POST /api/subscriptions/create-checkout` - 创建 Stripe 支付会话
- `POST /api/subscriptions/webhook` - Stripe webhook 回调

#### 用户管理（管理员）
- `GET /api/users` - 获取用户列表
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 中间件
- **认证中间件**: 验证 JWT token，注入用户信息到上下文
- **角色中间件**: 检查用户角色（admin/user）
- **租户隔离**: 所有查询自动按 user_id 过滤
- **订阅检查**: 验证订阅额度（app 数量、翻译次数）

---

## Development Workflow

### 1. 环境设置
```bash
# 克隆仓库
git clone https://github.com/fdddf/opentrans.git
cd opentrans

# 配置环境变量
cp .env.sample .env
# 编辑 .env 文件

# 安装前端依赖
npm --prefix web install

# 运行数据库迁移
bash backend/scripts/db-create.sh
bash backend/scripts/db-migrate.sh
```

### 2. 开发模式
```bash
# 终端 1: 启动后端
cd backend && go run main.go serve

# 终端 2: 启动前端（另一个终端）
npm --prefix web run dev
```

访问 `http://localhost:5173` 查看前端。

### 3. 代码提交
遵循 Conventional Commits 规范：
- `feat: 新功能`
- `fix: 修复 bug`
- `refactor: 重构`
- `docs: 文档更新`
- `test: 测试相关`
- `chore: 构建/工具相关`

示例：
```bash
git add .
git commit -m "feat: add Apple Store Connect synchronization"
git push origin next
```

### 4. Pull Request
- 从 `next` 分支创建功能分支
- 提交 PR 到 `next` 分支
- PR 描述应包含：
  - 功能简述
  - 测试说明（运行的命令）
  - UI 变化截图（如适用）
  - 相关 issue 链接

---

## Current Development Status

### 已完成功能 ✅
- 用户认证与授权（JWT）
- 应用 CRUD 操作
- 项目与翻译管理
- 多翻译提供商支持
- 订阅系统（Stripe 集成）
- 用户活动日志
- 数据库迁移系统
- Web UI 基础界面
- Apple Connect 基础客户端

### 进行中/待完成 🚧
详见 `docs/plan.md` 和 `docs/appstore-connect-next-steps.md`：

1. **App Store Connect 集成**
   - [x] 基础客户端（JWT 认证、获取应用、获取本地化）
   - [ ] 获取 App 最新版本
   - [ ] 更新本地化（PATCH）
   - [ ] 双向同步（pull/push）
   - [ ] 同步冲突处理策略

2. **前端完善**
   - [ ] 移除 AppLocalizations 硬编码
   - [ ] 完善编辑界面
   - [ ] 同步状态可视化
   - [ ] 配置选择器

3. **权限与租户隔离**
   - [ ] 完善所有接口的 user_id 过滤
   - [ ] App 成员角色管理
   - [ ] 请求速率限制

4. **翻译流程优化**
   - [ ] 翻译状态统计
   - [ ] 逐语言翻译进度显示
   - [ ] 翻译质量评分

### 已知问题 🐛
- `pkg/appleconnect/client.go` 中 `getAppLatestVersionID` 为占位符
- 前端 `AppLocalizations.vue` appId 硬编码
- 同步仅支持从 Apple 拉取，不支持推送
- 缺少同步冲突处理机制

---

## Troubleshooting

### 常见问题

**1. 数据库连接失败**
- 检查 `.env` 中的数据库配置
- 确保 PostgreSQL 服务正在运行
- 验证数据库用户权限

**2. 前端 API 调用失败**
- 检查 Vite 代理配置（`vite.config.ts`）
- 确保后端服务器正在运行（`:3000`）
- 查看浏览器控制台 CORS 错误

**3. 翻译失败**
- 验证翻译提供商 API 密钥
- 检查网络连接
- 查看后端日志获取详细错误

**4. Apple Connect 同步失败**
- 验证 JWT 配置（Issuer ID, Key ID, 私钥）
- 检查 Apple 开发者账号权限
- 确认 App Store Connect API 启用

### 调试
- 后端日志: 查看终端输出或配置日志文件
- 前端调试: 使用浏览器开发者工具
- 数据库查询: 使用 `backend/internal/dao/query` 生成的 DAO

---

## Contributing

欢迎贡献！请遵循以下流程：

1. Fork 仓库
2. 创建功能分支（`git checkout -b feature/amazing-feature`）
3. 提交更改（`git commit -m 'feat: add amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 创建 Pull Request

### 代码审查要点
- 遵循项目编码规范
- 添加必要的测试
- 更新相关文档
- 确保没有安全漏洞

---

## License

本项目采用 MIT 许可证 - 详见 LICENSE.md 文件。

---

## Contact & Support

- GitHub Issues: https://github.com/fdddf/opentrans/issues
- 文档: 查看项目 README.md 和 docs/ 目录

---

## Additional Resources

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Apple App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi)