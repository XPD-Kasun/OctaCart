# OctaCart Front — Agent Guide

> **Stack:** Next.js 16 (App Router) · React 19 · Tailwind CSS 4 · shadcn/ui (via `@base-ui/react`) · TypeScript · pnpm

This document is the canonical reference for any AI agent or developer working on this codebase.
It maps every meaningful entry point, explains the three-zone shell architecture, and lists the component contracts that every page must respect.

---

## 1. Repository Layout

```
front/
├── app/                        # Next.js App Router root
│   ├── layout.tsx              # ← ROOT ENTRY POINT — wraps every page
│   └── page.tsx                # ← DEFAULT PAGE — "/" route (Dashboard stub)
├── components/
│   ├── layout/                 # Shell components (add new layout primitives here)
│   │   ├── Shell.tsx           # AdminShell — three-zone layout orchestrator
│   │   ├── AppHeader.tsx       # Top bar: sidebar trigger, breadcrumb, actions
│   │   ├── AppSidebar.tsx      # Left navigation sidebar (collapsible icon rail)
│   │   └── DetailSidebar.tsx   # Right activity panel (optional, toggleable)
│   └── ui/                     # shadcn/ui primitives (generated — avoid manual edits)
├── lib/
│   └── utils.ts                # cn() Tailwind class-merge helper
├── styles/
│   └── style.css               # Global stylesheet (Tailwind directives + CSS vars)
├── hooks/                      # Custom React hooks (add shared hooks here)
├── scripts/                    # Build/tooling scripts
├── next.config.js              # Next.js config (standalone output)
├── components.json             # shadcn/ui CLI config
└── package.json
```

---

## 2. Init Entry Points

### 2.1 `app/layout.tsx` — Root Layout *(server component)*

**File:** `app/layout.tsx`

This is the **first file executed** for every request. It:

- Loads the **Inter** font via `next/font/google` and exposes it as the `--font-sans` CSS variable.
- Imports the global stylesheet (`styles/style.css`).
- Renders the `<html>` and `<body>` tags with `lang="en"`.
- Mounts `<AdminShell>` as the application-wide shell, passing `defaultRightPanelOpen={true}` and `rightPanelOpen={true}` so the right panel is visible on load.
- Injects `{children}` (the active page) into the shell's main content area.

```tsx
// app/layout.tsx — simplified
export default function Layout({ children }) {
    return (
        <html lang="en" className={cn("font-sans", ibmPlexSans.variable)}>
            <body>
                <AdminShell defaultRightPanelOpen={true} rightPanelOpen={true}>
                    {children}
                </AdminShell>
            </body>
        </html>
    )
}
```

> **Agent rule:** Never bypass `AdminShell` here. All new pages automatically receive the three-zone shell via this layout.

---

### 2.2 `app/page.tsx` — Default Route (`/`)

**File:** `app/page.tsx`

The default page rendered at the root path. Currently a placeholder (`"Dashboard Merchant"`).

```tsx
// app/page.tsx
export default function () {
    return <div className="">Dashboard Merchant</div>
}
```

> **Agent rule:** This is where the main dashboard content should be built out. It renders inside `<main>` within `AdminShell`.

---

## 3. Shell Architecture — `AdminShell`

**File:** `components/layout/Shell.tsx`

`AdminShell` is the **three-zone layout container**. All pages share this exact structure:

```
┌──────────────────────────────────────────────────────────────┐
│  Zone 1 — AppSidebar (zero-level, blends with page bg)       │
│  ┌────────────────────────────────┐  ┌──────────────────┐    │
│  │  Zone 2 — Elevated Main Panel  │  │ Zone 3 — Right   │    │
│  │  ┌──────────────────────────┐  │  │ Panel            │    │
│  │  │  AppHeader (pinned top)  │  │  │ (optional,       │    │
│  │  ├──────────────────────────┤  │  │  toggleable)     │    │
│  │  │  <main> — page content   │  │  │                  │    │
│  │  │  (scrollable)            │  │  │                  │    │
│  │  └──────────────────────────┘  │  └──────────────────┘    │
│  └────────────────────────────────┘                          │
└──────────────────────────────────────────────────────────────┘
```

### Props

| Prop | Type | Default | Description |
|---|---|---|---|
| `rightPanelOpen` | `boolean` | — | Controlled open state for the right panel |
| `onRightPanelChange` | `(open: boolean) => void` | — | Controlled change handler |
| `defaultRightPanelOpen` | `boolean` | `false` | Uncontrolled initial state |
| `children` | `ReactNode` | — | Page content rendered in `<main>` |

### Key behaviours

- Left sidebar width is capped at **15rem** (`--sidebar-width` CSS var override).
- The outer canvas is `h-svh` with `overflow-hidden` — the **viewport never scrolls**; only `<main>` scrolls internally.
- The right panel is hidden/shown by toggling `rightPanelOpen`; the header's `PanelRight` button is the trigger.

---

## 4. Layout Components

### 4.1 `AppSidebar` — Left Navigation

**File:** `components/layout/AppSidebar.tsx`

Collapsible icon-rail sidebar built on the shadcn `<Sidebar>` primitive.

**Navigation groups (hardcoded, to be wired to router):**

| Group | Items |
|---|---|
| **Overview** | Dashboard, Analytics, Inbox (badge: 12) |
| **Operations** | Orders, Customers, Billing |
| **System** | Settings |

- Active item is tracked with local `useState(selectedNav)`.
- User account menu is in `<SidebarFooter>` via a `<DropdownMenu>` (Account, Billing, Notifications, Log out).
- Currently uses placeholder user: **Mira Kapoor / mira@acme.io**.

> **Agent rule:** To add a new nav section, append an entry to the `navMain` array. To wire routing, replace `url: '#'` with Next.js `href` and wrap `<SidebarMenuButton>` in a `<Link>`.

---

### 4.2 `AppHeader` — Top Bar

**File:** `components/layout/AppHeader.tsx`

Fixed to the top of the elevated main panel. Contains:

- `<SidebarTrigger>` — collapses/expands the left sidebar.
- `<Breadcrumb>` — currently hardcoded to `Console > Dashboard`.
- **Search input** (hidden on mobile, visible `lg+`).
- **Notifications** bell button.
- **Right-panel toggle** button — switches between `PanelRight` and `PanelRightClose` icons; calls `onToggleRightPanel`.

### Props

| Prop | Type | Description |
|---|---|---|
| `rightPanelOpen` | `boolean` | Drives the toggle button icon and `aria-pressed` |
| `onToggleRightPanel` | `() => void` | Called when the panel toggle button is clicked |

---

### 4.3 `RightPanel` — Detail / Activity Sidebar

**File:** `components/layout/DetailSidebar.tsx`

An optional right-hand panel rendered alongside the main content card. Currently displays:

- **Quick stats** strip: Today's revenue ($4,820, +12.4%) and Pending orders (8, -2).
- **Activity feed**: timestamped list of team actions with success/warning/neutral tone icons.
- **"View all activity"** footer button (not yet wired).

Returns `null` when `open === false` — no DOM is rendered when hidden.

### Props

| Prop | Type | Description |
|---|---|---|
| `open` | `boolean` | Whether the panel is visible |
| `onClose` | `() => void` | Called by the × close button |

---

## 5. UI Primitives (`components/ui/`)

These are **shadcn/ui** generated components. Do **not** edit them manually unless patching a bug — regenerate via the shadcn CLI instead.

| File | Primitive |
|---|---|
| `sidebar.tsx` | `Sidebar`, `SidebarProvider`, `SidebarContent`, `SidebarMenu`, … |
| `avatar.tsx` | `Avatar`, `AvatarFallback`, `AvatarImage` |
| `badge.tsx` | `Badge` |
| `breadcrumb.tsx` | `Breadcrumb`, `BreadcrumbList`, `BreadcrumbItem`, … |
| `button.tsx` | `Button` |
| `card.tsx` | `Card`, `CardHeader`, `CardContent`, … |
| `dialog.tsx` | `Dialog`, `DialogContent`, `DialogHeader`, … |
| `dropdown-menu.tsx` | `DropdownMenu`, `DropdownMenuContent`, `DropdownMenuItem`, … |
| `input.tsx` | `Input` |
| `input-group.tsx` | `InputGroup` |
| `sheet.tsx` | `Sheet`, `SheetContent`, … |
| `tooltip.tsx` | `Tooltip`, `TooltipContent`, … |
| `progress.tsx` | `Progress` |
| `separator.tsx` | `Separator` |
| `skeleton.tsx` | `Skeleton` |
| `textarea.tsx` | `Textarea` |

---

## 6. Utilities

### `lib/utils.ts`

Exports `cn(...inputs)` — a `clsx` + `tailwind-merge` helper used throughout every component for conditional class composition.

---

## 7. Backend

The headless e-commerce backend lives at `../backend` (relative to this workspace root). API integration work (data fetching, server actions, route handlers) should bridge from that service into the pages under `app/`.

---

## 8. Development Commands

```bash
pnpm dev        # Start Next.js dev server
pnpm build      # Production build (standalone output)
pnpm npmlock    # Sync package-lock.json from pnpm lockfile
```

---

## 9. Agent Conventions

| Concern | Convention |
|---|---|
| **New pages** | Add a folder under `app/` with a `page.tsx`. The root layout applies automatically. |
| **New nav items** | Edit the `navMain` array in `AppSidebar.tsx`; wire `url` to the new route. |
| **New layout zones** | Extend `AdminShellProps` and `AdminShell` in `Shell.tsx`. |
| **New UI primitives** | Run `pnpm dlx shadcn@latest add <component>` — do not hand-write into `components/ui/`. |
| **Tailwind classes** | Always use `cn()` from `lib/utils.ts` for conditional/merged class strings. |
| **Client components** | Mark with `'use client'` when using hooks or browser APIs. Layout components are all client components. |
| **Data fetching** | Prefer React Server Components (`app/` pages without `'use client'`) for data fetching from `../backend`. |
