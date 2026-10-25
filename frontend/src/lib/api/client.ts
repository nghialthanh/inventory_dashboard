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
};

export class ApiError extends Error {
  readonly status: number;
  readonly errorCode: string;

  constructor(status: number, errorCode: string) {
    super(errorCode || "request failed");
    this.status = status;
    this.errorCode = errorCode;
  }
}

export function asList<T>(data: T[] | null | undefined): T[] {
  return data ?? [];
}

export async function getData<T>(path: string): Promise<T> {
  const response = await fetch(path);
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
    throw new ApiError(response.status, errorCode);
  }
  return body.data as T;
}
