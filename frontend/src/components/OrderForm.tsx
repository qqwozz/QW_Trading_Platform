import { useState } from 'react'
import { createOrder, type CreateOrderParams } from '../api/orders'

interface OrderFormProps {
  symbol: string
  onSuccess?: () => void
}

export default function OrderForm({ symbol, onSuccess }: OrderFormProps) {
  const [side, setSide] = useState<'buy' | 'sell'>('buy')
  const [orderType, setOrderType] = useState<'limit' | 'market'>('limit')
  const [price, setPrice] = useState('')
  const [quantity, setQuantity] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const quantityPct = quantity && price ? (parseFloat(quantity) * parseFloat(price)) : 0

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)

    const params: CreateOrderParams = {
      symbol,
      side,
      type: orderType,
      quantity: parseFloat(quantity),
    }

    if (orderType === 'limit' && price) {
      params.price = parseFloat(price)
    }

    try {
      await createOrder(params)
      setPrice('')
      setQuantity('')
      onSuccess?.()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Order failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="order-form" onSubmit={handleSubmit}>
      <div className="order-side-tabs">
        <button
          type="button"
          className={`side-tab buy${side === 'buy' ? ' active' : ''}`}
          onClick={() => setSide('buy')}
        >
          Buy
        </button>
        <button
          type="button"
          className={`side-tab sell${side === 'sell' ? ' active' : ''}`}
          onClick={() => setSide('sell')}
        >
          Sell
        </button>
      </div>

      <div className="order-type-tabs">
        <button
          type="button"
          className={`type-tab${orderType === 'limit' ? ' active' : ''}`}
          onClick={() => setOrderType('limit')}
        >
          Limit
        </button>
        <button
          type="button"
          className={`type-tab${orderType === 'market' ? ' active' : ''}`}
          onClick={() => setOrderType('market')}
        >
          Market
        </button>
      </div>

      {orderType === 'limit' && (
        <div className="form-group">
          <label>Price (USDT)</label>
          <input
            type="number"
            step="any"
            placeholder="0.00"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
            required
          />
        </div>
      )}

      <div className="form-group">
        <label>Quantity</label>
        <input
          type="number"
          step="any"
          placeholder="0.000000"
          value={quantity}
          onChange={(e) => setQuantity(e.target.value)}
          required
        />
      </div>

      <div className="quantity-slider">
        {[25, 50, 75, 100].map((pct) => (
          <button
            key={pct}
            type="button"
            className="pct-btn"
            onClick={() => {
              // Slider fills percentage of available balance (placeholder logic)
              const qty = quantity ? parseFloat(quantity) : 0
              setQuantity(String((qty * pct) / 100 || ''))
            }}
          >
            {pct}%
          </button>
        ))}
      </div>

      {orderType === 'limit' && price && quantity && (
        <div className="order-summary">
          <span>Total</span>
          <span className="summary-value">{quantityPct.toFixed(2)} USDT</span>
        </div>
      )}

      {error && <div className="order-error">{error}</div>}

      <button
        type="submit"
        className={`btn-submit ${side}`}
        disabled={loading}
      >
        {loading ? 'Placing...' : `${side === 'buy' ? 'Buy' : 'Sell'} ${symbol.split('/')[0]}`}
      </button>
    </form>
  )
}
