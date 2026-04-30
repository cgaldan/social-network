import { getToken } from "@/lib/auth";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8000/ws";

export type WebSocketHandlers = {
  onMessage?: (message: unknown) => void;
  onOpen?: () => void;
  onClose?: () => void;
  onError?: () => void;
};

export function connectWebSocket(handlers: WebSocketHandlers = {}) {
  const token = getToken();
  if (!token) {
    return null;
  }

  const socket = new WebSocket(`${WS_URL}?token=${encodeURIComponent(token)}`);

  socket.onopen = () => handlers.onOpen?.();
  socket.onclose = () => handlers.onClose?.();
  socket.onerror = () => handlers.onError?.();
  socket.onmessage = (event) => {
    try {
      handlers.onMessage?.(JSON.parse(event.data));
    } catch {
      handlers.onMessage?.(event.data);
    }
  };

  return socket;
}
