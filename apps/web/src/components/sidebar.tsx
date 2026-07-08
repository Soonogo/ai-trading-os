'use client';

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
import { useUI, type AppPage } from '@/stores/ui-store';
import { cn } from '@/lib/utils';

const nav: { page: AppPage; label: string; icon: LucideIcon }[] = [
  { page: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { page: 'research', label: 'Research', icon: Search },
  { page: 'markets', label: 'Markets', icon: BarChart3 },
  { page: 'portfolio', label: 'Portfolio', icon: Briefcase },
  { page: 'strategies', label: 'Strategies', icon: GitBranch },
  { page: 'agents', label: 'Agents', icon: Bot },
  { page: 'competitions', label: 'Competitions', icon: Trophy },
  { page: 'reports', label: 'Reports', icon: FileText },
  { page: 'prompt-studio', label: 'Prompt Studio', icon: Terminal },
  { page: 'model-router', label: 'Model Router', icon: GitCompare },
  { page: 'settings', label: 'Settings', icon: Settings },
];

export function Sidebar() {
  const { sidebarOpen, currentPage, setPage, toggleSidebar } = useUI();

  return (
    <aside
      className={cn(
        'fixed left-0 top-0 z-40 h-full glass-elevated border-r-0 transition-[width] duration-300 ease-out',
        sidebarOpen ? 'w-60' : 'w-16'
      )}
    >
      <div className="flex h-16 items-center gap-3 px-4 border-b border-white/[0.06]">
        <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center text-sm font-bold text-white">
          A
        </div>
        {sidebarOpen && (
          <span className="text-lg font-semibold tracking-tight">AITOS</span>
        )}
      </div>
      <nav className="flex flex-col gap-1 p-2">
        {nav.map((item) => (
          <button
            key={item.page}
            onClick={() => setPage(item.page)}
            className={cn(
              'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
              currentPage === item.page
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-white/[0.04] hover:text-foreground'
            )}
          >
            <item.icon className="h-5 w-5 shrink-0" />
            {sidebarOpen && <span>{item.label}</span>}
          </button>
        ))}
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
