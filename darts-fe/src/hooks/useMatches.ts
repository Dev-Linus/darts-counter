import { useEffect, useState } from "react";
import type { Match } from "../types";
import type { ApiClient } from "../lib/api";

export function useMatches(api: ApiClient) {
  const [matches, setMatches] = useState<Match[]>([]);
  const refresh = async () => setMatches(await api.call<Match[]>("/listMatches"));
  useEffect(() => { refresh().catch(() => {}); }, []);
  return { matches, refresh };
}