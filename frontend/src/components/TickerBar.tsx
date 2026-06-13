import { useEffect, useState } from 'react'
import { getTickers, type Ticker } from '../api/market'

export default function TickerBar() {
  const [tickers, setTickers] = useState<Ticker[]>([])

  useEffect(() => {
    getTickers().then(setTickers).catch(() => {})
    const interval = setInterval(() => {
      getTickers().then(setTickers).catch(() => {})
    }, 5000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="ticker-bar">
      {tickers.map((t) => (
        <div key={t.symbol} className="ticker-item">
          <span className="ticker-symbol">{t.symbol}</span>
          <span className="ticker-price">{formatPrice(t.last_price)}</span>
          <span className={`ticker-change ${t.change_pct_24h >= 0 ? 'positive' : 'negative'}`}>
            {t.change_pct_24h >= 0 ? '+' : ''}{t.change_pct_24h.toFixed(2)}%
          </span>
          <span className="ticker-volume">Vol: {formatVolume(t.volume_24h)}</span>
        </div>
      ))}
    </div>
  )
}

function formatPrice(price: number): string {
  if (price >= 1000) return price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  if (price >= 1) return price.toFixed(4)
  return price.toFixed(6)
}

function formatVolume(vol: number): string {
  if (vol >= 1_000_000) return `${(vol / 1_000_000).toFixed(2)}M`
  if (vol >= 1_000) return `${(vol / 1_000).toFixed(2)}K`
  return vol.toFixed(2)
}
