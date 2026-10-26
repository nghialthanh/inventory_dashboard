import { useState, type FormEvent } from "react";

import { ApiError, type ApiVehicle } from "../../lib/api/client";
import { useStockOutVehicle } from "../../lib/data";

type FieldErrors = {
  stockedOutOn?: string;
};

export function StockOutDialog({ vehicle }: { vehicle: ApiVehicle }) {
  const dialogId = `stock-out-${vehicle.id}`;
  const [open, setOpen] = useState(false);
  const [stockedOutOn, setStockedOutOn] = useState("");
  const [errors, setErrors] = useState<FieldErrors>({});
  const [toast, setToast] = useState<string | null>(null);
  const stockOut = useStockOutVehicle(vehicle.id);

  function close() {
    setOpen(false);
    setErrors({});
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (stockedOutOn === "") {
      setErrors({ stockedOutOn: "Enter stocked-out on." });
      return;
    }
    setErrors({});
    try {
      await stockOut.mutateAsync({ stocked_out_on: stockedOutOn });
      setOpen(false);
      setToast("Vehicle stocked out.");
    } catch (error) {
      const message =
        error instanceof ApiError ? error.message : "Could not stock out the vehicle.";
      setToast(message);
    }
  }

  return (
    <>
      <button
        type="button"
        className="mt-4 border border-[#1c2b24] px-3 py-1 text-sm text-[#1c2b24]"
        onClick={() => {
          setStockedOutOn("");
          setErrors({});
          setToast(null);
          setOpen(true);
        }}
      >
        Stock out
      </button>
      {open ? (
        <div className="fixed inset-0 z-10 flex items-start justify-center bg-[#1c1915]/40 px-4 py-16">
          <form
            aria-labelledby={dialogId}
            role="dialog"
            className="w-full max-w-md border border-[#e3dccf] bg-[#f6f1e7] p-6 shadow-lg"
            onSubmit={submit}
          >
            <h2 id={dialogId} className="text-2xl">
              Stock out
            </h2>
            <p className="mt-2 text-sm text-[#4e463c]">
              Record the calendar day this vehicle stocked out. The record stays on
              the list.
            </p>
            <label className="mt-4 flex flex-col gap-1 text-sm text-[#4e463c]">
              Stocked-out on
              <input
                type="date"
                className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                value={stockedOutOn}
                onChange={(event) => setStockedOutOn(event.target.value)}
              />
              {errors.stockedOutOn ? (
                <span className="text-[#8a3b2c]">{errors.stockedOutOn}</span>
              ) : null}
            </label>
            <div className="mt-6 flex justify-end gap-3">
              <button
                type="button"
                className="px-4 py-2 text-sm text-[#4e463c]"
                onClick={close}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="border border-[#1c2b24] bg-[#1c2b24] px-4 py-2 text-sm text-[#f6f1e7]"
                disabled={stockOut.isPending}
              >
                Stock out
              </button>
            </div>
          </form>
        </div>
      ) : null}
      {toast ? (
        <p
          className="fixed right-4 bottom-4 z-20 border border-[#d9d1c3] bg-white px-4 py-3 text-sm text-[#1c1915] shadow-lg"
          role="status"
        >
          {toast}
        </p>
      ) : null}
    </>
  );
}
