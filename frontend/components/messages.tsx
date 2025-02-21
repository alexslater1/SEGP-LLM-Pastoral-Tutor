import { Message, Status } from "@/types/message";
import { ErrorMessage, PreviewMessage, ThinkingMessage } from "./message";
import { useScrollToBottom } from "./use-scroll-to-bottom";
import { memo, useEffect } from "react";
import equal from "fast-deep-equal/es6/react";

interface MessagesProps {
  chatId: string | null;
  isLoading: boolean;
  messages: Array<Message>;
  isReadonly: boolean;
}

function PureMessages({
  chatId,
  isLoading,
  messages,
  isReadonly,
}: MessagesProps) {
  const [messagesContainerRef, messagesEndRef] =
    useScrollToBottom<HTMLDivElement>();

  useEffect(() => {
    if (chatId) {
      window.history.replaceState({}, "", `/chat/${chatId}`);
    }
  }, [chatId]);

  return (
    <div
      ref={messagesContainerRef}
      className="flex flex-col min-w-0 gap-6 flex-1 overflow-y-scroll pt-4"
    >
      {messages.map((message, index) =>
        message.status === Status.COMPLETED ? (
          <PreviewMessage
            key={message.id}
            chatId={chatId}
            message={message}
            isLoading={isLoading && messages.length - 1 === index}
            isReadonly={isReadonly}
          />
        ) : message.status === Status.PENDING &&
          message.actions.length !== 0 ? (
          <ThinkingMessage
            key={message.id}
            message={message.actions[message.actions.length - 1]}
          />
        ) : message.status === Status.PENDING &&
          message.actions.length === 0 ? (
          <ThinkingMessage key={message.id} />
        ) : (
          <ErrorMessage key={message.id} error={message.content} />
        )
      )}

      {isLoading &&
        messages.length > 0 &&
        messages[messages.length - 1].role === "user" && <ThinkingMessage />}

      <div
        ref={messagesEndRef}
        className="shrink-0 min-w-[24px] min-h-[24px]"
      />
    </div>
  );
}

export const Messages = PureMessages;
