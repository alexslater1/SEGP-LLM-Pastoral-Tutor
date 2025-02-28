import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getVotesByChatID, vote, removeVote, Vote } from "@/lib/supabase/vote";
import { useEffect, useState } from "react";
import { User } from "@/lib/supabase/user";

export type MessageVotesItems = {
  downvotedMessages: Vote[];
  upvotedMessages: Vote[];
  isLoading: boolean;
  error: Error | null;
  downvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
  upvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeUpvoteMessage: (messageId: string) => Promise<void>;
}

export function useMessagesVotes({chatId, user}: {chatId: string | null, user: User | null}): MessageVotesItems {
  const queryClient = useQueryClient();
  const [downvotedMessages, setDownvotedMessages] = useState<Vote[]>([]);
  const [upvotedMessages, setUpvotedMessages] = useState<Vote[]>([]);

  const { data, isLoading, error } = useQuery({
    queryKey: ['messages-votes', chatId],
    queryFn: async () => {
      if (!chatId) {
        return [];
      }
      const databaseDownvotes = await getVotesByChatID('downvote', chatId);
      const databaseUpvotes = await getVotesByChatID('upvote', chatId);
      setDownvotedMessages(databaseDownvotes);
      setUpvotedMessages(databaseUpvotes);

      return databaseDownvotes;
    },
    staleTime: Infinity,
  });

  useEffect(() => {
    return () => {
      queryClient.invalidateQueries({ queryKey: ['messages-votes', chatId] });
    };
  }, []);

  const downvoteMessage = async (messageId: string, reason?: string) => {
    if (!user) {
      throw new Error("User not found");
    }
    const data = await vote('downvote', messageId, reason ?? "", user);
    if (upvotedMessages.find((vote) => vote.request_id === messageId)) {
      removeVote('upvote', messageId);
      setUpvotedMessages((prev) => prev.filter((vote) => vote.request_id !== messageId));
    }
    setDownvotedMessages((prev) => [...prev, data]);
  }

  const removeDownvoteMessage = async (messageId: string) => {
    const data = await removeVote('downvote', messageId);
    setDownvotedMessages((prev) => prev.filter((vote) => vote.id !== data.id));
  }

  const upvoteMessage = async (messageId: string, reason?: string) => {
    if (!user) {
      throw new Error("User not found");
    }
    const data = await vote('upvote', messageId, reason ?? "", user);
    if (downvotedMessages.find((vote) => vote.request_id === messageId)) {
      removeVote('downvote', messageId);
      setDownvotedMessages((prev) => prev.filter((vote) => vote.request_id !== messageId));
    }
    setUpvotedMessages((prev) => [...prev, data]);
  }

  const removeUpvoteMessage = async (messageId: string) => {
    const data = await removeVote('upvote', messageId);
    setUpvotedMessages((prev) => prev.filter((vote) => vote.id !== data.id));
  }

  return { 
    downvotedMessages,
    upvotedMessages,
    isLoading, 
    error, 
    downvoteMessage, 
    removeDownvoteMessage, 
    upvoteMessage, 
    removeUpvoteMessage 
  };
}
