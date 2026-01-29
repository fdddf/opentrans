# App Store Connect 对接与 App Meta 多语言支持：下一步改进计划

> 目标：当前项目未完成 App Store Connect 对接与 App Meta 多语言管理，本计划基于现有前端 `web/` 与后端 `backend/` 现状整理可落地的下一步改进方向。

## 一、现状盘点（基于代码扫描）

### 前端（`web/`）
- 已有页面与文案占位：
  - Apple Connect 配置页：`web/src/views/AppleConnectConfig.vue`
  - App 本地化管理页：`web/src/views/AppLocalizations.vue`
  - App 工作台页中存在 Apple Connect 同步入口：`web/src/views/AppWorkspace.vue`
- 现有明显占位/假数据：
  - `AppLocalizations.vue` 中 `appId = 1` 的硬编码、同步调用未带配置参数。
  - 语言列表是静态数组，未与后端语言/Locale 体系对齐。
- API 调用已预留：`web/src/composables/useApi.ts` 包含 `getAppLocalizations`、`syncAppleAppLocalizations` 等方法。

### 后端（`backend/`）
- 已有 Apple Connect 客户端雏形：`backend/pkg/appleconnect/client.go`
  - 已实现：JWT 生成、获取 apps、本地化列表、创建本地化。
  - 明显缺口：`getAppLatestVersionID` 为 placeholder。
- 已有 App Localizations 基本 CRUD：
  - DB Model：`backend/internal/database/models.go` 中 `AppLocalization` 表已定义多语言字段。
  - REST 路由：`backend/internal/server/server.go` 中 `/apps/:id/localizations` 等路由已存在。
- 已有 Apple Connect 同步入口：
  - `handleSyncAppleApps` / `handleSyncAppleAppLocalizations` 已存在，但仅拉取同步到本地。
  - 缺乏“反向同步”（将本地 meta 写回 Apple Connect）。

## 二、关键缺口（阻塞 App Store Connect 对接与 Meta 多语言能力）

1. **App Store Connect API 不完整**
   - 缺少 `appStoreVersions` 查询，导致无法正确创建/更新本地化。
   - 缺少更新本地化（PATCH）和删除本地化等写入能力。

2. **配置与权限关联不足**
   - ProviderConfig 只支持配置保存，但与 App 之间没有绑定关系（当前按用户取配置）。
   - 前端同步流程未指定配置 ID，后端也未校验配置是否为目标 app 授权。 

3. **本地化数据管理不完善**
   - 前端 appId 仍是硬编码，路由和 app context 未接入。
   - 前端编辑能力缺失（仅显示 + 新增 + 删除）。
   - Locale 校验与支持语言清单未统一。

4. **同步策略不清晰**
   - 当前只支持“从 Apple Connect 拉取到本地”，没有“本地变更推送到 Apple”。
   - 缺少同步状态（待同步/已同步/冲突）与审计信息。

## 三、下一步改进计划（分阶段落地）

### 阶段 1：前后端打通 App 上下文与配置选择（1-2 周）
- 前端
  - 在 `AppLocalizations.vue` 中从路由获取 `appId`，移除硬编码。
  - 同步/新增/删除接口统一从路由参数获取 appId。
  - 配置选择（Apple Connect config）改为显式选择并传给后端。
- 后端
  - `sync-localizations` 接口支持 `configId` 入参并校验该配置归属用户。
  - `App` 与 `ProviderConfig` 增加绑定关系（可新建 `AppProviderConfig` 关联表）。

### 阶段 2：完善 App Store Connect API 能力（2-3 周）
- 新增 API 能力
  - 获取 App 最新版本（补齐 `getAppLatestVersionID`）。
  - 更新本地化（PATCH /v1/appStoreVersionLocalizations/{id}）。
  - 删除本地化（DELETE）。
- 服务层抽象
  - 在 `backend/internal/services/` 中增加 `appleconnect_service.go`，封装：
    - 同步 App 列表
    - 同步本地化
    - 推送本地化更新

### 阶段 3：双向同步与冲突策略（2-4 周）
- 增加“同步方向”概念：
  - `pull`（Apple -> 本地）
  - `push`（本地 -> Apple）
- AppLocalization 增加字段：
  - `SyncedAt`、`Source`、`SyncStatus`（pending/synced/failed）。
- 同步冲突策略：
  - 以 Apple 为准 / 以本地为准 / 手动确认。

### 阶段 4：前端编辑体验与状态可视化（2 周）
- AppLocalizations UI：
  - 编辑弹窗/抽屉
  - 显示同步状态、更新时间
  - 显示字段校验（长度、必填）
- 同步按钮支持“推送到 Apple”与“从 Apple 拉取”。

### 阶段 5：Meta 字段扩展与校验（2 周）
- AppLocalization 结构补齐 App Store Connect 支持字段（如 `promotionalText`、`description`、`releaseNotes` 等最新字段）。
- 校验规则（长度限制/必填）放在后端，前端同步展示。

## 四、建议的交付顺序

1. **先打通 appId + configId**（解决当前硬编码与配置不可控问题）
2. **补齐 Apple Connect API**（能正确创建/更新本地化）
3. **实现双向同步**（本地->Apple + Apple->本地）
4. **完善 UI 编辑与状态展示**
5. **增强字段与校验规则**

---

如需我进一步：
- 输出详细任务拆解表（含工时与负责人维度）
- 直接开始实现阶段 1 或 2
- 生成 API 设计文档/数据库变更方案

可以直接告诉我下一步优先级。
