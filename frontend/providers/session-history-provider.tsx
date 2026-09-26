'use client';

import { useContext } from "react";
import { createContext } from "react";
import { ChatHistoryItem, useChatSessionHistory } from "@/hooks/use-chat-history";


const SessionHistoryContext = createContext<ChatHistoryItem>({} as ChatHistoryItem);

export function SessionHistoryProvider({children}: {children: React.ReactNode}) {
  const sessionHistory = useChatSessionHistory();

  return (
    <SessionHistoryContext.Provider value={sessionHistory}>
      {children}
    </SessionHistoryContext.Provider>
  );
}

export function useSessionHistory (): ChatHistoryItem {
  return useContext(SessionHistoryContext);
};
