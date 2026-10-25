export function AgingBadge({ aging }: { aging: boolean }) {
  if (!aging) {
    return <span>No</span>;
  }
  return (
    <span className="inline-block bg-[#f3e2c4] px-2 py-0.5 text-sm text-[#6b4a16]">
      Aging stock
    </span>
  );
}
