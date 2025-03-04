"use client";

import { useChat } from "@/hooks/use-chat";
import { motion } from "framer-motion";
import { ChatHeader } from "@/components/chat-header";
import { MultimodalInput } from "./multimodal-input";
import { Messages } from "./messages";
import { useUser } from "@/providers/user-provider";
import { useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { useSessionHistory } from "@/providers/session-history-provider";
import { useChatUrl } from "@/providers/chat-url-provider";

type MessageType = "initial" | "loading" | "error";

export function Chat({
  id,
  isReadonly,
}: {
  id: string | null;
  isReadonly: boolean;
}) {
  const userContext = useUser();
  const query = useSearchParams().get("query");
  const capitalize = (name: string) => name.charAt(0).toUpperCase() + name.slice(1);
  const firstAgentMessage = "Hi Anshul, I am the Imperial College tutor agent! How are you doing? I see you have a big Graphics coursework due soon, how is that going?";

  const [_, setChatUrl] = useChatUrl();
  const { refresh: refreshHistory } = useSessionHistory();

  const {
    messages,
    handleSubmit,
    isLoading,
    stop,
    error,
    id: current_id,
    isInitialLoad,
  } = useChat({
    id,
    firstAgentMessage: !isReadonly ? firstAgentMessage : undefined
  });

  useEffect(() => {
    if (current_id && !isReadonly) {
      const newUrl = `/chat/${current_id}`;
      window.history.replaceState({ ...window.history.state, as: newUrl, url: newUrl }, "", newUrl);

      setChatUrl(newUrl);
      refreshHistory();
    }
    // setChatUrl and refreshHistory will not change, so we can disable the linting
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [current_id]);

  return isInitialLoad ? (
    <FullScreenMessage messageType="loading" />
  ) : error ? (
    <FullScreenMessage messageType="error" error={error} />
  ) : (
    <div className="flex flex-col h-full overflow-hidden bg-background">
      <ChatHeader
        isReadonly={isReadonly}
      />

      {messages.length == 0 ? (
        <FullScreenMessage messageType="initial" />
      ) : (
        <Messages
          chatId={current_id}
          isLoading={isLoading}
          messages={messages}
          isReadonly={isReadonly}
        />
      )}

      {!isReadonly && (
      <form className="flex mx-auto px-4 bg-background pb-4 md:pb-6 gap-2 w-full md:max-w-3xl">
          <MultimodalInput
            messages={messages}
            handleSubmit={handleSubmit}
            isLoading={isLoading}
            stop={stop}
            user={userContext.user}
            query={query}
          />
      </form>
    )}
    </div>
  );
}

function FadeInWrapper({ children }: { children: React.ReactNode }) {
  return (
    <motion.div
      className="h-full flex items-center justify-center"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      {children}
    </motion.div>
  );
}

function FullScreenMessage({
  messageType,
  error,
}: {
  messageType: MessageType;
  error?: string;
}) {
  switch (messageType) {
    case "initial":
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center p-16">
            <div className="text-center text-foreground">
              <p className="text-5xl font-bold">
                Hi, I&apos;m Amanda, the Imperial College tutor agent!
              </p>
              <p className="text-4xl font-bold p-8">Ask me anything</p>
            </div>
          </div>
        </FadeInWrapper>
      );
    case "loading":
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center h-dvh p-16">
            <div className="text-center text-primary">
              <p className="text-5xl font-bold">Loading...</p>
            </div>
          </div>
        </FadeInWrapper>
      );
    case "error":
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center h-dvh p-16">
            <div className="text-center text-red-500">
              <p className="text-5xl font-bold">Error loading chat</p>
              <p className="text-4xl font-bold pt-8">Please try again</p>
              <p className="text-1xl p-3 text-muted-foreground">{error}</p>
            </div>
          </div>
        </FadeInWrapper>
      );
    default:
      return null;
  }
}
