import { useState } from "react";

import type { ApiMeta, ApiVehicle } from "../../lib/api/client";
import { useVehicle } from "../../lib/data";
import { proposedActionLabel } from "../../lib/proposedAction";
import { AgingBadge } from "./AgingBadge";
import { CorrectVehicleDialog } from "./CorrectVehicleDialog";
import { ProposedActionControl } from "./ProposedActionControl";
import { StockOutDialog } from "./StockOutDialog";

function stockStateLabel(state: ApiVehicle["stock_state"]): string {
  return state === "in_stock" ? "In stock" : "Stocked out";
}

function rowClass(vehicle: ApiVehicle): string | undefined {
  if (vehicle.stock_state === "stocked_out") {
    return "text-[#8a8175]";
  }
  if (vehicle.is_aging) {
    return "bg-[#fbf6ee]";
  }
  return undefined;
}

type VehicleTableProps = {
  vehicles: ApiVehicle[];
  meta: ApiMeta;
  page: number;
  loading: boolean;
  failed: boolean;
  onPage: (page: number) => void;
};

export function VehicleTable({
  vehicles,
  meta,
  page,
  loading,
  failed,
  onPage,
}: VehicleTableProps) {
  const [openId, setOpenId] = useState<string | null>(null);
  const pageCount = Math.max(1, Math.ceil(meta.total / meta.page_size));

  return (
    <section
      aria-label="Inventory list"
      className="mt-6 min-h-80 border border-[#e3dccf] bg-white shadow-[0_1px_0_rgba(28,25,21,0.04)]"
    >
      <div className="flex items-center justify-between border-b border-[#efeae1] px-5 py-4">
        <h2 className="text-xl">Inventory list</h2>
        {loading || failed ? null : (
          <p className="text-sm text-[#4e463c]">{meta.total} vehicles</p>
        )}
      </div>
      {failed ? (
        <p className="px-5 py-8 text-[#8a3b2c]" role="alert">
          Could not load the inventory.
        </p>
      ) : loading ? (
        <p className="px-5 py-8 text-[#4e463c]">Loading inventory.</p>
      ) : vehicles.length === 0 ? (
        <p className="px-5 py-8 text-[#4e463c]">No vehicles</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full min-w-[56rem] border-collapse text-left text-sm">
            <thead className="text-[#4e463c]">
              <tr className="border-b border-[#efeae1]">
                <th className="px-4 py-3 font-medium">VIN</th>
                <th className="px-4 py-3 font-medium">Make</th>
                <th className="px-4 py-3 font-medium">Model</th>
                <th className="px-4 py-3 font-medium">Entered on</th>
                <th className="px-4 py-3 font-medium">Stocked-out on</th>
                <th className="px-4 py-3 font-medium">Stock state</th>
                <th className="px-4 py-3 font-medium">Inventory age</th>
                <th className="px-4 py-3 font-medium">Aging stock</th>
                <th className="px-4 py-3 font-medium">Proposed action</th>
              </tr>
            </thead>
            <tbody>
              {vehicles.map((vehicle) => (
                <VehicleRows
                  key={vehicle.id}
                  vehicle={vehicle}
                  open={openId === vehicle.id}
                  onToggle={() =>
                    setOpenId((current) => (current === vehicle.id ? null : vehicle.id))
                  }
                />
              ))}
            </tbody>
          </table>
        </div>
      )}
      {meta.total > meta.page_size ? (
        <div className="flex items-center justify-between border-t border-[#efeae1] px-5 py-3 text-sm">
          <p>
            Page {page} of {pageCount}
          </p>
          <div className="flex gap-2">
            <button
              type="button"
              className="border border-[#d9d1c3] px-3 py-1 disabled:opacity-40"
              disabled={page <= 1}
              onClick={() => onPage(page - 1)}
            >
              Previous
            </button>
            <button
              type="button"
              className="border border-[#d9d1c3] px-3 py-1 disabled:opacity-40"
              disabled={page >= pageCount}
              onClick={() => onPage(page + 1)}
            >
              Next
            </button>
          </div>
        </div>
      ) : null}
    </section>
  );
}

function VehicleRows({
  vehicle,
  open,
  onToggle,
}: {
  vehicle: ApiVehicle;
  open: boolean;
  onToggle: () => void;
}) {
  return (
    <>
      <tr className={rowClass(vehicle)}>
        <td className="px-4 py-3">
          <button type="button" className="underline" onClick={onToggle}>
            {vehicle.vin}
          </button>
        </td>
        <td className="px-4 py-3">{vehicle.make.name}</td>
        <td className="px-4 py-3">{vehicle.model.name}</td>
        <td className="px-4 py-3">{vehicle.entered_on}</td>
        <td className="px-4 py-3">{vehicle.stocked_out_on ?? ""}</td>
        <td className="px-4 py-3">{stockStateLabel(vehicle.stock_state)}</td>
        <td className="px-4 py-3">{vehicle.inventory_age}</td>
        <td className="px-4 py-3">
          <AgingBadge aging={vehicle.is_aging} />
        </td>
        <td className="px-4 py-3">
          {proposedActionLabel(vehicle.proposed_action)}
        </td>
      </tr>
      {open ? (
        <tr>
          <td colSpan={9} className="bg-[#f6f1e7] px-4 py-4">
            <VehicleView id={vehicle.id} />
          </td>
        </tr>
      ) : null}
    </>
  );
}

function VehicleView({ id }: { id: string }) {
  const vehicle = useVehicle(id);
  if (vehicle.isError) {
    return (
      <p role="alert" className="text-[#8a3b2c]">
        Could not load the vehicle.
      </p>
    );
  }
  if (!vehicle.data) {
    return <p className="text-[#4e463c]">Loading vehicle.</p>;
  }
  const row = vehicle.data;
  return (
    <div>
      <dl className="grid grid-cols-2 gap-x-6 gap-y-2 text-sm sm:grid-cols-3">
        <div>
          <dt className="text-[#4e463c]">VIN</dt>
          <dd>{row.vin}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Make</dt>
          <dd>{row.make.name}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Model</dt>
          <dd>{row.model.name}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Entered on</dt>
          <dd>{row.entered_on}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Stocked-out on</dt>
          <dd>{row.stocked_out_on ?? ""}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Stock state</dt>
          <dd>{stockStateLabel(row.stock_state)}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Inventory age</dt>
          <dd>{row.inventory_age}</dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Aging stock</dt>
          <dd>
            <AgingBadge aging={row.is_aging} />
          </dd>
        </div>
        <div>
          <dt className="text-[#4e463c]">Proposed action</dt>
          <dd>{proposedActionLabel(row.proposed_action)}</dd>
        </div>
      </dl>
      {row.stock_state === "stocked_out" ? (
        <p className="mt-4 text-sm text-[#4e463c]">
          Model and entered on cannot be corrected because the vehicle is
          stocked out. The last proposed action stays as recorded.
        </p>
      ) : null}
      <div className="flex flex-wrap gap-3">
        {row.stock_state === "in_stock" ? (
          <CorrectVehicleDialog vehicle={row} />
        ) : null}
        <ProposedActionControl vehicle={row} />
        <StockOutDialog vehicle={row} />
      </div>
    </div>
  );
}
