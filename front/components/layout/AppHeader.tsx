'use client'

import { Bell, PanelRight, PanelRightClose, Search } from 'lucide-react'

import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { SidebarTrigger } from '@/components/ui/sidebar'

export interface AppHeaderProps {
    rightPanelOpen: boolean
    onToggleRightPanel: () => void
}

export function AppHeader({
    rightPanelOpen,
    onToggleRightPanel,
}: AppHeaderProps) {
    return (
        <header className="flex h-11 shrink-0 items-center gap-2 border-b bg-card/80 px-3 backdrop-blur supports-[backdrop-filter]:bg-card/60 md:px-4">
            <SidebarTrigger className="-ml-1" />
            {/* <Separator orientation="vertical" className="mr-1 h-4" /> */}

            <Breadcrumb>
                <BreadcrumbList>
                    <BreadcrumbItem className="hidden md:block">
                        <span className="text-muted-foreground">Console</span>
                    </BreadcrumbItem>
                    <BreadcrumbSeparator className="hidden md:block" />
                    <BreadcrumbItem>
                        <BreadcrumbPage>Dashboard</BreadcrumbPage>
                    </BreadcrumbItem>
                </BreadcrumbList>
            </Breadcrumb>

            <div className="ml-auto flex items-center gap-1.5">
                <div className="relative hidden lg:block">
                    <Search className="pointer-events-none absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
                    <Input
                        placeholder="Search…"
                        className="h-7 w-56 rounded-lg pl-8 text-sm"
                    />
                </div>

                <Button
                    variant="ghost"
                    size="icon"
                    className="size-7 rounded-lg"
                    aria-label="Notifications"
                >
                    <Bell className="size-4" />
                </Button>

                {/* <ThemeToggle /> */}

                {/* <Separator orientation="vertical" className="mx-1 h-4" /> */}

                <Button
                    variant={rightPanelOpen ? 'secondary' : 'ghost'}
                    size="icon"
                    className="size-7 rounded-lg"
                    onClick={onToggleRightPanel}
                    aria-label={rightPanelOpen ? 'Close details panel' : 'Open details panel'}
                    aria-pressed={rightPanelOpen}
                >
                    {rightPanelOpen ? (
                        <PanelRightClose className="size-4" />
                    ) : (
                        <PanelRight className="size-4" />
                    )}
                </Button>

                <Avatar className="size-7 rounded-lg md:hidden">
                    <AvatarFallback className="rounded-lg text-xs">MK</AvatarFallback>
                </Avatar>
            </div>
        </header>
    )
}
