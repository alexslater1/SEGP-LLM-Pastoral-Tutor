import { useEffect, useState } from "react";
import { Session } from "@supabase/supabase-js";
import { getUserSession } from "@/lib/supabase/client";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export type ChatVisibility = "public" | "private";

export type Chat = {
  id: string;
  title: string;
  createdAt: Date;
  userId: string;
  visibility: ChatVisibility;
};

export type ChatHistoryItem = {
  history: Chat[];
  isLoading: boolean;
  refresh: () => void;
  error: string | null;
};

type BackendUserSession = {
  id: string;
  created_at: string;
  created_by_request_id: string;
  user_id: string;
  name?: string;
};

type BackendUserSessionsResponse = BackendUserSession[];

export function useSessionName(id?: string) {
  console.log("useSessionName called with id", id);

  const queryClient = useQueryClient();
  const { isPending, data, error } = useQuery({
    queryKey: ["session-name", id],
    queryFn: () => fetchSessionNameLoop(id!),
    enabled: !!id,
  });

  const mutate = () => {
    queryClient.invalidateQueries({ queryKey: ["session-name", id] });
  };

  const fetchSessionNameLoop = async (id: string) => {
    do {
      console.log("Fetching session name for", id);
      const session = await fetchSessionName(id);
      if (session?.name) {
        console.log("Session name found for", id);
        queryClient.invalidateQueries({ queryKey: ["chat-history"] });
        return session.name;
      }
      await new Promise((resolve) => setTimeout(resolve, 500));
    } while (true);
  };

  const fetchSessionName = async (id: string) => {
    const login_session = await getUserSession();
    if (!login_session) {
      return;
    }

    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + "/sessions/" + id,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${login_session.access_token}`,
        },
      }
    );

    return (await response.json()) as BackendUserSession;
  };

  return {
    isPending,
    data,
    error,
    mutate,
  };
}

export function useChatHistory(): ChatHistoryItem {
  const queryClient = useQueryClient();

  const { isPending, data, error } = useQuery({
    queryKey: ["chat-history"],
    queryFn: () => fetchUIChatHistory(),
  });

  const fetchUIChatHistory = async () => {
    const login_session = await getUserSession();
    if (!login_session) {
      return;
    }

    const chatHistory = await fetchChatHistory(login_session);
    const history = chatHistory.map((session) => {
      return {
        id: session.id,
        title: session.name ?? "Loading...",
        createdAt: new Date(session.created_at),
        userId: session.user_id,
        visibility: "public", // TODO: Add visibility
      } as Chat;
    });
    return history.reverse();
  };

  const refresh = () => {
    queryClient.invalidateQueries({ queryKey: ["chat-history"] });
  };

  return {
    history: data ?? [],
    isLoading: isPending,
    refresh: refresh,
    error: error?.message ?? null,
  };
}

async function fetchChatHistory(session: Session) {
  try {
    const completionEndpoint = "/sessions";
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + completionEndpoint,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${session.access_token}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return (await response.json()) as BackendUserSessionsResponse;
  } catch (currentError) {
    console.error("Error:", currentError);
    return [];
  }
}
