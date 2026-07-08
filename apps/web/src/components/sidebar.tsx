'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import {
  LayoutDashboard,
  Search,
  BarChart3,
  Briefcase,
  GitBranch,
  Bot,
  Trophy,
  FileText,
  Terminal,
  GitCompare,
  Settings,
  type LucideIcon,
} from 'lucide-react';
import { useUI } from '@/stores/ui-store';
import { cn } from '@/lib/utils';

const nav: { href: string; label: string; icon: LucideIcon }[] = [
  { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/research', label: 'Research', icon: Search },
  { href: '/markets', label: 'Markets', icon: BarChart3 },
  { href: '/portfolio', label: 'Portfolio', icon: Briefcase },
  { href: '/strategies', label: 'Strategies', icon: GitBranch },
  { href: '/agents', label: 'Agents', icon: Bot },
  { href: '/competitions', label: 'Competitions', icon: Trophy },
  { href: '/reports', label: 'Reports', icon: FileText },
  { href: '/prompt-studio', label: 'Prompt Studio', icon: Terminal },
  { href: '/model-router', label: 'Model Router', icon: GitCompare },
  { href: '/settings', label: 'Settings', icon: Settings },
];

export function Sidebar() {
  const { sidebarOpen, toggleSidebar } = useUI();
  const pathname = usePathname();

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 z-40 h-full glass-elevated border-r-0 transition-[width] duration-300 ease-out',
        sidebarOpen ? 'w-60' : 'w-16'
      )}
    >
      <Link href="/dashboard" className="flex h-16 items-center gap-3 px-4 border-b border-white/[0.06]">
        <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center text-sm font-bold text-white">
          A
        </div>
        {sidebarOpen && (
          <span className="text-lg font-semibold tracking-tight">AITOS</span>
        )}
      </Link>
      <nav className="flex flex-col gap-1 p-2">
        {nav.map((item) => {
          const active = pathname === item.href || pathname.startsWith(`${item.href}/`);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                active
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted-foreground hover:bg-white/[0.04] hover:text-foreground'
              )}
            >
              <item.icon className="h-5 w-5 shrink-0" />
              {sidebarOpen && <span>{item.label}</span>}
            </Link>
          );
        })}
      </nav>
      <button
        onClick={toggleSidebar}
        className="absolute bottom-4 left-1/2 -translate-x-1/2 rounded-md p-2 text-muted-foreground hover:bg-white/[0.04] hover:text-foreground"
        aria-label="Toggle sidebar"
      >
        <span className="sr-only">Toggle</span>
        <div
          className={cn(
            'h-1 w-6 rounded-full bg-current transition-all',
            sidebarOpen ? 'w-4' : 'w-6'
          )}
        />
      </button>
    </aside>
  );
}
