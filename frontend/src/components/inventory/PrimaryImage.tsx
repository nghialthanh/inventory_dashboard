import { useState, type FormEvent } from "react";

import { ApiError, type ApiVehicle } from "../../lib/api/client";
import { vehiclePrimaryImageUrl } from "../../lib/api/vehicles";
import { useDeletePrimaryImage, useUploadPrimaryImage } from "../../lib/data";

type PrimaryImageProps = {
  vehicle: ApiVehicle;
  compact?: boolean;
};

export function PrimaryImage({ vehicle, compact = false }: PrimaryImageProps) {
  const dialogId = `primary-image-${vehicle.id}`;
  const inStock = vehicle.stock_state === "in_stock";
  const [open, setOpen] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [fieldError, setFieldError] = useState<string | null>(null);
  const [toast, setToast] = useState<string | null>(null);
  const upload = useUploadPrimaryImage(vehicle.id);
  const remove = useDeletePrimaryImage(vehicle.id);

  function close() {
    setOpen(false);
    setFile(null);
    setFieldError(null);
  }

  async function submit(event: FormEvent) {
    event.preventDefault();
    if (!file) {
      setFieldError("Choose a JPEG, PNG, or WebP file.");
      return;
    }
    setFieldError(null);
    try {
      await upload.mutateAsync(file);
      close();
      setToast(
        vehicle.has_primary_image
          ? "Primary image replaced."
          : "Primary image attached.",
      );
    } catch (error) {
      setToast(
        error instanceof ApiError
          ? error.message
          : "Could not attach the primary image.",
      );
    }
  }

  async function onRemove() {
    try {
      await remove.mutateAsync();
      close();
      setToast("Primary image removed.");
    } catch (error) {
      setToast(
        error instanceof ApiError
          ? error.message
          : "Could not remove the primary image.",
      );
    }
  }

  const frame = compact
    ? "h-12 w-16 border border-[#d9d1c3] bg-[#efeae1]"
    : "h-40 w-56 border border-[#d9d1c3] bg-[#efeae1]";

  return (
    <div className={compact ? undefined : "mt-4"}>
      {vehicle.has_primary_image ? (
        <img
          src={vehiclePrimaryImageUrl(vehicle.id)}
          alt={`Primary image for ${vehicle.vin}`}
          className={`${frame} object-cover`}
        />
      ) : (
        <div
          className={`${frame} flex items-center justify-center text-xs text-[#8a8175]`}
          aria-label="No primary image"
        >
          No primary image
        </div>
      )}
      {compact ? null : (
        <>
          <button
            type="button"
            className="mt-3 border border-[#1c2b24] px-3 py-1 text-sm text-[#1c2b24] disabled:opacity-40"
            disabled={!inStock}
            onClick={() => {
              setFile(null);
              setFieldError(null);
              setToast(null);
              setOpen(true);
            }}
          >
            {vehicle.has_primary_image ? "Replace primary image" : "Attach primary image"}
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
                  Primary image
                </h2>
                <p className="mt-2 text-sm text-[#4e463c]">
                  Upload one JPEG, PNG, or WebP file of at most 5 MiB.
                </p>
                <label className="mt-4 flex flex-col gap-1 text-sm text-[#4e463c]">
                  Image file
                  <input
                    type="file"
                    accept="image/jpeg,image/png,image/webp"
                    className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
                    onChange={(event) =>
                      setFile(event.target.files?.[0] ?? null)
                    }
                  />
                  {fieldError ? (
                    <span className="text-[#8a3b2c]">{fieldError}</span>
                  ) : null}
                </label>
                <div className="mt-6 flex flex-wrap justify-end gap-3">
                  {vehicle.has_primary_image ? (
                    <button
                      type="button"
                      className="mr-auto px-4 py-2 text-sm text-[#8a3b2c]"
                      disabled={remove.isPending}
                      onClick={onRemove}
                    >
                      Remove
                    </button>
                  ) : null}
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
                    disabled={upload.isPending}
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
      )}
    </div>
  );
}
