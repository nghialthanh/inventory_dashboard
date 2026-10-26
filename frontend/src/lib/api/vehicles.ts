import {
  getData,
  getList,
  patchJSON,
  postJSON,
  type ApiVehicle,
} from "./client";

export type VehicleListQuery = {
  makeId: string;
  modelId: string;
  agingOnly: boolean;
  page: number;
};

export type RecordVehicleInput = {
  vin: string;
  model_id: string;
  entered_on: string;
};

export type CorrectVehicleInput = {
  vin?: string;
  model_id?: string;
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
  if (query.agingOnly) {
    params.set("aging", "true");
  }
  params.set("page", String(query.page));
  return getList<ApiVehicle>(`/api/v1/vehicles?${params.toString()}`).then(
    (page) => ({ vehicles: page.items, meta: page.meta }),
  );
}

export function getVehicle(id: string): Promise<ApiVehicle> {
  return getData<ApiVehicle>(`/api/v1/vehicles/${encodeURIComponent(id)}`);
}

export function createVehicle(input: RecordVehicleInput): Promise<ApiVehicle> {
  return postJSON<ApiVehicle>("/api/v1/vehicles", input);
}

export function correctVehicle(
  id: string,
  input: CorrectVehicleInput,
): Promise<ApiVehicle> {
  return patchJSON<ApiVehicle>(
    `/api/v1/vehicles/${encodeURIComponent(id)}`,
    input,
  );
}

export function stockOutVehicle(
  id: string,
  input: StockOutVehicleInput,
): Promise<ApiVehicle> {
  return postJSON<ApiVehicle>(
    `/api/v1/vehicles/${encodeURIComponent(id)}/stock-out`,
    input,
  );
}
