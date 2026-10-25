import { useState } from "react";

import { useMakes, useModels } from "../../lib/data";

export function InventoryFilters() {
  const [makeId, setMakeId] = useState("");
  const [modelId, setModelId] = useState("");
  const makes = useMakes();
  const models = useModels(makeId);

  function chooseMake(nextMakeId: string) {
    setMakeId(nextMakeId);
    setModelId("");
  }

  const failed = makes.isError || models.isError;

  return (
    <form
      className="mt-6 flex flex-wrap items-end gap-4"
      onSubmit={(event) => event.preventDefault()}
    >
      <label className="flex min-w-44 flex-col gap-1 text-sm text-[#4e463c]">
        Make
        <select
          className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
          value={makeId}
          disabled={makes.isPending}
          onChange={(event) => chooseMake(event.target.value)}
        >
          <option value="">All makes</option>
          {(makes.data ?? []).map((make) => (
            <option key={make.id} value={make.id}>
              {make.name}
            </option>
          ))}
        </select>
      </label>
      <label className="flex min-w-44 flex-col gap-1 text-sm text-[#4e463c]">
        Model
        <select
          className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
          value={modelId}
          disabled={models.isPending}
          onChange={(event) => setModelId(event.target.value)}
        >
          <option value="">All models</option>
          {(models.data ?? []).map((model) => (
            <option key={model.id} value={model.id}>
              {model.name}
            </option>
          ))}
        </select>
      </label>
      {failed ? (
        <p
          className="fixed right-4 bottom-4 border border-[#e7c3ba] bg-[#fff7f4] px-4 py-3 text-sm text-[#8a3b2c] shadow-lg"
          role="status"
        >
          Could not load the catalog.
        </p>
      ) : null}
    </form>
  );
}
