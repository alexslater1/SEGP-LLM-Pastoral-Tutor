import { useChatSessionHistory } from "./use-chat-history";
import { ITEMS_PER_PAGE } from "@/components/paged-object-table";

export function useUserChatHistory(userID: string | null, page: number) {
  const { history, isLoading, error } = useChatSessionHistory(userID);

  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  const slicedChats = history.slice(start, end);

  return {
    data: {
      data: slicedChats,
      totalPages: Math.ceil((history.length || 0) / ITEMS_PER_PAGE),
    },
    error: error ? new Error(error) : null,
    isLoading,
  }
}

