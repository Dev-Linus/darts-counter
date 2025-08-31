import Header from "../common/Header";
import DeleteButton from "../common/DeleteButton";
import { useState } from "react";
import type { ApiClient } from "../../lib/api";
import { usePlayers } from "../../hooks/usePlayers";
import type { Player } from "../../types";
import { Plus } from "lucide-react";

export default function PlayersScreen({
  api,
  playersState
}: {
  api: ApiClient;
  playersState: ReturnType<typeof usePlayers>;
}) {
  const { players, refresh } = playersState;
  const [name, setName] = useState("");

  const create = async () => {
    await api.call<Player>("/createPlayer", {
      method: "POST",
      body: JSON.stringify({ name })
    });
    setName("");
    refresh();
  };

  return (
    <div className="px-4 pb-24">
      <Header title="Spieler" />
      <div className="flex gap-2">
        <input
          className="flex-1 px-3 py-2 rounded-xl bg-zinc-800 border border-zinc-700"
          placeholder="Name"
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
        <button
          onClick={create}
          className="px-3 py-2 rounded-xl bg-green-700 hover:bg-green-700/90 flex items-center"
        >
          <Plus className="mr-2" />
          Erstellen
        </button>
      </div>

      <div className="mt-4 divide-y divide-zinc-800">
        {players.map((p) => (
          <div key={p.id} className="flex items-center py-3">
            <div className="font-semibold">{p.name}</div>
            <div className="ml-auto text-sm opacity-70">
              Matches: {p.matches} • Würfe: {p.throws} • Punkte: {p.totalScore}
            </div>
            <DeleteButton
              onClick={async () => {
                await api.call(`/deletePlayer?playerId=${p.id}`, {
                  method: "DELETE"
                });
                refresh();
              }}
            />
          </div>
        ))}
      </div>
    </div>
  );
}
