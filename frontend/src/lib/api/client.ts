export type ApiMake = {
  id: string;
  name: string;
};

export type ApiModel = {
  id: string;
  make_id: string;
  name: string;
};

type Envelope<T> = {
  code: number;
  message: string;
  data: T;
  meta?: ApiMeta;
};

export type ApiMeta = {
  page: number;
  page_size: number;
  total: number;
};

export type ApiRef = {
  id: string;
  name: string;
};

export type ApiVehicle = {
  id: string;
  vin: string;
  make: ApiRef;
  model: ApiRef;
  entered_on: string;
  stocked_out_on: string | null;
  stock_state: "in_stock" | "stocked_out";
  inventory_age: number;
  is_aging: boolean;
  proposed_action: string | null;
  has_primary_image: boolean;
};

export class ApiError extends Error {
  readonly status: number;
  readonly errorCode: string;

  constructor(status: number, errorCode: string, message = "request failed") {
    super(message || errorCode || "request failed");
    this.status = status;
    this.errorCode = errorCode;
  }
}

export function asList<T>(data: T[] | null | undefined): T[] {
  return data ?? [];
}

async function readEnvelope<T>(response: Response): Promise<Envelope<T>> {
  const body = (await response.json()) as Envelope<T> & {
    data?: { error_code?: string } | T;
  };
  if (!response.ok) {
    const errorCode =
      body.data &&
      typeof body.data === "object" &&
      "error_code" in body.data &&
      typeof body.data.error_code === "string"
        ? body.data.error_code
        : "";
    throw new ApiError(response.status, errorCode, body.message || "request failed");
  }
  return body;
}

export async function getData<T>(path: string): Promise<T> {
  const response = await fetch(path);
  const body = await readEnvelope<T>(response);
  return body.data;
}

export async function postJSON<T>(path: string, payload: unknown): Promise<T> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const body = await readEnvelope<T>(response);
  return body.data;
}

export async function getList<T>(
  path: string,
): Promise<{ items: T[]; meta: ApiMeta }> {
  const response = await fetch(path);
  const body = await readEnvelope<T[] | null>(response);
  if (!body.meta) {
    throw new ApiError(response.status, "", "request failed");
  }
  return { items: asList(body.data), meta: body.meta };
}
