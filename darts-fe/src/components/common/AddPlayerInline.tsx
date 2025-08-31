import { useState } from "react";
import type { ApiClient } from "../../lib/api";
import type { Player } from "../../types";
import { UserPlus } from "lucide-react";

export default function AddPlayerInline({
  api,
  onCreated
}: {
  api: ApiClient;
  onCreated: () => void;
}) {
  const [name, setName] = useState("");
  const create = async () => {
    if (!name.trim()) return;
    await api.call<Player>("/createPlayer", {
      method: "POST",
      body: JSON.stringify({ name })
    });
    setName("");
    onCreated();
  };
  return (
    <div className="flex items-center gap-2">
      <input
        className="w-40 px-3 py-2 rounded-xl bg-zinc-800 border border-zinc-700"
        placeholder="Name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <button
        onClick={create}
        className="px-3 py-2 rounded-xl bg-green-700 hover:bg-green-700/90 flex items-center"
      >
        <UserPlus className="mr-2" size={18} />
        Hinzufügen
      </button>
    </div>
  );
}
