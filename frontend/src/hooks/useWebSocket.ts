import { useEffect, useRef, useCallback, useState } from 'react'
import type { Ticker, OrderBook } from '../api/market'

type MessageHandler = (data: unknown) => void

interface UseWebSocketReturn {
  connected: boolean
  subscribe: (symbol: string) => void
  unsubscribe: (symbol: string) => void
  onTickerUpdate: (handler: (ticker: Ticker) => void) => () => void
  onOrderBookUpdate: (handler: (book: OrderBook & { symbol: string }) => void) => () => void
}

const WS_URL = 'ws://localhost:8080/market/ws'

export function useWebSocket(): UseWebSocketReturn {
  const wsRef = useRef<WebSocket | null>(null)
  const [connected, setConnected] = useState(false)
  const handlersRef = useRef<Map<string, Set<MessageHandler>>>(new Map())
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  const connect = useCallback(() => {
    const ws = new WebSocket(WS_URL)

    ws.onopen = () => {
      setConnected(true)
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data as string) as { type: string; symbol?: string }
        // Notify all handlers registered for the message type
        const handlers = handlersRef.current.get(msg.type)
        if (handlers) {
          for (const handler of handlers) {
            handler(msg)
          }
        }
        // Also notify type-specific handlers (e.g., "ticker.BTC/USDT")
        if (msg.symbol) {
          const specific = handlersRef.current.get(`${msg.type}.${msg.symbol}`)
          if (specific) {
            for (const handler of specific) {
              handler(msg)
            }
          }
        }
      } catch {
        // ignore parse errors
      }
    }

    ws.onclose = () => {
      setConnected(false)
      // Auto-reconnect after 3 seconds
      reconnectTimer.current = setTimeout(connect, 3000)
    }

    ws.onerror = () => {
      ws.close()
    }

    wsRef.current = ws
  }, [])

  useEffect(() => {
    connect()
    return () => {
      clearTimeout(reconnectTimer.current)
      wsRef.current?.close()
    }
  }, [connect])

  const send = useCallback((msg: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg))
    }
  }, [])

  const subscribe = useCallback((symbol: string) => {
    send({ type: 'subscribe', symbol })
  }, [send])

  const unsubscribe = useCallback((symbol: string) => {
    send({ type: 'unsubscribe', symbol })
  }, [send])

  const registerHandler = useCallback((type: string, handler: MessageHandler) => {
    if (!handlersRef.current.has(type)) {
      handlersRef.current.set(type, new Set())
    }
    handlersRef.current.get(type)!.add(handler)
    return () => {
      handlersRef.current.get(type)?.delete(handler)
    }
  }, [])

  const onTickerUpdate = useCallback((handler: (ticker: Ticker) => void) => {
    return registerHandler('ticker', handler as MessageHandler)
  }, [registerHandler])

  const onOrderBookUpdate = useCallback((handler: (book: OrderBook & { symbol: string }) => void) => {
    return registerHandler('orderbook', handler as MessageHandler)
  }, [registerHandler])

  return { connected, subscribe, unsubscribe, onTickerUpdate, onOrderBookUpdate }
}
