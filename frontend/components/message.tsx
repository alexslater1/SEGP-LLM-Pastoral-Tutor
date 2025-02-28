"use client";

import { Message, Status } from "@/types/message";
import cx from "classnames";
import { AnimatePresence, motion } from "framer-motion";
import { memo, useState } from "react";

import { CrossIcon, SparklesIcon } from "./icons";
import { Markdown } from "./markdown";
import { MessageActions } from "./message-actions";
import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "./ui/collapsible";
import { ArrowUpIcon } from "lucide-react";

const PurePreviewMessage = ({
  chatId,
  message,
  isLoading,
  isReadonly,
  downvoteMessage,
  removeDownvoteMessage,
  messageDownvoted,
}: {
  chatId: string | null;
  message: Message;
  isLoading: boolean;
  isReadonly: boolean;
  downvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
  messageDownvoted: boolean;
}) => {
  const [mode, setMode] = useState<"view" | "edit">("view");
  const [isCollapsibleOpen, setIsCollapsibleOpen] = useState(false);

  function getLastStatusOrMessage(): string {
    if (
      message.status === Status.COMPLETED ||
      message.status === Status.ERROR
    ) {
      return message.content;
    } else if (message.actions.length > 0) {
      return message.actions[message.actions.length - 1];
    } else {
      return "Thinking...";
    }
  }

  return (
    <AnimatePresence>
      <motion.div
        className="w-full mx-auto max-w-3xl px-4 group/message"
        initial={{ y: 5, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        data-role={message.role}
      >
        <Collapsible open={isCollapsibleOpen} onOpenChange={setIsCollapsibleOpen}>
          <div
            className={cn(
              "flex gap-4 w-full group-data-[role=user]/message:ml-auto group-data-[role=user]/message:max-w-2xl",
              {
                "w-full": mode === "edit",
                "group-data-[role=user]/message:w-fit": mode !== "edit",
              }
            )}
          >

            {message.role === "assistant" && (
              <div className="size-8 flex items-center rounded-full justify-center ring-1 shrink-0 ring-border bg-background">
                <div className="translate-y-px text-primary">
                  <SparklesIcon size={14} />
                </div>
              </div>
            )}

            <div>
              <div className="flex flex-col gap-2 w-full">
                <motion.div
                  className={cn("w-full mx-auto max-w-3xl group/message", {
                    "pl-4": message.role === "user",
                    "pr-4": message.role === "assistant",
                  })}
                  initial={{ y: 5, opacity: 0 }}
                  animate={{ y: 0, opacity: 1 }}
                >
                  {message.content && mode === "view" && (
                    <div className="flex flex-row gap-2 items-start ">
                      <div
                        className={cn("flex flex-col gap-4", {
                          "bg-chat-user text-chat-user-foreground px-3 py-2 rounded-xl":
                            message.role === "user",
                          "bg-chat-assistant text-chat-assistant-foreground px-3 py-2 rounded-xl":
                            message.role === "assistant",
                        })}
                      >
                        <Markdown>{getLastStatusOrMessage()}</Markdown>
                      </div>
                    </div>
                  )}
                </motion.div>

                <MessageActions
                  key={`action-${message.id}`}
                  chatId={chatId}
                  message={message}
                  isLoading={isLoading}
                  isReadonly={isReadonly}
                  isCollapsibleOpen={isCollapsibleOpen}
                  setIsCollapsibleOpen={setIsCollapsibleOpen}
                  downvoteMessage={downvoteMessage}
                  removeDownvoteMessage={removeDownvoteMessage}
                  messageDownvoted={messageDownvoted}
                />
              </div>
            </div>
          </div>
          <CollapsibleContent className="py-3 pr-4">
            <div className="bg-muted rounded-2xl p-4 ml-11">
              {message.actions.length > 0 ? (
                <div className="space-y-2">
                  {message.actions.map((action, index) => (
                    <div key={index} className="flex gap-2">
                      <span className="text-primary font-medium">{index + 1}.</span>
                      <span>{action}</span>
                    </div>
                  ))}
                </div>
              ) : (
                'No actions found'
              )}
            </div>
          </CollapsibleContent>
        </Collapsible>
      </motion.div>
    </AnimatePresence>
  );
};

export const PreviewMessage = PurePreviewMessage;

export const ThinkingMessage = ({
  message = "Thinking...",
}: {
  message?: string;
}) => {
  const role = "assistant";

  return (
    <AnimatePresence>
      <motion.div
        className="w-full mx-auto max-w-3xl px-4 group/message "
        initial={{ y: 5, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        data-role={role}
      >
        <div
          className={cx(
            "flex gap-4 group-data-[role=user]/message:px-3 w-full group-data-[role=user]/message:w-fit group-data-[role=user]/message:ml-auto group-data-[role=user]/message:max-w-2xl group-data-[role=user]/message:py-2 rounded-xl",
            {
              "group-data-[role=user]/message:bg-muted": true,
            }
          )}
        >
          <div className="size-8 flex items-center rounded-full justify-center ring-1 shrink-0 ring-border">
            <div className="text-primary">
              <SparklesIcon size={14} />
            </div>
          </div>

          <motion.div
            key={message.length}
            className="w-full mx-auto max-w-3xl px-4 group/message "
            initial={{ y: 5, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            data-role={role}
          >
            <div className="flex flex-col gap-2 w-full">
              <div className="flex flex-col gap-4 text-muted-foreground">
                {message}
              </div>
            </div>
          </motion.div>
        </div>
      </motion.div>
    </AnimatePresence>
  );
};

export const ErrorMessage = ({ error }: { error: string }) => {
  const role = "assistant";

  return (
    <motion.div
      className="w-full mx-auto max-w-3xl px-4 group/message "
      initial={{ y: 5, opacity: 0 }}
      animate={{ y: 0, opacity: 1, transition: { delay: 1 } }}
      data-role={role}
    >
      <div
        className={cx(
          "flex gap-4 group-data-[role=user]/message:px-3 w-full group-data-[role=user]/message:w-fit group-data-[role=user]/message:ml-auto group-data-[role=user]/message:max-w-2xl group-data-[role=user]/message:py-2 rounded-xl",
          {
            "group-data-[role=user]/message:bg-muted": true,
          }
        )}
      >
        <div className="size-8 flex items-center rounded-full justify-center ring-1 shrink-0 ring-border">
          <CrossIcon size={14} color="red" />
        </div>

        <div className="flex flex-col gap-2 w-full">
          <div className="flex flex-col gap-4 text-red-500">{error}</div>
        </div>
      </div>
    </motion.div>
  );
};
