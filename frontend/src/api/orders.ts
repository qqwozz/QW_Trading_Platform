import { api } from './client'

export interface Order {
  id: string
  symbol: string
  side: 'buy' | 'sell'
  type: 'limit' | 'market'
  price?: number
  quantity: number
  filled_quantity: number
  status: string
  time_in_force: string
  created_at: string
}

export interface CreateOrderParams {
  symbol: string
  side: 'buy' | 'sell'
  type: 'limit' | 'market'
  price?: number
  quantity: number
  time_in_force?: string
}

export interface OrderListResponse {
  data: Order[]
  total: number
}

export async function createOrder(params: CreateOrderParams): Promise<Order> {
  return api.post<Order>('/orders', params)
}

export async function listOrders(params?: {
  symbol?: string
  status?: string
  limit?: number
  offset?: number
}): Promise<OrderListResponse> {
  const query = new URLSearchParams()
  if (params?.symbol) query.set('symbol', params.symbol)
  if (params?.status) query.set('status', params.status)
  if (params?.limit) query.set('limit', String(params.limit))
  if (params?.offset) query.set('offset', String(params.offset))
  const qs = query.toString()
  return api.get<OrderListResponse>(`/orders${qs ? `?${qs}` : ''}`)
}

export async function cancelOrder(id: string): Promise<unknown> {
  return api.delete(`/orders/${id}`)
}
