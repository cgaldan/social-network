"use client";

import { useEffect, useRef, useState } from "react";
import { getStoredToken, getWebSocketUrl } from "../../../lib/api";

export default function RealtimePage() {
  const token = getStoredToken();
  const [status, setStatus] = useState("disconnected");
  const [log, setLog] = useState([]);
  const wsRef = useRef(null);

  useEffect(() => {
    if (!token) return undefined;

    let ws;
    try {
      ws = new WebSocket(getWebSocketUrl(token));
    } catch (e) {
      setStatus("error");
      setLog((l) => [...l, `Failed: ${e.message}`]);
      return undefined;
    }

    wsRef.current = ws;
    setStatus("connecting");

    ws.onopen = () => {
      setStatus("connected");
      setLog((l) => [...l, "Connected to /ws"]);
    };

    ws.onmessage = (ev) => {
      setLog((l) => [...l.slice(-80), ev.data]);
    };

    ws.onerror = () => {
      setLog((l) => [...l, "Socket error"]);
    };

    ws.onclose = () => {
      setStatus("disconnected");
      setLog((l) => [...l, "Disconnected"]);
    };

    const ping = setInterval(() => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type: "ping", payload: null }));
      }
    }, 30000);

    return () => {
      clearInterval(ping);
      ws.close();
    };
  }, [token]);

  return (
    <div className="page-stack">
      <section className="surface-card">
        <h1>Live connection</h1>
        <p className="helper-text">
          WebSocket endpoint <code>/ws?token=…</code> — used for presence and server pushes. Status:{" "}
          <strong>{status}</strong>
        </p>
        <pre className="ws-log">{log.join("\n") || "No messages yet."}</pre>
      </section>
    </div>
  );
}
