export function InventoryPage() {
  return (
    <div className="min-h-screen bg-[#f3efe6] text-[#1c1915]">
      <header className="border-b border-[#14201b] bg-[#1c2b24] text-[#f6f1e7]">
        <div className="mx-auto flex max-w-5xl items-end justify-between px-6 py-6">
          <div>
            <p className="text-xs font-medium tracking-[0.22em] text-[#c8bba6] uppercase">
              Dashboard
            </p>
            <h1 className="mt-1 text-4xl leading-none">Inventory</h1>
          </div>
        </div>
      </header>
      <main className="mx-auto max-w-5xl px-6 py-8">
        <p className="max-w-xl text-[15px] leading-6 text-[#4e463c]">
          Vehicles the dealership holds, including vehicles that have stocked
          out.
        </p>
        <section
          aria-label="Inventory list"
          className="mt-6 min-h-80 border border-[#e3dccf] bg-white shadow-[0_1px_0_rgba(28,25,21,0.04)]"
        >
          <div className="border-b border-[#efeae1] px-5 py-4">
            <h2 className="text-xl">Inventory list</h2>
          </div>
        </section>
      </main>
    </div>
  );
}
