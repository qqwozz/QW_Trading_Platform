import { api } from './client'

export interface Portfolio {
  total_balance: number
  total_frozen: number
  total_market_value: number
  total_unrealized_pnl: number
  total_equity: number
  positions: Position[]
}

export interface Position {
  symbol: string
  quantity: number
  avg_price: number
  market_price: number
  unrealized_pnl: number
  frozen_quantity: number
}

export interface Balance {
  currency: string
  available: number
  frozen: number
  total: number
}

export async function getPortfolio(): Promise<Portfolio> {
  return api.get<Portfolio>('/portfolio')
}

export async function getBalances(): Promise<Balance[]> {
  return api.get<Balance[]>('/balances')
}

export async function getPositions(): Promise<Position[]> {
  return api.get<Position[]>('/positions')
}
