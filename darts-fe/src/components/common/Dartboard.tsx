// Standard dartboard sector order starting at 20 on top and going clockwise
const DART_ORDER = [20, 1, 18, 4, 13, 6, 10, 15, 2, 17, 3, 19, 7, 16, 8, 11, 14, 9, 12, 5];

// Map ring + base number to backend ThrowType enum value
function ttFrom(n: number, ring: "S" | "D" | "T"): number {
  if (ring === "S") return n; // 1..20
  if (ring === "D") return 20 + n; // 21..40
  return 40 + n; // 41..60
}

export default function Dartboard({ onPick }: { onPick: (tt: number) => void }) {
  // Sizes (responsive-ish via tailwind max sizes)
  const size = 320;
  const ringGap = 12; // gap between ring buttons in a slice
  const sectorRadius = size / 2 - 20;

  return (
    <div
      className="relative"
      style={{ width: size, height: size }}
      aria-label="Dartboard"
    >
      {/* Outer circle visual */}
      <div className="absolute inset-0 rounded-full border-8 border-zinc-700 bg-zinc-800" />

      {/* Bulls */}
      <button
        className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-yellow-700 hover:bg-yellow-600 border border-yellow-400 text-xs font-bold"
        style={{ width: 44, height: 44 }}
        onClick={() => onPick(62)}
        aria-label="Bull (50)"
        title="BULL (50)"
      >
        50
      </button>
      <button
        className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full bg-yellow-900 hover:bg-yellow-800 border border-yellow-600 text-[10px] font-bold"
        style={{ width: 80, height: 80 }}
        onClick={() => onPick(61)}
        aria-label="Single Bull (25)"
        title="SBULL (25)"
      >
        25
      </button>

      {/* 20 sectors */}
      {DART_ORDER.map((num, idx) => {
        const angle = idx * 18 - 90; // start at top (20)
        // radial placement for the 3 buttons from outer to inner (S, D, T)
        const base = sectorRadius;
        const positions = [
          { ring: "S" as const, r: base },
          { ring: "D" as const, r: base - ringGap - 22 },
          { ring: "T" as const, r: base - 2 * ringGap - 44 }
        ];

        return (
          <div key={num} className="absolute left-1/2 top-1/2" style={{ transform: `rotate(${angle}deg)` }}>
            {positions.map((p) => (
              <button
                key={p.ring}
                onClick={() => onPick(ttFrom(num, p.ring))}
                className={
                  "absolute -translate-x-1/2 -translate-y-1/2 px-2 py-1 rounded-md border text-[10px] font-bold hover:opacity-90" +
                  (p.ring === "S"
                    ? " bg-zinc-900 border-zinc-600"
                    : p.ring === "D"
                    ? " bg-red-900/70 border-red-600"
                    : " bg-green-900/70 border-green-600")
                }
                style={{
                  left: `${Math.cos((angle * Math.PI) / 180) * p.r + size / 2}px`,
                  top: `${Math.sin((angle * Math.PI) / 180) * p.r + size / 2}px`,
                  transform: `translate(-50%, -50%) rotate(${-angle}deg)`
                }}
                aria-label={`${p.ring}${num}`}
                title={`${p.ring}${num}`}
              >
                {p.ring}
              </button>
            ))}

            {/* number label near rim */}
            <div
              className="absolute -translate-x-1/2 -translate-y-1/2 text-[10px] opacity-80"
              style={{
                left: `${Math.cos((angle * Math.PI) / 180) * (sectorRadius + 8) + size / 2}px`,
                top: `${Math.sin((angle * Math.PI) / 180) * (sectorRadius + 8) + size / 2}px`,
                transform: `translate(-50%, -50%) rotate(${-angle}deg)`
              }}
            >
              {num}
            </div>
          </div>
        );
      })}
    </div>
  );
}
