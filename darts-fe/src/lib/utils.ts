export function shuffle<T>(arr: T[]) {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [arr[i], arr[j]] = [arr[j], arr[i]];
  }
  return arr;
}

export function cycle<T>(list: T[], current: T): T {
  const i = list.findIndex((x) => x === current);
  return list[(i + 1) % list.length];
}

export const START_AT_OPTIONS = [101, 201, 301, 401, 501, 701, 1001] as const;

// Build throw type options for UI selection
export type ThrowTypeOption = { value: number; label: string };

export const THROW_TYPE_OPTIONS: ThrowTypeOption[] = (() => {
  const opts: ThrowTypeOption[] = [];
  // Singles S1..S20 map to 1..20
  for (let n = 1; n <= 20; n++) {
    const value = n; // Sn
    opts.push({ value, label: `S${n} (${value})` });
  }
  // Doubles D1..D20 map to 21..40
  for (let n = 1; n <= 20; n++) {
    const value = 20 + n; // Dn
    const points = n * 2;
    opts.push({ value, label: `D${n} (${points})` });
  }
  // Triples T1..T20 map to 41..60
  for (let n = 1; n <= 20; n++) {
    const value = 40 + n; // Tn
    const points = n * 3;
    opts.push({ value, label: `T${n} (${points})` });
  }
  // Bulls
  opts.push({ value: 61, label: "SBULL (25)" });
  opts.push({ value: 62, label: "BULL (50)" });

  // Order by usefulness (keep as pushed: S then D then T then Bulls) or could sort by points desc.
  return opts;
})();