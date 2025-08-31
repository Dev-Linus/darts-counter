import { useState } from "react";
import type { ApiLog } from "../../types";
import { ChevronDown, ChevronUp } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export default function DebugDrawer({ logs }: { logs: ApiLog[] }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="fixed right-3 bottom-16">
      <button
        onClick={() => setOpen((o) => !o)}
        className="rounded-full shadow-xl px-3 py-2 bg-zinc-800 border border-zinc-700"
      >
        {open ? <ChevronDown /> : <ChevronUp />} Debug
      </button>
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: 8 }}
            className="mt-2 w-[380px] max-h-[50vh] overflow-auto bg-zinc-900 border border-zinc-800 rounded-2xl p-3 text-xs text-zinc-200"
          >
            {logs.length === 0 && (
              <div className="opacity-60">No requests yet…</div>
            )}
            {logs.map((l, i) => (
              <div key={i} className="mb-3">
                <div className="font-semibold">{l.time}</div>
                <pre className="whitespace-pre-wrap break-words">
                  {JSON.stringify(l.request, null, 2)}
                </pre>
                {l.response && (
                  <pre className="whitespace-pre-wrap break-words mt-1">
                    {JSON.stringify(l.response, null, 2)}
                  </pre>
                )}
                {l.error && (
                  <div className="text-red-400">Error: {l.error}</div>
                )}
                <div className="h-px bg-zinc-800 mt-2" />
              </div>
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
