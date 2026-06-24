# 前端技术栈规范与 AI 编码上下文

> 本文件用于约束 AI 编程助手、开发者及自动化代码生成工具。  
> 任何新建页面、修改页面、生成组件、重构代码或补充交互的行为，都必须遵守本规范。
> 
> 本项目的核心原则是：
> 
> **优先复用，禁止臆造；优先组合，避免重复；先检索组件，再编写代码。**

---

# 1. 技术栈

## 1.1 运行时与包管理

- **JavaScript Runtime / Package Manager：Bun**

- 安装依赖、运行脚本和执行 CLI 时，默认使用 Bun。

- 除非项目已有明确约定，否则不得混用 npm、pnpm、Yarn。

常用命令：

```bash
bun install
bun add <package>
bun add -d <package>
bun run dev
bun run build
bun run lint
bun run typecheck
bunx <command>
```

禁止在同一项目中无理由生成或提交以下混合锁文件：

```text
package-lock.json
pnpm-lock.yaml
yarn.lock
bun.lock
bun.lockb
```

项目应只保留 Bun 对应的锁文件。

---

## 1.2 React 框架

默认技术体系：

- **React**

- **TanStack Start**

- **TanStack Router**

- **TanStack Query**

默认优先采用 **TanStack Start 项目结构**。

只有在当前项目明确是纯前端 SPA，且不存在 SSR、服务端函数、服务端渲染或全栈路由需求时，才使用：

- Vite

- React

- TanStack Router

- TanStack Query

不得在同一个项目中同时混用：

- TanStack Router

- React Router

- Next.js Router

不得为了实现简单跳转，绕过 TanStack Router 自行操作：

```ts
window.location.href
window.history.pushState
```

外部链接、文件下载、跨域跳转等合理场景除外。

---

## 1.3 样式与组件体系

- **Tailwind CSS**

- **shadcn/ui**

- **Lucide React**

- 项目已有设计令牌和 CSS 变量

shadcn/ui 组件源码默认存放于：

```text
@/components/ui/
```

业务组件默认存放于：

```text
@/components/
```

业务模块私有组件可以存放于对应 feature 目录，例如：

```text
@/features/auth/components/
@/features/users/components/
@/features/orders/components/
```

---

## 1.4 状态、表单与数据校验

- **TanStack Query**：服务端状态

- **Zustand**：客户端全局状态

- **React Hook Form**：表单状态

- **Zod**：运行时校验与数据 Schema

- **Hey API OpenAPI Client**：根据 OpenAPI 生成接口客户端

必须严格区分不同状态的职责，不得用一个工具包揽所有状态。

---

# 2. 总体编码原则

AI 在生成代码时，必须遵循以下优先级：

```text
现有业务组件
    ↓
本地已安装的 shadcn/ui 组件
    ↓
shadcn/ui 官方组件
    ↓
shadcn/ui 官方 Blocks
    ↓
项目已允许的 shadcn 生态组件
    ↓
基于现有原子组件进行业务组合
    ↓
最后才考虑编写项目专用基础组件
```

这里的“自己编写”只允许用于：

- 项目独有的业务展示

- 项目独有的业务流程

- shadcn 官方确实不存在的视觉表现

- 不属于通用 UI 原子的静态布局

- 通过已有 shadcn 原子组件组合出的业务组件

不得自行重复实现已有的通用 UI 原子组件。

---

# 3. AI 开始编码前的强制检查

任何 UI 代码生成任务开始前，AI 必须先完成以下检查。

这不是建议，而是强制前置流程。

---

## 3.1 第一步：读取项目技术上下文

首先检查以下文件是否存在：

```text
package.json
bun.lock
bun.lockb
tsconfig.json
components.json
vite.config.ts
app.config.ts
src/router.tsx
src/routes/
src/components/
src/components/ui/
src/lib/
src/hooks/
src/features/
```

需要确认：

1. 当前项目是 TanStack Start 还是 Vite SPA。

2. 路径别名 `@/` 指向哪里。

3. shadcn 组件目录配置在哪里。

4. 使用 Tailwind CSS 哪个版本。

5. 当前图标库是什么。

6. 已安装哪些状态管理、表单和请求依赖。

7. 项目是否已有相同或类似业务组件。

8. 项目现有目录结构、命名方式和代码风格。

不得脱离现有项目结构，凭空创建另一套架构。

---

## 3.2 第二步：检索现有业务组件

在创建新组件之前，必须搜索：

```text
@/components/
@/features/
@/routes/
```

检索范围至少包括：

- 组件文件名

- 组件导出名称

- 相似业务关键词

- 相似交互行为

- 相似页面区块

- 已有 hooks

- 已有 Schema

- 已有 API 查询封装

例如用户要求“用户筛选侧栏”，必须先搜索：

```text
user-filter
filter-sheet
filter-panel
sidebar-filter
UserFilter
FilterSheet
```

若已有可复用组件，应优先复用或扩展，不得重新创建重复组件。

---

## 3.3 第三步：扫描本地 shadcn/ui 组件

必须检查：

```text
@/components/ui/
```

确认其中现有的 `.tsx` 文件，例如：

```text
accordion.tsx
alert-dialog.tsx
avatar.tsx
badge.tsx
button.tsx
calendar.tsx
card.tsx
checkbox.tsx
command.tsx
dialog.tsx
drawer.tsx
dropdown-menu.tsx
form.tsx
input.tsx
popover.tsx
select.tsx
sheet.tsx
sidebar.tsx
skeleton.tsx
table.tsx
tabs.tsx
toast.tsx
tooltip.tsx
```

不得仅凭记忆判断某个组件是否已安装。

不得因为当前文件没有导入某个组件，就错误地认为项目没有安装该组件。

---

## 3.4 第四步：读取组件源码和现有用法

发现组件存在后，必须继续确认：

- 实际导出的组件名称

- 当前组件封装是否经过项目定制

- Props 是否与官方默认版本不同

- 项目是否增加了 variant

- 项目是否统一了 size

- 项目是否修改了 className

- 项目是否采用 Radix UI 或 Base UI 版本

- 项目中其他页面如何使用该组件

不得仅根据记忆猜测导出名称或 Props。

例如使用 `Select` 前，应先确认本地源码实际导出了哪些成员，再按项目现有方式组合：

```tsx
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
```

---

## 3.5 第五步：映射到 shadcn 官方组件名称

当用户使用中文或业务语言描述组件时，必须先转换为 shadcn/ui 标准组件名称。

常见映射如下：

| 用户描述     | 优先使用的 shadcn 组件               |
| -------- | ----------------------------- |
| 普通弹窗、模态框 | `Dialog`                      |
| 危险操作确认框  | `Alert Dialog`                |
| 右侧滑出面板   | `Sheet`                       |
| 桌面端侧边抽屉  | `Sheet`                       |
| 手机端底部弹窗  | `Drawer`                      |
| 下拉菜单     | `Dropdown Menu`               |
| 右键菜单     | `Context Menu`                |
| 标签页      | `Tabs`                        |
| 下拉选择框    | `Select`                      |
| 可搜索选择框   | `Combobox` 组合                 |
| 命令面板     | `Command`                     |
| 悬浮提示     | `Tooltip`                     |
| 点击触发浮层   | `Popover`                     |
| 鼠标悬停卡片   | `Hover Card`                  |
| 日期选择     | `Calendar` 或 `Date Picker` 组合 |
| 日期范围选择   | `Calendar` + `Popover`        |
| 折叠内容     | `Accordion`                   |
| 显示或隐藏区域  | `Collapsible`                 |
| 开关       | `Switch`                      |
| 单选项      | `Radio Group`                 |
| 多选项      | `Checkbox`                    |
| 滑块       | `Slider`                      |
| 分段调整布局   | `Resizable`                   |
| 进度状态     | `Progress`                    |
| 骨架加载     | `Skeleton`                    |
| 空状态      | `Empty`                       |
| 通知消息     | `Sonner`                      |
| 面包屑      | `Breadcrumb`                  |
| 分页       | `Pagination`                  |
| 数据表格     | `Table` 或 `Data Table`        |
| 左侧导航栏    | `Sidebar`                     |
| 卡片       | `Card`                        |
| 用户头像     | `Avatar`                      |
| 标签、状态标识  | `Badge`                       |
| 表单字段     | `Field` 或项目现有表单封装             |
| 一次性验证码   | `Input OTP`                   |
| 抽屉式导航菜单  | `Sheet` + 导航组件                |
| 搜索命令菜单   | `Command` + `Dialog`          |
| 可滚动区域    | `Scroll Area`                 |
| 分隔线      | `Separator`                   |
| 固定宽高比容器  | `Aspect Ratio`                |
| 提示信息块    | `Alert`                       |

不得将所有浮层都写成 `Dialog`。

不得将所有侧滑面板都写成 `Drawer`。

必须根据设备、方向和交互语义选择正确组件。

---

# 4. shadcn/ui 组件发现协议

## 4.1 本地存在组件时

如果组件已存在于：

```text
@/components/ui/
```

必须直接导入使用：

```tsx
import { Button } from "@/components/ui/button"
```

禁止：

- 用原生 `<button>` 重写 Button 的视觉和交互

- 重新复制一个 `CustomButton`

- 绕过本地组件直接从 Radix 导入

- 因为需要调整样式就另写一套组件

- 从网上复制另一版本覆盖项目现有版本

允许通过以下方式调整：

```tsx
<Button variant="outline" size="sm" className="...">
  保存
</Button>
```

如需项目级统一改造，应修改：

```text
@/components/ui/button.tsx
```

而不是在大量页面中重复覆盖相同样式。

---

## 4.2 本地不存在组件时

如果本地不存在对应组件，必须先检查 shadcn/ui 官方组件库。

优先使用 Bun 执行：

```bash
bunx shadcn@latest add <component-name>
```

例如：

```bash
bunx shadcn@latest add dialog
bunx shadcn@latest add sheet
bunx shadcn@latest add drawer
bunx shadcn@latest add select
bunx shadcn@latest add command
```

需要查看官方使用方式时：

```bash
bunx shadcn@latest docs <component-name>
```

需要检查当前 shadcn 项目状态时：

```bash
bunx shadcn@latest info
```

AI 必须明确提示：

```text
我检查了 @/components/ui/，当前缺少 Sheet 组件。

为了保持 shadcn/ui 组件规范，请先运行：

bunx shadcn@latest add sheet
```

不得在缺失组件的情况下用原生元素临时伪造同类组件。

---

## 4.3 缺失组件时的代码生成边界

当关键组件缺失时：

- 不得生成伪造版本。

- 不得假装组件已经安装。

- 不得编写错误 import。

- 不得引入未经允许的第三方组件库。

- 可以继续完成不依赖该组件的数据模型、Schema、Query 和静态内容。

- 涉及缺失组件的 UI 部分必须明确标记为等待安装。

- 安装后再基于本地源码生成最终代码。

禁止生成：

```tsx
import { Sheet } from "@/components/ui/sheet"
```

但项目中实际上没有：

```text
@/components/ui/sheet.tsx
```

---

## 4.4 shadcn 官方没有直接对应组件时

必须按以下顺序处理：

1. 搜索 shadcn 官方组件。

2. 搜索 shadcn 官方 Blocks。

3. 检查是否可由已有 shadcn 原子组件组合。

4. 检查项目中是否已有类似业务组件。

5. 检查允许使用的 shadcn 生态 Registry。

6. 最后才编写项目专用组件。

专用组件仍应尽量组合已有原子组件，例如：

```tsx
<Card>
  <CardHeader />
  <CardContent>
    <Tabs />
  </CardContent>
</Card>
```

而不是重新实现 Card 和 Tabs。

---

# 5. 禁止自行重写的组件类型

以下组件只要 shadcn 官方存在，就必须优先采用 shadcn 版本或基于其组合：

- Button

- Input

- Textarea

- Checkbox

- Radio Group

- Switch

- Slider

- Select

- Combobox

- Dialog

- Alert Dialog

- Sheet

- Drawer

- Popover

- Tooltip

- Hover Card

- Dropdown Menu

- Context Menu

- Menubar

- Navigation Menu

- Tabs

- Accordion

- Collapsible

- Command

- Calendar

- Date Picker

- Card

- Badge

- Avatar

- Alert

- Empty

- Skeleton

- Progress

- Separator

- Breadcrumb

- Pagination

- Sidebar

- Table

- Data Table

- Scroll Area

- Resizable

- Carousel

- Input OTP

- Sonner

- Toggle

- Toggle Group

禁止使用原生 HTML 和 Tailwind CSS 重新模拟这些组件的完整交互。

例如，以下代码不允许作为 Dialog 的替代：

```tsx
{open && (
  <div className="fixed inset-0 z-50 bg-black/50">
    <div className="rounded-xl bg-white p-6">
      ...
    </div>
  </div>
)}
```

原因包括但不限于：

- 缺少焦点锁定

- 缺少键盘关闭行为

- 缺少正确的可访问性语义

- 缺少焦点恢复

- 缺少统一动画

- 缺少遮罩层管理

- 容易产生层级冲突

正确做法是使用 `Dialog`、`AlertDialog`、`Sheet` 或 `Drawer`。

---

# 6. 允许自行编写的内容

“禁止乱写组件”不等于“禁止写 JSX”。

以下内容允许自行实现：

## 6.1 普通页面布局

例如：

```tsx
<div className="grid gap-6 lg:grid-cols-[240px_minmax(0,1fr)]">
  <aside>{/* ... */}</aside>
  <main>{/* ... */}</main>
</div>
```

普通的：

- `div`

- `section`

- `main`

- `header`

- `footer`

- `nav`

- `article`

- CSS Grid

- Flexbox

可以正常使用。

但这些标签不得被用于重新实现已有的交互组件。

---

## 6.2 业务展示组件

例如：

- 用户统计卡

- 学习进度区块

- 订单摘要

- 单词记忆状态

- 项目专属图表容器

- 业务时间线

- 课程信息卡

- 任务完成面板

这些组件应优先组合：

- Card

- Badge

- Avatar

- Progress

- Tooltip

- Button

- Separator

---

## 6.3 项目专属复合组件

例如：

```text
UserFilterSheet
OrderStatusCard
VocabularyReviewPanel
DailyTaskSummary
ExamProgressChart
```

复合组件可以自行编写，但其基础交互必须建立在现有 shadcn 原子组件之上。

---

# 7. 禁止引入其他 UI 体系

未经用户明确批准，禁止引入：

- Ant Design

- Material UI

- Chakra UI

- Mantine

- Arco Design

- Element Plus

- PrimeReact

- NextUI / HeroUI

- Bootstrap

- Semantic UI

- 其他完整 UI 组件库

同样禁止为了一个简单动效直接引入大型动画或 UI 依赖。

新增依赖前必须确认：

1. shadcn 是否已有对应能力。

2. 项目现有依赖是否已经能够实现。

3. 是否可以通过现有组件组合完成。

4. 新依赖是否会造成重复能力。

5. 新依赖是否会明显增加包体积。

6. 是否与 SSR 兼容。

7. 是否会引入另一套设计语言。

---

# 8. shadcn 生态扩展规范

当用户明确要求：

- 高完成度

- 高级感

- 视觉冲击力

- 高端动效

- 营销页面

- 精致 Dashboard

- 复杂交互区块

可以参考：

- shadcn/ui Blocks

- Shadcnblocks

- Magic UI

- Aceternity UI

但必须遵守以下规则：

1. 优先查看项目当前允许的 Registry。

2. 优先安装源码级组件，而不是黑盒依赖。

3. 安装后检查源代码和依赖。

4. 必须适配项目现有设计令牌。

5. 不得引入与现有 shadcn 体系冲突的基础组件。

6. 不得直接照搬演示站中的品牌色、背景和文案。

7. 不得为了一个简单效果引入大量依赖。

8. 必须考虑 SSR、移动端和无障碍支持。

9. 动效必须支持 `prefers-reduced-motion`。

10. 不得让装饰动效破坏核心操作。

---

# 9. Tailwind CSS 使用规范

## 9.1 使用原则

优先使用：

- 标准 Tailwind 工具类

- 项目 CSS 变量

- shadcn 语义化颜色

- 响应式断点

- `cn()` 合并类名

推荐：

```tsx
<div
  className={cn(
    "flex items-center gap-3 rounded-lg border bg-card p-4 text-card-foreground",
    selected && "border-primary",
    className,
  )}
/>
```

---

## 9.2 语义化颜色

优先使用：

```text
bg-background
bg-card
bg-muted
bg-primary
bg-secondary
bg-accent
bg-destructive

text-foreground
text-muted-foreground
text-primary-foreground
text-destructive

border-border
border-input
ring-ring
```

不得无理由硬编码：

```text
bg-[#ffffff]
text-[#111111]
border-[#e5e5e5]
```

只有品牌色、数据可视化色或设计稿明确要求时，才允许使用额外颜色。

---

## 9.3 避免重复任意值

禁止在多个文件中重复：

```text
rounded-[13px]
shadow-[0_8px_30px_rgba(...)]
w-[317px]
text-[15px]
```

如果相同值多次出现，应考虑：

- Tailwind 设计令牌

- CSS 变量

- 组件 variant

- 公共业务组件

- 全局主题配置

---

## 9.4 条件类名

统一使用项目现有 `cn()` 工具：

```tsx
import { cn } from "@/lib/utils"
```

禁止大量使用字符串拼接：

```tsx
className={"base " + (active ? "active" : "")}
```

---

## 9.5 响应式设计

默认采用移动优先：

```tsx
<div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
```

每个页面至少检查：

- 小屏是否横向溢出

- Sheet 和 Dialog 在手机端是否合适

- 表格是否需要 ScrollArea

- 操作按钮是否方便触摸

- 文本是否发生不可控截断

- 固定宽度是否导致布局破坏

---

# 10. 图标规范

默认使用项目已有的 **Lucide React**。

例如：

```tsx
import { Plus, Search, Settings } from "lucide-react"
```

禁止：

- 使用 Emoji 充当正式操作图标

- 混用多套线性图标库

- 手写常见 SVG 图标

- 为常见图标引入新的图标依赖

- 在同一按钮里混用不同视觉风格图标

纯图标按钮必须提供可访问名称：

```tsx
<Button variant="ghost" size="icon" aria-label="打开设置">
  <Settings aria-hidden="true" />
</Button>
```

必要时配合 Tooltip。

---

# 11. TanStack Router 规范

## 11.1 路由优先

所有站内页面跳转必须优先使用 TanStack Router：

```tsx
import { Link } from "@tanstack/react-router"

<Link to="/users/$userId" params={{ userId }}>
  查看用户
</Link>
```

命令式导航使用：

```tsx
const navigate = useNavigate()

navigate({
  to: "/users/$userId",
  params: { userId },
})
```

---

## 11.2 路由参数必须类型安全

不得手动解析 URL：

```tsx
const id = window.location.pathname.split("/").pop()
```

应使用：

```tsx
const { userId } = Route.useParams()
```

搜索参数应通过路由 Schema 校验：

```tsx
const searchSchema = z.object({
  page: z.coerce.number().int().positive().catch(1),
  keyword: z.string().catch(""),
})
```

---

## 11.3 页面状态优先进入 URL

以下状态优先放在 URL Search Params 中：

- 页码

- 排序

- 搜索关键词

- 筛选条件

- 当前标签页

- 可分享的视图模式

- 数据范围

以下状态通常不需要进入 URL：

- 临时 hover 状态

- 未提交的简单弹窗状态

- 一次性动画状态

- 局部输入过程状态

---

## 11.4 路由文件职责

路由文件主要负责：

- 路由定义

- 参数校验

- loader

- 页面级数据预取

- 权限守卫

- 页面组件挂载

复杂业务 UI 应拆分到：

```text
@/features/<feature>/
```

不得把几百行页面逻辑全部堆在路由文件中。

---

# 12. TanStack Query 规范

## 12.1 服务端状态必须使用 Query

以下数据属于服务端状态：

- 用户列表

- 用户详情

- 订单数据

- 后端配置

- 分页结果

- 权限数据

- 报表数据

- 字典数据

- API 返回的业务实体

这些数据必须优先由 TanStack Query 管理。

禁止：

```tsx
useEffect(() => {
  fetch("/api/users")
    .then(...)
    .then(setUsers)
}, [])
```

禁止把 API 返回数据长期复制进 Zustand。

---

## 12.2 Query Key 必须稳定

推荐集中管理：

```ts
export const userKeys = {
  all: ["users"] as const,
  lists: () => [...userKeys.all, "list"] as const,
  list: (params: UserListParams) =>
    [...userKeys.lists(), params] as const,
  details: () => [...userKeys.all, "detail"] as const,
  detail: (id: string) =>
    [...userKeys.details(), id] as const,
}
```

不得在不同页面随意创建含义不一致的 Query Key。

---

## 12.3 使用 Query Options 工厂

推荐：

```ts
export const userListQueryOptions = (params: UserListParams) =>
  queryOptions({
    queryKey: userKeys.list(params),
    queryFn: () => getUsers({ query: params }),
  })
```

路由 loader、组件预取和页面查询应尽量复用同一套 options。

---

## 12.4 Mutation 后正确维护缓存

Mutation 成功后，根据业务情况选择：

- `invalidateQueries`

- `setQueryData`

- 乐观更新

- 路由重新验证

不得在修改成功后简单执行：

```ts
window.location.reload()
```

---

## 12.5 必须处理完整状态

所有主要数据区域必须考虑：

- Pending

- Error

- Empty

- Success

- Refetching

- Pagination loading

优先使用：

- `Skeleton`

- `Alert`

- `Empty`

- `Button`

- `Sonner`

不得只处理成功状态。

---

# 13. Zustand 使用规范

Zustand 仅用于真正的客户端共享状态，例如：

- 全局侧栏展开状态

- 跨页面编辑器状态

- 本地工作区状态

- 播放器状态

- 未持久化的多步骤流程

- 需要跨层级共享的界面偏好

不得用于保存：

- 用户列表

- 订单详情

- 后端分页结果

- 查询缓存

- API loading 状态

- 可以由 URL 表达的筛选条件

- React Hook Form 已管理的表单字段

Store 应按领域拆分：

```text
stores/
  app-store.ts
  player-store.ts
  workspace-store.ts
```

避免创建无边界的：

```text
global-store.ts
```

选择状态时应使用最小选择器：

```tsx
const isSidebarOpen = useAppStore((state) => state.isSidebarOpen)
```

不要直接订阅整个 Store：

```tsx
const store = useAppStore()
```

---

# 14. React Hook Form 与 Zod 规范

## 14.1 表单必须使用 React Hook Form

以下表单优先使用 React Hook Form：

- 登录

- 注册

- 创建

- 编辑

- 搜索条件较多的筛选表单

- 多步骤表单

- 动态字段表单

- 带复杂校验的业务表单

简单的单字段即时搜索可以使用局部状态，但不得把复杂表单拆成大量 `useState`。

---

## 14.2 Zod 是校验真源

表单校验统一由 Zod Schema 定义：

```ts
export const userFormSchema = z.object({
  name: z.string().trim().min(1, "请输入姓名"),
  email: z.email("请输入有效邮箱"),
  role: z.enum(["admin", "member"]),
})

export type UserFormValues = z.infer<typeof userFormSchema>
```

不要同时手写一份重复 TypeScript 类型：

```ts
interface UserFormValues {
  name: string
  email: string
  role: string
}
```

除非生成类型和表单输入类型确实存在差异。

---

## 14.3 表单必须显示错误信息

每个字段至少需要考虑：

- Label

- Control

- Description

- Error message

- Disabled

- Required

- Loading

必须使用项目现有 shadcn 表单方案。

不得仅通过红色边框表达错误。

---

## 14.4 提交状态

提交按钮必须处理：

- 重复提交

- Loading

- 成功反馈

- 服务端错误

- 字段级错误

- 表单级错误

示例：

```tsx
<Button type="submit" disabled={form.formState.isSubmitting}>
  {form.formState.isSubmitting ? "保存中…" : "保存"}
</Button>
```

不得只改变按钮文案而不禁用重复提交。

---

## 14.5 前后端 Schema 边界

需要区分：

- API 响应 Schema

- 表单输入 Schema

- 路由 Search Schema

- 环境变量 Schema

- 本地持久化 Schema

不得将一个后端响应 Schema 强行用于所有场景。

---

# 15. Hey API OpenAPI Client 规范

## 15.1 生成代码是 API 类型真源

API 客户端、请求类型和响应类型应优先由 OpenAPI 生成。

禁止手写重复的：

- API Response 类型

- Request Body 类型

- Query 参数类型

- Path 参数类型

- Endpoint URL

- 重复 fetch 封装

例如已经生成：

```ts
getUsers
createUser
UserDto
CreateUserRequest
```

就不得再创建：

```ts
fetchUsers
UserResponse
CreateUserPayload
```

除非是在生成类型之上建立明确的业务适配层。

---

## 15.2 生成目录不得手工修改

建议将生成代码放在：

```text
@/lib/api/generated/
```

或：

```text
@/api/generated/
```

生成目录必须标明：

```text
DO NOT EDIT
```

禁止直接修改生成文件。

若生成结果不符合需求，应调整：

- OpenAPI 文档

- Hey API 配置

- 生成插件

- 业务适配层

---

## 15.3 API 调用层次

推荐结构：

```text
OpenAPI generated client
        ↓
Query options / mutation options
        ↓
Feature hooks
        ↓
Page or business component
```

例如：

```text
@/lib/api/generated/
@/features/users/api/user-queries.ts
@/features/users/hooks/use-users.ts
@/features/users/components/user-table.tsx
```

不得在展示组件内部散落大量底层请求配置。

---

## 15.4 错误处理

API 错误必须统一处理：

- 网络错误

- 未授权

- 权限不足

- 参数错误

- 业务错误

- 服务端异常

不得在每个页面中复制不同版本的错误解析逻辑。

统一错误解析可以放在：

```text
@/lib/api/error.ts
```

但不得丢失后端返回的字段级错误信息。

---

# 16. React 组件设计规范

## 16.1 默认使用函数组件

```tsx
export function UserCard() {
  return <Card />
}
```

不使用 class component。

---

## 16.2 优先组合，避免万能组件

不要创建包含大量布尔 Props 的万能组件：

```tsx
<Panel
  isUser
  isCompact
  showHeader
  showFooter
  isEditable
  useDialog
  useSheet
/>
```

应拆分为明确的业务组件或使用组合模式。

---

## 16.3 控制组件职责

一个组件应尽量只承担一个清晰职责。

以下情况应考虑拆分：

- 文件超过约 250 至 350 行

- 存在多个独立业务区块

- 状态和展示逻辑严重混合

- 表单、列表、弹窗全部堆在同一文件

- 多处出现相同 JSX

- 一个组件有大量无关 Props

行数不是绝对标准，职责边界优先。

---

## 16.4 不要过度抽象

仅使用一次、逻辑简单的 JSX 不需要强行抽成组件。

不得为了“看起来工程化”创建：

```text
PageWrapper
PageInner
PageContentWrapper
PageContentInner
PageSectionWrapper
```

抽象必须带来至少一种实际价值：

- 复用

- 一致性

- 可测试性

- 可维护性

- 明确语义

- 隔离复杂逻辑

---

## 16.5 服务端与客户端边界

在 TanStack Start 项目中：

- 服务端数据获取优先放在 loader 或服务端函数。

- 不应把密钥和服务端环境变量暴露给浏览器。

- 浏览器组件不得直接引用仅服务端模块。

- SSR 阶段不得直接访问 `window`、`document`、`localStorage`。

- 客户端专属逻辑必须进行环境判断或放在合适生命周期中。

---

# 17. 页面与目录结构规范

推荐采用 feature-first 结构：

```text
src/
├── components/
│   ├── ui/
│   └── shared/
├── features/
│   ├── users/
│   │   ├── api/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── schemas/
│   │   ├── stores/
│   │   └── utils/
│   └── orders/
├── lib/
│   ├── api/
│   ├── query/
│   ├── utils/
│   └── validation/
├── routes/
├── router.tsx
└── styles/
```

规则：

- 通用 UI 原子组件放 `components/ui`。

- 跨业务通用组件放 `components/shared`。

- 业务组件放 `features/<feature>`。

- 路由文件负责路由边界，不承担全部业务实现。

- 通用工具放 `lib`。

- 不得把所有内容都堆进 `components`。

- 不得创建含义不清的 `common`、`misc`、`helpers2`。

---

# 18. 加载、错误、空状态规范

每个异步页面或数据模块必须设计完整状态。

## 18.1 加载状态

优先使用与最终布局相似的 `Skeleton`。

禁止所有页面都只显示：

```tsx
<div>Loading...</div>
```

---

## 18.2 空状态

优先使用 shadcn 的 `Empty` 组件或项目统一空状态组件。

空状态应根据场景说明：

- 为什么为空

- 用户可以做什么

- 是否提供主要操作

- 是否可以清除筛选条件

---

## 18.3 错误状态

错误区域优先使用：

- `Alert`

- 重试 `Button`

- 统一错误文案

- 必要的错误详情

不得把技术堆栈、请求头、Token 或敏感信息直接暴露给用户。

---

## 18.4 局部刷新状态

后台重新获取数据时，不应无条件清空已有内容。

应合理显示：

- 小型加载指示

- 禁用相关操作

- 保留旧数据

- 刷新提示

避免页面闪烁。

---

# 19. 可访问性规范

所有交互组件必须满足基本可访问性要求。

必须遵守：

- 表单控件有对应 Label。

- 图标按钮有 `aria-label`。

- 图片有合理 `alt`。

- 装饰图标使用 `aria-hidden="true"`。

- 不只通过颜色传达状态。

- 键盘可以完成主要操作。

- 使用正确的语义化 HTML。

- 不破坏 shadcn/Radix/Base UI 内置焦点行为。

- Dialog 必须有标题。

- 必要时提供描述。

- 加载和动态状态应考虑屏幕阅读器反馈。

- 动效支持 reduced motion。

不得为了视觉效果移除可见焦点样式。

---

# 20. 动效规范

动效应服务于：

- 状态反馈

- 层级变化

- 元素进入或退出

- 操作结果

- 页面空间关系

禁止：

- 所有元素都使用大幅弹跳

- 无意义地持续旋转或漂浮

- 动画阻塞点击

- 页面加载时出现长时间开场动画

- 使用动画掩盖性能问题

- 为简单淡入引入大型动画库

优先使用：

- Tailwind transition

- shadcn 内置动画

- CSS keyframes

- 项目已有动画工具

只有复杂编排、手势或物理动画确实需要时，才考虑额外动画方案。

---

# 21. 性能规范

必须注意：

- 避免不必要的全局状态订阅。

- 避免把巨大对象放入 Query Key。

- 列表元素必须有稳定 key。

- 大列表考虑虚拟化。

- 图片提供合理尺寸。

- 避免无意义的 `useMemo` 和 `useCallback`。

- 避免重复请求。

- 避免在 render 中执行昂贵计算。

- 避免把整个大型页面强制变成客户端组件。

- 合理使用路由级代码拆分。

- 不要为了微小复用制造大量包装组件。

不得为了所谓“优化”提前写复杂缓存逻辑。

应先保证正确性，再针对实际瓶颈优化。

---

# 22. TypeScript 规范

必须：

- 开启严格类型检查。

- 避免 `any`。

- 优先使用推导。

- 为公共组件 Props 定义清晰类型。

- 使用判别联合表达复杂状态。

- 使用生成的 API 类型。

- 使用 Zod 推导表单与校验类型。

- 正确处理 `undefined` 和 `null`。

禁止用断言掩盖类型错误：

```ts
value as any
value as unknown as User
```

除非有明确、可解释且无法避免的边界。

---

# 23. AI 修改现有代码时的规则

AI 修改项目时必须：

1. 先读目标文件。

2. 读取关联组件。

3. 检查现有 import。

4. 检查本地 shadcn 组件。

5. 搜索相似实现。

6. 保持现有命名和格式风格。

7. 尽量做局部修改。

8. 不得无关重构。

9. 不得擅自升级依赖。

10. 不得擅自替换技术栈。

11. 不得删除未知用途代码。

12. 不得覆盖本地定制的 shadcn 组件。

13. 不得假设某个依赖已安装。

14. 不得输出无法编译的伪代码冒充成品。

当上下文不完整时，应基于已知项目结构做最保守、最兼容的实现。

---

# 24. AI 输出代码前的强制决策流程

AI 必须在内部依次回答：

```text
1. 项目中是否已有这个业务组件？
2. @/components/ui 中是否已有对应 shadcn 组件？
3. 本地组件实际导出了什么？
4. shadcn 官方是否提供对应组件？
5. 是否存在官方 Block 或推荐组合？
6. 是否需要安装组件？
7. 是否能用已有原子组件组合？
8. 这个状态属于 URL、Query、Zustand、Form 还是局部状态？
9. API 类型是否已由 Hey API 生成？
10. 页面是否处理 loading、error、empty、success？
11. 是否满足移动端与键盘操作？
12. 是否引入了不必要的新依赖？
```

任何一项未确认，不得直接凭记忆输出大段代码。

---

# 25. 组件缺失时的标准回复模板

```text
我检查了项目的 @/components/ui/ 目录，目前没有发现 Sheet 组件。

该需求属于右侧滑出面板，应使用 shadcn/ui 的 Sheet，而不是用原生 div 和 useState 自行模拟。

请在包含 components.json 的前端工作区运行：

bunx shadcn@latest add sheet

安装完成后，应按项目本地生成的 Sheet、SheetContent、SheetHeader、
SheetTitle、SheetDescription 等实际导出进行组合。
```

若缺少多个组件，应一次性列出：

```bash
bunx shadcn@latest add sheet select checkbox
```

但在给出批量安装命令前，必须确认这些组件本地确实都不存在。

---

# 26. 明确禁止的错误做法

## 26.1 未检查组件就直接手写

```tsx
const [open, setOpen] = useState(false)

{open && <div className="fixed inset-0">...</div>}
```

在 Dialog、Sheet、Drawer 已存在或官方可安装时禁止。

---

## 26.2 伪造不存在的 import

```tsx
import { FancySelect } from "@/components/ui/fancy-select"
```

禁止假设项目中存在该组件。

---

## 26.3 猜测 Props

```tsx
<Select searchable clearable animation="smooth" />
```

若本地组件没有这些 Props，则禁止使用。

---

## 26.4 使用 Zustand 保存服务端列表

```ts
const useUserStore = create(() => ({
  users: [],
  loading: false,
}))
```

用户列表应优先使用 TanStack Query。

---

## 26.5 手写重复 API 类型

```ts
interface UserResponse {
  id: number
  name: string
}
```

若 Hey API 已生成对应类型，则禁止重复定义。

---

## 26.6 使用 useEffect 发起普通页面查询

```tsx
useEffect(() => {
  loadUsers()
}, [])
```

应优先使用 TanStack Query、Route loader 或预取机制。

---

## 26.7 使用 div 冒充按钮

```tsx
<div onClick={handleSave}>保存</div>
```

必须使用 `Button` 或语义化 `button`。  
若项目已安装 shadcn Button，则优先使用 shadcn Button。

---

## 26.8 用颜色模拟所有状态

```tsx
<span className="text-green-500">正常</span>
```

应结合：

- Badge

- 图标

- 文本

- 可访问语义

---

## 26.9 大量硬编码设计值

```tsx
<div className="rounded-[17px] bg-[#F7F5FF] px-[19px] py-[13px]">
```

应先检查项目令牌和现有组件风格。

---

## 26.10 无理由添加依赖

禁止仅为了：

- 一个 debounce

- 一个日期格式化

- 一个 className 合并

- 一个简单动画

- 一个浅拷贝

就直接增加新包。

必须先检查项目现有能力。

---

# 27. 交付前强制自检

生成或修改代码后，必须检查：

## 组件

- 是否先检查了现有业务组件？

- 是否扫描了 `@/components/ui/`？

- 是否优先使用了 shadcn 组件？

- 是否存在手写 Dialog、Sheet、Select 等行为？

- 是否引入了不存在的组件？

- 是否猜测了组件 Props？

## 状态

- 服务端数据是否交给 TanStack Query？

- URL 状态是否进入 TanStack Router Search Params？

- 表单是否交给 React Hook Form？

- 校验是否由 Zod 管理？

- Zustand 是否只保存客户端共享状态？

- 是否存在重复状态源？

## API

- 是否使用 Hey API 生成客户端？

- 是否重复定义生成类型？

- 是否直接修改了生成文件？

- Query Key 是否稳定？

- Mutation 后是否正确更新缓存？

## UI

- 是否处理加载状态？

- 是否处理错误状态？

- 是否处理空状态？

- 是否处理禁用和提交状态？

- 是否适配移动端？

- 是否存在横向溢出？

- 图标按钮是否有可访问名称？

- Dialog 是否具有标题？

## 工程

- 是否使用 Bun？

- 是否引入不必要依赖？

- 是否保持现有目录结构？

- 是否出现 `any`？

- 是否存在无关重构？

- 是否通过 typecheck？

- 是否通过 lint？

- 是否能够构建？

建议执行：

```bash
bun run typecheck
bun run lint
bun run build
```

---

# 28. 最终执行准则

AI 必须始终遵守以下规则：

> **不要先写，再看有没有组件。必须先查，再决定怎么写。**

> **只要 shadcn/ui 已有对应组件，就不得用原生 HTML 和 Tailwind 重新实现它的完整交互。**

> **只要本地已经安装组件，就必须使用本地版本，而不是凭记忆生成另一套版本。**

> **组件不存在时，先提供准确安装命令；不得伪造文件、导出或 Props。**

> **TanStack Query 管服务端状态，Zustand 管客户端共享状态，React Hook Form 管表单状态，TanStack Router 管路由与 URL 状态。**

> **Hey API 生成的类型和客户端是接口层真源，不得重复手写。**

> **允许编写业务组件，但业务组件必须优先由现有 shadcn 原子组件组合而成。**

> **不允许为了赶进度而写一个“临时版”通用组件，因为临时版通常会成为永久技术债。**
