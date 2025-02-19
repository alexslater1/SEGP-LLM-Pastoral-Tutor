import { useEffect, useState } from "react";
import { Session } from "@supabase/supabase-js";
import { getUserSession } from "@/lib/supabase/client";

export type ChatVisibility = "public" | "private";

export type Chat = {
  id: string;
  title: string;
  createdAt: Date;
  userId: string;
  visibility: ChatVisibility;
}

export type ChatHistoryItem = {
  history: Chat[];
  isLoading: boolean;
  refresh: () => void;
  error: string | null;
}

type BackendUserSession = {
    id: string,
    created_at: string,
    created_by_request_id: string,
    user_id: string,
}

type BackendUserSessionsResponse = BackendUserSession[];

type ChatHistoryState = {
  history: Chat[];
  isLoading: boolean;
  error: string | null;
  attemptedInitialFetch: boolean;
}

export function useChatHistory(): ChatHistoryItem {
  const [chatHistoryState, setChatHistoryState] = useState<ChatHistoryState>({
    history: [],
    isLoading: false,
    error: null,
    attemptedInitialFetch: false,
  });

  let newChatHistoryState = chatHistoryState;

  const setHistory = (history: Chat[]) => {
    newChatHistoryState = {
      ...newChatHistoryState,
      history: history,
    };
  }

  const setIsLoading = (isLoading: boolean) => {
    newChatHistoryState = {
      ...newChatHistoryState,
      isLoading: isLoading,
    };
  }

  const setError = (error: string | null) => {
    newChatHistoryState = {
      ...newChatHistoryState,
      error: error,
    };
  }

  const setAttemptedInitialFetch = (attemptedInitialFetch: boolean) => {
    newChatHistoryState = {
      ...newChatHistoryState,
      attemptedInitialFetch: attemptedInitialFetch,
    };
  }

  const updateChatHistoryState = () => {
    setChatHistoryState(newChatHistoryState);
  }

  useEffect(() => {
    if (chatHistoryState.history.length === 0 &&
        !chatHistoryState.isLoading &&
        !chatHistoryState.attemptedInitialFetch) {
      setAttemptedInitialFetch(true);
      setIsLoading(true);
      updateChatHistoryState();
    } else if (chatHistoryState.isLoading) {
      fetchAndUpdateChatHistory();
    }
  }, [chatHistoryState]);

  const refresh = () => {
    setIsLoading(true);
    updateChatHistoryState();
  }

  const fetchAndUpdateChatHistory = async () => {
    const login_session = await getUserSession();
    if (!login_session) {
      setError("User not logged in");
      setIsLoading(false);
      setHistory([]);
      updateChatHistoryState();
      return;
    }

    const chatHistory = await fetchChatHistory(login_session);
    const history = chatHistory.map((session) => {
      return {
        id: session.id,
        title: "Chat " + session.id.slice(0, 8), // TODO: Add title
        createdAt: new Date(session.created_at),
        userId: session.user_id,
        visibility: "public", // TODO: Add visibility
      } as Chat;
    })
    setHistory(history);
    setIsLoading(false);
    setError(null);
    updateChatHistoryState();
  }

  return {
    history: structuredClone(chatHistoryState.history).reverse(),
    isLoading: chatHistoryState.isLoading,
    refresh: refresh,
    error: chatHistoryState.error,
  };
}

async function fetchChatHistory(session: Session) {
  try {
    const completionEndpoint = "/sessions"
    const response = await fetch(process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + completionEndpoint, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${session.access_token}`,
        }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return await response.json() as BackendUserSessionsResponse;
  } catch (currentError) {
    console.error('Error:', currentError);
    return [];
  }
}
