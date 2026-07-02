import { useState, useEffect } from 'react';
import { api } from '../lib/api';
import CandlestickChart from '../components/CandlestickChart';
import { ArrowDownUp } from 'lucide-react';

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
      <div className="flex-1 flex flex-col">
        <div className="p-4 border-b flex items-center gap-4" style={{ borderColor: '#2b3139', background: '#1e2329' }}>
          <div className="flex items-center gap-2">
            <ArrowDownUp size={16} style={{ color: '#fcd535' }} />
            <span className="text-sm font-bold">{symbol}</span>
          </div>
          <div className="flex gap-2">
            {['BTC/USDT', 'ETH/USDT', 'SOL/USDT'].map(s => (
              <button
                key={s}
                onClick={() => setSymbol(s)}
                className="px-3 py-1 rounded text-xs font-medium transition-all duration-200"
                style={{
                  background: symbol === s ? 'rgba(252,213,53,0.15)' : 'transparent',
                  color: symbol === s ? '#fcd535' : '#848e9c',
                  border: symbol === s ? '1px solid rgba(252,213,53,0.3)' : '1px solid transparent',
                }}
              >
                {s}
              </button>
            ))}
          </div>
        </div>
        <div className="flex-1 p-4" style={{ minHeight: 0 }}>
          <CandlestickChart symbol={symbol} height={340} />
        </div>
        <div className="p-4 border-t" style={{ borderColor: '#2b3139' }}>
          <div className="grid grid-cols-4 gap-4">
            <div className="p-3 rounded-lg" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
              <div className="text-xs mb-1" style={{ color: '#848e9c' }}>24h High</div>
              <div className="text-sm font-bold font-mono" style={{ color: '#0ecb81' }}>--</div>
            </div>
            <div className="p-3 rounded-lg" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
              <div className="text-xs mb-1" style={{ color: '#848e9c' }}>24h Low</div>
              <div className="text-sm font-bold font-mono" style={{ color: '#f6465d' }}>--</div>
            </div>
            <div className="p-3 rounded-lg" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
              <div className="text-xs mb-1" style={{ color: '#848e9c' }}>24h Volume</div>
              <div className="text-sm font-bold font-mono">--</div>
            </div>
            <div className="p-3 rounded-lg" style={{ background: '#1e2329', border: '1px solid #2b3139' }}>
              <div className="text-xs mb-1" style={{ color: '#848e9c' }}>Spread</div>
              <div className="text-sm font-bold font-mono">--</div>
            </div>
          </div>
        </div>
      </div>

      <div className="w-72 flex flex-col border-l" style={{ borderColor: '#2b3139', background: '#1e2329' }}>
        <div className="p-4 border-b" style={{ borderColor: '#2b3139' }}>
          <div className="text-xs font-semibold mb-3" style={{ color: '#848e9c' }}>Place Order</div>
          <div className="flex gap-2 mb-3">
            <button
              onClick={() => setSide('BUY')}
              className="flex-1 py-2.5 rounded-lg font-semibold text-xs transition-all duration-200"
              style={{
                background: side === 'BUY' ? 'linear-gradient(135deg, #0ecb81, #0db574)' : '#2b3139',
                color: side === 'BUY' ? '#fff' : '#848e9c',
                boxShadow: side === 'BUY' ? '0 2px 8px rgba(14,203,129,0.3)' : 'none',
              }}
            >
              Buy
            </button>
            <button
              onClick={() => setSide('SELL')}
              className="flex-1 py-2.5 rounded-lg font-semibold text-xs transition-all duration-200"
              style={{
                background: side === 'SELL' ? 'linear-gradient(135deg, #f6465d, #e53950)' : '#2b3139',
                color: side === 'SELL' ? '#fff' : '#848e9c',
                boxShadow: side === 'SELL' ? '0 2px 8px rgba(246,70,93,0.3)' : 'none',
              }}
            >
              Sell
            </button>
          </div>
          <div className="flex gap-2">
            {(['LIMIT', 'MARKET'] as const).map(t => (
              <button
                key={t}
                onClick={() => setType(t)}
                className="px-3 py-1.5 rounded text-xs font-medium transition-all duration-200"
                style={{
                  background: type === t ? 'rgba(252,213,53,0.12)' : 'transparent',
                  color: type === t ? '#fcd535' : '#848e9c',
                  border: type === t ? '1px solid rgba(252,213,53,0.3)' : '1px solid #2b3139',
                }}
              >
                {t}
              </button>
            ))}
          </div>
        </div>

        <form onSubmit={handleSubmit} className="flex-1 p-4 flex flex-col">
          {type === 'LIMIT' && (
            <div className="mb-3">
              <label className="block text-xs mb-1" style={{ color: '#848e9c' }}>Price (USDT)</label>
              <input
                type="number"
                step="0.01"
                value={price}
                onChange={e => setPrice(e.target.value)}
                className="w-full px-3 py-2.5 rounded-lg text-sm outline-none font-mono transition-all duration-200"
                style={{
                  background: '#0b0e11',
                  border: '1px solid #2b3139',
                  color: '#eaecef',
                }}
                onFocus={(e) => { e.currentTarget.style.borderColor = '#fcd535'; }}
                onBlur={(e) => { e.currentTarget.style.borderColor = '#2b3139'; }}
                placeholder="0.00"
                required
              />
            </div>
          )}
          <div className="mb-3">
            <label className="block text-xs mb-1" style={{ color: '#848e9c' }}>Quantity</label>
            <input
              type="number"
              step="0.0001"
              value={quantity}
              onChange={e => setQuantity(e.target.value)}
              className="w-full px-3 py-2.5 rounded-lg text-sm outline-none font-mono transition-all duration-200"
              style={{
                background: '#0b0e11',
                border: '1px solid #2b3139',
                color: '#eaecef',
              }}
              onFocus={(e) => { e.currentTarget.style.borderColor = '#fcd535'; }}
              onBlur={(e) => { e.currentTarget.style.borderColor = '#2b3139'; }}
              placeholder="0.0000"
              required
            />
          </div>
          {type === 'LIMIT' && total > 0 && (
            <div className="mb-3 p-2 rounded-lg text-xs flex justify-between" style={{ background: '#0b0e11' }}>
              <span style={{ color: '#848e9c' }}>Total</span>
              <span className="font-mono font-medium">{total.toFixed(2)} USDT</span>
            </div>
          )}
          {error && (
            <div className="mb-3 p-2.5 rounded-lg text-xs" style={{ background: 'rgba(246,70,93,0.1)', color: '#f6465d' }}>
              {error}
            </div>
          )}
          {success && (
            <div className="mb-3 p-2.5 rounded-lg text-xs" style={{ background: 'rgba(14,203,129,0.1)', color: '#0ecb81' }}>
              {success}
            </div>
          )}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 rounded-lg font-semibold text-sm transition-all duration-200 mt-auto"
            style={{
              background: side === 'BUY'
                ? 'linear-gradient(135deg, #0ecb81, #0db574)'
                : 'linear-gradient(135deg, #f6465d, #e53950)',
              color: '#fff',
              opacity: loading ? 0.7 : 1,
              boxShadow: side === 'BUY'
                ? '0 4px 12px rgba(14,203,129,0.3)'
                : '0 4px 12px rgba(246,70,93,0.3)',
            }}
          >
            {loading ? 'Placing...' : `${side} ${symbol.split('/')[0]}`}
          </button>
        </form>

        <div className="p-4 border-t" style={{ borderColor: '#2b3139' }}>
          <div className="text-xs font-semibold mb-3 flex items-center justify-between" style={{ color: '#848e9c' }}>
            <span>Order Book</span>
            <span className="text-xs" style={{ color: '#eaecef' }}>
              {orderBook?.asks?.[0]?.price?.toLocaleString() || '--'}
            </span>
          </div>
          <div className="mb-2">
            {orderBook?.asks?.slice(0, 6).reverse().map((ask: any, i: number) => {
              const maxQty = Math.max(...(orderBook?.asks?.map((a: any) => a.quantity) || [1]));
              const width = (ask.quantity / maxQty) * 100;
              return (
                <div key={i} className="relative flex justify-between text-xs py-0.5 px-1">
                  <div
                    className="absolute right-0 top-0 bottom-0 rounded-sm"
                    style={{ background: 'rgba(246,70,93,0.08)', width: `${width}%` }}
                  />
                  <span className="relative font-mono" style={{ color: '#f6465d' }}>{ask.price.toLocaleString()}</span>
                  <span className="relative font-mono" style={{ color: '#848e9c' }}>{ask.quantity}</span>
                </div>
              );
            })}
          </div>
          <div className="text-center py-1.5 mb-2 border-y" style={{ borderColor: '#2b3139' }}>
            <span className="text-base font-bold font-mono" style={{ color: orderBook?.asks?.[0] ? '#f6465d' : '#848e9c' }}>
              {orderBook?.asks?.[0]?.price?.toLocaleString() || '--'}
            </span>
          </div>
          <div>
            {orderBook?.bids?.slice(0, 6).map((bid: any, i: number) => {
              const maxQty = Math.max(...(orderBook?.bids?.map((b: any) => b.quantity) || [1]));
              const width = (bid.quantity / maxQty) * 100;
              return (
                <div key={i} className="relative flex justify-between text-xs py-0.5 px-1">
                  <div
                    className="absolute right-0 top-0 bottom-0 rounded-sm"
                    style={{ background: 'rgba(14,203,129,0.08)', width: `${width}%` }}
                  />
                  <span className="relative font-mono" style={{ color: '#0ecb81' }}>{bid.price.toLocaleString()}</span>
                  <span className="relative font-mono" style={{ color: '#848e9c' }}>{bid.quantity}</span>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
