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