import { Message } from '@/types/message';
import { toast } from 'sonner';
import { useSWRConfig } from 'swr';
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
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
  DialogPortal,
  DialogOverlay,
  DialogClose,
} from './ui/dialog';
import { memo, useState } from 'react';
import { InfoIcon } from 'lucide-react';
import { downvote } from '@/lib/supabase/vote';

export function PureMessageActions({
  chatId,
  message,
  isLoading,
  isCollapsibleOpen,
  setIsCollapsibleOpen,
}: {
  chatId: string | null;
  message: Message;
  isLoading: boolean;
  isCollapsibleOpen: boolean;
  setIsCollapsibleOpen: (open: boolean) => void;
}) {
  const { mutate } = useSWRConfig();
  const [_, copyToClipboard] = useCopyToClipboard();

  if (isLoading) return null;
  if (message.role === 'user') return null;
  if (!chatId) return null;

  const [vote, setVote] = useState<boolean>(false);

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

        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              className="py-1 px-2 h-fit text-muted-foreground !pointer-events-auto"
              variant="outline"
              disabled={vote}
              onClick={async () => {
                toast.promise(downvote(message.requestID), {
                  loading: 'Downvoting Response...',
                  success: () => {
                    return 'Downvoted Response!';
                  },
                  error: 'Failed to downvote response.',
                });
              }}
            >
              <ThumbDownIcon />
            </Button>
          </TooltipTrigger>
          <TooltipContent>Downvote Response</TooltipContent>
        </Tooltip>
      </div>
    </TooltipProvider>
  );
}

export const MessageActions = memo(
  PureMessageActions,
  (prevProps, nextProps) => {
    if (prevProps.isLoading !== nextProps.isLoading) return false;
    if (prevProps.isCollapsibleOpen !== nextProps.isCollapsibleOpen) return false;

    return true;
  },
);
