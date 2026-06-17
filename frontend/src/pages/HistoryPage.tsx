import { useState, useEffect } from 'react';
import { api } from '../lib/api';

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
      <h1 className="text-xl font-bold mb-6">History</h1>
      <div className="flex gap-4 mb-6">
        <button onClick={() => setTab('orders')}
          className="px-4 py-2 rounded-lg text-sm"
          style={{ background: tab === 'orders' ? '#2b3139' : 'transparent', color: tab === 'orders' ? '#eaecef' : '#848e9c' }}>
          Orders
        </button>
        <button onClick={() => setTab('trades')}
          className="px-4 py-2 rounded-lg text-sm"
          style={{ background: tab === 'trades' ? '#2b3139' : 'transparent', color: tab === 'trades' ? '#eaecef' : '#848e9c' }}>
          Trades
        </button>
      </div>
      <div className="rounded-xl overflow-hidden" style={{ background: '#1e2329' }}>
        {tab === 'orders' ? (
          <table className="w-full text-sm">
            <thead><tr style={{ color: '#848e9c' }}>
              <th className="text-left px-4 py-3 font-normal">Time</th>
              <th className="text-left px-4 py-3 font-normal">Pair</th>
              <th className="text-left px-4 py-3 font-normal">Side</th>
              <th className="text-left px-4 py-3 font-normal">Type</th>
              <th className="text-right px-4 py-3 font-normal">Price</th>
              <th className="text-right px-4 py-3 font-normal">Qty</th>
              <th className="text-right px-4 py-3 font-normal">Filled</th>
              <th className="text-right px-4 py-3 font-normal">Status</th>
            </tr></thead>
            <tbody>
              {orders.length === 0 ? (
                <tr><td colSpan={8} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>No orders</td></tr>
              ) : orders.map((o: any, i: number) => (
                <tr key={i} className="border-t" style={{ borderColor: '#2b3139' }}>
                  <td className="px-4 py-3" style={{ color: '#848e9c' }}>{new Date(o.created_at).toLocaleString()}</td>
                  <td className="px-4 py-3 font-semibold">{o.symbol}</td>
                  <td className="px-4 py-3" style={{ color: o.side === 'BUY' ? '#0ecb81' : '#f6465d' }}>{o.side}</td>
                  <td className="px-4 py-3" style={{ color: '#848e9c' }}>{o.type}</td>
                  <td className="px-4 py-3 text-right">{o.price || 'Market'}</td>
                  <td className="px-4 py-3 text-right">{o.quantity}</td>
                  <td className="px-4 py-3 text-right">{o.filled_quantity}</td>
                  <td className="px-4 py-3 text-right">
                    <span className="px-2 py-0.5 rounded text-xs" style={{
                      background: o.status === 'FILLED' ? 'rgba(14,203,129,0.15)' : o.status === 'CANCELLED' ? 'rgba(132,142,156,0.15)' : 'rgba(252,213,53,0.15)',
                      color: o.status === 'FILLED' ? '#0ecb81' : o.status === 'CANCELLED' ? '#848e9c' : '#fcd535'
                    }}>{o.status}</span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <table className="w-full text-sm">
            <thead><tr style={{ color: '#848e9c' }}>
              <th className="text-left px-4 py-3 font-normal">Time</th>
              <th className="text-left px-4 py-3 font-normal">Pair</th>
              <th className="text-right px-4 py-3 font-normal">Price</th>
              <th className="text-right px-4 py-3 font-normal">Qty</th>
              <th className="text-right px-4 py-3 font-normal">Fee</th>
            </tr></thead>
            <tbody>
              {trades.length === 0 ? (
                <tr><td colSpan={5} className="px-4 py-8 text-center" style={{ color: '#848e9c' }}>No trades</td></tr>
              ) : trades.map((t: any, i: number) => (
                <tr key={i} className="border-t" style={{ borderColor: '#2b3139' }}>
                  <td className="px-4 py-3" style={{ color: '#848e9c' }}>{new Date(t.executed_at).toLocaleString()}</td>
                  <td className="px-4 py-3 font-semibold">{t.symbol}</td>
                  <td className="px-4 py-3 text-right">{t.price}</td>
                  <td className="px-4 py-3 text-right">{t.quantity}</td>
                  <td className="px-4 py-3 text-right" style={{ color: '#f6465d' }}>{t.buyer_fee || t.seller_fee || 0}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
