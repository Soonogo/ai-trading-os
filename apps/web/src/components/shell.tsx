'use client';

import { Sidebar } from './sidebar';
import { AIChat } from './ai-chat';
import { useUI } from '@/stores/ui-store';
import { cn } from '@/lib/utils';

export function Shell({ children }: { children: React.ReactNode }) {
  const { sidebarOpen } = useUI();

  return (
    <div className="flex h-screen w-screen overflow-hidden">
      <Sidebar />
      <main
        className={cn(
          'flex-1 overflow-hidden transition-[margin-left] duration-300 ease-out',
          sidebarOpen ? 'ml-60' : 'ml-16'
        )}
      >
        <div className="grid h-full grid-cols-[1fr_320px]">
          <div className="overflow-y-auto p-6">
            {children}
          </div>
          <AIChat />
        </div>
      </main>
    </div>
  );
}
