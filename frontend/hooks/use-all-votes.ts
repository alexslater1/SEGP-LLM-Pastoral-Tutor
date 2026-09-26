import { fetchVotes } from "@/app/(admin)/actions";
import { useQuery } from "@tanstack/react-query";

export function useAllDownvotes(page: number) {
  return useQuery({
    queryKey: ['downvotes', page],
    queryFn: () => fetchVotes('downvote', page),
  });
}

export function useAllUpvotes(page: number) {
  return useQuery({
    queryKey: ['upvotes', page],
    queryFn: () => fetchVotes('upvote', page),
  });
}

export function useAllVotes(page: number) {
  return useQuery({
    queryKey: ['votes', page],
    queryFn: () => fetchVotes('both', page),
  });
}
