// Standard dartboard sector order starting at 20 on top and going clockwise
const DART_ORDER = [
  20, 1, 18, 4, 13, 6, 10, 15, 2, 17, 3, 19, 7, 16, 8, 11, 14, 9, 12, 5
];

// Map ring + base number to backend ThrowType enum value
function ttFrom(n: number, ring: "S" | "D" | "T"): number {
  if (ring === "S") return n; // 1..20
  if (ring === "D") return 20 + n; // 21..40
  return 40 + n; // 41..60
}

// build a donut sector path from inner radius r1 to outer radius r2 between angles a0..a1 (deg)
function donutSlice(
  cx: number,
  cy: number,
  r1: number,
  r2: number,
  a0: number,
  a1: number
) {
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

export default function Dartboard({
  onPick
}: {
  onPick: (tt: number) => void;
}) {
  // Larger geometry and responsive width; scales via viewBox
  const size = 600; // bigger base size
  const cx = size / 2;
  const cy = size / 2;

  // Radii (proportional)
  const rBullInner = 15; // 50
  const rBullOuter = 30; // 25
  const rInnerSingleOuter = 160;
  const rTripleInner = 160;
  const rTripleOuter = 175;
  const rOuterSingleOuter = 255;
  const rDoubleInner = 255;
  const rDoubleOuter = 270; // outermost playable ring
  const rNumbers = rDoubleOuter + 20; // where to draw numbers

  const sectorAngle = 360 / 20; // 18 degrees

  return (
    <svg
      role="img"
      aria-label="Dartboard"
      viewBox={`0 0 ${size} ${size}`}
      className="w-[360px] sm:w-[420px] md:w-[460px] h-auto drop-shadow-[0_0_8px_rgba(0,0,0,0.6)]"
    >
      {/* Backing circle */}
      <circle
        cx={cx}
        cy={cy}
        r={rDoubleOuter + 8}
        className="fill-zinc-900 stroke-zinc-700"
        strokeWidth={8}
      />

      {/* 20 sectors with Singles/Doubles/Triples */}
      {DART_ORDER.map((num, idx) => {
        const start = -100 + idx * sectorAngle;
        const end = start + sectorAngle;

        // Colors: beige singles; triples/doubles alternate red/green by sector
        const alternate = idx % 2 === 0; // alternate around the board
        const singleFill = alternate ? "#1a1a1a" : "#e7e3c6"; // black/beige
        const tripleFill = alternate ? "#c81e1e" : "#1f8a3a"; // red/green alternating
        const doubleFill = alternate ? "#c81e1e" : "#1f8a3a"; // red/green alternating

        const paths = [
          {
            r1: rBullOuter,
            r2: rInnerSingleOuter,
            fill: singleFill,
            ring: "S" as const
          }, // inner single
          {
            r1: rTripleInner,
            r2: rTripleOuter,
            fill: tripleFill,
            ring: "T" as const
          }, // triple
          {
            r1: rTripleOuter,
            r2: rOuterSingleOuter,
            fill: singleFill,
            ring: "S" as const
          }, // outer single
          {
            r1: rDoubleInner,
            r2: rDoubleOuter,
            fill: doubleFill,
            ring: "D" as const
          } // double
        ];

        const mid = start + sectorAngle / 2;
        const rad = (Math.PI / 180) * mid;
        const nx = cx + rNumbers * Math.cos(rad);
        const ny = cy + rNumbers * Math.sin(rad);

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
            {/* Rim number (not clickable) */}
            <text
              x={nx}
              y={ny}
              textAnchor="middle"
              dominantBaseline="middle"
              fill="#ffffff"
              fontSize={18}
              fontWeight={700}
              style={{ pointerEvents: "none" }}
            >
              {num}
            </text>
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
