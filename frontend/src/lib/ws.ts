import type { WsEnvelope } from "@/types/api";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8000/ws";

export type WsListener = (msg: WsEnvelope) => void;

export interface WsClient {
  close: () => void;
  send: (msg: WsEnvelope) => void;
  on: (listener: WsListener) => () => void;
}

export function connectWs(token: string): WsClient {
  const url = `${WS_URL}?token=${encodeURIComponent(token)}`;
  const listeners = new Set<WsListener>();
  let socket: WebSocket | null = null;
  let pingInterval: ReturnType<typeof setInterval> | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let attempt = 0;
  let closed = false;

  const clearPing = () => {
    if (pingInterval) {
      clearInterval(pingInterval);
      pingInterval = null;
    }
  };

  const scheduleReconnect = () => {
    if (closed) return;
    const delay = Math.min(30000, 1000 * Math.pow(2, attempt));
    attempt += 1;
    reconnectTimer = setTimeout(connect, delay);
  };

  function connect() {
    if (closed) return;
    reconnectTimer = null;
    const ws = new WebSocket(url);
    socket = ws;

    ws.addEventListener("open", () => {
      attempt = 0;
      pingInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({ type: "ping", payload: {} }));
        }
      }, 30000);
    });

    ws.addEventListener("message", (ev) => {
      try {
        const data = JSON.parse(ev.data) as WsEnvelope;
        listeners.forEach((l) => l(data));
      } catch {
        // ignore non-JSON
      }
    });

    ws.addEventListener("close", () => {
      clearPing();
      if (socket === ws) socket = null;
      scheduleReconnect();
    });

    ws.addEventListener("error", () => {
      // close event will follow; reconnect handled there
    });
  }

  connect();

  return {
    close: () => {
      closed = true;
      clearPing();
      if (reconnectTimer) {
        clearTimeout(reconnectTimer);
        reconnectTimer = null;
      }
      if (socket) {
        const s = socket;
        socket = null;
        if (s.readyState === WebSocket.CONNECTING) {
          s.addEventListener("open", () => s.close(), { once: true });
        } else {
          s.close();
        }
      }
    },
    send: (msg) => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(msg));
      }
    },
    on: (listener) => {
      listeners.add(listener);
      return () => {
        listeners.delete(listener);
      };
    },
  };
}
