import { asList, getData, type ApiMake, type ApiModel } from "./client";

export type { ApiMake, ApiModel };

export function listMakes(): Promise<ApiMake[]> {
  return getData<ApiMake[] | null>("/api/v1/catalog/makes").then(asList);
}

export function listModels(makeId: string): Promise<ApiModel[]> {
  const query = makeId ? `?make_id=${encodeURIComponent(makeId)}` : "";
  return getData<ApiModel[] | null>(`/api/v1/catalog/models${query}`).then(asList);
}
