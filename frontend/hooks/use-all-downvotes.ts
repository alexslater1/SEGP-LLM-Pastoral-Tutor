import { fetchDownvotes } from "@/app/(admin)/actions";
import { useQuery } from "@tanstack/react-query";

export function useAllDownvotes(page: number) {
  return useQuery({
    queryKey: ['downvotes', page],
    queryFn: () => fetchDownvotes(page),
  });
}

