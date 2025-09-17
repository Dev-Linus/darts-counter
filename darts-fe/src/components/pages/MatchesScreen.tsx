import Header from "../common/Header";
import DeleteButton from "../common/DeleteButton";
import { useState } from "react";
import type { ApiClient } from "../../lib/api";
import { useMatches } from "../../hooks/useMatches";
import type { Player } from "../../types";
import PlayScreen from "./PlayScreen";

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

  const [playMid, setPlayMid] = useState<string | null>(null);

  if (playMid) {
    return (
      <PlayScreen
        api={api}
        players={players}
        matches={matches}
        matchId={playMid}
        onClose={() => setPlayMid(null)}
        onRefreshMatches={refresh}
      />
    );
  }

  return (
    <div className="flex flex-col min-h-screen">
      <Header title="Alle Spiele" />
      <div className="flex-1 overflow-y-auto px-6 pt-4 pb-32">
        <div
          className="
          grid gap-6 w-full max-w-none
          [grid-template-columns:repeat(auto-fit,minmax(260px,1fr))]
        "
        >
          {matches.map((m) => (
            <div
              key={m.id}
              className="bg-zinc-900 border border-zinc-800 rounded-2xl p-4 flex flex-col gap-3"
            >
              <div className="flex items-center gap-2">
                <div className="text-xs opacity-70">Match</div>
                <div className="font-mono text-sm truncate max-w-[140px]">
                  {m.id}
                </div>
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
                  <div
                    key={pid}
                    className={`rounded-xl p-3 ${
                      pid === m.currentPlayer
                        ? "bg-green-900/30 border border-green-700"
                        : "bg-zinc-800/60"
                    }`}
                  >
                    <div className="text-sm opacity-80">{nameOf(pid)}</div>
                    <div className="text-2xl font-extrabold">
                      {m.scores?.[pid] ?? 0}
                    </div>
                  </div>
                ))}
              </div>

              <div className="flex">
                <button
                  onClick={() => setPlayMid(m.id)}
                  className="ml-auto px-3 py-2 rounded-xl bg-green-700 hover:bg-green-700/90"
                >
                  spielen
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
