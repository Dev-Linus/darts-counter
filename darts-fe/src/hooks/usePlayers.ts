import { useEffect, useState } from "react";
import type { Player } from "../types";
import type { ApiClient } from "../lib/api";

export function usePlayers(api: ApiClient) {
  const [players, setPlayers] = useState<Player[]>([]);
  const refresh = async () => setPlayers(await api.call<Player[]>("/listPlayers"));
  useEffect(() => { refresh().catch(() => {}); }, []);
  return { players, refresh };
}