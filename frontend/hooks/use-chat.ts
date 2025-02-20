import { Message, Status, Role } from "@/types/message";
import { generateUUID } from "@/lib/utils";
import { useState, useRef, useEffect, useMemo } from "react";
import { getUserSession } from "@/lib/supabase/client";
import { Session } from "@supabase/supabase-js";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

const STATUS_QUERY_INTERVAL_SECONDS = 1;
const IGNORED_ACTIONS: string[] = []; //["Thinking", "Thinking..."]
const DEFAULT_FIRST_AGENT_MESSAGE =
  "Hi, I'm the Imperial College tutor agent. Ask me anything!";

export type ChatItem = {
  messages: Message[];
  handleSubmit: (query: string) => void;
  isLoading: boolean;
  stop: () => void;
  error: string | null;
  id: string | null;
  isInitialLoad: boolean;
};

export type ChatItemProps = {
  id: string | null;
  firstAgentMessage?: string;
};

export function useChat({
  id,
  firstAgentMessage = DEFAULT_FIRST_AGENT_MESSAGE,
}: ChatItemProps): ChatItem {
  const [messages, setMessages] = useState<Message[]>([]);
  const [requestIDPollingKey, setRequestIDPollingKey] = useState<string | null>(
    null
  );
  const [chatSessionID, setChatSessionID] = useState<string | null>(id);

  const {
    data: allMessages,
    isPending: isAllMessagesPending,
    error: allMessagesError,
    fetchStatus: allMessagesFetchStatus,
  } = useQuery({
    queryKey: ["all-messages", chatSessionID],
    queryFn: async () => loadAllMessages(),
    staleTime: Infinity,
    // If we were given an id when first creating the chat hook, we should attempt to load all messages
    enabled: !!id,
  });

  if (messages.length === 0 && allMessages) {
    setMessages(allMessages);
  }

  const {
    isPending: isStatusPending,
    error: statusError,
    fetchStatus: statusFetchStatus,
  } = useQuery({
    queryKey: ["status", requestIDPollingKey],
    queryFn: async () => {
      let statusResponse = await checkStatus(requestIDPollingKey as string);
      updateLastMessage(statusResponse);
      return statusResponse;
    },
    enabled: !!requestIDPollingKey,
    refetchInterval: requestIDPollingKey
      ? STATUS_QUERY_INTERVAL_SECONDS * 1000
      : false,
    refetchIntervalInBackground: false,
  });

  const newLastMessageFrom = (statusResponse: BackendStatusResponse) => {
    const currentLastMessage = messages[messages.length - 1];

    if (statusResponse.type === Status.COMPLETED) {
      currentLastMessage.status = Status.COMPLETED;
      currentLastMessage.content = statusResponse.answer as string;
      stop();
      return currentLastMessage;
    }

    if (statusResponse.type === Status.FAILED) {
      currentLastMessage.status = Status.FAILED;
      currentLastMessage.content = statusResponse.error as string;
      stop();
      return currentLastMessage;
    }

    currentLastMessage.status = Status.PENDING;
    const newCurrentAction = statusResponse.current_action as string;
    if (
      currentLastMessage.actions.length > 0 &&
      newCurrentAction ==
        currentLastMessage.actions[currentLastMessage.actions.length - 1]
    ) {
      return currentLastMessage;
    }
    currentLastMessage.actions.push(statusResponse.current_action as string);
    return currentLastMessage;
  };

  const updateLastMessage = (statusResponse: BackendStatusResponse) => {
    if (
      messages.length > 0 &&
      messages[messages.length - 1].role === Role.ASSISTANT
    ) {
      const lastMessage = newLastMessageFrom(statusResponse);

      setMessages((messages) => {
        let newMessages = [...messages.slice(0, -1), lastMessage];
        return newMessages;
      });
    }
  };

  const {
    mutate: sendMessage,
    isPending: isSendPending,
    error: sendError,
  } = useMutation({
    mutationFn: ({
      query,
      newMessages,
    }: {
      query: string;
      newMessages: Message[];
    }) => sendMessageAndGetResponse(query),
    onSuccess: (newMessage, { newMessages }) => {
      setMessages([...newMessages, newMessage]);
      setRequestIDPollingKey(newMessage.requestID);
    },
  });

  const handleSubmit = (query: string) => {
    const id = generateUUID();
    const newMessages = [
      ...messages,
      {
        id: id + Role.USER,
        requestID: id,
        content: query,
        role: Role.USER,
        actions: [],
        status: Status.COMPLETED,
      },
    ];

    // add user's message to the messages array
    setMessages(newMessages);
    sendMessage({ query, newMessages });
  };

  const stop = () => {
    setRequestIDPollingKey(null);
  };

  const loadAllMessages = async (): Promise<Message[]> => {
    const login_session = await getUserSession();
    if (!login_session) {
      throw new Error("User not logged in");
    } else if (chatSessionID) {
      const { messages, error: loadError } = await fetchAllMessagesByID(
        chatSessionID,
        login_session
      );
      if (loadError) {
        throw new Error(loadError);
      } else {
        return messages;
      }
    } else {
      return [];
    }
  };

  const sendMessageAndGetResponse = async (query: string): Promise<Message> => {
    const loginSession = await getUserSession();
    if (!loginSession) {
      throw new Error("User not logged in");
    }

    const {
      requestID,
      sessionID: newSessionID,
      error: sendError,
    } = await sendMessageToBackend(chatSessionID, query, loginSession);

    if (sendError) {
      let id = generateUUID();
      return {
        id: id + Role.ASSISTANT,
        requestID: id,
        content: sendError,
        role: Role.ASSISTANT,
        actions: [],
        status: Status.FAILED,
      };
    } else {
      setChatSessionID(newSessionID as string);
      return {
        id: (requestID as string) + Role.ASSISTANT,
        requestID: requestID as string,
        content: "Thinking..",
        role: Role.ASSISTANT,
        actions: [],
        status: Status.PENDING,
      };
    }
  };

  const checkStatus = async (requestID: string) => {
    const login_session = await getUserSession();
    if (!login_session) {
      throw new Error("User not logged in");
    }

    return await checkStatusByRequestID(requestID, login_session);
  };

  const firstMessage = useMemo(() => {
    const id = generateUUID();
    return {
      id: id + Role.ASSISTANT,
      requestID: id,
      content: firstAgentMessage,
      role: Role.ASSISTANT,
      actions: [],
      status: Status.COMPLETED,
    };
  }, [id]);

  return {
    messages: [firstMessage, ...messages],
    handleSubmit,
    isLoading:
      isSendPending || (isStatusPending && statusFetchStatus !== "idle"),
    stop,
    error:
      sendError?.message ||
      statusError?.message ||
      allMessagesError?.message ||
      null,
    id: chatSessionID,
    isInitialLoad: isAllMessagesPending && allMessagesFetchStatus !== "idle",
  };
}

type BackendResponsePastQueryAndAnswer = {
  type: Status;
  answer?: string;
  current_action?: string;
  error?: string;
  query: string;
  request_id: string;
  actions: string[];
};

type BackendReponsePastQueriesAndAnswers = BackendResponsePastQueryAndAnswer[];

type AllBackendMessages = {
  messages: Message[];
  error: string | null;
};

function parseBackendResponse(
  response: BackendReponsePastQueriesAndAnswers
): Message[] {
  let messages: Message[] = [];
  response.forEach((item) => {
    messages.push({
      id: item.request_id + Role.USER,
      requestID: item.request_id,
      content: item.query,
      role: Role.USER,
      actions: [],
      status: Status.COMPLETED,
    });
    messages.push({
      id: item.request_id + Role.ASSISTANT,
      requestID: item.request_id,
      content:
        item.type === Status.COMPLETED
          ? (item.answer as string)
          : item.type === Status.PENDING
          ? "Thinking..."
          : (item.error as string),
      role: Role.ASSISTANT,
      actions: item.actions
        .filter((action) => !IGNORED_ACTIONS.includes(action))
        .concat(
          item.type === Status.PENDING ? [item.current_action as string] : []
        ),
      status: item.type,
    });
  });
  return messages;
}

async function fetchAllMessagesByID(
  id: string,
  session: Session
): Promise<AllBackendMessages> {
  try {
    const completionEndpoint = "/sessions/" + id + "/history";
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

    let backendResponse =
      (await response.json()) as BackendReponsePastQueriesAndAnswers;
    let messages = parseBackendResponse(backendResponse);

    return {
      messages,
      error: null,
    };
  } catch (currentError) {
    console.error("Error:", (currentError as Error).message);
    return {
      messages: [],
      error: (currentError as Error).message,
    };
  }
}

type RawBackendResponseCompletion = {
  request_id: string;
  session_id: string;
};

type BackendResponseCompletion = {
  requestID?: string;
  sessionID?: string;
  error: string | null;
};

async function sendMessageToBackend(
  id: string | null,
  message: string,
  session: Session
): Promise<BackendResponseCompletion> {
  try {
    const completionEndpoint = "/completion/v2";
    console.log({ id, message, session: session.access_token });
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + completionEndpoint,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.access_token}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(
          id
            ? {
                query: message,
                session_id: id,
              }
            : {
                query: message,
              }
        ),
      }
    );

    if (!response.ok) {
      throw new Error(
        `HTTP error! status: ${response.status} ${response.statusText} ${response.body}`
      );
    }

    let backendResponse =
      (await response.json()) as RawBackendResponseCompletion;
    return {
      requestID: backendResponse.request_id,
      sessionID: backendResponse.session_id,
      error: null,
    };
  } catch (currentError) {
    console.error("Error:", (currentError as Error).message);
    return {
      error: (currentError as Error).message,
    };
  }
}

type BackendStatusResponse = {
  type: Status;
  error?: string;
  answer?: string;
  current_action?: string;
};

async function checkStatusByRequestID(
  requestID: string,
  session: Session
): Promise<BackendStatusResponse> {
  try {
    const completionEndpoint = "/completion/v2/status/" + requestID;
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
      throw new Error(
        `HTTP error! status: ${response.status} ${response.statusText}`
      );
    }

    let backendResponse = (await response.json()) as BackendStatusResponse;
    return backendResponse;
  } catch (currentError) {
    console.error("Error:", currentError);
    return {
      type: Status.FAILED,
      error: (currentError as Error).message,
    };
  }
}
