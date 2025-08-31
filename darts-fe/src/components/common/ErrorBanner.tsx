import { AlertCircle } from "lucide-react";

export default function ErrorBanner({
  message,
  onClose
}: {
  message: string | null;
  onClose: () => void;
}) {
  if (!message) return null;
  return (
    <div className="mx-4 my-2 p-3 rounded-2xl bg-red-900/40 border border-red-700 text-red-200 flex items-center gap-2">
      <AlertCircle className="shrink-0" />
      <div className="text-sm whitespace-pre-wrap break-all">{message}</div>
      <div className="flex-1" />
      <button
        onClick={onClose}
        className="px-3 py-1 rounded-lg bg-red-800/60 hover:bg-red-800/80"
      >
        OK
      </button>
    </div>
  );
}
