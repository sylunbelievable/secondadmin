import { createContext, useContext, useState } from "react"
import { AnimatePresence, motion, useReducedMotion } from "motion/react"
import { Menu, X } from "lucide-react"

import { Button } from "#/components/ui/button"
import { cn } from "#/lib/utils"

// Adapted from Aceternity UI's Sidebar registry for this project's tokens and accessibility.
type SidebarContextValue = {
  open: boolean
  desktopOpen: boolean
  mobileOpen: boolean
  setOpen: (open: boolean) => void
  setMobileOpen: (open: boolean) => void
  animate: boolean
}

const SidebarContext = createContext<SidebarContextValue | undefined>(undefined)

export function useSidebar() {
  const context = useContext(SidebarContext)
  if (!context) throw new Error("useSidebar must be used within Sidebar")
  return context
}

export function Sidebar({
  children,
  animate = true,
  open: controlledOpen,
  onOpenChange,
}: {
  children: React.ReactNode
  animate?: boolean
  open?: boolean
  onOpenChange?: (open: boolean) => void
}) {
  const [localOpen, setLocalOpen] = useState(true)
  const [mobileOpen, setMobileOpen] = useState(false)
  const open = controlledOpen ?? localOpen
  const setOpen = onOpenChange ?? setLocalOpen
  const reduceMotion = useReducedMotion()
  return (
    <SidebarContext.Provider value={{ open: open || mobileOpen, desktopOpen: open, mobileOpen, setOpen, setMobileOpen, animate: animate && !reduceMotion }}>
      {children}
    </SidebarContext.Provider>
  )
}

type SidebarBodyProps = Omit<React.ComponentProps<typeof motion.aside>, "children"> & {
  children: React.ReactNode
}

export function SidebarBody(props: SidebarBodyProps) {
  return (
    <>
      <DesktopSidebar {...props} />
      <MobileSidebar className={props.className}>{props.children}</MobileSidebar>
    </>
  )
}

function DesktopSidebar({
  className,
  children,
  ...props
}: SidebarBodyProps) {
  const { desktopOpen, animate } = useSidebar()
  return (
    <motion.aside
      aria-label="主导航"
      data-state={desktopOpen ? "expanded" : "collapsed"}
      className={cn(
        "hidden h-svh shrink-0 flex-col overflow-hidden border-r bg-sidebar px-3 py-4 text-sidebar-foreground data-[state=collapsed]:[&_[data-sidebar-label]]:hidden md:flex",
        className,
      )}
      animate={{ width: animate ? (desktopOpen ? 280 : 68) : 280 }}
      transition={{ duration: 0.24, ease: "easeInOut" }}
      {...props}
    >
      {children}
    </motion.aside>
  )
}

function MobileSidebar({
  className,
  children,
}: {
  className?: string
  children: React.ReactNode
}) {
  const { mobileOpen, setMobileOpen, animate } = useSidebar()
  return (
    <>
      <div className="flex h-14 w-full items-center justify-between border-b bg-background px-4 md:hidden">
        <span className="font-semibold">Second Admin</span>
        <Button variant="ghost" size="icon" aria-label="打开主导航" onClick={() => setMobileOpen(true)}>
          <Menu aria-hidden="true" />
        </Button>
      </div>
      <AnimatePresence>
        {mobileOpen && (
          <motion.aside
            aria-label="主导航"
            initial={animate ? { x: "-100%", opacity: 0 } : false}
            animate={{ x: 0, opacity: 1 }}
            exit={animate ? { x: "-100%", opacity: 0 } : undefined}
            transition={{ duration: 0.24, ease: "easeInOut" }}
            className={cn(
              "fixed inset-0 z-50 flex h-svh w-full flex-col bg-sidebar p-5 text-sidebar-foreground md:hidden",
              className,
            )}
          >
            <Button
              variant="ghost"
              size="icon"
              className="absolute right-4 top-4"
              aria-label="关闭主导航"
              onClick={() => setMobileOpen(false)}
            >
              <X aria-hidden="true" />
            </Button>
            {children}
          </motion.aside>
        )}
      </AnimatePresence>
    </>
  )
}

export function SidebarLabel({ children, className }: { children: React.ReactNode; className?: string }) {
  const { open, animate } = useSidebar()
  return (
    <motion.span
      data-sidebar-label
      aria-hidden={!open}
      animate={{ opacity: animate ? (open ? 1 : 0) : 1 }}
      className={cn(
        "min-w-0 truncate whitespace-nowrap transition-transform duration-150 group-hover/sidebar-link:translate-x-1",
        !open && "hidden",
        className,
      )}
    >
      {children}
    </motion.span>
  )
}
