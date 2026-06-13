import { useState, useEffect, useCallback } from 'react'
import { getTickers, getOrderBook, type Ticker, type OrderBook } from '../api/market'
import { listOrders, cancelOrder, type Order } from '../api/orders'
import TickerBar from '../components/TickerBar'
import OrderBookComponent from '../components/OrderBook'
import OrderForm from '../components/OrderForm'

export default function TradingPage() {
  const [selectedSymbol, setSelectedSymbol] = useState('BTC/USDT')
  const [tickers, setTickers] = useState<Ticker[]>([])
  const [orderBook, setOrderBook] = useState<OrderBook>({ bids: [], asks: [] })
  const [openOrders, setOpenOrders] = useState<Order[]>([])
  const [loadingOrders, setLoadingOrders] = useState(false)

  // Fetch tickers periodically
  useEffect(() => {
    const fetchTickers = () => getTickers().then(setTickers).catch(() => {})
    fetchTickers()
    const interval = setInterval(fetchTickers, 5000)
    return () => clearInterval(interval)
  }, [])

  // Fetch order book for selected symbol
  useEffect(() => {
    const fetchBook = () => getOrderBook(selectedSymbol).then(setOrderBook).catch(() => {})
    fetchBook()
    const interval = setInterval(fetchBook, 2000)
    return () => clearInterval(interval)
  }, [selectedSymbol])

  const fetchOrders = useCallback(() => {
    setLoadingOrders(true)
    listOrders({ status: 'open', limit: 50 })
      .then((res) => setOpenOrders(Array.isArray(res.data) ? res.data : []))
      .catch(() => setOpenOrders([]))
      .finally(() => setLoadingOrders(false))
  }, [])

  useEffect(() => {
    fetchOrders()
  }, [fetchOrders])

  async function handleCancelOrder(id: string) {
    try {
      await cancelOrder(id)
      fetchOrders()
    } catch {
      // ignore
    }
  }

  const currentTicker = tickers.find((t) => t.symbol === selectedSymbol)

  return (
    <div className="trading-page">
      <TickerBar />

      <div className="trading-header">
        <div className="symbol-selector">
          {tickers.map((t) => (
            <button
              key={t.symbol}
              className={`symbol-btn${t.symbol === selectedSymbol ? ' active' : ''}`}
              onClick={() => setSelectedSymbol(t.symbol)}
            >
              {t.symbol}
            </button>
          ))}
        </div>
        {currentTicker && (
          <div className="current-ticker-info">
            <span className={`big-price ${currentTicker.change_pct_24h >= 0 ? 'positive' : 'negative'}`}>
              {currentTicker.last_price.toLocaleString('en-US', { minimumFractionDigits: 2 })}
            </span>
            <span className={`change ${currentTicker.change_pct_24h >= 0 ? 'positive' : 'negative'}`}>
              {currentTicker.change_pct_24h >= 0 ? '+' : ''}{currentTicker.change_pct_24h.toFixed(2)}%
            </span>
            <div className="ticker-details">
              <div>
                <span className="detail-label">24h High</span>
                <span className="detail-value">{currentTicker.high_24h.toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
              </div>
              <div>
                <span className="detail-label">24h Low</span>
                <span className="detail-value">{currentTicker.low_24h.toLocaleString('en-US', { minimumFractionDigits: 2 })}</span>
              </div>
              <div>
                <span className="detail-label">24h Volume</span>
                <span className="detail-value">{currentTicker.volume_24h.toLocaleString('en-US', { maximumFractionDigits: 2 })}</span>
              </div>
            </div>
          </div>
        )}
      </div>

      <div className="trading-grid">
        <div className="trading-left">
          <div className="panel-header">
            <span className="panel-title">Order Book</span>
          </div>
          <OrderBookComponent orderBook={orderBook} />
        </div>

        <div className="trading-center">
          {/* Chart placeholder */}
          <div className="chart-placeholder">
            <div className="chart-area">
              <span className="chart-label">Price Chart</span>
              <span className="chart-sublabel">{selectedSymbol}</span>
              {currentTicker && (
                <span className={`chart-price ${currentTicker.change_pct_24h >= 0 ? 'positive' : 'negative'}`}>
                  {currentTicker.last_price.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
                </span>
              )}
            </div>
          </div>

          {/* Open Orders Table */}
          <div className="open-orders-panel">
            <div className="panel-header">
              <span className="panel-title">Open Orders ({openOrders.length})</span>
              <button className="btn-refresh" onClick={fetchOrders} disabled={loadingOrders}>
                Refresh
              </button>
            </div>
            <div className="table-container">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Time</th>
                    <th>Symbol</th>
                    <th>Side</th>
                    <th>Type</th>
                    <th>Price</th>
                    <th>Qty</th>
                    <th>Filled</th>
                    <th>Status</th>
                    <th>Action</th>
                  </tr>
                </thead>
                <tbody>
                  {openOrders.length === 0 ? (
                    <tr><td colSpan={9} className="empty-state">No open orders</td></tr>
                  ) : (
                    openOrders.map((o) => (
                      <tr key={o.id}>
                        <td className="mono">{new Date(o.created_at).toLocaleTimeString()}</td>
                        <td>{o.symbol}</td>
                        <td className={`side-${o.side}`}>{o.side}</td>
                        <td>{o.type}</td>
                        <td className="mono">{o.price?.toFixed(2) ?? 'Market'}</td>
                        <td className="mono">{o.quantity.toFixed(6)}</td>
                        <td className="mono">{o.filled_quantity.toFixed(6)}</td>
                        <td>{o.status}</td>
                        <td>
                          <button
                            className="btn-cancel"
                            onClick={() => handleCancelOrder(o.id)}
                          >
                            Cancel
                          </button>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <div className="trading-right">
          <div className="panel-header">
            <span className="panel-title">Place Order</span>
          </div>
          <OrderForm symbol={selectedSymbol} onSuccess={fetchOrders} />
        </div>
      </div>
    </div>
  )
}
