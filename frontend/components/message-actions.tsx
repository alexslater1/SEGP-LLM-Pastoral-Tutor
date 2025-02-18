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

  return (
    <TooltipProvider delayDuration={0}>
      <div className="flex flex-row gap-2">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              className="py-1 px-2 h-fit text-muted-foreground"
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

              
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button 
                className="py-1 px-2 h-fit text-muted-foreground"
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
            <DialogContent>
              <DialogTitle>Message Actions</DialogTitle>
              <DialogDescription>
                {message.actions.length > 0 ?
                (message.actions.reduce((prev, action, index) => (
                  <>{prev}<br/>{index + 1}. {action}</>
                ), <></>)) :
                  'No actions found'
                }
              </DialogDescription>
              <DialogClose>Close</DialogClose>
            </DialogContent>
          </DialogPortal>
        </Dialog>
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
