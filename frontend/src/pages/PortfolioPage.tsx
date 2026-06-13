import { useState, useEffect } from 'react'
import { getPortfolio, type Portfolio, type Position, type Balance } from '../api/portfolio'
import { getBalances } from '../api/portfolio'

export default function PortfolioPage() {
  const [portfolio, setPortfolio] = useState<Portfolio | null>(null)
  const [balances, setBalances] = useState<Balance[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([getPortfolio(), getBalances()])
      .then(([p, b]) => {
        setPortfolio(p)
        setBalances(b)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="page-loading">Loading portfolio...</div>

  return (
    <div className="portfolio-page">
      <h1 className="page-title">Portfolio</h1>

      {portfolio && (
        <div className="portfolio-summary">
          <div className="summary-card">
            <span className="summary-label">Total Equity</span>
            <span className="summary-value large">
              {portfolio.total_equity.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
            </span>
          </div>
          <div className="summary-card">
            <span className="summary-label">Total Balance</span>
            <span className="summary-value">
              {portfolio.total_balance.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
            </span>
          </div>
          <div className="summary-card">
            <span className="summary-label">Market Value</span>
            <span className="summary-value">
              {portfolio.total_market_value.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
            </span>
          </div>
          <div className="summary-card">
            <span className="summary-label">Unrealized PnL</span>
            <span className={`summary-value ${portfolio.total_unrealized_pnl >= 0 ? 'positive' : 'negative'}`}>
              {portfolio.total_unrealized_pnl >= 0 ? '+' : ''}
              {portfolio.total_unrealized_pnl.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
            </span>
          </div>
          <div className="summary-card">
            <span className="summary-label">Frozen</span>
            <span className="summary-value">
              {portfolio.total_frozen.toLocaleString('en-US', { minimumFractionDigits: 2 })} USDT
            </span>
          </div>
        </div>
      )}

      <div className="portfolio-sections">
        <div className="section">
          <h2 className="section-title">Positions</h2>
          <div className="table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Symbol</th>
                  <th>Qty</th>
                  <th>Avg Price</th>
                  <th>Market Price</th>
                  <th>PnL</th>
                  <th>Frozen</th>
                </tr>
              </thead>
              <tbody>
                {!portfolio?.positions?.length ? (
                  <tr><td colSpan={6} className="empty-state">No positions</td></tr>
                ) : (
                  portfolio.positions.map((p: Position) => (
                    <tr key={p.symbol}>
                      <td>{p.symbol}</td>
                      <td className="mono">{p.quantity.toFixed(6)}</td>
                      <td className="mono">{p.avg_price.toFixed(2)}</td>
                      <td className="mono">{p.market_price.toFixed(2)}</td>
                      <td className={`mono ${p.unrealized_pnl >= 0 ? 'positive' : 'negative'}`}>
                        {p.unrealized_pnl >= 0 ? '+' : ''}{p.unrealized_pnl.toFixed(2)}
                      </td>
                      <td className="mono">{p.frozen_quantity.toFixed(6)}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>

        <div className="section">
          <h2 className="section-title">Balances</h2>
          <div className="table-container">
            <table className="data-table">
              <thead>
                <tr>
                  <th>Currency</th>
                  <th>Available</th>
                  <th>Frozen</th>
                  <th>Total</th>
                </tr>
              </thead>
              <tbody>
                {balances.length === 0 ? (
                  <tr><td colSpan={4} className="empty-state">No balances</td></tr>
                ) : (
                  balances.map((b: Balance) => (
                    <tr key={b.currency}>
                      <td>{b.currency}</td>
                      <td className="mono">{b.available.toFixed(8)}</td>
                      <td className="mono">{b.frozen.toFixed(8)}</td>
                      <td className="mono">{b.total.toFixed(8)}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}
