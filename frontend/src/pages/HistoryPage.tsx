import { useState, useEffect, useCallback } from 'react'
import {
  getOrderHistory,
  getTradeHistory,
  getBalanceHistory,
  type OrderHistoryItem,
  type TradeHistoryItem,
  type BalanceHistoryItem,
} from '../api/history'

type Tab = 'orders' | 'trades' | 'balance'

export default function HistoryPage() {
  const [activeTab, setActiveTab] = useState<Tab>('orders')
  const [orders, setOrders] = useState<OrderHistoryItem[]>([])
  const [trades, setTrades] = useState<TradeHistoryItem[]>([])
  const [balances, setBalances] = useState<BalanceHistoryItem[]>([])
  const [loading, setLoading] = useState(false)

  const fetchData = useCallback(async (tab: Tab) => {
    setLoading(true)
    try {
      if (tab === 'orders') {
        const res = await getOrderHistory({ limit: 100 })
        setOrders(Array.isArray(res.data) ? res.data : [])
      } else if (tab === 'trades') {
        const res = await getTradeHistory({ limit: 100 })
        setTrades(Array.isArray(res.data) ? res.data : [])
      } else {
        const res = await getBalanceHistory({ limit: 100 })
        setBalances(Array.isArray(res.data) ? res.data : [])
      }
    } catch {
      // ignore
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchData(activeTab)
  }, [activeTab, fetchData])

  return (
    <div className="history-page">
      <h1 className="page-title">History</h1>

      <div className="history-tabs">
        <button
          className={`tab-btn${activeTab === 'orders' ? ' active' : ''}`}
          onClick={() => setActiveTab('orders')}
        >
          Order History
        </button>
        <button
          className={`tab-btn${activeTab === 'trades' ? ' active' : ''}`}
          onClick={() => setActiveTab('trades')}
        >
          Trade History
        </button>
        <button
          className={`tab-btn${activeTab === 'balance' ? ' active' : ''}`}
          onClick={() => setActiveTab('balance')}
        >
          Balance History
        </button>
      </div>

      {loading ? (
        <div className="page-loading">Loading...</div>
      ) : (
        <div className="table-container">
          {activeTab === 'orders' && (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Time</th>
                  <th>Symbol</th>
                  <th>Side</th>
                  <th>Type</th>
                  <th>Price</th>
                  <th>Quantity</th>
                  <th>Filled</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                {orders.length === 0 ? (
                  <tr><td colSpan={8} className="empty-state">No order history</td></tr>
                ) : (
                  orders.map((o) => (
                    <tr key={o.id}>
                      <td className="mono">{new Date(o.created_at).toLocaleString()}</td>
                      <td>{o.symbol}</td>
                      <td className={`side-${o.side}`}>{o.side}</td>
                      <td>{o.type}</td>
                      <td className="mono">{o.price?.toFixed(2) ?? 'Market'}</td>
                      <td className="mono">{o.quantity.toFixed(6)}</td>
                      <td className="mono">{o.filled_quantity.toFixed(6)}</td>
                      <td>{o.status}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          )}

          {activeTab === 'trades' && (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Time</th>
                  <th>Symbol</th>
                  <th>Side</th>
                  <th>Price</th>
                  <th>Quantity</th>
                  <th>Fee</th>
                  <th>Fee Currency</th>
                </tr>
              </thead>
              <tbody>
                {trades.length === 0 ? (
                  <tr><td colSpan={7} className="empty-state">No trade history</td></tr>
                ) : (
                  trades.map((t) => (
                    <tr key={t.id}>
                      <td className="mono">{new Date(t.created_at).toLocaleString()}</td>
                      <td>{t.symbol}</td>
                      <td className={`side-${t.side}`}>{t.side}</td>
                      <td className="mono">{t.price.toFixed(2)}</td>
                      <td className="mono">{t.quantity.toFixed(6)}</td>
                      <td className="mono">{t.fee.toFixed(6)}</td>
                      <td>{t.fee_currency}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          )}

          {activeTab === 'balance' && (
            <table className="data-table">
              <thead>
                <tr>
                  <th>Time</th>
                  <th>Currency</th>
                  <th>Amount</th>
                  <th>Balance After</th>
                  <th>Type</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                {balances.length === 0 ? (
                  <tr><td colSpan={6} className="empty-state">No balance history</td></tr>
                ) : (
                  balances.map((b) => (
                    <tr key={b.id}>
                      <td className="mono">{new Date(b.created_at).toLocaleString()}</td>
                      <td>{b.currency}</td>
                      <td className={`mono ${b.amount >= 0 ? 'positive' : 'negative'}`}>
                        {b.amount >= 0 ? '+' : ''}{b.amount.toFixed(8)}
                      </td>
                      <td className="mono">{b.balance_after.toFixed(8)}</td>
                      <td>{b.type}</td>
                      <td>{b.description}</td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  )
}
