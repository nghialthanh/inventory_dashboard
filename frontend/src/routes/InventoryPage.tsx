import { useState } from "react";

import { InventoryFilters } from "../components/inventory/InventoryFilters";
import { RecordVehicleDialog } from "../components/inventory/RecordVehicleDialog";
import { VehicleTable } from "../components/inventory/VehicleTable";
import { useVehicles } from "../lib/data";

export function InventoryPage() {
  const [makeId, setMakeId] = useState("");
  const [modelId, setModelId] = useState("");
  const [stockState, setStockState] = useState<"" | "in_stock" | "stocked_out">(
    "",
  );
  const [stockedOutOn, setStockedOutOn] = useState("");
  const [agingOnly, setAgingOnly] = useState(false);
  const [page, setPage] = useState(1);
  const vehicles = useVehicles({
    makeId,
    modelId,
    stockState,
    stockedOutOn,
    agingOnly,
    page,
  });
  const meta = vehicles.data?.meta ?? { page, page_size: 20, total: 0 };

  return (
    <div className="min-h-screen bg-[#f3efe6] text-[#1c1915]">
      <header className="border-b border-[#14201b] bg-[#1c2b24] text-[#f6f1e7]">
        <div className="mx-auto flex max-w-7xl items-end justify-between px-6 py-6">
          <div>
            <p className="text-xs font-medium tracking-[0.22em] text-[#c8bba6] uppercase">
              Dashboard
            </p>
            <h1 className="mt-1 text-4xl leading-none">Inventory</h1>
          </div>
          <RecordVehicleDialog />
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-6 py-8">
        <p className="max-w-xl text-[15px] leading-6 text-[#4e463c]">
          Vehicles the dealership holds, including vehicles that have stocked
          out.
        </p>
        <InventoryFilters
          makeId={makeId}
          modelId={modelId}
          stockState={stockState}
          stockedOutOn={stockedOutOn}
          agingOnly={agingOnly}
          onMakeId={(next) => {
            setMakeId(next);
            setModelId("");
            setPage(1);
          }}
          onModelId={(next) => {
            setModelId(next);
            setPage(1);
          }}
          onStockState={(next) => {
            setStockState(next);
            setPage(1);
          }}
          onStockedOutOn={(next) => {
            setStockedOutOn(next);
            setPage(1);
          }}
          onAgingOnly={(next) => {
            setAgingOnly(next);
            setPage(1);
          }}
        />
        <VehicleTable
          vehicles={vehicles.data?.vehicles ?? []}
          meta={meta}
          page={page}
          loading={vehicles.isLoading}
          failed={
            vehicles.isError ||
            vehicles.fetchStatus === "paused" ||
            (vehicles.isPending && vehicles.fetchStatus === "idle")
          }
          onPage={setPage}
        />
      </main>
    </div>
  );
}
