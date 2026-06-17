import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export default function PortfolioPage() {
  const [portfolio, setPortfolio] = useState<any>(null);
  const [balances, setBalances] = useState<any[]>([]);

  useEffect(() => {
    api.portfolio.get().then((res: any) => setPortfolio(res.data)).catch(() => {});
    api.portfolio.balances().then((res: any) => setBalances(res.data?.balances || [])).catch(() => {});
  }, []);

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold mb-6">Portfolio</h1>
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="p-4 rounded-xl" style={{ background: '#1e2329' }}>
          <div className="text-xs mb-1" style={{ color: '#848e9c' }}>Total Balance</div>
          <div className="text-xl font-bold" style={{ color: '#fcd535' }}>${(portfolio?.total_balance || 0).toFixed(2)}</div>
        </div>
        <div className="p-4 rounded-xl" style={{ background: '#1e2329' }}>
          <div className="text-xs mb-1" style={{ color: '#848e9c' }}>Market Value</div>
          <div className="text-xl font-bold">${(portfolio?.total_market_value || 0).toFixed(2)}</div>
        </div>
        <div className="p-4 rounded-xl" style={{ background: '#1e2329' }}>
          <div className="text-xs mb-1" style={{ color: '#848e9c' }}>Unrealized PnL</div>
          <div className="text-xl font-bold" style={{ color: (portfolio?.total_unrealized_pnl || 0) >= 0 ? '#0ecb81' : '#f6465d' }}>
            ${(portfolio?.total_unrealized_pnl || 0).toFixed(2)}
          </div>
        </div>
      </div>
      <div className="rounded-xl overflow-hidden mb-6" style={{ background: '#1e2329' }}>
        <div className="px-4 py-3 border-b" style={{ borderColor: '#2b3139' }}>
          <span className="text-sm font-semibold">Balances</span>
        </div>
        <table className="w-full text-sm">
          <thead><tr style={{ color: '#848e9c' }}>
            <th className="text-left px-4 py-3 font-normal">Currency</th>
            <th className="text-right px-4 py-3 font-normal">Balance</th>
            <th className="text-right px-4 py-3 font-normal">Frozen</th>
            <th className="text-right px-4 py-3 font-normal">Available</th>
          </tr></thead>
          <tbody>
            {balances.map((b, i) => (
              <tr key={i} className="border-t" style={{ borderColor: '#2b3139' }}>
                <td className="px-4 py-3 font-semibold">{b.currency}</td>
                <td className="px-4 py-3 text-right">{b.balance.toFixed(2)}</td>
                <td className="px-4 py-3 text-right" style={{ color: '#f6465d' }}>{b.frozen.toFixed(2)}</td>
                <td className="px-4 py-3 text-right" style={{ color: '#0ecb81' }}>{b.available.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329' }}>
        <div className="px-4 py-3 border-b" style={{ borderColor: '#2b3139' }}>
          <span className="text-sm font-semibold">Positions</span>
        </div>
        <table className="w-full text-sm">
          <thead><tr style={{ color: '#848e9c' }}>
            <th className="text-left px-4 py-3 font-normal">Symbol</th>
            <th className="text-right px-4 py-3 font-normal">Quantity</th>
            <th className="text-right px-4 py-3 font-normal">Avg Price</th>
            <th className="text-right px-4 py-3 font-normal">PnL</th>
          </tr></thead>
          <tbody>
            {(!portfolio?.positions || portfolio.positions.length === 0) ? (
              <tr><td colSpan={4} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>No positions</td></tr>
            ) : portfolio.positions.map((p: any, i: number) => (
              <tr key={i} className="border-t" style={{ borderColor: '#2b3139' }}>
                <td className="px-4 py-3 font-semibold">{p.symbol}</td>
                <td className="px-4 py-3 text-right">{p.quantity}</td>
                <td className="px-4 py-3 text-right">{p.average_price.toFixed(2)}</td>
                <td className="px-4 py-3 text-right" style={{ color: p.unrealized_pnl >= 0 ? '#0ecb81' : '#f6465d' }}>
                  ${p.unrealized_pnl.toFixed(2)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
