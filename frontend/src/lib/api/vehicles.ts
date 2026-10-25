import {
  getData,
  getList,
  postJSON,
  type ApiVehicle,
} from "./client";

export type VehicleListQuery = {
  makeId: string;
  modelId: string;
  page: number;
};

export type RecordVehicleInput = {
  vin: string;
  model_id: string;
  entered_on: string;
};

export function listVehicles(query: VehicleListQuery) {
  const params = new URLSearchParams();
  if (query.makeId) {
    params.set("make_id", query.makeId);
  }
  if (query.modelId) {
    params.set("model_id", query.modelId);
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
