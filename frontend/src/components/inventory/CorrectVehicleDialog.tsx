import { useState, type FormEvent } from "react";

import { ApiError, type ApiVehicle } from "../../lib/api/client";
import { useCorrectVehicle, useMakes, useModels } from "../../lib/data";

type FieldErrors = {
  vin?: string;
  model?: string;
  enteredOn?: string;
};

export function CorrectVehicleDialog({ vehicle }: { vehicle: ApiVehicle }) {
  const dialogId = `correct-vehicle-${vehicle.id}`;
  const [open, setOpen] = useState(false);
  const [vin, setVin] = useState(vehicle.vin);
  const [makeId, setMakeId] = useState(vehicle.make.id);
  const [modelId, setModelId] = useState(vehicle.model.id);
  const [enteredOn, setEnteredOn] = useState(vehicle.entered_on);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [toast, setToast] = useState<string | null>(null);
  const makes = useMakes();
  const models = useModels(makeId);
  const correct = useCorrectVehicle(vehicle.id);

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

    const input: { vin?: string; model_id?: string; entered_on?: string } = {};
    if (vin.trim() !== vehicle.vin) {
      input.vin = vin;
    }
    if (modelId !== vehicle.model.id) {
      input.model_id = modelId;
    }
    if (enteredOn !== vehicle.entered_on) {
      input.entered_on = enteredOn;
    }
    if (Object.keys(input).length === 0) {
      setToast("Change the VIN, model, or entered on.");
      return;
    }

    try {
      await correct.mutateAsync(input);
      setOpen(false);
      setToast("Vehicle corrected.");
    } catch (error) {
      const message =
        error instanceof ApiError ? error.message : "Could not correct the vehicle.";
      setToast(message);
    }
  }

  return (
    <>
      <button
        type="button"
        className="mt-4 border border-[#1c2b24] px-3 py-1 text-sm text-[#1c2b24]"
        onClick={() => {
          setVin(vehicle.vin);
          setMakeId(vehicle.make.id);
          setModelId(vehicle.model.id);
          setEnteredOn(vehicle.entered_on);
          setErrors({});
          setToast(null);
          setOpen(true);
        }}
      >
        Correct
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
              Correct vehicle
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
                    <option key={make.id} value={make.id}>
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
                    <option key={model.id} value={model.id}>
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
                disabled={correct.isPending}
              >
                Correct
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
