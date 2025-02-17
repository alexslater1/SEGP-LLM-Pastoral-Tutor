import { Message, Status } from '@/types/message';
import { ErrorMessage, PreviewMessage, ThinkingMessage } from './message';
import { useScrollToBottom } from './use-scroll-to-bottom';
import { memo, useEffect } from 'react';
import { Vote } from '@/lib/db/schema';
import equal from 'fast-deep-equal/es6/react';

interface MessagesProps {
  chatId: string | null;
  isLoading: boolean;
  votes: Array<Vote> | undefined;
  messages: Array<Message>;
  setMessages: (
    messages: Message[] | ((messages: Message[]) => Message[]),
  ) => void;
  reload: () => void;
  isReadonly: boolean;
}

function PureMessages({
  chatId,
  isLoading,
  votes,
  messages,
  setMessages,
  reload,
  isReadonly,
}: MessagesProps) {
  const [messagesContainerRef, messagesEndRef] =
    useScrollToBottom<HTMLDivElement>();

  useEffect(() => {
    if (chatId) {
      window.history.replaceState({}, '', `/chat/${chatId}`);
    }
  }, [chatId]);

  return (
    <div
      ref={messagesContainerRef}
      className="flex flex-col min-w-0 gap-6 flex-1 overflow-y-scroll pt-4"
    >
      {messages.map((message, index) => (
        message.status === Status.COMPLETED ?
          <PreviewMessage
            // This key is used to force the component to re-render (with the nice animation) when the message is updated
            key={message.id + message.actions.length + message.content.length}
            chatId={chatId}
            message={message}
            isLoading={isLoading && messages.length - 1 === index}
            vote={
              votes
                ? votes.find((vote) => vote.messageId === message.id)
                : undefined
            }
            setMessages={setMessages}
            reload={reload}
            isReadonly={isReadonly}
          /> :
        message.status === Status.PENDING && message.actions.length !== 0 ?
          <ThinkingMessage 
            key={message.id + message.actions.length + message.content.length} 
            message={message.actions[message.actions.length - 1]} 
          /> :
        message.status === Status.PENDING && message.actions.length === 0 ?
          <ThinkingMessage 
            key={message.id + message.actions.length + message.content.length} 
          /> :
          <ErrorMessage error={message.content} />
      ))}

      {isLoading &&
        messages.length > 0 &&
        messages[messages.length - 1].role === 'user' && <ThinkingMessage />}

      <div
        ref={messagesEndRef}
        className="shrink-0 min-w-[24px] min-h-[24px]"
      />
    </div>
  );
}

export const Messages = memo(PureMessages, (prevProps, nextProps) => {
  if (prevProps.isLoading !== nextProps.isLoading) return false;
  if (prevProps.isLoading && nextProps.isLoading) return false;
  if (prevProps.messages.length !== nextProps.messages.length) return false;
  if (!equal(prevProps.messages, nextProps.messages)) return false;
  if (!equal(prevProps.votes, nextProps.votes)) return false;

  return true;
});
