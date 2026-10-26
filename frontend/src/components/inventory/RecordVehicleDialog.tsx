import { useState, type FormEvent } from "react";

import { ApiError } from "../../lib/api/client";
import { useMakes, useModels, useCreateVehicle } from "../../lib/data";

type FieldErrors = {
  vin?: string;
  model?: string;
  enteredOn?: string;
};

export function RecordVehicleDialog() {
  const dialogId = "record-vehicle";
  const [open, setOpen] = useState(false);
  const [vin, setVin] = useState("");
  const [makeId, setMakeId] = useState("");
  const [modelId, setModelId] = useState("");
  const [enteredOn, setEnteredOn] = useState("");
  const [errors, setErrors] = useState<FieldErrors>({});
  const [toast, setToast] = useState<string | null>(null);
  const makes = useMakes();
  const models = useModels(makeId);
  const create = useCreateVehicle();

  function close() {
    setOpen(false);
    setErrors({});
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    const next: FieldErrors = {};
    if (vin.trim() === "") {
      next.vin = "Enter a VIN.";
    }
    if (modelId === "") {
      next.model = "Choose a model.";
    }
    if (enteredOn === "") {
      next.enteredOn = "Enter entered on.";
    }
    setErrors(next);
    if (Object.keys(next).length > 0) {
      return;
    }

    try {
      await create.mutateAsync({
        vin,
        model_id: Number(modelId),
        entered_on: enteredOn,
      });
      setVin("");
      setMakeId("");
      setModelId("");
      setEnteredOn("");
      setOpen(false);
      setToast("Vehicle recorded.");
    } catch (error) {
      const message =
        error instanceof ApiError ? error.message : "Could not record the vehicle.";
      setToast(message);
    }
  }

  return (
    <>
      <button
        type="button"
        className="border border-[#f6f1e7] bg-[#f6f1e7] px-4 py-2 text-sm text-[#1c2b24]"
        onClick={() => {
          setToast(null);
          setOpen(true);
        }}
      >
        Record a vehicle
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
              Record a vehicle
            </h2>
            <div className="mt-4 flex flex-col gap-4">
              <label className="flex flex-col gap-1 text-sm text-[#4e463c]">
                VIN
                <input
                  className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                  value={vin}
                  autoComplete="off"
                  onChange={(event) => setVin(event.target.value)}
                />
                {errors.vin ? (
                  <span className="text-[#8a3b2c]">{errors.vin}</span>
                ) : null}
              </label>
              <label className="flex flex-col gap-1 text-sm text-[#4e463c]">
                Make
                <select
                  className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                  value={makeId}
                  onChange={(event) => {
                    setMakeId(event.target.value);
                    setModelId("");
                  }}
                >
                  <option value="">Choose a make</option>
                  {(makes.data ?? []).map((make) => (
                    <option key={make.id} value={String(make.id)}>
                      {make.name}
                    </option>
                  ))}
                </select>
              </label>
              <label className="flex flex-col gap-1 text-sm text-[#4e463c]">
                Model
                <select
                  className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                  value={modelId}
                  onChange={(event) => setModelId(event.target.value)}
                >
                  <option value="">Choose a model</option>
                  {(models.data ?? []).map((model) => (
                    <option key={model.id} value={String(model.id)}>
                      {model.name}
                    </option>
                  ))}
                </select>
                {errors.model ? (
                  <span className="text-[#8a3b2c]">{errors.model}</span>
                ) : null}
              </label>
              <label className="flex flex-col gap-1 text-sm text-[#4e463c]">
                Entered on
                <input
                  type="date"
                  className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                  value={enteredOn}
                  onChange={(event) => setEnteredOn(event.target.value)}
                />
                {errors.enteredOn ? (
                  <span className="text-[#8a3b2c]">{errors.enteredOn}</span>
                ) : null}
              </label>
            </div>
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
                disabled={create.isPending}
              >
                Record
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
