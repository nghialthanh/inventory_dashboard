import { useMakes, useModels } from "../../lib/data";

type InventoryFiltersProps = {
  makeId: string;
  modelId: string;
  stockState: "" | "in_stock" | "stocked_out";
  stockedOutOn: string;
  agingOnly: boolean;
  onMakeId: (makeId: string) => void;
  onModelId: (modelId: string) => void;
  onStockState: (stockState: "" | "in_stock" | "stocked_out") => void;
  onStockedOutOn: (stockedOutOn: string) => void;
  onAgingOnly: (agingOnly: boolean) => void;
};

export function InventoryFilters({
  makeId,
  modelId,
  stockState,
  stockedOutOn,
  agingOnly,
  onMakeId,
  onModelId,
  onStockState,
  onStockedOutOn,
  onAgingOnly,
}: InventoryFiltersProps) {
  const makes = useMakes();
  const models = useModels(makeId);

  function chooseMake(nextMakeId: string) {
    onMakeId(nextMakeId);
    onModelId("");
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
          onChange={(event) => onModelId(event.target.value)}
        >
          <option value="">All models</option>
          {(models.data ?? []).map((model) => (
            <option key={model.id} value={model.id}>
              {model.name}
            </option>
          ))}
        </select>
      </label>
      <label className="flex min-w-44 flex-col gap-1 text-sm text-[#4e463c]">
        Stock state
        <select
          className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
          value={stockState}
          onChange={(event) =>
            onStockState(event.target.value as "" | "in_stock" | "stocked_out")
          }
        >
          <option value="">All stock states</option>
          <option value="in_stock">In stock</option>
          <option value="stocked_out">Stocked out</option>
        </select>
      </label>
      <label className="flex min-w-44 flex-col gap-1 text-sm text-[#4e463c]">
        Stocked-out on
        <input
          type="date"
          className="border border-[#d9d1c3] bg-white px-3 py-2 text-base text-[#1c1915]"
          value={stockedOutOn}
          onChange={(event) => onStockedOutOn(event.target.value)}
        />
      </label>
      <label className="flex items-center gap-2 pb-2 text-sm text-[#4e463c]">
        <input
          type="checkbox"
          className="size-4 accent-[#1c2b24]"
          checked={agingOnly}
          onChange={(event) => onAgingOnly(event.target.checked)}
        />
        Aging stock only
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
