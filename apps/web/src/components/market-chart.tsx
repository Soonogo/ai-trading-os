'use client';

import { useEffect, useRef } from 'react';
import {
  createChart,
  ColorType,
  CrosshairMode,
  CandlestickSeries,
  type IChartApi,
  type CandlestickData,
} from 'lightweight-charts';

function generateSampleData(): CandlestickData[] {
  const data: CandlestickData[] = [];
  let price = 150;
  const now = new Date();
  for (let i = 0; i < 120; i++) {
    const open = price;
    const close = price + (Math.random() - 0.45) * 4;
    const high = Math.max(open, close) + Math.random() * 2;
    const low = Math.min(open, close) - Math.random() * 2;
    const time = new Date(now);
    time.setDate(time.getDate() - (120 - i));
    data.push({
      time: time.toISOString().split('T')[0],
      open,
      high,
      low,
      close,
    });
    price = close;
  }
  return data;
}

export function MarketChart() {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<IChartApi | null>(null);

  useEffect(() => {
    if (!chartContainerRef.current) return;

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: 'transparent' },
        textColor: 'hsl(220 10% 70%)',
      },
      grid: {
        vertLines: { color: 'rgba(255,255,255,0.03)' },
        horzLines: { color: 'rgba(255,255,255,0.03)' },
      },
      crosshair: {
        mode: CrosshairMode.Magnet,
      },
      rightPriceScale: {
        borderColor: 'rgba(255,255,255,0.06)',
      },
      timeScale: {
        borderColor: 'rgba(255,255,255,0.06)',
      },
      autoSize: true,
    });

    const series = chart.addSeries(CandlestickSeries, {
      upColor: 'hsl(142 60% 48%)',
      downColor: 'hsl(4 75% 55%)',
      wickUpColor: 'hsl(142 60% 48%)',
      wickDownColor: 'hsl(4 75% 55%)',
      borderUpColor: 'hsl(142 60% 48%)',
      borderDownColor: 'hsl(4 75% 55%)',
    });

    series.setData(generateSampleData());
    chart.timeScale().fitContent();

    chartRef.current = chart;

    return () => {
      chart.remove();
    };
  }, []);

  return <div ref={chartContainerRef} className="h-full w-full" />;
}
