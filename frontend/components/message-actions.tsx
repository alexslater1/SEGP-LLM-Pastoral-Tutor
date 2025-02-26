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
import equal from 'fast-deep-equal';
import { InfoIcon } from 'lucide-react';
import { downvote } from '@/lib/supabase/vote';

export function PureMessageActions({
  chatId,
  message,
  isLoading,
}: {
  chatId: string | null;
  message: Message;
  isLoading: boolean;
}) {
  const { mutate } = useSWRConfig();
  const [_, copyToClipboard] = useCopyToClipboard();

  const [isDialogOpen, setIsDialogOpen] = useState(false);

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

              
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button 
                className="py-1 px-2 h-fit text-muted-foreground hover:text-primary"
                variant="outline"
                onClick={() => setIsDialogOpen(true)}
              >
                <InfoIcon />
              </Button>
            </TooltipTrigger>
            <TooltipContent>Message Actions</TooltipContent>
          </Tooltip>
          <DialogPortal>
            <DialogOverlay />
            <DialogContent className="sm:max-w-[425px]">
              <DialogTitle className="text-primary text-xl font-semibold">Message Actions</DialogTitle>
              <DialogDescription className="text-foreground">
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
              </DialogDescription>
              <DialogClose asChild>
                <Button className="w-full mt-6 font-semibold bg-primary text-primary-foreground hover:bg-primary/90">
                  Close
                </Button>
              </DialogClose>
            </DialogContent>
          </DialogPortal>
        </Dialog>

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

    return true;
  },
);
