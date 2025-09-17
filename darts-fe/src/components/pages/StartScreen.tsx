import Header from "../common/Header";
import OptionBox from "../common/OptionBox";
import AddPlayerInline from "../common/AddPlayerInline";
import DeleteButton from "../common/DeleteButton";
import { useState } from "react";
import type { ApiClient } from "../../lib/api";
import { usePlayers } from "../../hooks/usePlayers";
import { useMatches } from "../../hooks/useMatches";
import { START_AT_OPTIONS, cycle, shuffle } from "../../lib/utils";
import type { Match } from "../../types";
import { Bot, Target } from "lucide-react";

export default function StartScreen({
  api,
  playersState,
  matchesState
}: {
  api: ApiClient;
  playersState: ReturnType<typeof usePlayers>;
  matchesState: ReturnType<typeof useMatches>;
}) {
  const { players } = playersState;
  const [startAt, setStartAt] = useState<number>(301);
  const [startMode, setStartMode] = useState<0 | 1 | 2>(2);
  const [endMode, setEndMode] = useState<0 | 1 | 2>(2);
  const [sets, setSets] = useState<number>(1);
  const [legs, setLegs] = useState<number>(1);
  const [randomOrder, setRandomOrder] = useState(false);
  const [selected, setSelected] = useState<string[]>([]);

  const toggle = (id: string) =>
    setSelected((s) =>
      s.includes(id) ? s.filter((x) => x !== id) : [...s, id]
    );
  const startModeLabel = (m: 0 | 1 | 2) =>
    m === 0 ? "Straight In" : m === 1 ? "Double In" : "Master In";
  const endModeLabel = (m: 0 | 1 | 2) =>
    m === 0 ? "Straight Out" : m === 1 ? "Double Out" : "Master Out";

  const create = async () => {
    const chosen = randomOrder ? shuffle([...selected]) : selected;
    if (chosen.length === 0) throw new Error("Bitte Spieler auswählen");
    await api.call<Match>("/createMatch", {
      method: "POST",
      body: JSON.stringify({
        Pids: chosen,
        StartAt: startAt,
        StartMode: startMode,
        EndMode: endMode
      })
    });
    await matchesState.refresh();
  };

  return (
    <div className="px-4 pb-24">
      <Header title="Dart Zähler" />

      <div className="grid grid-cols-3 gap-3 mt-4">
        <OptionBox
          label="Punkte"
          value={String(startAt)}
          color="green"
          onClick={() => setStartAt(cycle([...START_AT_OPTIONS], startAt))}
        />
        <OptionBox
          label="Check-Out"
          value={endModeLabel(endMode)}
          color="red"
          onClick={() => setEndMode(endMode === 0 ? 1 : endMode === 1 ? 2 : 0)}
        />
        <OptionBox
          label="Sätzen"
          value={String(sets)}
          color="green"
          onClick={() => setSets((n) => (n % 5) + 1)}
        />
        <OptionBox label="Satz/Leg" value="First to" color="green" />
        <OptionBox
          label="Check-In"
          value={startModeLabel(startMode)}
          color="red"
          onClick={() =>
            setStartMode(startMode === 0 ? 1 : startMode === 1 ? 2 : 0)
          }
        />
        <OptionBox
          label="Legs"
          value={String(legs)}
          color="green"
          onClick={() => setLegs((n) => (n % 5) + 1)}
        />
      </div>

      <button
        onClick={create}
        className="w-full mt-6 bg-red-600 hover:bg-red-700 rounded-2xl py-3 text-lg font-bold"
      >
        START
      </button>

      <div className="mt-5 flex items-center gap-3">
        <input
          id="rand"
          type="checkbox"
          checked={randomOrder}
          onChange={(e) => setRandomOrder(e.target.checked)}
        />
        <label htmlFor="rand" className="text-sm">
          Zufällige Reihenfolge
        </label>
      </div>

      <div className="mt-3 flex gap-3">
        <button className="rounded-xl px-3 py-2 bg-green-700 hover:bg-green-700/90 flex items-center opacity-70 cursor-not-allowed">
          <Bot className="mr-2" size={18} />
          KI hinzufügen
        </button>
        <AddPlayerInline api={api} onCreated={playersState.refresh} />
      </div>

      <div className="mt-4 text-sm opacity-80">Spieler</div>
      <div className="divide-y divide-zinc-800">
        {players.map((p) => (
          <div key={p.id} className="flex items-center py-3">
            <Target
              className={`mr-3 ${
                selected.includes(p.id) ? "text-green-500" : "text-zinc-500"
              }`}
            />
            <button
              onClick={() => toggle(p.id)}
              className="text-base font-semibold mr-auto text-left"
            >
              {p.name}
            </button>
            <DeleteButton
              onClick={async () => {
                await api.call(`/deletePlayer?playerId=${p.id}`, {
                  method: "DELETE"
                });
                playersState.refresh();
              }}
            />
          </div>
        ))}
      </div>
    </div>
  );
}
