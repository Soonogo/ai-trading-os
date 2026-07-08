'use client';

import { Activity, DollarSign, Brain, Shield } from 'lucide-react';
import { motion } from 'framer-motion';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@aitos/design-system';
import { MarketChart } from './market-chart';
import { cn } from '@/lib/utils';

const stats = [
  { label: 'Portfolio P&L', value: '+3.24%', sub: '$124,500', trend: 'up' },
  { label: 'Active Trades', value: '12', sub: '4 pending fills', trend: 'neutral' },
  { label: 'AI Confidence', value: '78%', sub: 'Multi-model debate', trend: 'up' },
  { label: 'Risk Heat', value: '12%', sub: 'Within limits', trend: 'down' },
];

const agents = [
  { name: 'Research', status: 'scanning', icon: Brain },
  { name: 'Risk', status: 'idle', icon: Shield },
  { name: 'Execution', status: 'live', icon: Activity },
  { name: 'Macro', status: 'analyzing', icon: DollarSign },
];

export function Dashboard() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground">Autonomous market intelligence for {new Date().toLocaleDateString()}.</p>
        </div>
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span className="relative flex h-2 w-2">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-chart-up opacity-75"></span>
            <span className="relative inline-flex h-2 w-2 rounded-full bg-chart-up"></span>
          </span>
          Live
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((s, i) => (
          <motion.div
            key={s.label}
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: i * 0.05 }}
          >
            <Card className="p-0">
              <CardHeader className="pb-2">
                <CardDescription className="text-muted-foreground">{s.label}</CardDescription>
                <CardTitle className="text-2xl">
                  <span className={s.trend === 'up' ? 'text-chart-up' : s.trend === 'down' ? 'text-chart-down' : 'text-foreground'}>
                    {s.value}
                  </span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-xs text-muted-foreground">{s.sub}</p>
              </CardContent>
            </Card>
          </motion.div>
        ))}
      </div>

      <div className="grid h-[460px] grid-cols-1 gap-4 lg:grid-cols-3">
        <Card className="col-span-2 p-0 overflow-hidden">
          <CardHeader className="pb-0">
            <CardTitle className="text-base font-medium">Market View</CardTitle>
            <CardDescription>Price, volume, and agent signals</CardDescription>
          </CardHeader>
          <CardContent className="h-[calc(100%-80px)]">
            <MarketChart />
          </CardContent>
        </Card>

        <Card className="p-0">
          <CardHeader>
            <CardTitle className="text-base font-medium">Agent Monitor</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {agents.map((a) => {
              const Icon = a.icon;
              return (
                <div key={a.name} className="flex items-center justify-between rounded-lg bg-white/[0.03] px-3 py-2">
                  <div className="flex items-center gap-2">
                    <Icon className="h-4 w-4 text-primary" />
                    <span className="text-sm font-medium">{a.name}</span>
                  </div>
                  <span className={cn(
                    'rounded-full px-2 py-0.5 text-[10px] font-medium uppercase',
                    a.status === 'live' && 'bg-chart-up/20 text-chart-up',
                    a.status === 'analyzing' && 'bg-primary/20 text-primary',
                    a.status === 'scanning' && 'bg-accent/20 text-accent',
                    a.status === 'idle' && 'bg-muted text-muted-foreground'
                  )}>
                    {a.status}
                  </span>
                </div>
              );
            })}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
