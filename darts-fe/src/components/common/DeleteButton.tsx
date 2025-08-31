import { Trash2 } from "lucide-react";
export default function DeleteButton({ onClick }: { onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className="text-zinc-400 hover:text-red-400 p-2"
      title="Löschen"
    >
      <Trash2 />
    </button>
  );
}
