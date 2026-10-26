import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { listMakes, listModels } from "./api/catalog";
import {
  correctVehicle,
  createVehicle,
  getVehicle,
  listVehicles,
  setProposedAction,
  stockOutVehicle,
  type CorrectVehicleInput,
  type RecordVehicleInput,
  type SetProposedActionInput,
  type StockOutVehicleInput,
  type VehicleListQuery,
} from "./api/vehicles";

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

export function useVehicles(query: VehicleListQuery) {
  return useQuery({
    queryKey: ["vehicles", query],
    queryFn: () => listVehicles(query),
    retry: false,
  });
}

export function useVehicle(id: string) {
  return useQuery({
    queryKey: ["vehicles", id],
    queryFn: () => getVehicle(id),
    enabled: id !== "",
    retry: false,
  });
}

export function useCorrectVehicle(id: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CorrectVehicleInput) => correctVehicle(id, input),
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: ["vehicles"] });
    },
  });
}

export function useCreateVehicle() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: RecordVehicleInput) => createVehicle(input),
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: ["vehicles"] });
    },
  });
}

export function useStockOutVehicle(id: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: StockOutVehicleInput) => stockOutVehicle(id, input),
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: ["vehicles"] });
    },
  });
}

export function useSetProposedAction(id: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: SetProposedActionInput) => setProposedAction(id, input),
    onSuccess: async () => {
      await client.invalidateQueries({ queryKey: ["vehicles"] });
    },
  });
}
