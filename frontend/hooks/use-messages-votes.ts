import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getDownvotesByChatID, downvote, removeDownvote, Downvote } from "@/lib/supabase/vote";
import { useEffect, useState } from "react";
import { User } from "@/lib/supabase/user";

export type MessageVotesItems = {
  downvotedMessages: Downvote[];
  isLoading: boolean;
  error: Error | null;
  downvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
}

export function useMessagesVotes({chatId, user}: {chatId: string | null, user: User | null}): MessageVotesItems {
  const queryClient = useQueryClient();
  const [downvotedMessages, setDownvotedMessages] = useState<Downvote[]>([]);

  const { data, isLoading, error } = useQuery({
    queryKey: ['messages-votes', chatId],
    queryFn: async () => {
      if (!chatId) {
        return [];
      }
      const databaseDownvotes = await getDownvotesByChatID(chatId);
      setDownvotedMessages(mapDownvotedMessages(databaseDownvotes));

      return databaseDownvotes;
    },
    staleTime: Infinity,
  });

  useEffect(() => {
    return () => {
      queryClient.invalidateQueries({ queryKey: ['messages-votes', chatId] });
    };
  }, []);

  const mapDownvotedMessage = (vote: Downvote) => {
    return {
      id: vote.id,
      created_at: vote.created_at,
      request_id: vote.request_id,
      user_id: vote.user_id,
      reason: vote.reason,
    } as Downvote;
  }

  const mapDownvotedMessages = (data: Downvote[]) => data.map((vote) => {
    return mapDownvotedMessage(vote);
  });

  const downvoteMessage = async (messageId: string, reason?: string) => {
    if (!user) {
      throw new Error("User not found");
    }
    const data = await downvote(messageId, reason ?? "", user);
    setDownvotedMessages((prev) => [...prev, mapDownvotedMessage(data)]);
  }

  const removeDownvoteMessage = async (messageId: string) => {
    const data = await removeDownvote(messageId);
    setDownvotedMessages((prev) => prev.filter((vote) => vote.id !== data.id));
  }

  return { downvotedMessages, isLoading, error, downvoteMessage, removeDownvoteMessage };
}
