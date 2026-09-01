'use client'

import * as React from 'react'

import { SidebarProvider } from '@/components/ui/sidebar'
import { AppHeader } from './AppHeader'
import { AppSidebar } from './AppSidebar'
import { RightPanel } from './DetailSidebar'

export interface AdminShellProps {
    /** Controlled right-panel state. If omitted, the shell manages it internally. */
    rightPanelOpen?: boolean
    onRightPanelChange?: (open: boolean) => void
    /** Page content rendered inside the elevated main panel. */
    children?: React.ReactNode
    /** When true, the right panel is shown on initial render. */
    defaultRightPanelOpen?: boolean
}

/**
 * AdminShell — three-zone layout:
 *
 *  1. Zero level: <AppSidebar /> — blends with the page background.
 *  2. Elevated main panel: a rounded, bordered, shadowed card that holds the
 *     header + page content. This is the "app area" the user sees most of the time.
 *  3. Right sticky panel: an optional rounded card that appears alongside the
 *     main panel when needed (details, activity, inspector, etc.). Hidden by
 *     default; toggled from the header.
 */
export function AdminShell({
    rightPanelOpen,
    onRightPanelChange,
    children,
    defaultRightPanelOpen = false,
}: AdminShellProps) {
    const [internalOpen, setInternalOpen] = React.useState(defaultRightPanelOpen)
    const open = rightPanelOpen ?? internalOpen
    const setOpen = onRightPanelChange ?? setInternalOpen

    return (
        <SidebarProvider
            style={{
                // Left sidebar: cap width at 250px (15.625rem) — overrides the
                // default 16rem (256px) baked into the shadcn sidebar primitive.
                '--sidebar-width': '15rem',
            } as React.CSSProperties}
        >
            <AppSidebar />

            {/*
        Outer canvas — fixed to the viewport height so the *page* itself never
        scrolls. Only the <main> inside the elevated card scrolls. This keeps
        the header pinned at the top of the card naturally (no sticky hacks).
      */}
            <div className="relative flex h-svh flex-1 flex-col overflow-hidden">
                {/* Inner padded frame that holds the elevated main panel + right panel */}
                <div className="flex h-full min-h-0 flex-1 gap-4 p-3 md:p-4">
                    {/* ── Elevated main panel ─────────────────────────────────────── */}
                    {/* Elevated main panel — border only, no shadow (per request) */}
                    <div className="flex min-w-0 flex-1 flex-col overflow-hidden rounded-2xl border bg-card">
                        {/* Header is a flex child of the card (not inside the scrolling
                main), so it stays pinned at the top while <main> scrolls. */}
                        <AppHeader
                            rightPanelOpen={open}
                            onToggleRightPanel={() => setOpen(!open)}
                        />

                        <main className="min-h-0 flex-1 overflow-y-auto p-4 md:p-6">
                            {children ?? <div>test</div>}
                        </main>
                    </div>

                    {/* ── Optional right panel ──────────────────────────────────────
              Fills the full height of the inner frame; its own internal
              activity list scrolls. */}
                    <RightPanel open={open} onClose={() => setOpen(false)} />
                </div>
            </div>
        </SidebarProvider>
    )
}
