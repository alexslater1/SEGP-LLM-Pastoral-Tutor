import { fetchDownvotes } from "@/app/(admin)/actions";
import { useQuery } from "@tanstack/react-query";

export function useAllDownvotes(page: number) {
  const { data, error, isLoading } = useQuery({
    queryKey: ['downvotes', page],
    queryFn: () => fetchDownvotes(page),
  });

  return {
    data,
    error,
    isLoading,
  };
}

