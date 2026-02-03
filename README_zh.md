# OpenTrans

现代化的 iOS/macOS 多语言字符串翻译管理平台，支持多种翻译提供商（Google、DeepL、百度、OpenAI、腾讯混元本地模型）。提供 CLI 命令行工具和 Web UI 界面。平台包含完整的用户认证、订阅管理、App Store Connect 集成和团队协作功能。

## 功能特性

- **多提供商翻译**
  - Google Translate API
  - DeepL API
  - 百度翻译 API
  - OpenAI 兼容 API（支持 GPT、Claude 等）
  - 腾讯混元本地模型

- **Web UI 管理**
  - 基于 Vue 3 + Vite 的现代化界面
  - 项目管理
  - 翻译状态追踪
  - 可视化翻译编辑工作区

- **用户与权限**
  - 完整的用户认证系统（JWT）
  - 角色管理（admin/user）
  - 订阅计划

- **App Store Connect 集成**
  - 同步 App 元数据
  - 本地化信息
  - 双向同步支持（部分完成）

- **协作功能**
  - 团队成员管理
  - 用量追踪
  - 订阅额度控制

- **CLI 工具**
  - 命令行批量翻译 xcstrings 文件

## 安装

### 环境要求

- Go 1.24.4 或更高版本
- Node.js 16+ 和 npm
- PostgreSQL 或 MySQL
- （可选）Docker

### 使用 Docker 快速启动

```bash
docker build -t opentrans .
docker run -p 3000:3000 opentrans
```

### 手动安装

1. 克隆仓库：
```bash
git clone https://github.com/fdddf/opentrans.git
cd opentrans
```

2. 配置环境变量：
```bash
cp .env.sample .env
# 编辑 .env 文件
```

3. 安装依赖：
```bash
npm --prefix web install
```

4. 初始化数据库：
```bash
bash backend/scripts/db-create.sh
bash backend/scripts/db-migrate.sh
```

5. 构建并运行：
```bash
make -C backend binary
./backend/opentrans serve
```

## 使用方法

### CLI 翻译

```bash
# 翻译单个文件
./backend/opentrans translate -i input.xcstrings -o output.xcstrings

# 配置翻译提供商
./backend/opentrans config set google.api_key "your-api-key"
./backend/opentrans config set deepl.api_key "your-api-key"
```

### Web UI

启动开发服务器：

```bash
# 后端
cd backend && go run main.go serve

# 前端（在另一个终端）
npm --prefix web run dev
```

访问 Web UI：`http://localhost:5173`

## 配置

### 环境变量

`.env` 中的关键环境变量：

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=xcstrings_user
DB_PASSWORD=your_password
DB_NAME=xcstrings_translator

# JWT
JWT_SECRET=your-secret-key

# 服务器
SERVER_PORT=3000
BASE_URL=http://localhost:8080

# 翻译提供商
GOOGLE_API_KEY=your-google-api-key
DEEPL_API_KEY=your-deepl-api-key
BAIDU_APP_ID=your-baidu-app-id
BAIDU_APP_SECRET=your-baidu-app-secret
OPENAI_API_KEY=your-openai-api-key
```

### CLI 配置文件

编辑 `backend/config.yaml` 进行 CLI 专用设置：

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

## 开发

### 后端

```bash
# 构建
make -C backend binary

# 运行测试
make -C backend test

# 安装到 GOPATH
make -C backend install
```

### 前端

```bash
# 开发模式
npm --prefix web run dev

# 生产构建
npm --prefix web run build

# 类型检查
npm --prefix web run lint
```

### 数据库迁移

```bash
# 创建数据库
bash backend/scripts/db-create.sh

# 运行迁移
bash backend/scripts/db-migrate.sh

# 重置数据库
bash backend/scripts/db-reset.sh
```

## API 文档

平台提供以下 RESTful API 端点：

- **认证**: `/api/auth/*` - 注册、登录、用户信息
- **应用**: `/api/apps/*` - 应用的 CRUD 操作
- **本地化**: `/api/apps/:id/localizations` - 管理应用本地化
- **Apple Connect**: `/api/apple-connect/*` - App Store Connect 同步
- **翻译**: `/api/translate/*` - 文件翻译和批量操作
- **订阅**: `/api/subscriptions/*` - Stripe 集成
- **用户**: `/api/users/*` - 用户管理（管理员）

## 架构

- **后端**: Go 1.24.4 + Fiber + GORM + Uber FX
- **前端**: Vue 3 + TypeScript + Vite + Tailwind CSS
- **数据库**: PostgreSQL（支持 MySQL）
- **认证**: JWT
- **支付**: Stripe

## 贡献

欢迎贡献！请按照以下步骤：

1. Fork 仓库
2. 创建功能分支（`git checkout -b feature/amazing-feature`）
3. 提交更改（`git commit -m 'feat: add amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE.md](LICENSE.md) 文件。

## 支持

- GitHub Issues: https://github.com/fdddf/opentrans/issues
- 文档: 查看 [AGENTS.md](AGENTS.md) 了解详细项目文档

## 开发路线图

- [ ] 完成 App Store Connect 双向同步
- [ ] 添加翻译质量评分
- [ ] 实现实时协作
- [ ] 添加更多翻译提供商
- [ ] 增强同步操作的冲突解决机制