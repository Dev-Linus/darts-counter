import Header from "../common/Header";
import type { Player } from "../../types";
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip
} from "recharts";

export default function StatsScreen({ players }: { players: Player[] }) {
  const rows = players.map((p) => ({
    name: p.name,
    avg: p.throws ? p.totalScore / p.throws : 0,
    total: p.totalScore,
    throws: p.throws,
    matches: p.matches
  }));

  return (
    <div className="px-4 pb-24">
      <Header title="Statistiken" />

      <Section title="Average und höchste Punktzahl">
        <ResponsiveContainer width="100%" height={220}>
          <BarChart data={rows}>
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip />
            <Bar dataKey="avg" />
          </BarChart>
        </ResponsiveContainer>
      </Section>

      <Section title="Punkte gesamt">
        <ResponsiveContainer width="100%" height={220}>
          <BarChart data={rows}>
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip />
            <Bar dataKey="total" />
          </BarChart>
        </ResponsiveContainer>
      </Section>

      <Table
        title="Spieler Werte"
        headers={["Spieler", "Matches", "Würfe", "Punkte", "Ø pro Wurf"]}
        rows={rows.map((r) => [
          r.name,
          r.matches,
          r.throws,
          r.total,
          r.avg.toFixed(2)
        ])}
      />
    </div>
  );
}

function Section({
  title,
  children
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-bold mb-2">{title}</h3>
      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-3">
        {children}
      </div>
    </div>
  );
}

function Table({
  title,
  headers,
  rows
}: {
  title: string;
  headers: string[];
  rows: React.ReactNode[][];
}) {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-bold mb-2">{title}</h3>
      <div className="overflow-auto rounded-2xl border border-zinc-800">
        <table className="min-w-full text-sm">
          <thead className="bg-zinc-800">
            <tr>
              {headers.map((h, i) => (
                <th key={i} className="text-left px-3 py-2 font-semibold">
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((r, ri) => (
              <tr key={ri} className="odd:bg-zinc-950 even:bg-zinc-900">
                {r.map((c, ci) => (
                  <td key={ci} className="px-3 py-2">
                    {c}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
