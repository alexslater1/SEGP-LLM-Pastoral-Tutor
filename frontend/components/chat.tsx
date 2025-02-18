'use client';

import { useChat } from '@/hooks/use-chat';
import useSWR, { useSWRConfig } from 'swr';
import { motion } from 'framer-motion';

import { ChatHeader } from '@/components/chat-header';
import type { Vote } from '@/lib/db/schema';
import { fetcher } from '@/lib/utils';

import { MultimodalInput } from './multimodal-input';
import { Messages } from './messages';
import { VisibilityType } from './visibility-selector';

type MessageType = 'initial' | 'loading' | 'error';

export function Chat({
  id,
  selectedVisibilityType,
  isReadonly,
}: {
  id: string | null;
  selectedVisibilityType: VisibilityType;
  isReadonly: boolean;
}) {
  const {
    messages,
    setMessages,
    handleSubmit,
    input,
    setInput,
    append,
    isLoading,
    stop,
    reload,
    error,
    attemptedInitialMessageLoad,
    id: current_id,
  } = useChat({
    id,
  });

  const { data: votes } = useSWR<Array<Vote>>(
    `/api/vote?chatId=${current_id}`,
    fetcher,
  );

  return (
    <>
      {!attemptedInitialMessageLoad ? 
        <InitialMessage messageType="loading" />
      : error ?
        <InitialMessage messageType="error" />
      :(
        <div className="flex flex-col min-w-0 h-dvh bg-background">
          <ChatHeader
            chatId={current_id}
            selectedVisibilityType={selectedVisibilityType}
            isReadonly={isReadonly}
          />

          {messages.length == 0 && (
            <InitialMessage messageType="initial" />
          )}

          <Messages
            chatId={current_id}
            isLoading={isLoading}
            votes={votes}
            messages={messages}
            setMessages={setMessages}
            reload={reload}
            isReadonly={isReadonly}
          />

          <form className="flex mx-auto px-4 bg-background pb-4 md:pb-6 gap-2 w-full md:max-w-3xl">
            {!isReadonly && (
              <MultimodalInput
                chatId={current_id}
                input={input}
                setInput={setInput}
                handleSubmit={handleSubmit}
                isLoading={isLoading}
                stop={stop}
                messages={messages}
                append={append}
              />
            )}
          </form>
        </div>
      )}
    </>
  );
}

function FadeInWrapper({children}: {children: React.ReactNode}) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
    >
      {children}
    </motion.div>
  )
}

function InitialMessage({messageType}: {messageType: MessageType}) {
  switch (messageType) {
    case 'initial':
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center h-dvh p-16">
            <div className="text-center text-white">
              <p className="text-5xl font-bold">
                Hi, I’m the Imperial College tutor agent
              </p>
              <p className="text-4xl font-bold p-8">
                Ask me anything
              </p>
            </div>
          </div>
        </FadeInWrapper>
      )
    case 'loading':
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center h-dvh p-16">
            <div className="text-center text-white">
              <p className="text-5xl font-bold">
                Loading...
              </p>
            </div>
          </div>
        </FadeInWrapper>
      )
    case 'error':
      return (
        <FadeInWrapper>
          <div className="flex flex-col items-center justify-center h-dvh p-16">
            <div className="text-center text-red-500">
              <p className="text-5xl font-bold">
                Error loading chat
              </p>
              <p className="text-4xl font-bold p-8">
                Please try again
              </p>
            </div>
          </div>
        </FadeInWrapper>
      )
    default:
      return null;
  }
}
