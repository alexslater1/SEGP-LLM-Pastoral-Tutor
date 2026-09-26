import { useQuery } from "@tanstack/react-query";
import { fetchAgentEvents, fetchAgentRequests } from "@/app/(admin)/actions";

export function useAgentEvents(page: number) {
  return useQuery({
    queryKey: ["agent-events", page],
    queryFn: () => fetchAgentEvents(page),
  });
}

export function useAgentRequests(page: number) {
  return useQuery({
    queryKey: ["agent-requests", page],
    queryFn: () => fetchAgentRequests(page),
  });
} 