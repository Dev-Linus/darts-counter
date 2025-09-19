import { useEffect, useState } from "react";
import type { ApiClient } from "../../lib/api";
import type { Match, MatchHistory, Player } from "../../types";
import { THROW_TYPE_OPTIONS } from "../../lib/utils";
import { ArrowLeft } from "lucide-react";
import Dartboard from "../common/Dartboard";

export default function PlayScreen({
  api,
  players,
  matchId,
  onClose,
  onRefreshMatches
}: {
  api: ApiClient;
  players: Player[];
  matches: Match[];
  matchId: string;
  onClose: () => void;
  onRefreshMatches: () => void;
}) {
  const [matchHistory, setMatchHistory] = useState<MatchHistory | null>(null);
  const [scores, setScores] = useState<Record<string, number>>({});
  const [currentPid, setCurrentPid] = useState<string | undefined>();
  const [turnThrows, setTurnThrows] = useState<number[]>([]);
  const [finishes, setFinishes] = useState<number[]>([]);
  const [useBoardView, setUseBoardView] = useState(false);

  const nameOf = (id: string) => players.find((p) => p.id === id)?.name || id;

  useEffect(() => {
    const loadMatchHistory = async () => {
      const resp = await api.call<MatchHistory>(
        `/getMatch?matchId=${matchId}`,
        { method: "GET" }
      );
      setMatchHistory(resp);
      setScores(resp.match.scores || {});
      setCurrentPid(resp.match.currentPlayer);
    };
    if (matchId) loadMatchHistory();
  }, [matchId]);

  const throwOnce = async (tt: number) => {
    if (!currentPid) return;

    const resp = await api.call<{
      Won: boolean;
      NotValid: boolean;
      NextThrowBy: string;
      Scores: Record<string, number>;
      PossibleFinish: number[];
    }>("/playerThrow", {
      method: "POST",
      body: JSON.stringify({ Mid: matchId, Pid: currentPid, Throw: tt })
    });

    setScores(resp.Scores || scores);
    setFinishes(resp.PossibleFinish || []);

    if (resp.NextThrowBy !== currentPid) {
      setTurnThrows([]);
    } else {
      setTurnThrows((prev) => [...prev, tt].slice(-3));
    }

    setCurrentPid(resp.NextThrowBy);
    onRefreshMatches();

    // refresh match + history after every throw
    const historyResp = await api.call<MatchHistory>(
      `/getMatch?matchId=${matchId}`,
      { method: "GET" }
    );
    setMatchHistory(historyResp);
  };

  const finishLabels = finishes
    .map(
      (v) => THROW_TYPE_OPTIONS.find((o) => o.value === v)?.label || String(v)
    )
    .slice(0, 3);

  const currentMatch = matchHistory?.match;

  return (
    <div className="px-4 pb-24">
      <div className="flex items-center gap-2 py-2">
        <button
          onClick={onClose}
          className="p-2 rounded-xl bg-zinc-800 hover:bg-zinc-700"
          aria-label="Zurück"
        >
          <ArrowLeft />
        </button>
        <h1 className="text-2xl font-extrabold ml-2">Spiel</h1>

        {/* Toggle button */}
        <button
          onClick={() => setUseBoardView((v) => !v)}
          className="ml-auto px-3 py-2 text-sm rounded-lg bg-blue-800 hover:bg-blue-700"
        >
          {useBoardView ? "Zahlen-Eingabe" : "Board-Eingabe"}
        </button>
      </div>

      <div className="grid md:grid-cols-[220px,1fr] gap-6 mt-4">
        {/* LEFT: Players */}
        <div className="space-y-2">
          {currentMatch?.players.map((pid) => {
            const throws = matchHistory?.history?.[pid] ?? [];

            return (
              <div key={pid} className="space-y-2">
                <div
                  className={`flex items-center justify-between rounded-xl px-3 py-2 border ${
                    pid === currentPid
                      ? "bg-green-900/30 border-green-700"
                      : "bg-zinc-900 border-zinc-800"
                  }`}
                >
                  <div className="font-semibold truncate mr-3">
                    {nameOf(pid)}
                  </div>
                  <div className="text-xl font-extrabold">
                    {scores?.[pid] ?? 0}
                  </div>
                </div>

                {/* Throw history grouped by turn */}
                {throws.length > 0 && (
                  <div className="flex flex-col gap-1 pl-1">
                    {Object.values(
                      throws.reduce((acc, t) => {
                        if (!acc[t.turnNumber]) acc[t.turnNumber] = [];
                        acc[t.turnNumber].push(t);
                        return acc;
                      }, {} as Record<number, typeof throws>)
                    ).map((turn, ti) => (
                      <div key={ti} className="flex gap-1">
                        {turn.map((t, i) => {
                          const label =
                            THROW_TYPE_OPTIONS.find((o) => o.value === t.throw)
                              ?.label ?? String(t.throw);
                          return (
                            <span
                              key={i}
                              className="px-2 py-1 rounded-lg bg-zinc-800 text-xs border border-zinc-700"
                            >
                              {label}
                            </span>
                          );
                        })}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {/* RIGHT: Board or Grid */}
        <div className="space-y-6">
          {useBoardView ? (
            <div className="flex justify-center items-center">
              <Dartboard onPick={(tt) => throwOnce(tt)} />
            </div>
          ) : (
            <div
              className="
          grid gap-1
          [grid-template-columns:repeat(auto-fit,minmax(1fr))]
        "
            >
              {THROW_TYPE_OPTIONS.filter((o) => o.value <= 60).map((opt) => (
                <button
                  key={opt.value}
                  className={`px-3 py-3 rounded-lg border text-base font-bold hover:opacity-90 ${
                    opt.value <= 20
                      ? "bg-zinc-800 border-zinc-600"
                      : opt.value <= 40
                      ? "bg-red-900/60 border-red-700"
                      : "bg-green-900/50 border-green-700"
                  }`}
                  onClick={() => throwOnce(opt.value)}
                >
                  {opt.label}
                </button>
              ))}
              {/* Bulls */}
              {THROW_TYPE_OPTIONS.filter((o) => o.value > 60).map((opt) => (
                <button
                  key={opt.value}
                  className="col-span-2 sm:col-span-4 md:col-span-5 px-3 py-4 rounded-lg bg-yellow-900/60 border border-yellow-700 text-lg font-bold"
                  onClick={() => throwOnce(opt.value)}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          )}

          {/* Throws of this turn */}
          <div className="grid grid-cols-3 gap-3">
            {[0, 1, 2].map((i) => {
              const val = turnThrows[i];
              const label = val
                ? THROW_TYPE_OPTIONS.find((o) => o.value === val)?.label
                : "-";
              return (
                <div
                  key={i}
                  className="rounded-xl bg-zinc-900 border border-zinc-800 px-3 py-3 text-center"
                >
                  <div className="text-sm opacity-80">Wurf {i + 1}</div>
                  <div className="mt-1 font-extrabold">{label}</div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
