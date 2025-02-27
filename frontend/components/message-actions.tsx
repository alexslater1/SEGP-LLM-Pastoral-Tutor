import { Message } from '@/types/message';
import { toast } from 'sonner';
import { useCopyToClipboard } from 'usehooks-ts';

import { CopyIcon, ThumbDownIcon, ThumbUpIcon } from './icons';
import { Button } from './ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from './ui/tooltip';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
  PopoverPortal,
  PopoverClose,
  PopoverArrow,
} from './ui/popover';
import { memo, useContext, useState } from 'react';
import { InfoIcon, XIcon } from 'lucide-react';
import { cx } from 'class-variance-authority';
import { Textarea } from './ui/textarea';
import { UserContext } from '@/lib/userContext';
import { cn } from '@/lib/utils';

export function PureMessageActions({
  chatId,
  message,
  isLoading,
  isCollapsibleOpen,
  setIsCollapsibleOpen,
  downvoteMessage,
  removeDownvoteMessage,
  messageDownvoted,
}: {
  chatId: string | null;
  message: Message;
  isLoading: boolean;
  isCollapsibleOpen: boolean;
  setIsCollapsibleOpen: (open: boolean) => void;
  downvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
  messageDownvoted: boolean;
}) {
  const [_, copyToClipboard] = useCopyToClipboard();
  const [isDownvotePopoverOpen, setIsDownvotePopoverOpen] = useState(false);
  const [isDownvoteHover, setIsDownvoteHover] = useState(false);
  const [input, setInput] = useState('');

  const user = useContext(UserContext);

  if (isLoading) return null;
  if (message.role === 'user') return null;
  if (!chatId) return null;


  const changePopoverState = (state: boolean) => {
    if (state && !messageDownvoted) {
      setIsDownvotePopoverOpen(state);
    } else if (!state) {
      setIsDownvotePopoverOpen(state);
      if (user) {
        toast.promise(downvoteMessage(message.requestID, input), {
          loading: 'Downvoting Response...',
          success: () => {
            return 'Downvoted Response!';
          },
          error: (error) => {
            if (error instanceof Error && error.message === 'Message already downvoted') {
              return 'Response already downvoted';
            } else {
              return 'Failed to downvote response';
            }
          },
        });
      } else {
        toast.error('Please login to downvote responses');
      }
    }

    setInput('');
  }

  return (
    <TooltipProvider delayDuration={0}>
      <div className="flex flex-row gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
          <Button
            className="py-1 px-2 h-fit text-muted-foreground hover:text-primary"
            variant="outline"
            onClick={async () => {
              await copyToClipboard(message.content as string);
              toast.success('Copied to clipboard!', {
                className: 'bg-success text-success-foreground'
              });
            }}
          >
            <CopyIcon />
          </Button>
          </TooltipTrigger>
          <TooltipContent>Copy</TooltipContent>
        </Tooltip>

        <Tooltip>
          <TooltipTrigger asChild>
            <Button 
              className="py-1 px-2 h-fit text-muted-foreground hover:text-primary"
              variant="outline"
              onClick={() => setIsCollapsibleOpen(!isCollapsibleOpen)}
            >
              <InfoIcon />
            </Button>
          </TooltipTrigger>
          <TooltipContent>Message Actions</TooltipContent>
        </Tooltip>

        <Popover open={isDownvotePopoverOpen} onOpenChange={changePopoverState}>
          {/* Separate state for hover to avoid tooltip opening when pressing enter on popover input */}
          <Tooltip open={isDownvoteHover} onOpenChange={() => {}}>
            <TooltipTrigger asChild>
              <PopoverTrigger asChild>
                <Button
                  className={cn(
                    "py-1 px-2 h-fit text-muted-foreground !pointer-events-auto",
                    messageDownvoted && "bg-zinc-300"
                  )}
                  variant="outline"
                  onClick={() => {
                    if (messageDownvoted) {
                      toast.promise(removeDownvoteMessage(message.requestID), {
                        loading: 'Removing downvote...',
                        success: () => {
                          return 'Downvote removed!';
                        },
                        error: (error) => {
                          return 'Failed to remove downvote';
                        },
                      });
                    } else {
                      changePopoverState(!isDownvotePopoverOpen);
                    }
                  }}
                  onMouseEnter={() => setIsDownvoteHover(true)}
                  onMouseLeave={() => setIsDownvoteHover(false)}
                >
                  <ThumbDownIcon />
                </Button>
              </PopoverTrigger>
            </TooltipTrigger>
            <TooltipContent>{messageDownvoted ? 'Remove Downvote' : 'Downvote Response'}</TooltipContent>
            <PopoverPortal>
              <PopoverContent className="rounded-2xl">
                <div className="flex gap-2">
                  <p className="text-base p-1">Downvote Reason</p>
                  <PopoverClose className="ml-auto">
                    <XIcon />
                  </PopoverClose>
                </div>
                <div className="flex flex-col gap-2 px-1 py-3">
                  <Textarea
                    placeholder="(Optional) Reason for downvote..."
                    value={input}
                    onChange={(event) => {
                      setInput(event.target.value);
                    }}
                    className={cx(
                      'min-h-[24px] max-h-[calc(75dvh)] overflow-hidden resize-none rounded-2xl !text-sm pb-10',
                      'bg-muted dark:bg-background',
                    )}
                    rows={2}
                    autoFocus
                    onKeyDown={(event) => {
                      if (event.key === 'Enter' && !event.shiftKey) {
                        event.preventDefault();
                        changePopoverState(false);
                      }
                    }}
                  />
                </div>
                <PopoverArrow />
              </PopoverContent>
            </PopoverPortal>
          </Tooltip>
        </Popover>
      </div>
    </TooltipProvider>
  );
}

export const MessageActions = memo(
  PureMessageActions,
  (prevProps, nextProps) => {
    if (prevProps.isLoading !== nextProps.isLoading) return false;
    if (prevProps.isCollapsibleOpen !== nextProps.isCollapsibleOpen) return false;
    if (prevProps.messageDownvoted !== nextProps.messageDownvoted) return false;

    return true;
  },
);

const DownVoteButton = ({
  messageDownvoted,
  removeDownvoteMessage,
  message,
  isDownvotePopoverOpen,
  setIsDownvoteHover,
  changePopoverState,
}: {
  messageDownvoted: boolean;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
  message: Message;
  isDownvotePopoverOpen: boolean;
  setIsDownvoteHover: (hover: boolean) => void;
  changePopoverState: (open: boolean) => void;
}) => {
  return (
    <Button
      className={cn(
        "py-1 px-2 h-fit text-muted-foreground !pointer-events-auto",
        messageDownvoted && "bg-muted"
      )}
      variant="outline"
        onClick={() => {
        if (messageDownvoted) {
          removeDownvoteMessage(message.requestID);
        } else {
          changePopoverState(!isDownvotePopoverOpen);
        }
      }}
      onMouseEnter={() => setIsDownvoteHover(true)}
      onMouseLeave={() => setIsDownvoteHover(false)}
    >
      <ThumbDownIcon />
    </Button>
  )
}
