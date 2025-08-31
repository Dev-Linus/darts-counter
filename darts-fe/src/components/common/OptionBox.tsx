export default function OptionBox({
  label,
  value,
  color = "green",
  onClick
}: {
  label: string;
  value: string;
  color?: "green" | "red";
  onClick?: () => void;
}) {
  const isGreen = color === "green";
  return (
    <button
      onClick={onClick}
      className={`rounded-2xl p-3 text-center ${
        isGreen
          ? "bg-green-800/60 border-green-700"
          : "bg-red-800/70 border-red-700"
      } border shadow-inner`}
    >
      <div className="text-xs opacity-80 mb-1">{label}</div>
      <div className="text-2xl font-extrabold">{value}</div>
    </button>
  );
}
