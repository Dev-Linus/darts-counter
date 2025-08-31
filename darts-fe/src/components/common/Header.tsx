export default function Header({ title }: { title: string }) {
  const now = new Date().toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit"
  });
  return (
    <div className="px-4 pt-4 pb-2">
      <div className="opacity-70 text-xs">{now}</div>
      <h1 className="text-3xl font-extrabold mt-1">{title}</h1>
    </div>
  );
}
