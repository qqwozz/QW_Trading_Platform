import { api } from './client'

export interface Ticker {
  symbol: string
  last_price: number
  best_bid: number
  best_ask: number
  volume_24h: number
  high_24h: number
  low_24h: number
  change_24h: number
  change_pct_24h: number
}

export interface OrderBook {
  bids: [number, number][]
  asks: [number, number][]
}

export async function getTickers(): Promise<Ticker[]> {
  return api.get<Ticker[]>('/market/tickers')
}

export async function getTicker(symbol: string): Promise<Ticker> {
  return api.get<Ticker>(`/market/tickers/${encodeURIComponent(symbol)}`)
}

export async function getOrderBook(symbol: string, depth = 20): Promise<OrderBook> {
  return api.get<OrderBook>(`/market/orderbook/${encodeURIComponent(symbol)}?depth=${depth}`)
}
