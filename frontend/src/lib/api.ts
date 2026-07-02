const API_BASE = '/v1';

let accessToken = localStorage.getItem('access_token');

export function setToken(token: string) {
  accessToken = token;
  localStorage.setItem('access_token', token);
}

export function clearToken() {
  accessToken = null;
  localStorage.removeItem('access_token');
}

export function getToken() {
  return accessToken;
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };
  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.error || 'Request failed');
  }
  return data;
}

export const api = {
  auth: {
    register: (email: string, username: string, password: string) =>
      request('/auth/register', { method: 'POST', body: JSON.stringify({ email, username, password }) }),
    login: (email: string, password: string) =>
      request<{ access_token: string; refresh_token: string; expires_in: number }>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      }),
    guest: (email: string, username: string, password: string) =>
      request<{ access_token: string; refresh_token: string; expires_in: number }>('/auth/guest', {
        method: 'POST',
        body: JSON.stringify({ email, username, password }),
      }),
    me: () => request<{ id: string; email: string; username: string; created_at: string }>('/users/me'),
  },
  accounts: {
    list: () => request<{ data: { accounts: Array<{ id: string; type: string; balance: number; frozen_balance: number; currency: string; status: string }> } }>('/accounts'),
    deposit: (currency: string, amount: number) =>
      request('/accounts/deposit', { method: 'POST', body: JSON.stringify({ currency, amount }) }),
    withdraw: (currency: string, amount: number) =>
      request('/accounts/withdraw', { method: 'POST', body: JSON.stringify({ currency, amount }) }),
  },
  orders: {
    create: (order: { symbol: string; side: string; type: string; price?: number; quantity: number; time_in_force?: string }) =>
      request('/orders', { method: 'POST', body: JSON.stringify(order) }),
    list: (params?: string) => request(`/orders${params ? '?' + params : ''}`),
    get: (id: string) => request(`/orders/${id}`),
    cancel: (id: string) => request(`/orders/${id}`, { method: 'DELETE' }),
  },
  portfolio: {
    get: () => request('/portfolio'),
    positions: () => request('/positions'),
    balances: () => request('/balances'),
  },
  market: {
    tickers: () => request('/market/tickers'),
    ticker: (symbol: string) => request(`/market/tickers/${symbol}`),
    orderbook: (symbol: string, depth?: number) =>
      request(`/market/orderbook/${symbol}${depth ? '?depth=' + depth : ''}`),
  },
  history: {
    orders: (params?: string) => request(`/history/orders${params ? '?' + params : ''}`),
    trades: (params?: string) => request(`/history/trades${params ? '?' + params : ''}`),
  },
};
