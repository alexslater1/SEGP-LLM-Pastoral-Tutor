import { Message, Status, Role } from "@/types/message";
import { generateUUID } from "@/lib/utils";
import { useState, useRef, useEffect } from "react";

const STATUS_QUERY_INTERVAL_SECONDS = 1;
const BACKEND_AGENT_URL="https://segp-backend-agent.serve.freemyip.com"
const IGNORED_ACTIONS = ["Thinking", "Thinking..."]

export type ChatItem = {
  messages: Message[];
  setMessages: (messages: Message[] | ((messages: Message[]) => Message[])) => void;
  handleSubmit: () => void;
  input: string, 
  setInput: (input: string) => void;
  append: (query: string) => void;
  isLoading: boolean;
  stop: () => void;
  reload: () => void;
  error: string | null;
  id: string | null;
}

export type ChatItemProps = {
  id: string | null;
}

type BackendResponsePastQueryAndAnswer = {
  type: Status,
  answer?: string,
  current_action?: string,
  error?: string,
  query: string,
  requestID: string,
  actions: string[],
}

type BackendReponsePastQueriesAndAnswers = BackendResponsePastQueryAndAnswer[];

type BackendResponseCompletion = {
  request_id: string,
  session_id: string,
}

type BackendStatusResponse = {
    type: Status,
    error?: string,
    answer?: string,
    current_action?: string,
}

type ChatState = {
  messages: Message[];
  isAwaitingResponse: boolean;
  checkStatus: boolean;
  attemptedIntitialMessageLoad: boolean;
  input: string;
  statusCheckedCount: number;
}

export function useChat({ id }: ChatItemProps): ChatItem {
  const [chatState, setChatState] = useState<ChatState>({
    messages: [],
    isAwaitingResponse: false,
    checkStatus: false,
    attemptedIntitialMessageLoad: false,
    input: '',
    statusCheckedCount: 0 // Incase 2 or more consecutive status checks are made which don't 
                          // update the state. We still want the next status check to be made
  });

  const error = useRef<string | null>(null);
  const session_id = useRef<string | null>(null);

  let newChatState = chatState;

  useEffect(() => {
    if (chatState.messages.length === 0 && 
        !chatState.isAwaitingResponse && 
        !chatState.attemptedIntitialMessageLoad) {
      setIsAwaitingResponse(true);
      updateChatState();
    } else if (chatState.messages.length === 0 && 
               chatState.isAwaitingResponse && 
               !chatState.attemptedIntitialMessageLoad) {
      loadAllMessages();
    } else if (chatState.attemptedIntitialMessageLoad && 
               chatState.isAwaitingResponse) {
      sendMessageAndAddResponse();
    } else if (chatState.messages.length > 0 && 
               !chatState.isAwaitingResponse && 
               chatState.checkStatus) {
      checkStatus();
    }
  }, [chatState]);


  const setMessages = (messages: Message[] | ((messages: Message[]) => Message[])) => {
    newChatState = { 
      ...newChatState, 
      messages: typeof messages === 'function' ? 
        messages(newChatState.messages) : messages
    };
  }

  const setIsAwaitingResponse = (isAwaitingResponse: boolean) => {
    newChatState = { 
      ...newChatState, 
      isAwaitingResponse: isAwaitingResponse 
    };
  }

  const setCheckStatus = (checkStatus: boolean) => {
    newChatState = { 
      ...newChatState, 
      checkStatus: checkStatus 
    };
  }

  const setAttemptedIntitialMessageLoad = (attemptedIntitialMessageLoad: boolean) => {
    newChatState = { 
      ...newChatState, 
      attemptedIntitialMessageLoad: attemptedIntitialMessageLoad 
    };
  }

  const setInput = (input: string) => {
    newChatState = { 
      ...newChatState, 
      input: input 
    };
  }

  const incrementStatusCheckedCount = () => {
    newChatState = { 
      ...newChatState, 
      statusCheckedCount: newChatState.statusCheckedCount + 1 
    };
  }

  const resetStatusCheckedCount = () => {
    newChatState = { 
      ...newChatState, 
      statusCheckedCount: 0 
    };
  }

  const updateChatState = () => {
    setChatState(newChatState);
  }

  const append = (query: string) => {
    setInput(query);
    updateChatState();
    handleSubmit();
  }

  const handleSubmit = () => {
    setIsAwaitingResponse(true);
    setCheckStatus(false);
    setMessages((messages) => {
      let id = generateUUID();
      return [...messages, {
        id: id + Role.USER,
        requestID: id,
        content: chatState.input,
        role: Role.USER,
        actions: [],
        status: Status.COMPLETED
      }];
    });
    updateChatState();
  }

  const stop = () => {
    setIsAwaitingResponse(false);
    setCheckStatus(false);
    updateChatState();
  }

  const reload = async () => {
    setIsAwaitingResponse(true);
    setMessages([]);
    setCheckStatus(false);
    updateChatState();
  }

  // Pre-condition: chatState.messages.length === 0 && chatState.isLoading === true
  const loadAllMessages = async () => {
    if (session_id.current) {
      const { messages, checkStatus, error: loadError } = await fetchAllMessagesByID(session_id.current);
      if (loadError) {
        error.current = loadError;
      } else {
        setMessages(messages);
        setCheckStatus(checkStatus);
      }
    } else {
      error.current = "Could not find chat ID";
    }
    
    setIsAwaitingResponse(false);
    setAttemptedIntitialMessageLoad(true);
    updateChatState();
  }

  const sendMessageAndAddResponse = async (query?: string) => {
    const { requestID, sessionID, checkStatus, error: sendError } = 
      await sendMessageToBackend(session_id.current, query ? query : chatState.input);

    if (sendError) {
      let id = generateUUID();
      let newMessages = [...chatState.messages, {
        id: id + Role.ASSISTANT,
        requestID: id,
        content: sendError,
        role: Role.ASSISTANT,
        actions: [],
        status: Status.FAILED,
      }]
      setMessages(newMessages);
      setCheckStatus(checkStatus);
    } else {
      session_id.current = sessionID;
      let newMessages = [...chatState.messages, {
        id: requestID as string + Role.ASSISTANT,
        requestID: requestID as string,
        content: "Thinking..",
        role: Role.ASSISTANT,
        actions: [],
        status: Status.PENDING,
      }]
      setMessages(newMessages);
      setCheckStatus(checkStatus);
    }
    
    setIsAwaitingResponse(false);
    setInput('');
    updateChatState();
  }

  const checkStatus = async () => {
    await new Promise(resolve => setTimeout(resolve, 1000 * STATUS_QUERY_INTERVAL_SECONDS));
    if (chatState.messages[chatState.messages.length - 1].role === Role.USER) {
      setCheckStatus(false);
      setIsAwaitingResponse(false);
      updateChatState();
      return;
    }
    const statusResponse = 
      await checkStatusByRequestID(chatState.messages[chatState.messages.length - 1].requestID);

    let lastMessage = structuredClone(chatState.messages[chatState.messages.length - 1]);
    if (statusResponse.type === Status.COMPLETED) {
      lastMessage.status = Status.COMPLETED;
      lastMessage.content = statusResponse.answer as string;
      setCheckStatus(false);
      resetStatusCheckedCount();
    } else if (
        statusResponse.type === Status.PENDING && 
        statusResponse.current_action !== lastMessage.actions[lastMessage.actions.length - 1] &&
        !IGNORED_ACTIONS.includes(statusResponse.current_action as string)
    ) {
      lastMessage.status = Status.PENDING;
      lastMessage.actions.push(statusResponse.current_action as string)
      incrementStatusCheckedCount();
    } else if (statusResponse.type === Status.FAILED) {
      lastMessage.status = Status.FAILED;
      lastMessage.content = statusResponse.error as string;
      setCheckStatus(false);
      resetStatusCheckedCount();
    }

    setMessages((messages) => {
      let newMessages = [...messages];
      newMessages[newMessages.length - 1] = lastMessage;
      return newMessages;
    });
    setIsAwaitingResponse(false);
    updateChatState();
  }

  return { 
    messages: chatState.messages, 
    setMessages,
    handleSubmit,
    input: chatState.input,
    setInput: (newInput: string) => { 
      setInput(newInput);
      updateChatState();
     },
    append,
    isLoading: chatState.isAwaitingResponse,
    stop,
    reload,
    error: error.current,
    id: session_id.current,
  };
}

function parseBackendResponse(response: BackendReponsePastQueriesAndAnswers): Message[] {
  let messages: Message[] = [];
  response.forEach((item) => {
    messages.push({
      id: item.requestID + Role.USER,
      requestID: item.requestID,
      content: item.query,
      role: Role.USER,
      actions: [],
      status: Status.COMPLETED,
    });
    messages.push({
      id: item.requestID + Role.ASSISTANT,
      requestID: item.requestID,
      content: item.type === Status.COMPLETED ? item.answer as string : 
               item.type === Status.PENDING ? "Thinking..." :
               item.error as string,
      role: Role.ASSISTANT,
      actions: item.actions.concat(item.type === Status.PENDING ? [item.current_action as string] : []),
      status: item.type,
    });
  });
  return messages;
}

async function fetchAllMessagesByID (id: string) {
  try {
    const completionEndpoint = "/sessions/" + id
    const response = await fetch(BACKEND_AGENT_URL + completionEndpoint, {
        method: 'GET',
    });

    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }

    let backendResponse = await response.json() as BackendReponsePastQueriesAndAnswers;
    let messages = parseBackendResponse(backendResponse);

    return {
      messages,
      checkStatus: messages.length > 0 && messages[messages.length - 1].status === Status.PENDING,
      error: null,
    }
  } catch (currentError) {
    console.error('Error:', currentError);
    return {
      messages: [],
      checkStatus: false,
      error: (currentError as Error).message,
    }
  } 
}

async function sendMessageToBackend (id: string | null, message: string) {
  try {
    const completionEndpoint = "/completion/v2"
    const response = await fetch(BACKEND_AGENT_URL + completionEndpoint, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(id ?{
          query: message,
          session_id: id,
        } : {
          query: message,
        }),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    let backendResponse = await response.json() as BackendResponseCompletion;
    return {
      requestID: backendResponse.request_id,
      sessionID: backendResponse.session_id,
      checkStatus: true,
      error: null,
    }
  } catch (currentError) {
    console.error('Error:', currentError);
    return {
      requestID: null,
      sessionID: null,
      checkStatus: false,
      error: (currentError as Error).message,
    }
  }
}

async function checkStatusByRequestID(requestID: string) {
  try {
    const completionEndpoint = "/completion/v2/status/" + requestID
    const response = await fetch(BACKEND_AGENT_URL + completionEndpoint, {
        method: 'GET',
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    let backendResponse = await response.json() as BackendStatusResponse;
    return backendResponse.type === Status.COMPLETED ? {
      type: backendResponse.type,
      answer: backendResponse.answer,
    } as BackendStatusResponse : backendResponse.type === Status.PENDING ? {
      type: backendResponse.type,
      current_action: backendResponse.current_action,
    } as BackendStatusResponse : {
      type: backendResponse.type,
      error: backendResponse.error,
    } as BackendStatusResponse;
  } catch (currentError) {
    console.error('Error:', currentError);
    return {
      type: Status.FAILED,
      error: (currentError as Error).message,
    }
  }
}