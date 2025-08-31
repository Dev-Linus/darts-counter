import Header from "../common/Header";
import DeleteButton from "../common/DeleteButton";
import { useState } from "react";
import type { ApiClient } from "../../lib/api";
import { useMatches } from "../../hooks/useMatches";
import type { Player } from "../../types";

export default function MatchesScreen({
  api,
  players,
  matchesState
}: {
  api: ApiClient;
  players: Player[];
  matchesState: ReturnType<typeof useMatches>;
}) {
  const { matches, refresh } = matchesState;
  const nameOf = (id: string) => players.find((p) => p.id === id)?.name || id;

  const [scoreBy, setScoreBy] = useState<
    Record<string, { playerId: string; amount: number }>
  >({});

  const submit = async (mid: string) => {
    const sb = scoreBy[mid];
    if (!sb || !sb.playerId) return;
    await api.call("/playerThrow", {
      method: "POST",
      body: JSON.stringify({
        matchId: mid,
        playerId: sb.playerId,
        amount: Number(sb.amount || 0)
      })
    });
    setScoreBy((m) => ({ ...m, [mid]: { playerId: sb.playerId, amount: 0 } }));
    refresh();
  };

  return (
    <div className="px-4 pb-24">
      <Header title="Alle Spiele" />
      <div className="space-y-4">
        {matches.map((m) => (
          <div
            key={m.id}
            className="bg-zinc-900 border border-zinc-800 rounded-2xl p-4 space-y-3"
          >
            <div className="flex items-center gap-2">
              <div className="text-xs opacity-70">Match</div>
              <div className="font-mono text-sm">{m.id}</div>
              <div className="ml-auto text-xs opacity-70">
                Start: {m.startAt}
              </div>
              <DeleteButton
                onClick={async () => {
                  await api.call(`/deleteMatch?matchId=${m.id}`, {
                    method: "DELETE"
                  });
                  refresh();
                }}
              />
            </div>
            <div className="grid grid-cols-2 gap-3">
              {m.players.map((pid) => (
                <div key={pid} className="rounded-xl p-3 bg-zinc-800/60">
                  <div className="text-sm opacity-80">{nameOf(pid)}</div>
                  <div className="text-2xl font-extrabold">
                    {m.scores?.[pid] ?? 0}
                  </div>
                </div>
              ))}
            </div>
            <div className="flex items-center gap-2">
              <select
                className="bg-zinc-800 rounded-xl px-3 py-2"
                value={scoreBy[m.id]?.playerId || ""}
                onChange={(e) =>
                  setScoreBy((s) => ({
                    ...s,
                    [m.id]: {
                      playerId: e.target.value,
                      amount: s[m.id]?.amount || 0
                    }
                  }))
                }
              >
                <option value="">Spieler wählen…</option>
                {m.players.map((pid) => (
                  <option key={pid} value={pid}>
                    {nameOf(pid)}
                  </option>
                ))}
              </select>
              <input
                type="number"
                placeholder="Punkte"
                className="w-28 px-3 py-2 rounded-xl bg-zinc-800 border border-zinc-700"
                value={scoreBy[m.id]?.amount ?? ""}
                onChange={(e) =>
                  setScoreBy((s) => ({
                    ...s,
                    [m.id]: {
                      playerId: s[m.id]?.playerId || "",
                      amount: Number(e.target.value)
                    }
                  }))
                }
              />
              <button
                onClick={() => submit(m.id)}
                className="px-3 py-2 rounded-xl bg-green-700 hover:bg-green-700/90"
              >
                Wurf speichern
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
