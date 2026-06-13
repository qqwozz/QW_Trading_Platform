import type { OrderBook as OrderBookType } from '../api/market'

interface OrderBookProps {
  orderBook: OrderBookType
}

export default function OrderBook({ orderBook }: OrderBookProps) {
  const maxQty = Math.max(
    ...orderBook.asks.map(([, q]) => q),
    ...orderBook.bids.map(([, q]) => q),
    1,
  )

  // Display asks reversed so lowest ask is at bottom
  const displayAsks = [...orderBook.asks].reverse().slice(0, 15)
  const displayBids = orderBook.bids.slice(0, 15)

  return (
    <div className="order-book">
      <div className="order-book-header">
        <span>Price (USDT)</span>
        <span>Quantity</span>
        <span>Total</span>
      </div>
      <div className="order-book-asks">
        {displayAsks.map(([price, qty], i) => {
          const depthPct = (qty / maxQty) * 100
          const total = displayAsks.slice(0, i + 1).reduce((sum, [, q]) => sum + q, 0)
          return (
            <div key={`ask-${price}-${i}`} className="order-book-row ask">
              <div className="depth-bar ask-bar" style={{ width: `${depthPct}%` }} />
              <span className="ob-price">{formatPrice(price)}</span>
              <span className="ob-qty">{qty.toFixed(6)}</span>
              <span className="ob-total">{total.toFixed(6)}</span>
            </div>
          )
        })}
      </div>
      <div className="order-book-spread">
        <span className="spread-label">Spread</span>
        <span className="spread-value">
          {orderBook.asks.length && orderBook.bids.length
            ? formatPrice(orderBook.asks[0][0] - orderBook.bids[0][0])
            : '--'}
        </span>
      </div>
      <div className="order-book-bids">
        {displayBids.map(([price, qty], i) => {
          const depthPct = (qty / maxQty) * 100
          const total = displayBids.slice(0, i + 1).reduce((sum, [, q]) => sum + q, 0)
          return (
            <div key={`bid-${price}-${i}`} className="order-book-row bid">
              <div className="depth-bar bid-bar" style={{ width: `${depthPct}%` }} />
              <span className="ob-price">{formatPrice(price)}</span>
              <span className="ob-qty">{qty.toFixed(6)}</span>
              <span className="ob-total">{total.toFixed(6)}</span>
            </div>
          )
        })}
      </div>
    </div>
  )
}

function formatPrice(price: number): string {
  if (price >= 1000) return price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  if (price >= 1) return price.toFixed(4)
  return price.toFixed(6)
}
