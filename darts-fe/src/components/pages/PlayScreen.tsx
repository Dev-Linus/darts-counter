import { useEffect, useState } from "react";
import type { ApiClient } from "../../lib/api";
import type { Match, MatchHistory, Player, HistoryElement } from "../../types";
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
      const cur = resp.match.currentPlayer;
      setCurrentPid(cur);
      // initialize current turn throws for the active player from history (last turn)
      const hmap = resp.history?.history || {};
      const curHist = hmap[cur] ? (hmap[cur] as any[]) : [];
      const curTurn = curHist.map((h: any) => h.throw).reverse().slice(0, 3);
      setTurnThrows(curTurn);
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
            const throws = (matchHistory as any)?.history?.history?.[pid] ?? [];

            return (
              <div
                key={pid}
                className={`rounded-xl px-3 py-2 border ${
                  pid === currentPid
                    ? "bg-green-900/30 border-green-700"
                    : "bg-zinc-900 border-zinc-800"
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="font-semibold truncate mr-3">{nameOf(pid)}</div>
                  <div className="text-xl font-extrabold">{scores?.[pid] ?? 0}</div>
                </div>

                {/* Possible Finish for active player (inside the box) */}
                {pid === currentPid && finishes.length > 0 && (
                  <div className="mt-2">
                    <div className="text-xs opacity-70 mb-1">Mögliches Finish</div>
                    <div className="flex gap-1 flex-wrap">
                      {finishLabels.map((label, i) => (
                        <span
                          key={i}
                          className="px-2 py-1 rounded-lg bg-yellow-900/50 text-yellow-100 text-xs border border-yellow-700"
                        >
                          {label}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Throw history (slot grid 3xN) only for non-current players */}
                {pid !== currentPid && throws.length > 0 && (
                  <div className="mt-2">
                    {Object.values(
                      (throws as HistoryElement[]).reduce<Record<number, HistoryElement[]>>(
                        (acc: Record<number, HistoryElement[]>, t: HistoryElement) => {
                          const tn = t.turn_number;
                          if (!acc[tn]) acc[tn] = [];
                          acc[tn].push(t);
                          return acc;
                        },
                        {} as Record<number, HistoryElement[]>
                      )
                    )
                      // latest turns first
                      .sort((a, b) => (b[0]?.turn_number ?? 0) - (a[0]?.turn_number ?? 0))
                      .map((turn, ti) => {
                        const cells = [turn[0], turn[1], turn[2]] as (HistoryElement | undefined)[];
                        return (
                          <div key={ti} className="grid grid-cols-3 gap-1 mb-1">
                            {cells.map((t, i) => {
                              const isEmpty = !t;
                              const label = t
                                ? (THROW_TYPE_OPTIONS.find((o) => o.value === t.throw)?.label ?? String(t.throw))
                                : "\u00A0"; // non-breaking space to keep height without showing '-'
                              const ended = !!t && (t as HistoryElement).ended_turn;
                              return (
                                <div
                                  key={i}
                                  className={`px-2 py-1 rounded-md text-xs text-center border ${ended ? "bg-red-900/60 border-red-700 text-red-100" : "bg-zinc-800 border-zinc-700"} ${isEmpty ? "opacity-40" : ""}`}
                                >
                                  {label}
                                </div>
                              );
                            })}
                          </div>
                        );
                      })}
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
              const isEmpty = !val;
              const label = val
                ? THROW_TYPE_OPTIONS.find((o) => o.value === val)?.label
                : "\u00A0"; // keep height without showing '-'
              return (
                <div
                  key={i}
                  className={`rounded-xl bg-zinc-900 border border-zinc-800 px-3 py-3 text-center ${isEmpty ? "opacity-60" : ""}`}
                >
                  <div className="text-sm opacity-80">Wurf {i + 1}</div>
                  <div className={`mt-1 font-extrabold ${isEmpty ? "opacity-60" : ""}`}>{label}</div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
