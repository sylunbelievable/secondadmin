# Second Admin

一个基于 Go + React 的前后端分离管理后台脚手架。功能方向参考
[Gin-Vue-Admin](https://www.gin-vue-admin.com/guide/introduce/project.html)，但后端使用 Iris，
前端使用 React/TanStack 技术栈，不直接复制其代码或目录。

当前 Phase 0～9 已完成，已具备发布一个完整脚手架版本的基础。

## 项目目标

- 提供可复用的用户、角色、菜单、按钮和 API 权限管理能力。
- 使用 JWT 完成多客户端认证，使用 Casbin 实现 RBAC 授权。
- 通过 OpenAPI 生成接口文档和前端 TypeScript Client。
- 默认使用 PostgreSQL 17，保留切换 MySQL 的能力。
- 支持 `dev`、`staging`、`prod` 三套环境配置。
- 同时支持 Web 后台、App、小程序和外部 API。
- 最终可通过 Docker Compose 部署。

## 非目标（首期不做）

- 微服务、API 网关、消息队列和分布式事务。
- 插件市场、AI 代码生成、可视化表单生成器。
- Kubernetes、多云对象存储、分片上传和断点续传。
- 分布式文件管理、定时任务、监控大屏和多租户。

这些能力在基础权限系统稳定后按真实需求增加，避免先维护一套暂时用不到的复杂框架。

## 目录规划

```text
.
├── src/
│   ├── server/                 # Go 后端
│   │   ├── cmd/server/         # 程序入口
│   │   ├── configs/            # dev/staging/prod 非敏感配置
│   │   ├── docs/               # OpenAPI 产物
│   │   ├── internal/
│   │   │   ├── bootstrap/      # 初始化顺序与资源释放
│   │   │   ├── config/         # Viper 配置结构
│   │   │   ├── dto/            # 与 Web 框架无关的请求/响应结构
│   │   │   ├── entity/         # GORM 数据库实体
│   │   │   ├── repository/     # 数据访问
│   │   │   ├── service/        # 业务逻辑
│   │   │   └── transport/
│   │   │       └── iris/       # Iris 路由、Handler 和中间件适配
│   │   ├── migrations/
│   │   │   ├── postgres/       # PostgreSQL Goose SQL 迁移
│   │   │   └── mysql/          # MySQL Goose SQL 迁移
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── go.sum
│   └── web/                    # React 管理后台
├── deploy/
│   ├── docker-compose.yml      # 本地依赖：PostgreSQL、Redis、可选 MySQL
│   ├── compose.yaml            # staging/prod 应用编排
│   ├── compose.staging.yaml
│   ├── compose.prod.yaml
│   └── nginx.conf
├── .env.example
├── LICENSE
└── readme.md
```

首期保持单体应用和单 Go module。`service`、`repository` 直接使用具体类型，不为单一实现预建接口。
业务逻辑不依赖 Iris，框架相关代码集中在 `transport/iris`。切换 Gin 时替换该目录即可，不设计一套
自制的通用 Router/Context 接口。

## 技术栈

### 后端

| 领域     | 选型                       | 用途                          |
| ------ | ------------------------ | --------------------------- |
| 语言     | Go                       | 后端开发                        |
| Web    | Iris                     | HTTP 服务、路由和中间件              |
| ORM    | GORM                     | 数据访问及 MySQL/PostgreSQL 方言切换 |
| 数据库    | PostgreSQL 17              | 默认主数据库                      |
| 可选数据库  | MySQL ≥ 5.7              | 使用 InnoDB                   |
| 缓存     | Redis                    | 登录会话、令牌撤销和必要缓存              |
| 认证     | JWT + Redis Session      | Cookie/Bearer 多客户端认证          |
| 授权     | Casbin RBAC              | API 权限校验                    |
| 权限同步   | Casbin Redis Watcher     | 多实例策略变更通知                   |
| 配置     | Viper + 环境变量             | 多环境配置                       |
| 日志     | Zap                      | 结构化日志                       |
| 日志轮转   | Lumberjack               | 非容器部署时的文件日志轮转               |
| 校验     | go-playground/validator  | 请求参数校验                      |
| API 文档 | OpenAPI/Swagger          | 文档及前端 Client 生成             |
| 数据迁移   | Goose                    | 可回滚 SQL 版本迁移                |
| 读写分离   | GORM dbresolver          | 后期按负载启用                     |

迁移工具首期只使用 Goose，不同时维护 Atlas。数据库结构由迁移文件负责，不在生产环境执行
`AutoMigrate`。PostgreSQL 和 MySQL 分别维护迁移目录，不尝试用一份 SQL 兼容两种方言。

### 前端

- Bun
- React
- TanStack Start / Router / Query
- Zustand
- Tailwind CSS
- shadcn/ui
- Aceternity UI Sidebar + Motion
- Zod
- React Hook Form
- Hey API OpenAPI Client

## 环境约定

环境名固定为：

| 环境  | `APP_ENV` | 用途                      |
| --- | --------- | ----------------------- |
| 开发  | `dev`     | 本地开发、调试日志、允许 Swagger    |
| 演示  | `staging` | 演示和上线前验证，使用独立数据库        |
| 生产  | `prod`    | 生产环境，关闭调试信息和 Swagger UI |

非敏感默认值放入 `src/server/configs/config.<env>.yaml`。密码、JWT 密钥、数据库地址等敏感配置只通过
环境变量或部署平台 Secret 注入，不提交到仓库。

建议的必要环境变量：

```dotenv
APP_ENV=dev
HTTP_ADDR=:8080
DB_DRIVER=postgres
DB_DSN=
DB_READ_WRITE_ENABLED=false
# 开启读写分离时才需要；不填 DB_WRITER_DSN 时默认复用 DB_DSN 作为写库
DB_WRITER_DSN=
DB_READER_DSNS=
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
JWT_SECRET=
COOKIE_SECURE=false
```

配置优先级：环境变量 > 当前环境配置文件 > 代码安全默认值。启动时必须校验必要配置，缺失即退出。
`JWT_SECRET` 不提供代码默认值。`staging`、`prod` 要求至少 32 字节；检测到空值、长度不足或
`secret`、`change_me` 等常见开发占位值时拒绝启动。该校验保留一个最小单元测试。

数据库默认单 DSN 运行，不启用读写分离。关闭时只配置 `database.dsn`/`DB_DSN`，不要把读库和写库
重复填成同一个地址。开启读写分离时使用 GORM dbresolver：一个 writer、多个 readers；事务、Casbin
策略和初始化脚本强制走 writer，普通查询/写入由 GORM 回调分发。

```yaml
database:
  driver: postgres
  dsn: postgres://app:password@writer:5432/secondadmin?sslmode=disable
  readWrite:
    enabled: true
    writerDsn: "" # 空值表示复用 database.dsn
    readerDsns:
      - postgres://app:password@reader-1:5432/secondadmin?sslmode=disable
      - postgres://app:password@reader-2:5432/secondadmin?sslmode=disable
```

认证相关的非敏感配置放入 YAML：

```yaml
auth:
  accessTokenTTL: 15m
  refreshTokenTTL: 168h
  maxDevices: 5 # 0 表示不限制
  cookie:
    secure: false
    sameSite: lax
  corsOrigins:
    - http://localhost:3000
```

`staging` 和 `prod` 的 `cookie.secure` 必须为 `true`。TTL 和设备数可按环境覆盖。配置修改后重启
服务，不使用 fsnotify 热更新运行中的数据库、日志或认证组件。

## 后端架构

```mermaid
flowchart LR
    Clients[Web / App / 小程序 / 外部 API] -->|/api/v1| Iris[Iris Transport]
    Iris --> MW[Recovery / Log / JWT / Casbin]
    MW --> Handler[Handler]
    Handler --> Service[Service]
    Service --> Repo[Repository]
    Repo --> DB[(PostgreSQL / MySQL)]
    Service --> Redis[(Redis Session / PubSub)]
    Casbin[Casbin Enforcer] --> DB
    Casbin <--> Watcher[Redis Watcher]
    Watcher <--> Redis
    MW --> Casbin
```

职责边界：

- Transport：解析和校验 HTTP 输入，适配 Iris 路由及中间件，调用 Service。
- Service：业务规则、事务边界，不依赖 Iris Context。
- DTO：只表达接口输入输出，不带 GORM Tag，不与数据库实体混用。
- Entity：只表达数据库表结构和关联关系。
- Repository：封装 GORM 查询和事务，不向 Service 暴露 `*gorm.DB`。
- 认证、授权规则放在框架无关的业务函数中，Iris 中间件只负责提取请求并返回结果。
- Bootstrap：按顺序初始化配置、日志、数据库、Redis、Casbin、路由和 HTTP 服务。

路由可替换的边界是 HTTP Transport，不是运行时插件。首期不同时实现 Iris 和 Gin；需要切换时新增
`transport/gin` 并删除 Iris 适配，DTO、Service、Repository、Entity 和数据库迁移保持不变。

使用一个具体的 `Repositories` 聚合结构持有各业务 Repository，并提供
`WithinTransaction(ctx, func(tx *Repositories) error)`。事务回调可同时操作用户、角色等多个
Repository；Service 决定事务边界，但不接触 `*gorm.DB`。不把事务对象隐式塞入
`context.Context`，也不增加只有一个实现的 `TransactionManager` 或 Unit of Work 接口。

## 首期功能范围

### 1. 基础设施

- 配置加载与启动校验。
- Zap 日志、请求 ID、访问日志和 panic recovery。
- PostgreSQL/MySQL 驱动选择及连接池配置。
- Redis 连接。
- PostgreSQL/MySQL 两套 Goose 数据库迁移。
- `/healthz` 存活检查和 `/readyz` 依赖就绪检查。
- 统一响应、错误码、分页参数和参数校验。
- OpenAPI 文档生成。
- 优雅停机。

### 2. 认证与账号

- 用户名密码登录、退出和当前用户信息。
- 密码使用 bcrypt 哈希，禁止保存明文密码。
- Web 后台支持 `HttpOnly` Cookie；App、小程序和外部 API 支持 `Authorization: Bearer`。
- 短期 Access Token + 可轮换 Refresh Token，Refresh Token 只保存哈希。
- Redis 是有效设备会话的唯一事实源，保存 Refresh Token 哈希、设备信息和 TTL。
- 每次认证只查询 Redis；Redis 不可用时安全失败并返回 503，不从 PostgreSQL 重建会话。
- PostgreSQL 只记录登录、刷新、退出和设备淘汰日志，不参与在线会话判定。
- 会话写入 Redis 成功后即可完成登录；审计日志写入失败只记录错误，不回滚已创建会话。
- 同账号最大在线设备数由 YAML 的 `auth.maxDevices` 配置；`0` 表示不限制。
- 超限时默认淘汰最早活跃设备，不额外增加交互流程。
- 登录失败限流和账号状态检查。
- 认证中间件优先检查 `Authorization: Bearer`：存在但无效时直接返回 401，不回退 Cookie。
- 没有 Bearer 时才读取 Cookie，并对 POST、PUT、PATCH、DELETE 强制校验 CSRF Token。
- CSRF 判断依据是本次成功认证的凭证来源，不依赖 User-Agent 或客户端自报类型。
- Web 跨域来源使用 YAML 白名单，禁止 `Access-Control-Allow-Origin: *` 与凭证同时启用。

### 3. RBAC 权限

- 用户管理。
- 角色管理。
- 用户分配角色。
- 菜单树和按钮权限标识。
- 角色分配菜单、按钮和 API 权限。
- Casbin 校验“角色 + HTTP 方法 + API 路径”。
- 权限写入数据库成功后，本实例刷新策略并通过 Redis Watcher 发布通知。
- 其他实例收到通知后执行 `LoadPolicy()`，保证多实例权限最终同步。
- Redis Watcher 断线后自动重连；重连成功后主动执行一次 `LoadPolicy()`。

### 4. 系统管理

- API 资源管理。
- 数据字典和字典项。
- 操作日志查询。
- 基础系统配置读取；首期不允许后台直接改写服务器配置文件。

## 首期数据模型

所有业务表使用统一的 `id`、`created_at`、`updated_at`，需要软删除的表增加 `deleted_at`。

| 表                      | 说明                    |
| ---------------------- | --------------------- |
| `sys_users`            | 用户、密码哈希、昵称、状态         |
| `sys_roles`            | 角色编码、名称、状态            |
| `sys_menus`            | 目录、菜单、按钮及父子关系         |
| `sys_role_menus`       | 角色与菜单/按钮关联            |
| `sys_apis`             | API 路径、HTTP 方法和分组     |
| `casbin_rule`          | Casbin 策略             |
| `sys_dictionaries`     | 字典类型                  |
| `sys_dictionary_items` | 字典项                   |
| `sys_operation_logs`   | 操作人、方法、路径、状态码、耗时和必要摘要 |
| `sys_login_logs`       | 登录、刷新、退出、淘汰结果及设备摘要 |

`sys_users.username`、`sys_roles.code`、字典类型编码建立唯一索引。外键是否由数据库强制执行在首次迁移时统一决定，不混用两种策略。

在线设备列表直接读取 Redis，不维护 PostgreSQL 会话副本，避免双写一致性问题。

## API 约定

- 基础路径：`/api/v1`
- JSON 字段：`camelCase`
- 时间：RFC 3339
- 分页：`page`、`pageSize`
- 排序字段使用白名单，不直接拼接客户端输入。
- 成功 HTTP 状态码保持语义正确；业务错误不统一伪装成 HTTP 200。

统一错误响应：

```json
{
  "code": "AUTH_INVALID_CREDENTIALS",
  "message": "用户名或密码错误",
  "requestId": "..."
}
```

首期路由：

```text
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
GET    /api/v1/auth/me
GET    /api/v1/auth/sessions
DELETE /api/v1/auth/sessions/:id

GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
PUT    /api/v1/users/:id/roles

GET    /api/v1/roles
POST   /api/v1/roles
PUT    /api/v1/roles/:id
DELETE /api/v1/roles/:id
PUT    /api/v1/roles/:id/menus
PUT    /api/v1/roles/:id/apis

GET    /api/v1/menus
POST   /api/v1/menus
PUT    /api/v1/menus/:id
DELETE /api/v1/menus/:id
GET    /api/v1/menus/current

GET    /api/v1/apis
POST   /api/v1/apis
PUT    /api/v1/apis/:id
DELETE /api/v1/apis/:id

GET    /api/v1/dictionaries
POST   /api/v1/dictionaries
PUT    /api/v1/dictionaries/:id
DELETE /api/v1/dictionaries/:id

GET    /api/v1/operation-logs
GET    /api/v1/login-logs
```

## 认证授权流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant W as Web / App / 小程序 / API
    participant A as Iris API
    participant R as Redis
    participant C as Casbin
    participant D as Database

    U->>W: 提交登录
    W->>A: POST /api/v1/auth/login
    A->>D: 查询用户并校验密码
    A->>R: 保存设备会话、Refresh Token 哈希和 TTL
    A->>D: 记录登录审计
    A-->>W: Cookie 或 Bearer Token

    W->>A: Cookie 或 Authorization Bearer
    A->>A: 校验 JWT
    A->>R: 校验会话状态
    A->>C: 校验角色、方法和路径
    C-->>A: 允许或拒绝
    A-->>W: 返回结果
```

## 初始化顺序

1. 读取 `APP_ENV` 并加载配置。
2. 校验必填配置。
3. 初始化 Zap。
4. 连接数据库。
5. 连接 Redis。
6. 初始化 Redis Watcher 并加载 Casbin 策略。
7. 注册中间件和路由。
8. 启动 Iris HTTP Server。

收到退出信号后的关闭顺序：

1. 标记服务为未就绪，停止接收新请求。
2. 在超时时间内等待 Iris 当前 HTTP 请求完成。
3. 停止 Casbin Redis Watcher。
4. 关闭 Redis 连接。
5. 关闭数据库连接。
6. 刷新并关闭日志。

首期没有 WebSocket。以后增加长连接时，必须先通知并关闭长连接，再等待 HTTP Shutdown，不能依赖
普通请求超时自动结束。

数据库迁移作为独立命令执行，不隐藏在服务启动流程中。

## 日志策略

- `dev`：彩色控制台日志，可选 Lumberjack 文件日志。
- `staging`：JSON 控制台日志。
- `prod`：JSON 输出到 stdout/stderr，由 Docker logging driver 收集。
- 单机非容器部署时才默认启用 Lumberjack。
- 操作日志只记录必要请求摘要，密码、Cookie、Token 和敏感字段必须脱敏或忽略。

## 实施阶段

当前状态：

- Phase 0：已完成
- Phase 1：已完成
- Phase 2：已完成
- Phase 3：已完成
- Phase 3.1：已完成
- Phase 4：已完成
- Phase 5：权限闭环打磨，已完成
- Phase 6：OpenAPI + Hey API 生成链路，已完成
- Phase 7：通用 CRUD 样板页，已完成
- Phase 8：登录日志与操作日志闭环，已完成
- Phase 9：初始化与部署体验打磨，已完成

### Phase 0：工程骨架

- 创建 `src/server` Go module 和最小目录。
- 完成配置、日志、Iris 启动、健康检查、优雅停机。
- 增加生产环境弱 `JWT_SECRET` 拒绝启动检查及一个最小测试。
- 接入数据库、Redis、Casbin Redis Watcher、双数据库 Goose 迁移和 OpenAPI。
- 提供 `dev` 环境的 Docker Compose 依赖。

验收：服务可启动；健康检查正常；PostgreSQL/MySQL 迁移均可执行和回滚；缺失关键配置时拒绝
启动；停止 HTTP 后才释放 Redis 和数据库连接。

### Phase 1：认证与核心权限

- 完成用户、角色和多设备登录会话。
- 完成 JWT Cookie/Bearer、刷新令牌、CSRF 和 Casbin 中间件。
- 完成用户分配角色、角色分配 API 权限。
- 初始化一个管理员账号和管理员角色。

验收：未登录返回 401；无权限返回 403；Bearer 请求跳过 CSRF；Cookie 写请求缺少 CSRF Token
时被拒绝；Redis 不可用时认证安全失败；退出或设备被淘汰后令牌立即失效；两个服务实例的权限变更
能够通过 Redis Watcher 同步。

#### Phase 1 实施边界

Phase 1 只完成认证和 API 级 RBAC 闭环：

- 实现用户、角色、API 资源和登录日志；用户角色关系直接使用 Casbin `g` 策略。
- 实现登录、刷新、退出、当前用户和在线设备管理。
- 实现用户/角色 CRUD、用户分配角色、角色分配 API 权限。
- 菜单、按钮、字典、操作日志查询和后台配置留在 Phase 2。
- 不实现 OAuth、短信登录、验证码、找回密码、MFA 和第三方登录。

#### Phase 1 数据表

PostgreSQL 和 MySQL 分别新增 `00002_phase1_auth.sql`：

| 表 | Phase 1 字段 |
| --- | --- |
| `sys_users` | `id`、`username`、`password_hash`、`nickname`、`status`、`password_changed_at`、时间戳 |
| `sys_roles` | `id`、`code`、`name`、`status`、时间戳 |
| `sys_apis` | `id`、`group`、`name`、`path`、`method`、时间戳 |
| `sys_login_logs` | `id`、`user_id`、`username`、`event`、`success`、`ip`、`user_agent`、`device_id`、`created_at` |

约束：

- `username` 和角色 `code` 保存前转为小写，建立唯一索引。
- `sys_apis(path, method)` 建立唯一索引，HTTP 方法统一保存为大写。
- 用户和角色首期不做软删除；有关联或审计记录时改为禁用，避免唯一索引和历史关系复杂化。
- 密码使用 bcrypt；数据库只保存哈希。
- 不创建 `sys_user_roles`：`casbin_rule` 中 `g, userID, roleCode` 是用户角色关系的唯一事实源，
  `p, roleCode, path, method` 是角色 API 权限的唯一事实源。

#### Token 与 Redis 会话

- Access Token 使用 JWT，只包含 `sub`（用户 ID）、`sid`（会话 ID）、`iat`、`exp`。
- Access Token 不保存角色，避免角色变更后旧 Token 保留旧权限。
- Refresh Token 使用 `crypto/rand` 生成的随机字符串，不使用 JWT。
- Redis 只保存 Refresh Token 的 SHA-256 哈希，不保存原文。
- 每次刷新后轮换 Refresh Token，旧 Token 立即失效。
- Casbin 以用户 ID 作为 subject，通过 `g(userID, roleCode)` 维护用户角色关系。

Redis Key 固定为：

```text
auth:session:{sessionID}       # Hash，会话、Refresh Token 哈希、CSRF 哈希，TTL=refreshTokenTTL
auth:user:{userID}:sessions   # ZSet，用户的会话 ID，score=创建时间
auth:login:{ip}:{username}    # 登录失败计数，短 TTL
```

登录时使用一段 Redis Lua 脚本原子写入新会话并淘汰最早会话，避免并发登录突破
`auth.maxDevices`。被淘汰、退出或用户禁用时删除对应 Session Key。

#### Cookie、Bearer 与 CSRF

登录请求显式传递 `authMode: cookie | bearer`，服务端不通过 User-Agent 猜客户端：

- `cookie`：Access Token 和 Refresh Token 写入 `HttpOnly` Cookie；CSRF Token 写入可读 Cookie，
  同时保存哈希到 Redis Session。
- `bearer`：Access Token 和 Refresh Token 返回 JSON，不设置认证 Cookie。
- Cookie 认证的 POST、PUT、PATCH、DELETE 请求必须携带 `X-CSRF-Token`，并与当前 Session 中的
  哈希匹配。
- 存在 `Authorization` 请求头时只使用 Bearer 认证；无效 Bearer 不回退 Cookie。
- Refresh Cookie 只允许发送到 `/api/v1/auth/refresh`。

#### Phase 1 代码增量

```text
src/server/
├── cmd/bootstrap-admin/       # 显式创建首个管理员
├── internal/
│   ├── dto/                   # Auth/User/Role 请求响应
│   ├── entity/                # User/Role/API/LoginLog
│   ├── repository/            # Repositories 聚合及事务
│   ├── service/               # Auth/User/Role 服务
│   └── transport/iris/
│       ├── middleware.go      # JWT、Session、CSRF、Casbin
│       └── routes.go          # Phase 1 路由
└── migrations/
    ├── postgres/00002_phase1_auth.sql
    └── mysql/00002_phase1_auth.sql
```

不按每个实体拆成大量小文件；文件变长到影响阅读时再拆分。

#### Phase 1 执行顺序

1. **1A 数据层**：双数据库迁移、Entity、Repositories 聚合、事务检查。
2. **1B 认证核心**：密码哈希、JWT、随机 Refresh Token、Redis Session 和设备限制。
3. **1C HTTP 认证**：登录、刷新、退出、`me`、设备列表、Cookie/Bearer、CSRF、CORS。
4. **1D 核心 RBAC**：用户、角色、用户角色、API 权限和 Casbin 中间件。
5. **1E 初始化与契约**：`bootstrap-admin`、OpenAPI 更新、最小集成测试。

每一步完成后保持服务可编译、迁移可回滚；不等到 Phase 1 末尾一次性集成。

#### Phase 1 验收清单

- PostgreSQL/MySQL 的 `00002` 均通过 `up → down → up`。
- 密码哈希、JWT 签名、Refresh Token 轮换各有一个最小测试。
- Bearer 和 Cookie 两种登录方式均可访问受保护接口。
- Cookie 写请求缺少或伪造 CSRF Token 时返回 403。
- Access Token 过期后可刷新；旧 Refresh Token 重放失败。
- 达到设备上限时最早会话失效。
- 用户禁用、退出或设备删除后对应会话立即失效。
- 未登录返回 401，无 API 权限返回 403。
- 用户角色或角色 API 权限变更后立即生效。
- 两个服务实例通过 Redis Watcher 同步 Casbin 策略。
- `bootstrap-admin` 重复执行不会创建重复管理员；不提供默认生产密码。

### Phase 2：动态菜单与系统管理

- 完成菜单、按钮权限和当前用户菜单树。
- 完成数据字典和操作日志。
- 固化 OpenAPI，并验证 Hey API 可生成前端 Client。

验收：不同角色得到不同菜单和按钮权限；关键写操作可审计；OpenAPI Client 可生成。

#### Phase 2 实施边界

Phase 2 只补齐后台导航和基础系统管理：

- 菜单类型固定为 `directory`、`menu`、`button`。
- 角色分配菜单和按钮；当前用户只获取已授权且启用的菜单树。
- 数据字典类型和字典项 CRUD。
- 记录已认证写请求的操作日志并提供分页查询。
- API 资源管理已在 Phase 1 完成，本阶段只补 OpenAPI 契约，不重复实现。
- 不实现菜单缓存、字典缓存、日志消息队列、日志导出和日志全文检索。

#### Phase 2 数据表

PostgreSQL 和 MySQL 分别新增 `00003_phase2_system.sql`：

| 表 | 核心字段 |
| --- | --- |
| `sys_menus` | `id`、`parent_id`、`type`、`name`、`path`、`component`、`icon`、`sort`、`visible`、`permission`、`status`、时间戳 |
| `sys_role_menus` | `role_id`、`menu_id`，联合主键 |
| `sys_dictionaries` | `id`、`code`、`name`、`status`、时间戳 |
| `sys_dictionary_items` | `id`、`dictionary_id`、`label`、`value`、`sort`、`status`、时间戳 |
| `sys_operation_logs` | `id`、`user_id`、`request_id`、`method`、`path`、`status_code`、`duration_ms`、`ip`、`user_agent`、`created_at` |

约束：

- `sys_menus.parent_id` 默认 `0` 表示根节点，不使用可空父级。
- 同一父级下菜单名称唯一；`permission` 使用可空字段，非空时全局唯一。
- `button` 不允许拥有子节点，也不返回路由字段。
- `sys_role_menus` 使用外键级联删除关联行，不级联删除角色或菜单。
- 字典 `code` 全局唯一，字典项 `(dictionary_id, value)` 唯一。
- 操作日志不保存请求体、响应体、Cookie、Token 或密码。
- 菜单、字典首期不做软删除；存在子节点或关联时拒绝删除，先解除关系。

#### 菜单与按钮权限

- 菜单和按钮授权使用 `sys_role_menus`，不塞进 Casbin；Casbin 继续只负责 API 权限。
- 当前用户角色仍从 Casbin `g(userID, roleCode)` 读取，再关联启用的 `sys_roles` 查询菜单。
- 多角色菜单取并集，按 `sort ASC, id ASC` 排序后在 Service 内一次构建树。
- `directory`、`menu` 返回给前端作为导航；`button.permission` 汇总为字符串数组。
- 新建或移动菜单时校验父节点存在、父节点不是 `button`、节点不能成为自己的后代。
- 删除只允许叶子节点；避免自动级联删除整棵菜单树。
- 首期每次请求直接查数据库，不加 Redis 缓存；菜单规模通常很小。

当前用户菜单响应：

```json
{
  "menus": [
    {
      "id": 1,
      "type": "directory",
      "name": "系统管理",
      "path": "/system",
      "children": []
    }
  ],
  "permissions": ["system:user:create"]
}
```

#### 数据字典

- 字典类型与字典项使用同一组 Repository/Service，不引入通用“配置中心”。
- 后台管理接口需要 Casbin API 权限。
- 提供一个已认证公共读取接口：
  `GET /api/v1/dictionaries/{code}/items`，只返回启用项。
- 字典项按 `sort ASC, id ASC` 返回。
- 字典禁用后，其公共读取接口返回空数组，不回退历史值。

#### 操作日志

在 Iris Transport 增加一个审计中间件：

- 放在认证之后、Casbin 授权之前，因此可记录成功、403 和业务失败的已认证写请求。
- 只记录 `/api/v1` 下的 POST、PUT、PATCH、DELETE。
- `/auth/login`、`/auth/refresh` 继续使用 `sys_login_logs`，不重复写操作日志。
- `ctx.Next()` 返回后读取最终状态码和耗时，同步写入数据库。
- 日志写入失败只写 Zap error，不修改已经生成的 HTTP 响应。
- 审计中间件复用请求 ID；同时修正当前错误响应从 Context 读取请求 ID，而不是读取请求头。

同步写入是首期的刻意简化；只有实测审计写入成为延迟瓶颈时，才增加批量队列。

#### Phase 2 API

```text
GET    /api/v1/menus
POST   /api/v1/menus
PUT    /api/v1/menus/:id
DELETE /api/v1/menus/:id
GET    /api/v1/menus/current
PUT    /api/v1/roles/:id/menus

GET    /api/v1/dictionaries
POST   /api/v1/dictionaries
PUT    /api/v1/dictionaries/:id
DELETE /api/v1/dictionaries/:id
POST   /api/v1/dictionaries/:id/items
PUT    /api/v1/dictionary-items/:id
DELETE /api/v1/dictionary-items/:id
GET    /api/v1/dictionaries/:code/items

GET    /api/v1/operation-logs
GET    /api/v1/login-logs
```

除 `menus/current` 和字典项公共读取外，以上接口都经过现有 Casbin 中间件。

#### Phase 2 代码增量

复用现有聚合文件，先不按实体拆包：

```text
src/server/
├── internal/
│   ├── dto/dto.go
│   ├── entity/entity.go
│   ├── repository/repository.go
│   ├── service/service.go
│   └── transport/iris/
│       ├── middleware.go      # 增加操作审计
│       └── routes.go          # 增加 Phase 2 路由
└── migrations/
    ├── postgres/00003_phase2_system.sql
    └── mysql/00003_phase2_system.sql
```

现有聚合文件明显影响阅读时才按领域拆分；不预建 `menu`、`dictionary`、`audit` 子包。

#### Phase 2 执行顺序

1. **2A 数据层**：双数据库迁移、Entity、Repository 和删除约束。
2. **2B 菜单权限**：菜单 CRUD、角色菜单、当前用户菜单树和按钮权限。
3. **2C 数据字典**：字典类型、字典项和公共读取接口。
4. **2D 操作审计**：请求 ID 修正、审计中间件和分页查询。
5. **2E 契约验收**：更新 OpenAPI，在临时目录运行 Hey API Client 生成，不提交前端产物。

每一步完成后保持双数据库迁移可回滚、服务可编译，不提前进入 Docker 生产部署。

#### Phase 2 验收清单

- PostgreSQL/MySQL 的 `00003` 均通过 `up → down → up`。
- 菜单父子校验、循环检测和树构建保留一个最小测试。
- 不同角色获得不同菜单；多角色菜单正确合并且无重复。
- 禁用菜单和禁用角色不出现在当前用户菜单中。
- 按钮只出现在 `permissions`，不混入前端路由树。
- 删除有子节点或仍被角色引用的菜单返回 409。
- 字典 code 和同字典 value 重复时返回 409。
- 公共字典接口只返回启用字典的启用项。
- 已认证写请求的 2xx、4xx、5xx 均能记录状态码和请求 ID。
- 操作日志不包含认证头、Cookie、请求体和响应体。
- 操作日志分页最大 `pageSize=100`。
- OpenAPI JSON 有效，Hey API 能在临时目录生成 TypeScript Client。

### Phase 3：Web 管理后台

使用既定的 React/TanStack 技术栈完成可用的管理后台。Phase 3 结束后再统一处理前后端镜像与部署。

#### Phase 3 实施边界

- 使用 Bun、React、TypeScript、TanStack Start/Router/Query、Tailwind CSS、shadcn/ui。
- 左侧导航使用 Aceternity UI Sidebar Registry，保留 Motion 悬停展开和移动端抽屉，并适配项目
  Lucide 图标、语义颜色、TanStack Router 与 `prefers-reduced-motion`。
- React Hook Form + Zod 只用于表单输入和客户端提示；后端仍是最终校验边界。
- Hey API 根据 OpenAPI 生成 Client，业务页面不手写重复的请求类型。
- Web 固定使用 `authMode: cookie`，请求统一携带 `credentials: include`。
- 后台不需要 SEO，首期不做 SSR 数据预取、服务端 Session 转发或 BFF。
- Zustand 作为基础依赖，使用一个轻量 UI Store 管理侧边栏折叠、主题和布局偏好；用户、菜单和列表
  数据交给 TanStack Query。
- 不实现大屏图表、拖拽布局、国际化、多标签页、主题市场和低代码页面。

#### Phase 3 前置修正：OpenAPI 契约

OpenAPI 已补齐以下请求/响应 Schema，并以生成 Client 作为 Web 接口类型真源：

- 登录、刷新、当前用户和会话响应。
- 用户、角色、API、菜单、字典、操作日志模型。
- 分页响应、错误响应和所有写请求 Body。
- Cookie/CSRF Header 描述。
- 每个 `operationId` 保持稳定且唯一。

生成命令固定写入 `src/web/package.json`：

```json
{
  "scripts": {
    "api:generate": "openapi-ts",
    "api:check": "bun scripts/check-openapi.mjs",
    "api:sync": "bun run api:generate && bun run api:check"
  }
}
```

Hey API 配置固定在 `src/web/openapi-ts.config.ts`。生成目录不手改；`api:check` 会校验所有接口都有
唯一 `operationId`，并确认 `src/api/generated` 内的 TypeScript 文件带有 Hey API 自动生成标记。
OpenAPI 变化后执行 `bun run api:sync`，再由 TypeScript 编译暴露不兼容变更。

#### Web 目录规划

```text
src/web/
├── src/
│   ├── api/
│   │   ├── client.ts          # credentials、CSRF 和一次刷新重试
│   │   └── generated/         # Hey API 生成代码
│   ├── components/
│   │   ├── layout/            # Header、Sidebar、Content
│   │   └── ui/                # shadcn/ui 按需生成
│   ├── features/
│   │   ├── auth/
│   │   ├── users/
│   │   ├── roles/
│   │   ├── apis/
│   │   ├── menus/
│   │   ├── dictionaries/
│   │   ├── operation-logs/
│   │   └── login-logs/
│   ├── lib/                   # queryClient、权限判断和小型通用函数
│   ├── routes/                # TanStack 文件路由
│   ├── store/ui.ts            # Zustand UI 状态
│   └── styles.css
├── components.json
├── package.json
├── tsconfig.json
└── vite.config.ts
```

按业务功能组织页面，不增加 `repositories`、前端 Service 层或每个请求一层 Hook 包装。只有被多个页面
复用的 Query 才提取为 Hook。

#### 认证与 API Client

- 浏览器不读取 Access/Refresh Token；两者保持 HttpOnly Cookie。
- 从 `csrf_token` Cookie 读取 CSRF Token，仅为写请求添加 `X-CSRF-Token`。
- 收到 401 时只允许一次 Refresh；并发 401 共用同一个刷新 Promise，避免刷新风暴。
- Refresh 成功后原请求重试一次；再次 401 时清空前端 Query Cache 并跳转登录页。
- 403 显示无权限，不自动刷新。
- 登录成功后获取 `auth/me` 和 `menus/current`，再进入后台布局。
- 退出后清空 Query Cache 并返回登录页。

首期不把用户、菜单、权限或业务列表复制到 Zustand；`auth/me`、`menus/current` 和业务接口是服务端
状态唯一来源。UI Store 需要持久化的字段通过 Zustand persist 保存到 localStorage。

#### 路由、菜单与按钮权限

TanStack Router 路由在构建时确定，后端菜单不能动态加载任意 React 文件：

- 后端 `component` 字段只作为页面标识，映射到前端本地组件注册表。
- 未在注册表中的 component 显示“页面未实现”，不执行动态字符串 import。
- `menus/current` 决定侧边栏可见项和可访问页面；前端路由守卫防止手工输入隐藏 URL。
- 后端 Casbin 仍是安全边界，前端守卫只改善体验。
- 按钮通过 `can("system:user:create")` 检查 `permissions` 后控制显示。
- 目录不对应业务页面；叶子菜单对应已注册路由。

首期页面注册表直接使用一个对象，不构建插件系统。

#### Phase 3 页面范围

1. **登录**
   - 用户名、密码、登录错误和限流提示。
2. **后台布局**
   - 动态响应式侧边栏、当前用户、主题切换和退出。
3. **用户管理**
   - 分页列表、新建、编辑、启停、重置密码、分配角色。
4. **角色管理**
   - 列表、新建、编辑、删除、分配 API、分配菜单。
5. **API 管理**
   - 列表、新建、编辑、删除。
6. **菜单管理**
   - 树形列表、新建子节点、编辑、删除；不做拖拽排序。
7. **数据字典**
   - 字典与字典项 CRUD。
8. **操作日志**
   - 分页只读表格。
9. **在线设备**
   - 当前账号会话列表和下线设备。

首期表格使用普通响应式 HTML Table/shadcn Table，不引入重型 DataGrid。

#### 表单和交互约定

- Zod Schema 与 React Hook Form 放在对应 feature 中。
- 创建使用页内表单，编辑使用 Dialog；角色授权使用 Dialog。
- 删除、禁用、下线设备必须二次确认。
- 提交中禁用按钮，成功后关闭表单并精准 invalidate 对应 Query。
- 409 显示重复或关联冲突；401 跳登录；403 显示无权限；其他错误展示 `requestId` 便于排查。
- Loading、空数据和错误状态必须明确，不使用全屏旋转动画覆盖整个后台。

#### Phase 3 执行顺序

1. **3A 契约与骨架**：补齐 OpenAPI Schema、创建 TanStack Start、Tailwind、shadcn、Zustand UI
   Store 和生成 Client。
2. **3B 认证闭环**：API Client、CSRF、单次刷新重试、登录、退出、当前用户和路由守卫。
3. **3C 后台壳层**：动态菜单、组件注册表、侧边栏、面包屑和按钮权限。
4. **3D 管理页面**：用户、角色、API、菜单、字典、操作日志和在线设备。
5. **3E 质量验收**：类型检查、构建、关键组件测试和真实后端端到端冒烟。

每一步保持 `bun run build` 可通过；不等所有页面完成后再处理类型错误。

#### Phase 3 最小检查

- `bun run api:sync`
- `bun run typecheck`
- `bun run test`
- `bun run build`
- 对 CSRF Header、401 单次刷新和权限判断保留最小测试。
- 使用真实后端验证登录、刷新、退出、动态菜单和一个完整 CRUD。

自动化检查已通过；当前执行环境禁止浏览器访问 `localhost:3000`，真实后端浏览器冒烟需要在本机手工
复核，不因此进入 Phase 3.1 或 Phase 4。

不在首期引入完整浏览器 E2E 框架；当页面流程稳定且 CI 需要回归保护时再增加 Playwright。

#### Phase 3 验收清单

- 未登录访问后台路由跳转登录页。
- Cookie 登录后刷新浏览器仍保持会话。
- 多个并发 401 只触发一次 Refresh 请求。
- CSRF Token 自动添加到写请求，缺失时能显示后端 403。
- 动态菜单顺序、层级和隐藏状态正确。
- 无权限按钮不显示，手工访问页面会被前端守卫拦截，直接请求 API 仍由后端返回 403。
- 用户、角色、API、菜单和字典 CRUD 可用。
- 角色 API 与菜单授权修改后界面立即刷新。
- 操作日志与在线设备分页/列表可用。
- OpenAPI Client 可重复生成且工作区无手工生成代码修改。
- `bun run typecheck` 和 `bun run build` 通过。

### Phase 3.1：Web 视觉完善

复用根目录 `Illustration/` 中现有 SVG，不引入额外插画库。只按实际状态挑选并复制到
`src/web/public/illustrations/`：

- 404、无权限、502 和网络断开页面。
- 用户、角色、API、菜单、字典、日志和在线设备的空数据状态。
- 登录或页面加载状态仅在确有等待时使用，不用插画替代普通骨架屏。

保留 SVG `viewBox`，提供可访问的替代文本，并统一尺寸、留白和深浅色背景表现；不要求一次使用全部
插画。验收以常见错误页和空状态视觉一致、移动端不溢出为准。

### Phase 4：部署

- API 使用 Go 多阶段构建和 distroless 非 root 运行镜像；镜像内只保留二进制与运行配置。
- Web 使用 Bun 多阶段构建并运行 TanStack Start 产物；最终镜像不包含源码、依赖缓存或构建工具链。
- Nginx 网关将 `/api`、`/healthz`、`/readyz` 转发到 API，其余请求转发到 Web，浏览器保持同源。
- 最终业务镜像不包含 Goose、迁移 SQL、编译器、Bun 缓存或源码。
- Web 与 API 镜像可独立发布，Compose 默认通过同源网关组合。
- 迁移由 CI/CD Runner 在发布业务容器前执行；首期不额外维护 migration-runner 镜像。
- 提供 `staging`、`prod` Compose 配置、健康检查、数据卷、备份和部署说明。

验收：全新环境可通过文档启动；CI/CD 迁移成功后才发布；Web Cookie、App/小程序 Bearer 均正常；
日志可由容器平台采集。

### Phase 5：权限闭环打磨

- 后端明确内置 `admin` 角色为超级管理员：拥有该角色的用户 API 授权直接放行。
- `admin` 角色不可删除、不可禁用，避免误操作把系统锁死。
- `admin` 当前菜单返回所有启用菜单和按钮权限，不依赖 `sys_role_menus` 是否完整。
- 前端在用户角色、角色状态、角色菜单和菜单节点变更后刷新 `menus/current` 缓存。
- 保留后端 Casbin 作为最终安全边界，前端按钮和路由守卫只改善体验。

验收：管理员不会因 API 策略或角色菜单关联误删而失去后台入口；普通用户仍按 Casbin API 权限和
菜单授权返回 403/隐藏菜单；权限变更后前端无需手动刷新即可更新可见菜单和按钮。

### Phase 6：OpenAPI + Hey API 生成链路

- 固化后端 OpenAPI 输出和 `operationId` 命名。
- 固化前端 Hey API 生成命令、生成目录和 “DO NOT EDIT” 边界。
- 前端业务请求只通过生成 client + query options，不手写 endpoint/type。
- `src/web/openapi-ts.config.ts` 是生成配置唯一入口；`bun run api:sync` 负责生成和检查。

验收：修改后端契约后可一键重新生成前端 API；生成代码无手工修改；类型错误能在 typecheck 暴露。

### Phase 7：通用 CRUD 样板页

- 用户、角色、API、菜单、字典页面在现有 OpenAPI 契约内补齐创建、编辑、删除或禁用的统一体验。
- 支持分页的列表统一使用 URL 页码和 `Pagination`；普通列表统一使用 `DataTable`。
- 不抽万能 CRUD 框架，只沉淀小型共享组件和 query/mutation 写法。
- 每个主要数据区覆盖 loading、error、empty、success 和 refetching。
- 搜索和高级筛选等后端查询参数进入 OpenAPI 后再补，不在前端硬做假筛选。

验收：后续新增普通管理页面可以复制现有 feature 结构落地，不需要重新设计表格、弹窗、确认操作、
分页、刷新提示和权限按钮。

### Phase 8：登录日志与操作日志闭环

- 登录成功/失败记录完整，且不保存密码、Token、Cookie 等敏感信息。
- 关键写操作写入操作日志，包含用户、请求 ID、方法、路径、状态码、耗时、IP 和 UA。
- 前端可查询登录日志和操作日志，支持基础分页。
- `GET /api/v1/login-logs` 和 `GET /api/v1/operation-logs` 均通过 Casbin 后台权限保护。

验收：常见安全审计问题能通过后台查询定位；日志字段不会泄露认证凭证。

### Phase 9：初始化与部署体验打磨

- 保证 `make dev`、`make migrate-up`、`make bootstrap-admin` 文档和实际行为一致。
- Docker/Compose 明确迁移、启动、初始化 admin、健康检查和备份路径。
- PostgreSQL 默认链路完整，MySQL 迁移至少保留可验证路径。
- 本地 Makefile 自动读取 `.env`；不存在 `.env` 时回退 `.env.example`。

验收：全新环境按 README 从零启动，不需要翻代码补环境变量；部署失败能通过健康检查和日志定位。

#### Staging / Production 部署

复制对应环境文件并填写强密码：

```bash
cp deploy/.env.staging.example deploy/.env.staging
# 或
cp deploy/.env.prod.example deploy/.env.prod
```

`POSTGRES_PASSWORD` 当前直接进入 PostgreSQL URL，使用 URL 安全字符；`JWT_SECRET` 至少 32 字节。
生产和演示环境强制使用 Secure Cookie，因此网关前必须由云负载均衡或反向代理提供 HTTPS，并传递
`X-Forwarded-Proto: https`。

部署顺序：

```bash
make deploy-deps DEPLOY_ENV=staging
make deploy-migrate DEPLOY_ENV=staging
make deploy-up DEPLOY_ENV=staging
```

生产环境将 `staging` 替换为 `prod`。常用维护命令：

```bash
make deploy-config DEPLOY_ENV=prod
make deploy-logs DEPLOY_ENV=prod
make deploy-backup DEPLOY_ENV=prod
make deploy-down DEPLOY_ENV=prod
```

迁移在宿主机或 CI Runner 执行，数据库只绑定到 `127.0.0.1`；API 镜像不包含 Goose 和迁移 SQL。
备份输出到未提交的 `backups/`，恢复前应先在隔离环境验证。

## 后续能力

以下功能不进入首期，确认有业务需求后单独设计：

- 文件管理和 S3 兼容对象存储。
- 定时任务。
- 数据导入导出。
- 代码生成器。
- 系统监控面板。
- 多租户。

## 已确认决策

- 默认数据库：PostgreSQL 17；MySQL 作为可选数据库。
- 环境名：`dev`、`staging`、`prod`。
- 客户端：Web 后台、App、小程序和外部 API。
- 认证传输：Web 支持安全 Cookie，非浏览器客户端支持 Bearer Token。
- 设备限制：由 YAML 配置最大在线设备数。
- 首期保留数据字典和操作日志。
- Iris 只存在于 Transport 层，不通过自制 Web 框架抽象制造额外耦合。

## 下一步

在 staging 环境完成 HTTPS、Cookie 登录、数据库迁移、备份恢复和滚动发布演练后再发布生产环境。

## Phase 1 管理员初始化

迁移完成后显式创建首个管理员，不提供默认账号或密码：

```bash
ADMIN_USERNAME=admin \
ADMIN_PASSWORD='replace-with-a-strong-password' \
make bootstrap-admin
```

该命令可重复执行，不会创建重复用户、角色或菜单；同时补齐管理员的系统菜单和按钮权限，不覆盖
后来新增的自定义菜单授权。

## Phase 0 本地运行

```bash
cp .env.example .env
make dev-deps
make migrate-up
ADMIN_USERNAME=admin \
ADMIN_PASSWORD='replace-with-a-strong-password' \
make bootstrap-admin
make dev
```

健康检查：

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/openapi.json
```

可选 MySQL 验证：

```bash
docker compose -f deploy/docker-compose.yml --profile mysql up -d mysql
DB_DRIVER=mysql \
DB_DSN='secondadmin:secondadmin@tcp(127.0.0.1:3306)/secondadmin?parseTime=true' \
make migrate-up
```

## License

本项目采用 [MIT License](https://opensource.org/license/mit)。
