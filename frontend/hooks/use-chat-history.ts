import { useEffect, useRef, useState } from "react";
import { Session } from "@supabase/supabase-js";
import { getUserSession } from "@/lib/supabase/client";
import { useQuery, useQueryClient } from "@tanstack/react-query";

export type ChatVisibility = "public" | "private";

const STATUS_QUERY_INTERVAL_SECONDS = 0.5;

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

type BackendUserSessions = BackendUserSession[];

export function useChatSessionHistory(): ChatHistoryItem {
  const queryClient = useQueryClient();
  const [history, setHistory] = useState<Chat[]>([]);
  const chatIDPollingKeys = useRef<string[]>([]);

  const { isPending, error, fetchStatus } = useQuery({
    queryKey: ["chat-history"],
    queryFn: async () => {
      const history = await fetchUIChatHistory();
      setHistory(history);
      chatIDPollingKeys.current = getSessionsWithNoNames(history);
      return history;
    },
    staleTime: Infinity,
  });

  useEffect(() => {
    return () => {
      queryClient.invalidateQueries({ queryKey: ["chat-history"] });
    };
  }, []);

  const {
    isPending: isSessionStatusPending,
    error: sessionStatusError,
  } = useQuery({
    queryKey: ["session-status", chatIDPollingKeys],
    queryFn: async () => {
      let newSessionNames = []
      for (let chatIDPollingKey of chatIDPollingKeys.current) {
        let sessionName = await fetchUIChatName(chatIDPollingKey);
        if (sessionName) {
          newSessionNames.push({id: chatIDPollingKey, name: sessionName});
          chatIDPollingKeys.current = chatIDPollingKeys.current.filter(key => key !== chatIDPollingKey);
        }
      }
      setSessionIDNames(newSessionNames);

      // Return something so TanStack doesn't complain
      return newSessionNames;
    },
    enabled: chatIDPollingKeys.current.length !== 0,
    refetchInterval: chatIDPollingKeys
      ? STATUS_QUERY_INTERVAL_SECONDS * 1000
      : false,
    refetchIntervalInBackground: false,
  });

  const setSessionIDNames = (sessionNames: {id: string, name: string}[]) => {
    setHistory(prev => prev.map(chat => {
      const session = sessionNames.find(sessionName => sessionName.id === chat.id)
      return session ? {...chat, title: session.name} : chat;
    }));
  }

  const getSessionsWithNoNames = (history: Chat[]) => {
    return history.filter(chat => chat.title === "Loading...").map(chat => chat.id);
  }

  const fetchUIChatHistory = async () => {
    const login_session = await getUserSession();
    if (!login_session) {
      return [];
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

  const fetchUIChatName = async (chatSessionID: string) => {
    const login_session = await getUserSession();
    if (!login_session) {
      throw new Error("No login session found");
    }

    const sessionStatus = await fetchSessionName(login_session, chatSessionID);
    return sessionStatus?.name ?? null;
  };

  const refresh = () => {
    queryClient.invalidateQueries({ queryKey: ["chat-history"] });
  };

  return {
    history: history,
    isLoading: isPending || (isSessionStatusPending && fetchStatus !== "idle"),
    refresh: refresh,
    error: error?.message ?? sessionStatusError?.message ?? null,
  };
}

async function fetchChatHistory(session: Session): Promise<BackendUserSessions> {
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

    return (await response.json()) as BackendUserSessions;
  } catch (currentError) {
    console.error("Fetch Chat History Error:", currentError);
    throw currentError;
  }
}

const fetchSessionName = async (loginSession: Session, chatSessionID: string) => {
  try {
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + "/sessions/" + chatSessionID,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${loginSession.access_token}`,
        },
      }
    );

    return (await response.json()) as BackendUserSession;
  } catch (currentError) {
    console.error("Fetch Session Name Error:", currentError);
    return null;
  }
};