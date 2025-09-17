import { useEffect, useMemo, useState } from "react";
import type { ApiClient } from "../../lib/api";
import type { Match, Player } from "../../types";
import { THROW_TYPE_OPTIONS } from "../../lib/utils";
import { ArrowLeft } from "lucide-react";

export default function PlayScreen({
  api,
  players,
  matches,
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
  const matchFromList = useMemo(
    () => matches.find((m) => m.id === matchId),
    [matches, matchId]
  );

  const nameOf = (id: string) => players.find((p) => p.id === id)?.name || id;

  const [scores, setScores] = useState<Record<string, number>>(
    matchFromList?.scores || {}
  );
  const [currentPid, setCurrentPid] = useState<string | undefined>(
    matchFromList?.currentPlayer
  );
  const [turnThrows, setTurnThrows] = useState<number[]>([]);
  const [finishes, setFinishes] = useState<number[]>([]);
  const [useBoardView, setUseBoardView] = useState(false); // <-- NEW toggle

  useEffect(() => {
    if (matchFromList) {
      setScores(matchFromList.scores || {});
      setCurrentPid(matchFromList.currentPlayer);
    }
  }, [matchFromList?.id]);

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
      setTurnThrows([tt]);
    } else {
      setTurnThrows((prev) => [...prev, tt].slice(-3));
    }

    setCurrentPid(resp.NextThrowBy);
    onRefreshMatches();
  };

  const finishLabels = finishes
    .map(
      (v) => THROW_TYPE_OPTIONS.find((o) => o.value === v)?.label || String(v)
    )
    .slice(0, 20);

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
          {matchFromList?.players.map((pid) => (
            <div
              key={pid}
              className={`flex items-center justify-between rounded-xl px-3 py-2 border ${
                pid === currentPid
                  ? "bg-green-900/30 border-green-700"
                  : "bg-zinc-900 border-zinc-800"
              }`}
            >
              <div className="font-semibold truncate mr-3">{nameOf(pid)}</div>
              <div className="text-xl font-extrabold">{scores?.[pid] ?? 0}</div>
            </div>
          ))}
        </div>

        {/* RIGHT: Board or Grid */}
        <div className="space-y-6">
          {useBoardView ? (
            <div className="flex justify-center items-center">
              {/* Dartboard placeholder */}
              <div className="w-[320px] h-[320px] rounded-full border-8 border-zinc-700 bg-zinc-800 flex items-center justify-center text-zinc-400">
                🎯 Dartboard (klickbare Felder später)
              </div>
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

          {/* Possible finish */}
          {finishLabels.length > 0 && (
            <div>
              <div className="text-sm opacity-80 mb-2">Mögliche Finishes</div>
              <div className="flex flex-wrap gap-2">
                {finishLabels.map((l, idx) => (
                  <span
                    key={idx}
                    className="px-2 py-1 rounded-full bg-emerald-900/30 border border-emerald-700 text-sm"
                  >
                    {l}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
