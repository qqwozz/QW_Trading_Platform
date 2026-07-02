import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import { History, FileText, ArrowLeftRight } from 'lucide-react';

export default function HistoryPage() {
  const [orders, setOrders] = useState<any[]>([]);
  const [trades, setTrades] = useState<any[]>([]);
  const [tab, setTab] = useState<'orders' | 'trades'>('orders');

  useEffect(() => {
    api.history.orders().then((res: any) => setOrders(res.data || [])).catch(() => {});
    api.history.trades().then((res: any) => setTrades(res.data || [])).catch(() => {});
  }, []);

  return (
    <div className="p-6">
      <div className="flex items-center gap-3 mb-6">
        <History size={22} style={{ color: '#fcd535' }} />
        <h1 className="text-xl font-bold">History</h1>
      </div>

      <div className="flex gap-3 mb-6">
        <button
          onClick={() => setTab('orders')}
          className="flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition-all duration-200"
          style={{
            background: tab === 'orders' ? 'linear-gradient(135deg, rgba(252,213,53,0.12), rgba(252,213,53,0.04))' : 'transparent',
            color: tab === 'orders' ? '#fcd535' : '#848e9c',
            border: tab === 'orders' ? '1px solid rgba(252,213,53,0.3)' : '1px solid #2b3139',
          }}
        >
          <FileText size={14} />
          Orders
        </button>
        <button
          onClick={() => setTab('trades')}
          className="flex items-center gap-2 px-4 py-2.5 rounded-lg text-sm font-medium transition-all duration-200"
          style={{
            background: tab === 'trades' ? 'linear-gradient(135deg, rgba(252,213,53,0.12), rgba(252,213,53,0.04))' : 'transparent',
            color: tab === 'trades' ? '#fcd535' : '#848e9c',
            border: tab === 'trades' ? '1px solid rgba(252,213,53,0.3)' : '1px solid #2b3139',
          }}
        >
          <ArrowLeftRight size={14} />
          Trades
        </button>
      </div>

      <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
        {tab === 'orders' ? (
          <table className="w-full text-sm">
            <thead>
              <tr style={{ color: '#848e9c' }}>
                <th className="text-left px-4 py-3 font-normal">Time</th>
                <th className="text-left px-4 py-3 font-normal">Pair</th>
                <th className="text-left px-4 py-3 font-normal">Side</th>
                <th className="text-left px-4 py-3 font-normal">Type</th>
                <th className="text-right px-4 py-3 font-normal">Price</th>
                <th className="text-right px-4 py-3 font-normal">Qty</th>
                <th className="text-right px-4 py-3 font-normal">Filled</th>
                <th className="text-right px-4 py-3 font-normal">Status</th>
              </tr>
            </thead>
            <tbody>
              {orders.length === 0 ? (
                <tr>
                  <td colSpan={8} className="px-4 py-12 text-center" style={{ color: '#848e9c' }}>
                    <FileText size={32} className="mx-auto mb-3" style={{ opacity: 0.3 }} />
                    <div>No orders yet</div>
                  </td>
                </tr>
              ) : (
                orders.map((o: any, i: number) => (
                  <tr
                    key={i}
                    className="border-t transition-colors duration-150"
                    style={{ borderColor: '#2b3139' }}
                    onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.02)'; }}
                    onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; }}
                  >
                    <td className="px-4 py-3 font-mono" style={{ color: '#848e9c' }}>
                      {new Date(o.created_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 font-semibold">{o.symbol}</td>
                    <td className="px-4 py-3">
                      <span
                        className="px-2 py-0.5 rounded text-xs font-medium"
                        style={{
                          background: o.side === 'BUY' ? 'rgba(14,203,129,0.12)' : 'rgba(246,70,93,0.12)',
                          color: o.side === 'BUY' ? '#0ecb81' : '#f6465d',
                        }}
                      >
                        {o.side}
                      </span>
                    </td>
                    <td className="px-4 py-3" style={{ color: '#848e9c' }}>{o.type}</td>
                    <td className="px-4 py-3 text-right font-mono">{o.price || 'Market'}</td>
                    <td className="px-4 py-3 text-right font-mono">{o.quantity}</td>
                    <td className="px-4 py-3 text-right font-mono">{o.filled_quantity}</td>
                    <td className="px-4 py-3 text-right">
                      <span
                        className="px-2 py-0.5 rounded text-xs font-medium"
                        style={{
                          background: o.status === 'FILLED' ? 'rgba(14,203,129,0.12)' : o.status === 'CANCELLED' ? 'rgba(132,142,156,0.12)' : 'rgba(252,213,53,0.12)',
                          color: o.status === 'FILLED' ? '#0ecb81' : o.status === 'CANCELLED' ? '#848e9c' : '#fcd535',
                        }}
                      >
                        {o.status}
                      </span>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr style={{ color: '#848e9c' }}>
                <th className="text-left px-4 py-3 font-normal">Time</th>
                <th className="text-left px-4 py-3 font-normal">Pair</th>
                <th className="text-right px-4 py-3 font-normal">Price</th>
                <th className="text-right px-4 py-3 font-normal">Qty</th>
                <th className="text-right px-4 py-3 font-normal">Fee</th>
              </tr>
            </thead>
            <tbody>
              {trades.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-4 py-12 text-center" style={{ color: '#848e9c' }}>
                    <ArrowLeftRight size={32} className="mx-auto mb-3" style={{ opacity: 0.3 }} />
                    <div>No trades yet</div>
                  </td>
                </tr>
              ) : (
                trades.map((t: any, i: number) => (
                  <tr
                    key={i}
                    className="border-t transition-colors duration-150"
                    style={{ borderColor: '#2b3139' }}
                    onMouseEnter={(e) => { e.currentTarget.style.background = 'rgba(255,255,255,0.02)'; }}
                    onMouseLeave={(e) => { e.currentTarget.style.background = 'transparent'; }}
                  >
                    <td className="px-4 py-3 font-mono" style={{ color: '#848e9c' }}>
                      {new Date(t.executed_at).toLocaleString()}
                    </td>
                    <td className="px-4 py-3 font-semibold">{t.symbol}</td>
                    <td className="px-4 py-3 text-right font-mono">{t.price}</td>
                    <td className="px-4 py-3 text-right font-mono">{t.quantity}</td>
                    <td className="px-4 py-3 text-right font-mono" style={{ color: '#f6465d' }}>
                      {t.buyer_fee || t.seller_fee || 0}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
