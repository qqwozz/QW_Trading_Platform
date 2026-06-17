import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { TrendingUp, TrendingDown, Wallet, BarChart3 } from 'lucide-react';

interface Ticker {
  id: string; symbol: string; last_price: number; best_bid: number; best_ask: number;
  volume_24h: number; high_24h: number; low_24h: number; change_24h: number; change_pct_24h: number;
}

export default function DashboardPage() {
  const [tickers, setTickers] = useState<Ticker[]>([]);
  const [portfolio, setPortfolio] = useState<any>(null);

  useEffect(() => {
    api.market.tickers().then((res: any) => setTickers(res.data || [])).catch(() => {});
    api.portfolio.get().then((res: any) => setPortfolio(res.data)).catch(() => {});
  }, []);

  const stats = [
    { label: 'Total Equity', value: `$${(portfolio?.total_equity || 0).toFixed(2)}`, icon: Wallet, color: '#fcd535' },
    { label: 'Available Balance', value: `$${(portfolio?.total_balance || 0).toFixed(2)}`, icon: BarChart3, color: '#0ecb81' },
    { label: 'Unrealized PnL', value: `$${(portfolio?.total_unrealized_pnl || 0).toFixed(2)}`, icon: TrendingUp, color: '#0ecb81' },
    { label: 'Positions', value: `${portfolio?.positions?.length || 0}`, icon: TrendingDown, color: '#848e9c' },
  ];

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold mb-6">Dashboard</h1>
      <div className="grid grid-cols-4 gap-4 mb-8">
        {stats.map(s => {
          const Icon = s.icon;
          return (
            <div key={s.label} className="p-4 rounded-xl" style={{ background: '#1e2329' }}>
              <div className="flex items-center gap-2 mb-2">
                <Icon size={16} style={{ color: s.color }} />
                <span className="text-xs" style={{ color: '#848e9c' }}>{s.label}</span>
              </div>
              <div className="text-xl font-bold" style={{ color: s.color }}>{s.value}</div>
            </div>
          );
        })}
      </div>
      <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329' }}>
        <div className="px-4 py-3 border-b" style={{ borderColor: '#2b3139' }}>
          <span className="text-sm font-semibold">Market Tickers</span>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr style={{ color: '#848e9c' }}>
                <th className="text-left px-4 py-3 font-normal">Pair</th>
                <th className="text-right px-4 py-3 font-normal">Last Price</th>
                <th className="text-right px-4 py-3 font-normal">24h Change</th>
                <th className="text-right px-4 py-3 font-normal">24h High</th>
                <th className="text-right px-4 py-3 font-normal">24h Low</th>
                <th className="text-right px-4 py-3 font-normal">Volume</th>
              </tr>
            </thead>
            <tbody>
              {tickers.length === 0 ? (
                <tr><td colSpan={6} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>No market data available</td></tr>
              ) : tickers.map(t => (
                <tr key={t.id} className="border-t" style={{ borderColor: '#2b3139' }}>
                  <td className="px-4 py-3 font-semibold">{t.symbol}</td>
                  <td className="px-4 py-3 text-right">{t.last_price.toLocaleString()}</td>
                  <td className="px-4 py-3 text-right" style={{ color: t.change_24h >= 0 ? '#0ecb81' : '#f6465d' }}>
                    {t.change_24h >= 0 ? '+' : ''}{t.change_pct_24h.toFixed(2)}%
                  </td>
                  <td className="px-4 py-3 text-right" style={{ color: '#848e9c' }}>{t.high_24h.toLocaleString()}</td>
                  <td className="px-4 py-3 text-right" style={{ color: '#848e9c' }}>{t.low_24h.toLocaleString()}</td>
                  <td className="px-4 py-3 text-right" style={{ color: '#848e9c' }}>{t.volume_24h.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
