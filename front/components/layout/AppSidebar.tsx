'use client'

import {
    BarChart3,
    ChevronsUpDown,
    Inbox,
    LayoutDashboard,
    Search,
    Settings,
    ShoppingCart,
    Users,
    Wallet
} from 'lucide-react';

import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarInput,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem
} from '@/components/ui/sidebar';
import { useState } from 'react';

const navMain = [
    {
        title: 'Overview',
        items: [
            { title: 'Dashboard', url: '#', icon: LayoutDashboard, active: true },
            { title: 'Analytics', url: '#', icon: BarChart3 },
            { title: 'Inbox', url: '#', icon: Inbox, badge: '12' },
        ],
    },
    {
        title: 'Operations',
        items: [
            { title: 'Orders', url: '#', icon: ShoppingCart },
            { title: 'Customers', url: '#', icon: Users },
            { title: 'Billing', url: '#', icon: Wallet },
        ],
    },
    {
        title: 'System',
        items: [{ title: 'Settings', url: '#', icon: Settings }],
    },
]

export function AppSidebar() {

    let [selectedNav, setSelectedNav] = useState('Dashboard');

    const onMenuSelected = (item) => {
        if (selectedNav != item.title) {
            setSelectedNav(item.title)
        }
    };


    return (
        <Sidebar collapsible="icon" variant='inset'>
            <SidebarHeader>
                <div className="flex items-center gap-2 px-2 py-1.5">
                    <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                        <span className="text-sm font-bold">A</span>
                    </div>
                    <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                        <span className="truncate font-semibold">Acme Inc.</span>
                        <span className="truncate text-xs text-muted-foreground">
                            Enterprise
                        </span>
                    </div>
                </div>
                <form className="group-data-[collapsible=icon]:hidden">
                    <label className="sr-only" htmlFor="sidebar-search">
                        Search
                    </label>
                    <SidebarInput
                        id="sidebar-search"
                        placeholder="Search…"
                        className="bg-background border-1 border-sidebar-border"
                    />
                    <span className="sr-only">
                        <Search />
                    </span>
                </form>
            </SidebarHeader>

            <SidebarContent>
                {navMain.map((group) => (
                    <SidebarGroup key={group.title}>
                        <SidebarGroupLabel>{group.title}</SidebarGroupLabel>
                        <SidebarGroupContent>
                            <SidebarMenu>
                                {group.items.map((item) => (
                                    <SidebarMenuItem key={item.title}>
                                        <SidebarMenuButton
                                            variant={item.title == selectedNav ? "outline" : "default"}
                                            isActive={item.title == selectedNav}
                                            tooltip={item.title}
                                            onClick={e => onMenuSelected(item)}
                                        >
                                            {item.icon && <item.icon className="size-4" />}
                                            <span>{item.title}</span>
                                        </SidebarMenuButton>
                                    </SidebarMenuItem>
                                ))}
                            </SidebarMenu>
                        </SidebarGroupContent>
                    </SidebarGroup>
                ))}
            </SidebarContent>

            <SidebarFooter>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <DropdownMenu>
                            <DropdownMenuTrigger
                                render={
                                    <SidebarMenuButton
                                        size="lg"
                                        className="sidebar-btn"
                                    />
                                }
                            >
                                <Avatar className="size-8 rounded-lg">
                                    <AvatarFallback className="rounded-lg">MK</AvatarFallback>
                                </Avatar>
                                <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                                    <span className="truncate font-semibold">Mira Kapoor</span>
                                    <span className="truncate text-xs text-muted-foreground">
                                        mira@acme.io
                                    </span>
                                </div>
                                <ChevronsUpDown className="ml-auto size-4 group-data-[collapsible=icon]:hidden" />
                            </DropdownMenuTrigger>
                            <DropdownMenuContent
                                className="w-[--radix-dropdown-menu-trigger-width] min-w-52 rounded-lg"
                                side="bottom"
                                align="end"
                                sideOffset={4}
                            >
                                <DropdownMenuGroup>
                                    <DropdownMenuLabel className="p-0 font-normal" >
                                        <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                                            <Avatar className="size-8 rounded-lg">
                                                <AvatarFallback className="rounded-lg">MK</AvatarFallback>
                                            </Avatar>
                                            <div className="grid flex-1 text-left text-sm leading-tight">
                                                <span className="truncate font-semibold">
                                                    Mira Kapoor
                                                </span>
                                                <span className="truncate text-xs text-muted-foreground">
                                                    mira@acme.io
                                                </span>
                                            </div>
                                        </div>
                                    </DropdownMenuLabel>
                                </DropdownMenuGroup>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>Account</DropdownMenuItem>
                                <DropdownMenuItem>Billing</DropdownMenuItem>
                                <DropdownMenuItem>Notifications</DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem>Log out</DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarFooter>

            {/* <SidebarRail /> */}
        </Sidebar>
    )
}
