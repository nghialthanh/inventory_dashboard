import type { ProposedAction } from "./api/vehicles";

export const PROPOSED_ACTION_OPTIONS: {
  value: ProposedAction;
  label: string;
}[] = [
  { value: "", label: "None" },
  { value: "price_reduction_planned", label: "Price reduction planned" },
  { value: "wholesale_planned", label: "Wholesale planned" },
  { value: "hold", label: "Hold" },
];

export function proposedActionLabel(value: string | null): string {
  const match = PROPOSED_ACTION_OPTIONS.find((option) => option.value === (value ?? ""));
  return match?.label ?? "None";
}

export function asProposedAction(value: string | null): ProposedAction {
  if (
    value === "price_reduction_planned" ||
    value === "wholesale_planned" ||
    value === "hold"
  ) {
    return value;
  }
  return "";
}
