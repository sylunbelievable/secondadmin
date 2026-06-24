import { useEffect, useState, type ReactNode, type RefObject } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate, useRouterState } from '@tanstack/react-router'
import { KeepAlive, useKeepAliveRef, type KeepAliveRef } from 'keepalive-for-react'
import {
  BookOpen,
  Braces,
  ChevronRight,
  FolderTree,
  LayoutDashboard,
  ListTree,
  LogOut,
  MonitorSmartphone,
  Moon,
  PanelLeftClose,
  PanelLeftOpen,
  ScrollText,
  ShieldCheck,
  Sun,
  Users,
  X,
  type LucideIcon,
} from 'lucide-react'
import { logout, responseData, type Menu } from '#/api'
import { Button } from '#/components/ui/button'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '#/components/ui/collapsible'
import { Separator } from '#/components/ui/separator'
import { Sidebar, SidebarBody, SidebarLabel, useSidebar } from '#/components/ui/sidebar'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '#/components/ui/tooltip'
import { authQueries } from '#/lib/auth'
import { cn } from '#/lib/utils'
import { useUIStore } from '#/store/ui'

const knownRoutes: Record<string, string> = {
  dashboard: '/dashboard',
  users: '/users',
  roles: '/roles',
  apis: '/apis',
  menus: '/menus',
  dictionaries: '/dictionaries',
  'operation-logs': '/operation-logs',
  'login-logs': '/login-logs',
  sessions: '/sessions',
}

const icons: Record<string, LucideIcon> = {
  dashboard: LayoutDashboard,
  users: Users,
  roles: ShieldCheck,
  apis: Braces,
  menus: ListTree,
  dictionaries: BookOpen,
  'operation-logs': ScrollText,
  'login-logs': ScrollText,
  sessions: MonitorSmartphone,
}

export function routeFor(item: Pick<Menu, 'component' | 'path'>) {
  return item.component ? knownRoutes[item.component] : item.path
}

function NavLink({ item, depth = 0 }: { item: Menu; depth?: number }) {
  const { open, setMobileOpen } = useSidebar()
  const openTab = useUIStore((state) => state.openTab)
  const to = routeFor(item)
  const Icon = icons[item.component ?? ''] ?? LayoutDashboard
  const content = (
    <>
      <Icon className={cn('size-5 shrink-0', depth > 0 && 'size-4')} aria-hidden="true" />
      <SidebarLabel>{to ? item.name : `${item.name}（未实现）`}</SidebarLabel>
    </>
  )
  if (!to) {
    return <span className="group/sidebar-link flex items-center gap-3 rounded-lg px-3 py-2 text-sm text-muted-foreground">{content}</span>
  }
  const link = (
    <Link
      to={to as never}
      onClick={() => {
        openTab({ path: to, title: item.name })
        setMobileOpen(false)
      }}
      className="group/sidebar-link flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
      activeProps={{ className: 'bg-sidebar-accent font-medium text-sidebar-accent-foreground' }}
      style={{ paddingLeft: open ? 12 + depth * 16 : 12 }}
    >
      {content}
    </Link>
  )
  return open ? link : (
    <Tooltip>
      <TooltipTrigger asChild>{link}</TooltipTrigger>
      <TooltipContent side="right">{item.name}</TooltipContent>
    </Tooltip>
  )
}

export function containsPath(item: Menu, pathname: string): boolean {
  return routeFor(item) === pathname || (item.children ?? []).some((child) => containsPath(child, pathname))
}

function NavItem({ item, depth, pathname }: { item: Menu; depth: number; pathname: string }) {
  const { open } = useSidebar()
  const children = (item.children ?? []).filter((child) => child.type !== 'button')
  const active = containsPath(item, pathname)
  const [expanded, setExpanded] = useState(active)

  useEffect(() => {
    if (active) setExpanded(true)
  }, [active])

  if (item.type !== 'directory' && !children.length) return <NavLink item={item} depth={depth} />

  return (
    <Collapsible open={expanded} onOpenChange={setExpanded} className="group/tree">
      <CollapsibleTrigger asChild>
        <Button
          variant="ghost"
          className="group/sidebar-link w-full justify-start gap-3 px-3 text-sidebar-foreground"
          style={{ paddingLeft: open ? 12 + depth * 16 : 12 }}
          aria-label={!open ? item.name : undefined}
          title={!open ? item.name : undefined}
        >
          <FolderTree className="size-5 shrink-0" aria-hidden="true" />
          <SidebarLabel className="flex-1 text-left">{item.name}</SidebarLabel>
          <SidebarLabel><ChevronRight className="size-4 transition-transform group-data-[state=open]/tree:rotate-90" aria-hidden="true" /></SidebarLabel>
        </Button>
      </CollapsibleTrigger>
      <CollapsibleContent className="grid gap-1">
        {children.map((child) => <NavItem key={child.id} item={child} depth={depth + 1} pathname={pathname} />)}
      </CollapsibleContent>
    </Collapsible>
  )
}

function NavItems({ items, pathname }: { items: Menu[]; pathname: string }) {
  return items.filter((item) => item.type !== 'button').map((item) => (
    <NavItem key={item.id} item={item} depth={0} pathname={pathname} />
  ))
}

function Navigation({ menus }: { menus: Menu[] }) {
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const { open, setOpen } = useSidebar()
  return (
    <SidebarBody className="justify-between">
      <div className="grid min-h-0 gap-6 overflow-y-auto overflow-x-hidden">
        <Link to="/dashboard" className="flex h-10 items-center gap-3 px-2">
          <span className="grid size-9 shrink-0 place-items-center rounded-xl bg-primary font-bold text-primary-foreground shadow-sm">S</span>
          <SidebarLabel className="font-semibold">Second Admin</SidebarLabel>
        </Link>
        <nav className="grid gap-1">
          <NavLink item={{ id: '0', parentId: '0', type: 'menu', name: '概览', component: 'dashboard', sort: 0, visible: true, status: 1 }} />
          <NavItems items={menus} pathname={pathname} />
        </nav>
      </div>
      <Button variant="ghost" className="group/sidebar-link w-full justify-start gap-3 px-3" onClick={() => setOpen(!open)}>
        {open ? <PanelLeftClose aria-hidden="true" /> : <PanelLeftOpen aria-hidden="true" />}
        <SidebarLabel>收起侧边栏</SidebarLabel>
      </Button>
    </SidebarBody>
  )
}

export function routeTitle(pathname: string, menus: Menu[]) {
  if (pathname === '/dashboard') return '概览'
  if (pathname === '/forbidden') return '暂无权限'
  const stack = [...menus]
  while (stack.length) {
    const item = stack.shift()!
    if (routeFor(item) === pathname) return item.name
    stack.push(...(item.children ?? []))
  }
  return pathname
}

function PageTabs({ menus, aliveRef }: { menus: Menu[]; aliveRef: RefObject<KeepAliveRef | null> }) {
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const navigate = useNavigate()
  const tabs = useUIStore((state) => state.tabs)
  const openTab = useUIStore((state) => state.openTab)
  const closeTab = useUIStore((state) => state.closeTab)

  useEffect(() => {
    openTab({ path: pathname, title: routeTitle(pathname, menus) })
  }, [menus, openTab, pathname])

  async function close(path: string) {
    const index = tabs.findIndex((tab) => tab.path === path)
    const next = tabs[index - 1] ?? tabs[index + 1] ?? { path: '/dashboard' }
    if (path === pathname) await navigate({ to: next.path as never })
    await aliveRef.current?.destroy(path)
    closeTab(path)
  }

  return (
    <div className="flex h-10 items-end gap-1 overflow-x-auto border-b bg-background px-2">
      {tabs.map((tab) => (
        <div key={tab.path} className={cn('flex h-8 shrink-0 items-center rounded-t-md border border-b-0 px-2 text-sm', tab.path === pathname ? 'bg-muted font-medium' : 'bg-background text-muted-foreground')}>
          <Link to={tab.path as never} className="px-1">{tab.title}</Link>
          {tab.path !== '/dashboard' && (
            <Button variant="ghost" size="icon-xs" aria-label={`关闭${tab.title}`} onClick={() => close(tab.path)}>
              <X aria-hidden="true" />
            </Button>
          )}
        </div>
      ))}
    </div>
  )
}

export function AppShell({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const theme = useUIStore((state) => state.theme)
  const toggleTheme = useUIStore((state) => state.toggleTheme)
  const sidebarExpanded = useUIStore((state) => state.sidebarExpanded)
  const setSidebarExpanded = useUIStore((state) => state.setSidebarExpanded)
  const resetTabs = useUIStore((state) => state.resetTabs)
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const aliveRef = useKeepAliveRef()
  const { data: user } = useQuery(authQueries.me())
  const { data: menus } = useQuery(authQueries.menus())

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme])

  async function signOut() {
    await responseData(logout({ throwOnError: true })).catch(() => undefined)
    queryClient.clear()
    await aliveRef.current?.destroyAll()
    resetTabs()
    await navigate({ to: '/login' })
  }

  return (
    <TooltipProvider>
      <Sidebar open={sidebarExpanded} onOpenChange={setSidebarExpanded}>
        <div className="flex h-svh w-full flex-col bg-muted/30 md:flex-row">
          <Navigation menus={menus?.menus ?? []} />
          <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
            <header className="z-10 flex h-14 shrink-0 items-center justify-end border-b bg-background/90 px-4 backdrop-blur">
              <div className="flex items-center gap-2">
                <Button variant="ghost" size="icon" aria-label="切换主题" onClick={toggleTheme}>
                  {theme === 'dark' ? <Sun aria-hidden="true" /> : <Moon aria-hidden="true" />}
                </Button>
                <span className="hidden text-sm sm:inline">{user?.nickname || user?.username}</span>
                <Button variant="ghost" size="icon" aria-label="退出" onClick={signOut}>
                  <LogOut aria-hidden="true" />
                </Button>
              </div>
            </header>
            <PageTabs menus={menus?.menus ?? []} aliveRef={aliveRef} />
            <Separator />
            <KeepAlive
              activeCacheKey={pathname}
              aliveRef={aliveRef}
              max={10}
              enableActivity
              containerClassName="min-h-0 flex-1"
              cacheNodeClassName="h-full overflow-y-auto"
            >
              <main className="p-4 md:p-6">{children}</main>
            </KeepAlive>
          </div>
        </div>
      </Sidebar>
    </TooltipProvider>
  )
}
