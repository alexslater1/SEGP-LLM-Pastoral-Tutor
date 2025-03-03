"use client";

import { useChat } from "@/hooks/use-chat";
import { motion } from "framer-motion";

import { ChatHeader } from "@/components/chat-header";

import { MultimodalInput } from "./multimodal-input";
import { Messages } from "./messages";
import { VisibilityType } from "./visibility-selector";
import { UserContext } from "@/lib/userContext";
import { useContext, useEffect } from "react";
import { useSearchParams } from "next/navigation";

type MessageType = "initial" | "loading" | "error";

export function Chat({
  id,
  selectedVisibilityType,
  isReadonly,
}: {
  id: string | null;
  selectedVisibilityType: VisibilityType;
  isReadonly: boolean;
}) {
  const user = useContext(UserContext);
  const query = useSearchParams().get("query");
  const capitalize = (name: string) => name ? name.charAt(0).toUpperCase() + name.slice(1) : "";
  const firstAgentMessage = `Hi${user ? (" " + capitalize(user.name)) : ""}, I'm the Imperial College tutor agent. Ask me anything!`;

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
    firstAgentMessage
  });

  useEffect(() => {
    if (current_id) {
      const newUrl = `/chat/${current_id}`;
      window.history.replaceState({ ...window.history.state, as: newUrl, url: newUrl }, "", newUrl);
    }
  }, [current_id]);

  return isInitialLoad ? (
    <FullScreenMessage messageType="loading" />
  ) : error ? (
    <FullScreenMessage messageType="error" error={error} />
  ) : (
    <div className="flex flex-col min-w-0 h-dvh bg-background">
      <ChatHeader
        chatId={current_id}
        selectedVisibilityType={selectedVisibilityType}
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

      <form className="flex mx-auto px-4 bg-background pb-4 md:pb-6 gap-2 w-full md:max-w-3xl">
        {!isReadonly && (
          <MultimodalInput
            messages={messages}
            handleSubmit={handleSubmit}
            isLoading={isLoading}
            stop={stop}
            user={user}
            query={query}
          />
        )}
      </form>
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
