import { Info } from "lucide-react";
import { TooltipTrigger, TooltipContent, Tooltip, TooltipProvider } from "./ui/tooltip";
import { Switch } from "./ui/switch";

export function SwitchCard({titleText, descriptionText, enabled, setEnabled, loading, error}: 
  {
    titleText: string, 
    descriptionText: string, 
    enabled: boolean, 
    setEnabled: (enabled: boolean) => void, 
    loading: boolean, 
    error: Error | null
  }) {
  return (
    <div className="rounded-md border px-4 py-2 h-10">
      <div className="flex flex-row items-left gap-2 text-sm font-medium">
        <div>
          {titleText}
        </div>
        <div className="flex flex-row items-center">
          <TooltipProvider delayDuration={200}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="size-4 text-muted-foreground hover:text-primary" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="max-w-xs font-normal text-foreground">
                  {descriptionText}
                </p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
        <div className="flex flex-row items-center gap-2">
          {error ? 
            <p className="text-error">{error.message}</p> :
            <Switch checked={enabled} onCheckedChange={setEnabled} disabled={loading} />
          }
        </div>
      </div>
    </div>
  )
}
