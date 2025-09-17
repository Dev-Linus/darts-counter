import { useState } from "react";

export default function PlayScreen() {
  const { matchId } = useParams();
  const [throws, setThrows] = useState(["", "", ""]);
  const [players, setPlayers] = useState([
    { id: "p1", name: "Lea", active: true, score: 301 },
    { id: "p2", name: "Linus", active: false, score: 301 }
  ]);
  const [possibleFinish, setPossibleFinish] = useState<string | null>(null);

  async function submitThrows() {
    const response = await fetch(`/playerThrow`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ matchId, throws: throws.map(Number) })
    });

    const data = await response.json();
    setPlayers(data.players);
    setPossibleFinish(data.possibleFinish);
  }

  return (
    <div className="flex min-h-screen">
      {/* Left player list */}
      <div className="w-1/4 bg-gray-900 text-white p-4">
        <h2 className="text-xl font-bold mb-4">Spieler</h2>
        <ul className="space-y-2">
          {players.map((p) => (
            <li
              key={p.id}
              className={`p-2 rounded-xl ${
                p.active ? "bg-green-600 font-bold" : "bg-gray-700"
              }`}
            >
              {p.name}: {p.score}
            </li>
          ))}
        </ul>

        {possibleFinish && (
          <div className="mt-6 p-3 rounded-xl bg-yellow-600 text-black font-bold">
            Möglicher Checkout: {possibleFinish}
          </div>
        )}
      </div>

      {/* Right main game board */}
      <div className="flex-1 p-6">
        <h1 className="text-2xl font-bold mb-6">Match {matchId}</h1>

        {/* 🎯 Dartboard placeholder */}
        <div className="bg-gray-800 rounded-full aspect-square w-96 mx-auto mb-8 flex items-center justify-center text-gray-500">
          Darts Board
        </div>

        {/* 🎯 Throw input */}
        <div className="flex justify-center gap-4 mb-6">
          {throws.map((t, i) => (
            <input
              key={i}
              value={t}
              onChange={(e) =>
                setThrows((arr) => {
                  const newArr = [...arr];
                  newArr[i] = e.target.value;
                  return newArr;
                })
              }
              type="number"
              className="w-20 text-center text-xl p-2 rounded-lg border border-gray-300"
              placeholder="0"
            />
          ))}
        </div>

        <button
          onClick={submitThrows}
          className="block mx-auto bg-green-600 hover:bg-green-700 text-white font-bold py-3 px-8 rounded-xl"
        >
          Würfe speichern
        </button>
      </div>
    </div>
  );
}
