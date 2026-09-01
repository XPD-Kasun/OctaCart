'use client'

import {
    Activity,
    AlertCircle,
    CheckCircle2,
    Clock,
    TrendingDown,
    TrendingUp,
    X,
} from 'lucide-react'

import {
    Avatar,
    AvatarFallback,
} from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

export interface RightPanelProps {
    open: boolean
    onClose: () => void
}

const activity = [
    {
        id: '1',
        who: 'OL',
        name: 'Olivia Lee',
        action: 'approved invoice #1042',
        time: '2m ago',
        tone: 'success' as const,
    },
    {
        id: '2',
        who: 'JM',
        name: 'Jack Maher',
        action: 'flagged order #8831',
        time: '14m ago',
        tone: 'warning' as const,
    },
    {
        id: '3',
        who: 'SR',
        name: 'Sofia Reyes',
        action: 'updated billing profile',
        time: '1h ago',
        tone: 'neutral' as const,
    },
    {
        id: '4',
        who: 'AK',
        name: 'Amir Khan',
        action: 'closed ticket #412',
        time: '3h ago',
        tone: 'success' as const,
    },
    {
        id: '5',
        who: 'TN',
        name: 'Téa Núñez',
        action: 'refunded $129.00',
        time: '5h ago',
        tone: 'warning' as const,
    },
]

const quickStats = [
    {
        label: 'Today',
        value: '$4,820',
        delta: '+12.4%',
        trend: 'up' as const,
    },
    {
        label: 'Pending',
        value: '8',
        delta: '-2',
        trend: 'down' as const,
    },
]

export function RightPanel({ open, onClose }: RightPanelProps) {
    if (!open) return null

    return (
        <aside
            className={cn(
                // Cap width at 300px (18.75rem) — was previously w-80 (320) on md, w-96 (384) on lg
                'hidden h-full w-[18.75rem] max-w-[18.75rem] shrink-0 flex-col overflow-hidden rounded-2xl border bg-card md:flex',
            )}
            aria-label="Details panel"
        >
            {/* Panel header */}
            <div className="flex items-center justify-between border-b px-4 py-3">
                <div className="flex items-center gap-2">
                    <Activity className="size-4 text-muted-foreground" />
                    <span className="text-sm font-semibold">Activity</span>
                    <Badge variant="secondary" className="rounded-full px-2 py-0 text-xs">
                        Live
                    </Badge>
                </div>
                <Button
                    variant="ghost"
                    size="icon"
                    className="size-7 rounded-md"
                    onClick={onClose}
                    aria-label="Close panel"
                >
                    <X className="size-4" />
                </Button>
            </div>

            {/* Quick stats */}
            <div className="grid grid-cols-2 gap-px border-b bg-border/40">
                {quickStats.map((s) => (
                    <div key={s.label} className="bg-card px-4 py-3">
                        <div className="text-xs text-muted-foreground">{s.label}</div>
                        <div className="mt-0.5 text-lg font-semibold tracking-tight">
                            {s.value}
                        </div>
                        <div
                            className={cn(
                                'mt-0.5 flex items-center gap-1 text-xs',
                                s.trend === 'up' ? 'text-emerald-600' : 'text-rose-600',
                            )}
                        >
                            {s.trend === 'up' ? (
                                <TrendingUp className="size-3" />
                            ) : (
                                <TrendingDown className="size-3" />
                            )}
                            <span>{s.delta}</span>
                        </div>
                    </div>
                ))}
            </div>

            {/* Activity feed */}
            <div className="flex-1 overflow-y-auto">
                <ul className="divide-y">
                    {activity.map((item) => (
                        <li
                            key={item.id}
                            className="flex items-start gap-3 px-4 py-3 transition-colors hover:bg-muted/40"
                        >
                            <Avatar className="size-8 rounded-full">
                                <AvatarFallback className="rounded-full bg-muted text-xs">
                                    {item.who}
                                </AvatarFallback>
                            </Avatar>
                            <div className="min-w-0 flex-1">
                                <p className="text-sm leading-snug">
                                    <span className="font-medium">{item.name}</span>{' '}
                                    <span className="text-muted-foreground">{item.action}</span>
                                </p>
                                <div className="mt-1 flex items-center gap-1.5 text-xs text-muted-foreground">
                                    <Clock className="size-3" />
                                    <span>{item.time}</span>
                                </div>
                            </div>
                            <span className="mt-0.5">
                                {item.tone === 'success' && (
                                    <CheckCircle2 className="size-4 text-emerald-600" />
                                )}
                                {item.tone === 'warning' && (
                                    <AlertCircle className="size-4 text-amber-600" />
                                )}
                                {item.tone === 'neutral' && (
                                    <span className="size-2 rounded-full bg-muted-foreground/40" />
                                )}
                            </span>
                        </li>
                    ))}
                </ul>
            </div>

            {/* Footer action */}
            <div className="border-t p-3">
                <Button variant="outline" className="w-full justify-center" size="sm">
                    View all activity
                </Button>
            </div>
        </aside>
    )
}
