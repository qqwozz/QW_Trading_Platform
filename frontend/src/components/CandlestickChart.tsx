import { useEffect, useRef, useState } from 'react';
import { createChart, ColorType, CandlestickSeries } from 'lightweight-charts';

interface Kline {
  openTime: number;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  closeTime: number;
}

interface Props {
  symbol?: string;
  height?: number;
}

const SYMBOL_MAP: Record<string, string> = {
  'BTC/USDT': 'BTCUSDT',
  'ETH/USDT': 'ETHUSDT',
  'SOL/USDT': 'SOLUSDT',
};

export default function CandlestickChart({ symbol = 'BTC/USDT', height = 400 }: Props) {
  const chartContainerRef = useRef<HTMLDivElement>(null);
  const chartRef = useRef<ReturnType<typeof createChart> | null>(null);
  const seriesRef = useRef<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!chartContainerRef.current) return;

    const chart = createChart(chartContainerRef.current, {
      layout: {
        background: { type: ColorType.Solid, color: '#1e2329' },
        textColor: '#848e9c',
      },
      grid: {
        vertLines: { color: '#2b3139' },
        horzLines: { color: '#2b3139' },
      },
      crosshair: {
        mode: 0,
      },
      timeScale: {
        borderColor: '#2b3139',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#2b3139',
      },
      width: chartContainerRef.current.clientWidth,
      height,
    });

    const series = chart.addSeries(CandlestickSeries, {
      upColor: '#0ecb81',
      downColor: '#f6465d',
      borderDownColor: '#f6465d',
      borderUpColor: '#0ecb81',
      wickDownColor: '#f6465d',
      wickUpColor: '#0ecb81',
    });

    chartRef.current = chart;
    seriesRef.current = series;

    const handleResize = () => {
      if (chartContainerRef.current && chartRef.current) {
        chartRef.current.resize(chartContainerRef.current.clientWidth, height);
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      chart.remove();
    };
  }, [height]);

  useEffect(() => {
    const fetchKlines = async () => {
      try {
        const binanceSymbol = SYMBOL_MAP[symbol] || symbol.replace('/', '');
        const res = await fetch(
          `https://api.binance.com/api/v3/klines?symbol=${binanceSymbol}&interval=1h&limit=100`
        );
        const data = await res.json();

        const klines: Kline[] = data.map((k: any[]) => ({
          openTime: k[0],
          open: parseFloat(k[1]),
          high: parseFloat(k[2]),
          low: parseFloat(k[3]),
          close: parseFloat(k[4]),
          volume: parseFloat(k[5]),
          closeTime: k[6],
        }));

        if (seriesRef.current) {
          seriesRef.current.setData(
            klines.map(k => ({
              time: Math.floor(k.openTime / 1000) as any,
              open: k.open,
              high: k.high,
              low: k.low,
              close: k.close,
            }))
          );
        }
        setLoading(false);
      } catch (err) {
        console.error('Failed to fetch klines:', err);
        setLoading(false);
      }
    };

    fetchKlines();
    const interval = setInterval(fetchKlines, 30000);
    return () => clearInterval(interval);
  }, [symbol]);

  return (
    <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329' }}>
      <div className="px-4 py-3 border-b flex items-center justify-between" style={{ borderColor: '#2b3139' }}>
        <span className="text-sm font-semibold">{symbol} - 1H Chart</span>
        <span className="text-xs" style={{ color: '#848e9c' }}>Last 100 candles</span>
      </div>
      {loading && (
        <div className="flex items-center justify-center" style={{ height }}>
          <span style={{ color: '#848e9c' }}>Loading chart data...</span>
        </div>
      )}
      <div ref={chartContainerRef} style={{ display: loading ? 'none' : 'block' }} />
    </div>
  );
}
