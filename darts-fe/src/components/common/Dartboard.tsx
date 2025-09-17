// Standard dartboard sector order starting at 20 on top and going clockwise
const DART_ORDER = [20, 1, 18, 4, 13, 6, 10, 15, 2, 17, 3, 19, 7, 16, 8, 11, 14, 9, 12, 5];

// Map ring + base number to backend ThrowType enum value
function ttFrom(n: number, ring: "S" | "D" | "T"): number {
  if (ring === "S") return n; // 1..20
  if (ring === "D") return 20 + n; // 21..40
  return 40 + n; // 41..60
}

// build a donut sector path from inner radius r1 to outer radius r2 between angles a0..a1 (deg)
function donutSlice(cx: number, cy: number, r1: number, r2: number, a0: number, a1: number) {
  const toRad = (a: number) => (Math.PI / 180) * a;
  const x1o = cx + r2 * Math.cos(toRad(a0));
  const y1o = cy + r2 * Math.sin(toRad(a0));
  const x2o = cx + r2 * Math.cos(toRad(a1));
  const y2o = cy + r2 * Math.sin(toRad(a1));
  const x2i = cx + r1 * Math.cos(toRad(a1));
  const y2i = cy + r1 * Math.sin(toRad(a1));
  const x1i = cx + r1 * Math.cos(toRad(a0));
  const y1i = cy + r1 * Math.sin(toRad(a0));
  const large = a1 - a0 > 180 ? 1 : 0;
  return `M ${x1o} ${y1o} A ${r2} ${r2} 0 ${large} 1 ${x2o} ${y2o} L ${x2i} ${y2i} A ${r1} ${r1} 0 ${large} 0 ${x1i} ${y1i} Z`;
}

export default function Dartboard({ onPick }: { onPick: (tt: number) => void }) {
  // Geometry based on a 400x400 board; scales via viewBox
  const size = 400;
  const cx = size / 2;
  const cy = size / 2;

  // Radii (approximate but proportionally correct)
  const rBullInner = 18;   // 50
  const rBullOuter = 36;   // 25
  const rInnerSingleOuter = 120; // up to before triple
  const rTripleInner = 130;
  const rTripleOuter = 145;
  const rOuterSingleOuter = 185;
  const rDoubleInner = 195;
  const rDoubleOuter = 210;

  const sectorAngle = 360 / 20; // 18 degrees

  return (
    <svg
      role="img"
      aria-label="Dartboard"
      viewBox={`0 0 ${size} ${size}`}
      className="max-w-full h-auto drop-shadow-[0_0_8px_rgba(0,0,0,0.6)]"
    >
      {/* Backing circle */}
      <circle cx={cx} cy={cy} r={rDoubleOuter + 6} className="fill-zinc-900 stroke-zinc-700" strokeWidth={6} />

      {/* 20 sectors with Singles/Doubles/Triples */}
      {DART_ORDER.map((num, idx) => {
        const start = -90 + idx * sectorAngle;
        const end = start + sectorAngle;

        // Colors to match the reference: beige singles, green triples, red doubles
        const singleFill = "#e7e3c6"; // beige
        const tripleFill = "#1f8a3a"; // green
        const doubleFill = "#c81e1e"; // red

        const paths = [
          { r1: rBullOuter, r2: rInnerSingleOuter, fill: singleFill, ring: "S" as const }, // inner single
          { r1: rTripleInner, r2: rTripleOuter, fill: tripleFill, ring: "T" as const }, // triple
          { r1: rTripleOuter, r2: rOuterSingleOuter, fill: singleFill, ring: "S" as const }, // outer single
          { r1: rDoubleInner, r2: rDoubleOuter, fill: doubleFill, ring: "D" as const }, // double
        ];

        return (
          <g key={num}>
            {paths.map((p, i) => (
              <path
                key={i}
                d={donutSlice(cx, cy, p.r1, p.r2, start, end)}
                fill={p.fill}
                stroke="#111827"
                strokeWidth={1}
                onClick={() => onPick(ttFrom(num, p.ring))}
                cursor="pointer"
              />
            ))}
          </g>
        );
      })}

      {/* Bulls */}
      <circle
        cx={cx}
        cy={cy}
        r={rBullOuter}
        className="fill-green-600 stroke-green-700"
        strokeWidth={2}
        onClick={() => onPick(61)}
        cursor="pointer"
      />
      <circle
        cx={cx}
        cy={cy}
        r={rBullInner}
        className="fill-red-600 stroke-red-700"
        strokeWidth={2}
        onClick={() => onPick(62)}
        cursor="pointer"
      />
    </svg>
  );
}
