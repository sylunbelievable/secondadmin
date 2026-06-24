# Second Admin Web

Second Admin 的 React 管理后台，使用 Bun、TanStack Start/Router/Query、Zustand、Tailwind CSS、shadcn/ui、React Hook Form、Zod 和 Hey API。

## 常用命令

```bash
bun install
bun run dev
bun run api:sync
bun run typecheck
bun run test
bun run build
```

## 目录约定

- `src/routes/`：TanStack Router 文件路由。
- `src/features/`：业务页面和业务组件。
- `src/components/ui/`：shadcn/ui 本地组件源码。
- `src/api/generated/`：Hey API 生成代码，不手改。
- `src/store/`：Zustand 客户端 UI 状态。

新增页面前先检查根目录 `AGENTS.md`，优先复用现有 feature、query options 和 shadcn/ui 组件。
