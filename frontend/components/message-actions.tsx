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
import { ArrowUpIcon, InfoIcon, XIcon } from 'lucide-react';
import { cx } from 'class-variance-authority';
import { Textarea } from './ui/textarea';
import { cn } from '@/lib/utils';
import { VoteType } from '@/lib/supabase/vote';
import { useUser } from '@/providers/user-provider';

export function PureMessageActions({
  chatId,
  message,
  isLoading,
  isCollapsibleOpen,
  setIsCollapsibleOpen,
  isReadonly,
  downvoteMessage,
  removeDownvoteMessage,
  upvoteMessage,
  removeUpvoteMessage,
  messageDownvoted,
  messageUpvoted,
}: {
  chatId: string | null;
  message: Message;
  isLoading: boolean;
  isCollapsibleOpen: boolean;
  isReadonly: boolean;
  setIsCollapsibleOpen: (open: boolean) => void;
  downvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeDownvoteMessage: (messageId: string) => Promise<void>;
  upvoteMessage: (messageId: string, reason?: string) => Promise<void>;
  removeUpvoteMessage: (messageId: string) => Promise<void>;
  messageDownvoted: boolean;
  messageUpvoted: boolean;
}) {
  const [_, copyToClipboard] = useCopyToClipboard();
  const [isDownvotePopoverOpen, setIsDownvotePopoverOpen] = useState(false);
  const [isDownvoteHover, setIsDownvoteHover] = useState(false);
  const [input, setInput] = useState('');
  const [isUpvotePopoverOpen, setIsUpvotePopoverOpen] = useState(false);
  const [isUpvoteHover, setIsUpvoteHover] = useState(false);

  const user = useUser();

  if (isLoading) return null;
  if (message.role === 'user') return null;
  if (!chatId) return null;

  const setPopoverState = (type: VoteType, state: boolean) => {
    if (state && !messageDownvoted && type === 'downvote') {
      setIsDownvotePopoverOpen(state);
    } else if (state && !messageUpvoted && type === 'upvote') {
      setIsUpvotePopoverOpen(state);
    } else if (!state) {
      type === 'downvote' ? setIsDownvotePopoverOpen(state) : setIsUpvotePopoverOpen(state);
    }
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
              toast.success('Copied to clipboard!');
            }}
          >
            <CopyIcon />
          </Button>
          </TooltipTrigger>
          <TooltipContent>Copy</TooltipContent>
        </Tooltip>

        {!isReadonly && (
          <>
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

            <VotePopover
              type="downvote"
              messageVoted={messageDownvoted}
              removeVoteMessage={removeDownvoteMessage}
              voteMessage={downvoteMessage}
              isVotePopoverOpen={isDownvotePopoverOpen}
              changePopoverState={setPopoverState}
              isVoteHover={isDownvoteHover}
              setIsVoteHover={setIsDownvoteHover}
              input={input}
              setInput={setInput}
              message={message}
            />

            <VotePopover
              type="upvote"
              messageVoted={messageUpvoted}
              removeVoteMessage={removeUpvoteMessage}
              voteMessage={upvoteMessage}
              isVotePopoverOpen={isUpvotePopoverOpen}
              changePopoverState={setPopoverState}
              isVoteHover={isUpvoteHover}
              setIsVoteHover={setIsUpvoteHover}
              input={input}
              setInput={setInput}
              message={message}
            />
          </>
        )}
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
    if (prevProps.messageUpvoted !== nextProps.messageUpvoted) return false;

    return true;
  },
);

const capitalize = (str: string) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
}

type VotePopoverProps = {
  type: VoteType,
  messageVoted: boolean,
  removeVoteMessage: (messageId: string) => Promise<void>,
  voteMessage: (messageId: string, reason?: string) => Promise<void>
  isVotePopoverOpen: boolean,
  changePopoverState: (type: VoteType, state: boolean) => void,
  isVoteHover: boolean,
  setIsVoteHover: (state: boolean) => void,
  input: string,
  setInput: (input: string) => void,
  message: Message
}

const VotePopover = ({
  type, 
  messageVoted, 
  removeVoteMessage, 
  voteMessage,
  isVotePopoverOpen,
  changePopoverState,
  isVoteHover,
  setIsVoteHover,
  input,
  setInput,
  message
}: VotePopoverProps) => {
  const submitVote = () => {
    changePopoverState(type, false);
    toast.promise(voteMessage(message.requestID, input), {
      loading: type === 'downvote' ? 'Downvoting Response...' : 'Upvoting Response...',
      success: () => {
        return type === 'downvote' ? 'Downvoted Response!' : 'Upvoted Response!';
      },
      error: (error) => {
        if (error instanceof Error && error.message === 'Message already ' + (type === 'downvote' ? 'downvoted' : 'upvoted')) {
          return 'Response already ' + (type === 'downvote' ? 'downvoted' : 'upvoted');
        } else {
          return 'Failed to ' + type + ' response';
        }
      },
    });
    setInput('');
  }

  return (
    <Popover open={isVotePopoverOpen} onOpenChange={(state) => changePopoverState(type, state)}>
      {/* Separate state for hover to avoid tooltip opening when pressing enter on popover input */} 
      <Tooltip open={isVoteHover} onOpenChange={() => {}}> 
        <TooltipTrigger asChild> 
          <PopoverTrigger asChild> 
            <Button 
              className={cn( "py-1 px-2 h-fit text-muted-foreground !pointer-events-auto", 
                             messageVoted && "bg-zinc-300")} 
              variant="outline" 
              onClick={() => { 
                if (messageVoted) { 
                  toast.promise(removeVoteMessage(message.requestID), 
                  { loading: 'Removing ' + type + '...', 
                    success: () => { return capitalize(type) + ' removed!'; }, 
                    error: (error) => { return 'Failed to remove ' + capitalize(type); }, 
                  }); 
                } else { 
                  changePopoverState(type, !isVotePopoverOpen); 
                } 
              }}
              onMouseEnter={() => setIsVoteHover(true)}
              onMouseLeave={() => setIsVoteHover(false)}
            >
              {type === 'downvote' ? <ThumbDownIcon /> : <ThumbUpIcon />}
            </Button>
          </PopoverTrigger>
        </TooltipTrigger>
        <TooltipContent>{messageVoted ? 'Remove ' + capitalize(type) : capitalize(type) + ' Response'}</TooltipContent>
        <PopoverPortal>
          <PopoverContent className="rounded-2xl">
            <div className="flex gap-2">
              <p className="text-base p-1">{capitalize(type) + ' Reason'}</p>
              <PopoverClose className="ml-auto">
                <XIcon />
              </PopoverClose>
            </div>
            <div className="flex flex-col gap-2 py-2 w-80">
              <Textarea
                placeholder={'(Optional) Reason for ' + type + '...'}
                value={input}
                onChange={(event) => {
                  setInput(event.target.value);
                }}
                className={cx(
                  'min-h-[24px] max-h-[calc(75dvh)] overflow-hidden resize-none rounded-2xl !text-sm pb-10',
                  'bg-muted dark:bg-background',
                )}
                rows={3}
                autoFocus
                onKeyDown={(event) => {
                  if (event.key === 'Enter' && !event.shiftKey) {
                    event.preventDefault();
                    submitVote();
                  }
                }}
              />

              <div className="absolute bottom-0 right-0 px-5 py-6 w-fit flex flex-row justify-end">
                <VoteSubmitButton submitVote={submitVote} input={input} />
              </div>
            </div>
            <PopoverArrow />
          </PopoverContent>
        </PopoverPortal>
      </Tooltip>
    </Popover>
  )
}

function VoteSubmitButton({submitVote, input}: {submitVote: () => void, input: string}) {
  return (
    <Button
      className="rounded-full p-1.5 h-fit border dark:border-zinc-600"
      onClick={(event) => {
        event.preventDefault();
        submitVote();
      }}
      disabled={input.length === 0}
    >
      <ArrowUpIcon size={14} />
    </Button>
  );
}