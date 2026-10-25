import { useQuery } from "@tanstack/react-query";

import { listMakes, listModels } from "./api/catalog";

export function useMakes() {
  return useQuery({
    queryKey: ["catalog", "makes"],
    queryFn: listMakes,
  });
}

export function useModels(makeId: string) {
  return useQuery({
    queryKey: ["catalog", "models", makeId],
    queryFn: () => listModels(makeId),
  });
}
