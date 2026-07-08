'use client';

import * as React from 'react';
import { Dashboard } from './dashboard';
import { useUI, type AppPage } from '@/stores/ui-store';

const pages: Record<AppPage, React.ReactNode> = {
  dashboard: <Dashboard />,
  research: <Placeholder title="Research" description="Multi-asset screen and RAG summaries" />,
  markets: <Placeholder title="Markets" description="Market structure, heatmaps, and cross-asset flow" />,
  portfolio: <Placeholder title="Portfolio" description="Holdings, P&L attribution, and drift" />,
  strategies: <Placeholder title="Strategies" description="Backtests, parameters, and live deployments" />,
  agents: <Placeholder title="Agents" description="Agent monitor, reasoning timeline, and logs" />,
  competitions: <Placeholder title="Competitions" description="Strategy leaderboards and paper-trading leagues" />,
  reports: <Placeholder title="Reports" description="Daily AI reports and reflection summaries" />,
  'prompt-studio': <Placeholder title="Prompt Studio" description="Version, test, and deploy prompts" />,
  'model-router': <Placeholder title="Model Router" description="Provider registry, debate, and confidence fusion" />,
  settings: <Placeholder title="Settings" description="Exchanges, brokers, and system preferences" />,
};

function Placeholder({ title, description }: { title: string; description: string }) {
  return (
    <div className="space-y-4">
      <h1 className="text-3xl font-semibold tracking-tight">{title}</h1>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}

export function PageRouter() {
  const { currentPage } = useUI();
  return pages[currentPage] ?? pages.dashboard;
}
