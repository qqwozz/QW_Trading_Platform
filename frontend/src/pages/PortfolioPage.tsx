import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { Briefcase, DollarSign, TrendingUp, TrendingDown, Layers } from 'lucide-react';

export default function PortfolioPage() {
  const [portfolio, setPortfolio] = useState<any>(null);
  const [balances, setBalances] = useState<any[]>([]);

  useEffect(() => {
    api.portfolio.get().then((res: any) => setPortfolio(res.data)).catch(() => {});
    api.portfolio.balances().then((res: any) => setBalances(res.data?.balances || [])).catch(() => {});
  }, []);

  const pnlPositive = (portfolio?.total_unrealized_pnl || 0) >= 0;

  return (
    <div className="p-6">
      <div className="flex items-center gap-3 mb-6">
        <Briefcase size={22} style={{ color: '#fcd535' }} />
        <h1 className="text-xl font-bold">Portfolio</h1>
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div
          className="p-5 rounded-xl"
          style={{
            background: 'linear-gradient(135deg, rgba(252,213,53,0.12), rgba(252,213,53,0.03))',
            border: '1px solid rgba(252,213,53,0.25)',
          }}
        >
          <div className="flex items-center gap-2 mb-2">
            <DollarSign size={16} style={{ color: '#fcd535' }} />
            <div className="text-xs font-medium" style={{ color: '#848e9c' }}>Total Balance</div>
          </div>
          <div className="text-2xl font-bold font-mono" style={{ color: '#fcd535' }}>
            ${(portfolio?.total_balance || 0).toFixed(2)}
          </div>
        </div>
        <div
          className="p-5 rounded-xl"
          style={{
            background: 'linear-gradient(135deg, rgba(132,142,156,0.12), rgba(132,142,156,0.03))',
            border: '1px solid rgba(132,142,156,0.25)',
          }}
        >
          <div className="flex items-center gap-2 mb-2">
            <Layers size={16} style={{ color: '#848e9c' }} />
            <div className="text-xs font-medium" style={{ color: '#848e9c' }}>Market Value</div>
          </div>
          <div className="text-2xl font-bold font-mono">
            ${(portfolio?.total_market_value || 0).toFixed(2)}
          </div>
        </div>
        <div
          className="p-5 rounded-xl"
          style={{
            background: pnlPositive
              ? 'linear-gradient(135deg, rgba(14,203,129,0.12), rgba(14,203,129,0.03))'
              : 'linear-gradient(135deg, rgba(246,70,93,0.12), rgba(246,70,93,0.03))',
            border: pnlPositive ? '1px solid rgba(14,203,129,0.25)' : '1px solid rgba(246,70,93,0.25)',
          }}
        >
          <div className="flex items-center gap-2 mb-2">
            {pnlPositive ? <TrendingUp size={16} style={{ color: '#0ecb81' }} /> : <TrendingDown size={16} style={{ color: '#f6465d' }} />}
            <div className="text-xs font-medium" style={{ color: '#848e9c' }}>Unrealized PnL</div>
          </div>
          <div className="text-2xl font-bold font-mono" style={{ color: pnlPositive ? '#0ecb81' : '#f6465d' }}>
            ${(portfolio?.total_unrealized_pnl || 0).toFixed(2)}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-6">
        <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
          <div className="px-4 py-3 border-b" style={{ borderColor: '#2b3139' }}>
            <span className="text-sm font-semibold">Balances</span>
          </div>
          <table className="w-full text-sm">
            <thead>
              <tr style={{ color: '#848e9c' }}>
                <th className="text-left px-4 py-3 font-normal">Currency</th>
                <th className="text-right px-4 py-3 font-normal">Balance</th>
                <th className="text-right px-4 py-3 font-normal">Frozen</th>
                <th className="text-right px-4 py-3 font-normal">Available</th>
              </tr>
            </thead>
            <tbody>
              {balances.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>
                    No balances
                  </td>
                </tr>
              ) : (
                balances.map((b, i) => (
                  <tr
                    key={i}
                    className="border-t transition-colors duration-150"
                    style={{ borderColor: '#2b3139' }}
                    onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.02)'; }}
                    onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; }}
                  >
                    <td className="px-4 py-3 font-semibold">{b.currency}</td>
                    <td className="px-4 py-3 text-right font-mono">{b.balance.toFixed(2)}</td>
                    <td className="px-4 py-3 text-right font-mono" style={{ color: '#f6465d' }}>{b.frozen.toFixed(2)}</td>
                    <td className="px-4 py-3 text-right font-mono" style={{ color: '#0ecb81' }}>{b.available.toFixed(2)}</td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
          <div className="px-4 py-3 border-b" style={{ borderColor: '#2b3139' }}>
            <span className="text-sm font-semibold">Positions</span>
          </div>
          <table className="w-full text-sm">
            <thead>
              <tr style={{ color: '#848e9c' }}>
                <th className="text-left px-4 py-3 font-normal">Symbol</th>
                <th className="text-right px-4 py-3 font-normal">Quantity</th>
                <th className="text-right px-4 py-3 font-normal">Avg Price</th>
                <th className="text-right px-4 py-3 font-normal">PnL</th>
              </tr>
            </thead>
            <tbody>
              {!portfolio?.positions || portfolio.positions.length === 0 ? (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>
                    No positions
                  </td>
                </tr>
              ) : (
                portfolio.positions.map((p: any, i: number) => (
                  <tr
                    key={i}
                    className="border-t transition-colors duration-150"
                    style={{ borderColor: '#2b3139' }}
                    onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.02)'; }}
                    onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; }}
                  >
                    <td className="px-4 py-3 font-semibold">{p.symbol}</td>
                    <td className="px-4 py-3 text-right font-mono">{p.quantity}</td>
                    <td className="px-4 py-3 text-right font-mono">{p.average_price.toFixed(2)}</td>
                    <td className="px-4 py-3 text-right font-mono font-medium" style={{ color: p.unrealized_pnl >= 0 ? '#0ecb81' : '#f6465d' }}>
                      ${p.unrealized_pnl.toFixed(2)}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
