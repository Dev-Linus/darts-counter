import type { ApiLog } from "../types";


const API_BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

export type ApiClient = {
  call: <T>(path: string, options?: RequestInit) => Promise<T>;
  logs: ApiLog[];
  lastError: string | null;
  clearError: () => void;
};

export function createApi(): ApiClient {
  let logs: ApiLog[] = [];
  let lastError: string | null = null;

  const notify = () => {
    listeners.forEach((l) => l());
  };
  const listeners = new Set<() => void>();
  const subscribe = (cb: () => void) => { listeners.add(cb); return () => listeners.delete(cb); };

  const api: ApiClient & { subscribe: (cb: () => void) => () => void } = {
    async call<T>(path: string, options?: RequestInit): Promise<T> {
      const url = `${API_BASE}${path}`;
      const reqBody = options?.body;
      const entry: ApiLog = {
        time: new Date().toISOString(),
        request: { method: options?.method || "GET", url, body: safeParse(reqBody) },
      };
      logs = [entry, ...logs].slice(0, 100);
      lastError = null;
      notify();

      try {
        const res = await fetch(url, {
          ...options,
          headers: { "Content-Type": "application/json", ...(options?.headers || {}) },
        });
        const text = await res.text();
        let body: any;
        try { body = text ? JSON.parse(text) : undefined; } catch { body = text; }
        const response = { status: res.status, body };
        logs = [{ ...entry, response }, ...logs.slice(1)];
        notify();

        if (!res.ok) {
          const msg = typeof body === "string" ? body : JSON.stringify(body);
          lastError = `${res.status}: ${msg}`;
          notify();
          throw new Error(msg || `HTTP ${res.status}`);
        }
        return body as T;
      } catch (e: any) {
        logs = [{ ...entry, error: e?.message || String(e) }, ...logs.slice(1)];
        lastError = e?.message || String(e);
        notify();
        throw e;
      }
    },
    get logs() { return logs; },
    get lastError() { return lastError; },
    clearError() { lastError = null; notify(); },
    subscribe,
  };

  // small pub/sub hook
  (api as any).use = function useApiState() {
    const [, setTick] = useState(0);
    useEffect(() => api.subscribe(() => setTick((x) => x + 1)), []);
    return api as ApiClient;
  };

  return api;
}

function safeParse(x: any) {
  try { return typeof x === "string" ? JSON.parse(x) : x; } catch { return x; }
}

// local React without circular deps
import { useEffect, useState } from "react";