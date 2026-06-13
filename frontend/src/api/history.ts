import { api } from './client'

export interface OrderHistoryItem {
  id: string
  symbol: string
  side: 'buy' | 'sell'
  type: string
  price: number
  quantity: number
  filled_quantity: number
  status: string
  created_at: string
  updated_at: string
}

export interface TradeHistoryItem {
  id: string
  order_id: string
  symbol: string
  side: 'buy' | 'sell'
  price: number
  quantity: number
  fee: number
  fee_currency: string
  created_at: string
}

export interface BalanceHistoryItem {
  id: string
  currency: string
  amount: number
  balance_after: number
  type: string
  description: string
  created_at: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
}

export async function getOrderHistory(params?: {
  symbol?: string
  limit?: number
  offset?: number
}): Promise<PaginatedResponse<OrderHistoryItem>> {
  const query = new URLSearchParams()
  if (params?.symbol) query.set('symbol', params.symbol)
  if (params?.limit) query.set('limit', String(params.limit))
  if (params?.offset) query.set('offset', String(params.offset))
  const qs = query.toString()
  return api.get<PaginatedResponse<OrderHistoryItem>>(`/history/orders${qs ? `?${qs}` : ''}`)
}

export async function getTradeHistory(params?: {
  symbol?: string
  limit?: number
  offset?: number
}): Promise<PaginatedResponse<TradeHistoryItem>> {
  const query = new URLSearchParams()
  if (params?.symbol) query.set('symbol', params.symbol)
  if (params?.limit) query.set('limit', String(params.limit))
  if (params?.offset) query.set('offset', String(params.offset))
  const qs = query.toString()
  return api.get<PaginatedResponse<TradeHistoryItem>>(`/history/trades${qs ? `?${qs}` : ''}`)
}

export async function getBalanceHistory(params?: {
  currency?: string
  limit?: number
  offset?: number
}): Promise<PaginatedResponse<BalanceHistoryItem>> {
  const query = new URLSearchParams()
  if (params?.currency) query.set('currency', params.currency)
  if (params?.limit) query.set('limit', String(params.limit))
  if (params?.offset) query.set('offset', String(params.offset))
  const qs = query.toString()
  return api.get<PaginatedResponse<BalanceHistoryItem>>(`/history/balance${qs ? `?${qs}` : ''}`)
}
