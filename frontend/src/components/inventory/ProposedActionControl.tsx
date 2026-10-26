import { useState, type FormEvent } from "react";

import { ApiError, type ApiVehicle } from "../../lib/api/client";
import type { ProposedAction } from "../../lib/api/vehicles";
import { useSetProposedAction } from "../../lib/data";
import {
  PROPOSED_ACTION_OPTIONS,
  asProposedAction,
} from "../../lib/proposedAction";

export function ProposedActionControl({ vehicle }: { vehicle: ApiVehicle }) {
  const dialogId = `proposed-action-${vehicle.id}`;
  const agingStock = vehicle.is_aging;
  const [open, setOpen] = useState(false);
  const [proposedAction, setProposedAction] = useState<ProposedAction>(
    asProposedAction(vehicle.proposed_action),
  );
  const [toast, setToast] = useState<string | null>(null);
  const save = useSetProposedAction(vehicle.id);

  function close() {
    setOpen(false);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    try {
      await save.mutateAsync({ proposed_action: proposedAction });
      setOpen(false);
      setToast(
        proposedAction === ""
          ? "Proposed action cleared."
          : "Proposed action recorded.",
      );
    } catch (error) {
      const message =
        error instanceof ApiError
          ? error.message
          : "Could not record the proposed action.";
      setToast(message);
    }
  }

  return (
    <>
      <button
        type="button"
        className="mt-4 border border-[#1c2b24] px-3 py-1 text-sm text-[#1c2b24] disabled:opacity-40"
        disabled={!agingStock}
        onClick={() => {
          setProposedAction(asProposedAction(vehicle.proposed_action));
          setToast(null);
          setOpen(true);
        }}
      >
        Proposed action
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
              Proposed action
            </h2>
            <p className="mt-2 text-sm text-[#4e463c]">
              Record the intention for this aging vehicle. Hold is a decision,
              distinct from having none. Clear with None.
            </p>
            <label className="mt-4 flex flex-col gap-1 text-sm text-[#4e463c]">
              Proposed action
              <select
                className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                value={proposedAction}
                onChange={(event) =>
                  setProposedAction(event.target.value as ProposedAction)
                }
              >
                {PROPOSED_ACTION_OPTIONS.map((option) => (
                  <option key={option.value || "none"} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
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
                disabled={save.isPending}
              >
                Save
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
