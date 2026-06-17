import { useState, useEffect } from 'react';
import { api } from '../lib/api';

export default function TradePage() {
  const [symbol, setSymbol] = useState('BTC/USDT');
  const [side, setSide] = useState<'BUY' | 'SELL'>('BUY');
  const [type, setType] = useState<'LIMIT' | 'MARKET'>('LIMIT');
  const [price, setPrice] = useState('');
  const [quantity, setQuantity] = useState('');
  const [orderBook, setOrderBook] = useState<any>(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = () => {
      api.market.orderbook(symbol, 10).then((res: any) => setOrderBook(res.data)).catch(() => {});
    };
    load();
    const interval = setInterval(load, 5000);
    return () => clearInterval(interval);
  }, [symbol]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setLoading(true);
    try {
      await api.orders.create({
        symbol,
        side,
        type,
        price: type === 'LIMIT' ? parseFloat(price) : undefined,
        quantity: parseFloat(quantity),
      });
      setSuccess('Order placed successfully');
      setQuantity('');
      setPrice('');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const total = (parseFloat(price) || 0) * (parseFloat(quantity) || 0);

  return (
    <div className="flex h-full">
      <div className="flex-1 p-6">
        <div className="flex items-center gap-4 mb-6">
          <select value={symbol} onChange={e => setSymbol(e.target.value)}
            className="px-4 py-2 rounded-lg text-sm font-semibold outline-none"
            style={{ background: '#1e2329', color: '#fcd535', border: '1px solid #2b3139' }}>
            <option>BTC/USDT</option><option>ETH/USDT</option><option>SOL/USDT</option>
          </select>
        </div>
        <div className="flex gap-4 mb-6">
          <button onClick={() => setSide('BUY')}
            className="flex-1 py-3 rounded-lg font-semibold text-sm"
            style={{ background: side === 'BUY' ? '#0ecb81' : '#2b3139', color: side === 'BUY' ? '#fff' : '#848e9c' }}>
            Buy
          </button>
          <button onClick={() => setSide('SELL')}
            className="flex-1 py-3 rounded-lg font-semibold text-sm"
            style={{ background: side === 'SELL' ? '#f6465d' : '#2b3139', color: side === 'SELL' ? '#fff' : '#848e9c' }}>
            Sell
          </button>
        </div>
        <div className="flex gap-4 mb-4">
          <button onClick={() => setType('LIMIT')}
            className="px-4 py-2 rounded-lg text-xs"
            style={{ background: type === 'LIMIT' ? '#2b3139' : 'transparent', color: type === 'LIMIT' ? '#eaecef' : '#848e9c' }}>
            Limit
          </button>
          <button onClick={() => setType('MARKET')}
            className="px-4 py-2 rounded-lg text-xs"
            style={{ background: type === 'MARKET' ? '#2b3139' : 'transparent', color: type === 'MARKET' ? '#eaecef' : '#848e9c' }}>
            Market
          </button>
        </div>
        <form onSubmit={handleSubmit}>
          {type === 'LIMIT' && (
            <div className="mb-4">
              <label className="block text-xs mb-1" style={{ color: '#848e9c' }}>Price (USDT)</label>
              <input type="number" step="0.01" value={price} onChange={e => setPrice(e.target.value)}
                className="w-full px-4 py-3 rounded-lg text-sm outline-none"
                style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }}
                required />
            </div>
          )}
          <div className="mb-4">
            <label className="block text-xs mb-1" style={{ color: '#848e9c' }}>Quantity</label>
            <input type="number" step="0.0001" value={quantity} onChange={e => setQuantity(e.target.value)}
              className="w-full px-4 py-3 rounded-lg text-sm outline-none"
              style={{ background: '#0b0e11', border: '1px solid #2b3139', color: '#eaecef' }}
              required />
          </div>
          {type === 'LIMIT' && total > 0 && (
            <div className="mb-4 text-xs" style={{ color: '#848e9c' }}>
              Total: <span style={{ color: '#eaecef' }}>${total.toFixed(2)} USDT</span>
            </div>
          )}
          {error && <div className="mb-4 p-3 rounded-lg text-xs" style={{ background: 'rgba(246,70,93,0.1)', color: '#f6465d' }}>{error}</div>}
          {success && <div className="mb-4 p-3 rounded-lg text-xs" style={{ background: 'rgba(14,203,129,0.1)', color: '#0ecb81' }}>{success}</div>}
          <button type="submit" disabled={loading}
            className="w-full py-3 rounded-lg font-semibold text-sm"
            style={{ background: side === 'BUY' ? '#0ecb81' : '#f6465d', color: '#fff', opacity: loading ? 0.7 : 1 }}>
            {loading ? 'Placing...' : `${side} ${symbol.split('/')[0]}`}
          </button>
        </form>
      </div>
      <div className="w-80 p-4 border-l" style={{ borderColor: '#2b3139', background: '#1e2329' }}>
        <div className="text-xs font-semibold mb-3" style={{ color: '#848e9c' }}>Order Book</div>
        <div className="mb-4">
          <div className="flex justify-between text-xs mb-1" style={{ color: '#848e9c' }}>
            <span>Price</span><span>Qty</span>
          </div>
          {orderBook?.asks?.slice(0, 8).reverse().map((ask: any, i: number) => (
            <div key={i} className="flex justify-between text-xs py-0.5">
              <span style={{ color: '#f6465d' }}>{ask.price.toLocaleString()}</span>
              <span style={{ color: '#848e9c' }}>{ask.quantity}</span>
            </div>
          ))}
        </div>
        <div className="text-lg font-bold py-2 text-center" style={{ color: orderBook?.asks?.[0] ? '#f6465d' : '#848e9c' }}>
          {orderBook?.asks?.[0]?.price?.toLocaleString() || '--'}
        </div>
        <div>
          {orderBook?.bids?.slice(0, 8).map((bid: any, i: number) => (
            <div key={i} className="flex justify-between text-xs py-0.5">
              <span style={{ color: '#0ecb81' }}>{bid.price.toLocaleString()}</span>
              <span style={{ color: '#848e9c' }}>{bid.quantity}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
