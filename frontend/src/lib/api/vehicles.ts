import {
  getData,
  getList,
  patchJSON,
  postJSON,
  putJSON,
  type ApiVehicle,
} from "./client";

export type ProposedAction =
  | "price_reduction_planned"
  | "wholesale_planned"
  | "hold"
  | "";

export type VehicleListQuery = {
  makeId: string;
  modelId: string;
  stockState: "" | "in_stock" | "stocked_out";
  stockedOutOn: string;
  agingOnly: boolean;
  page: number;
};

export type RecordVehicleInput = {
  vin: string;
  model_id: number;
  entered_on: string;
};

export type CorrectVehicleInput = {
  vin?: string;
  model_id?: number;
  entered_on?: string;
};

export type StockOutVehicleInput = {
  stocked_out_on: string;
};

export function listVehicles(query: VehicleListQuery) {
  const params = new URLSearchParams();
  if (query.makeId) {
    params.set("make_id", query.makeId);
  }
  if (query.modelId) {
    params.set("model_id", query.modelId);
  }
  if (query.stockState) {
    params.set("stock_state", query.stockState);
  }
  if (query.stockedOutOn) {
    params.set("stocked_out_on", query.stockedOutOn);
  }
  if (query.agingOnly) {
    params.set("aging", "true");
  }
  params.set("page", String(query.page));
  return getList<ApiVehicle>(`/api/v1/vehicles?${params.toString()}`).then(
    (page) => ({ vehicles: page.items, meta: page.meta }),
  );
}

export function getVehicle(id: number): Promise<ApiVehicle> {
  return getData<ApiVehicle>(`/api/v1/vehicles/${id}`);
}

export function createVehicle(input: RecordVehicleInput): Promise<ApiVehicle> {
  return postJSON<ApiVehicle>("/api/v1/vehicles", input);
}

export function correctVehicle(
  id: number,
  input: CorrectVehicleInput,
): Promise<ApiVehicle> {
  return patchJSON<ApiVehicle>(`/api/v1/vehicles/${id}`, input);
}

export function stockOutVehicle(
  id: number,
  input: StockOutVehicleInput,
): Promise<ApiVehicle> {
  return postJSON<ApiVehicle>(`/api/v1/vehicles/${id}/stock-out`, input);
}

export type SetProposedActionInput = {
  proposed_action: ProposedAction;
};

export function setProposedAction(
  id: number,
  input: SetProposedActionInput,
): Promise<ApiVehicle> {
  return putJSON<ApiVehicle>(
    `/api/v1/vehicles/${id}/proposed-action`,
    input,
  );
}
