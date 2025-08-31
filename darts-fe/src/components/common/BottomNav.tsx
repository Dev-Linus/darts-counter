import { Home, UserPlus, LineChart, List, Settings } from "lucide-react";

export default function BottomNav({
  tab,
  setTab
}: {
  tab: string;
  setTab: (t: string) => void;
}) {
  const Item = ({
    id,
    Icon,
    label
  }: {
    id: string;
    Icon: any;
    label: string;
  }) => (
    <button
      onClick={() => setTab(id)}
      className={`flex-1 flex flex-col items-center py-2 ${
        tab === id ? "text-green-500" : "text-gray-400"
      }`}
    >
      <Icon size={22} />
      <span className="text-xs mt-1">{label}</span>
    </button>
  );
  return (
    <div className="fixed bottom-0 left-0 right-0 bg-zinc-900 border-t border-zinc-800 flex">
      <Item id="start" Icon={Home} label="Start" />
      <Item id="players" Icon={UserPlus} label="Spieler" />
      <Item id="stats" Icon={LineChart} label="Statistiken" />
      <Item id="matches" Icon={List} label="Alle Spiele" />
      <Item id="settings" Icon={Settings} label="Einstellungen" />
    </div>
  );
}
