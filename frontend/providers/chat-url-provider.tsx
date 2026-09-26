'use client';

import { useContext, useState } from "react";
import { createContext } from "react";


const ChatUrlContext = createContext<[string | null, (chatUrl: string | null) => void]>([null, () => {}]);

export function ChatUrlProvider({children}: {children: React.ReactNode}) {
  const [chatUrl, setChatUrl] = useState<string | null>(null);

  return (
    <ChatUrlContext.Provider value={[chatUrl, setChatUrl]}>
      {children}
    </ChatUrlContext.Provider>
  );
}

export function useChatUrl (): [string | null, (chatUrl: string | null) => void] {
  return useContext(ChatUrlContext);
};
