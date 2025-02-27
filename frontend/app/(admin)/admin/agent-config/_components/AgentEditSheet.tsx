import React, { useState, useEffect } from "react";
import { Info } from "lucide-react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { AgentsResponse } from "@/hooks/use-agent-config";

// Import the Agent interface or define it here

interface AgentEditSheetProps {
  isOpen: boolean;
  setIsOpen: (open: boolean) => void;
  agent: AgentsResponse | null;
  onSave: (agent: AgentsResponse) => void;
  setIsDeleteDialogOpen: (open: boolean) => void;
  availableTools: string[];
  availableApis: {
    abc_apis: string[];
    emarking_apis: string[];
  };
  isSubmitting?: boolean;
  isDeleting?: boolean;
}

const AgentEditSheet: React.FC<AgentEditSheetProps> = ({
  isOpen,
  setIsOpen,
  agent,
  onSave,
  setIsDeleteDialogOpen,
  availableTools,
  availableApis,
  isSubmitting = false,
  isDeleting = false,
}) => {
  const [editedAgent, setEditedAgent] = useState<AgentsResponse | null>(null);

  // Update local state when the agent prop changes
  useEffect(() => {
    setEditedAgent(agent);
  }, [agent]);

  const handleInputChange = (field: keyof AgentsResponse, value: string) => {
    if (editedAgent) {
      setEditedAgent({
        ...editedAgent,
        [field]: value,
      });
    }
  };

  const handleToolToggle = (tool: string) => {
    if (editedAgent) {
      const updatedTools = editedAgent.tools || [];
      if (updatedTools.includes(tool)) {
        // Remove tool if already selected
        setEditedAgent({
          ...editedAgent,
          tools: updatedTools.filter((t) => t !== tool),
        });
      } else {
        // Add tool if not selected
        setEditedAgent({
          ...editedAgent,
          tools: [...updatedTools, tool],
        });
      }
    }
  };

  const handleApiToggle = (api: string, type: "abc_apis" | "emarking_apis") => {
    if (editedAgent) {
      const apis = editedAgent.apis || { abc_apis: [], emarking_apis: [] };
      const currentTypeApis = apis[type] || [];

      if (currentTypeApis.includes(api)) {
        // Remove API if already selected
        setEditedAgent({
          ...editedAgent,
          apis: {
            ...apis,
            [type]: currentTypeApis.filter((a) => a !== api),
          },
        });
      } else {
        // Add API if not selected
        setEditedAgent({
          ...editedAgent,
          apis: {
            ...apis,
            [type]: [...currentTypeApis, api],
          },
        });
      }
    }
  };

  const handleSaveChanges = () => {
    if (editedAgent) {
      onSave(editedAgent);
    }
  };

  const handleDelete = () => {
    if (editedAgent) {
      setIsDeleteDialogOpen(true);
    }
  };

  if (!editedAgent) return null;

  return (
    <Sheet open={isOpen} onOpenChange={setIsOpen}>
      <SheetContent className="overflow-y-auto">
        <SheetHeader>
          <SheetTitle>Edit Agent</SheetTitle>
          <SheetDescription>
            Modify the agent&aposs configuration
          </SheetDescription>
        </SheetHeader>

        <div className="py-4 space-y-4">
          <div className="space-y-2">
            <label htmlFor="agentName" className="text-sm font-medium">
              Agent Name
            </label>
            <input
              id="agentName"
              className="w-full p-2 border rounded-md"
              value={editedAgent.name || ""}
              onChange={(e) => handleInputChange("name", e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label
              htmlFor="routerDescription"
              className="text-sm font-medium flex items-center gap-1"
            >
              Router Description
              <TooltipProvider delayDuration={200}>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Info className="size-4 text-muted-foreground" />
                  </TooltipTrigger>
                  <TooltipContent>
                    <p className="max-w-xs">
                      This is the description of the agent which the router LLM
                      uses to decide who to route the request to
                    </p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </label>
            <input
              id="routerDescription"
              className="w-full p-2 border rounded-md"
              value={editedAgent.description || ""}
              onChange={(e) => handleInputChange("description", e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="agentPrompt" className="text-sm font-medium">
              Agent Prompt
            </label>
            <textarea
              id="agentPrompt"
              className="w-full p-2 border rounded-md min-h-[150px]"
              value={editedAgent.prompt || ""}
              onChange={(e) => handleInputChange("prompt", e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <h3 className="text-sm font-medium">Tools:</h3>
            <div className="space-y-2">
              {availableTools.map((tool) => (
                <div key={tool} className="flex items-center">
                  <div
                    className={`w-full p-3 border rounded-md flex items-center gap-2 cursor-pointer ${
                      editedAgent.tools?.includes(tool)
                        ? "bg-blue-50 border-blue-500"
                        : "bg-white"
                    }`}
                    onClick={() => handleToolToggle(tool)}
                  >
                    <div className="shrink-0">
                      <div
                        className={`size-5 rounded-full border flex items-center justify-center ${
                          editedAgent.tools?.includes(tool)
                            ? "border-blue-500 bg-blue-500"
                            : "border-gray-300"
                        }`}
                      >
                        {editedAgent.tools?.includes(tool) && (
                          <div className="size-2  bg-white rounded-full" />
                        )}
                      </div>
                    </div>
                    <span>{tool}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <h3 className="text-sm font-medium">ABC APIs:</h3>
            <div className="grid grid-cols-2 gap-2">
              {availableApis.abc_apis?.map((api) => (
                <div key={api} className="flex items-center">
                  <div
                    className={`w-full p-2 border rounded-md flex items-center gap-2 cursor-pointer text-sm ${
                      editedAgent.apis?.abc_apis?.includes(api)
                        ? "bg-blue-50 border-blue-500"
                        : "bg-white"
                    }`}
                    onClick={() => handleApiToggle(api, "abc_apis")}
                  >
                    <div className="shrink-0">
                      <div
                        className={`size-5 rounded-full border flex items-center justify-center ${
                          editedAgent.apis?.abc_apis?.includes(api)
                            ? "border-blue-500 bg-blue-500"
                            : "border-gray-300"
                        }`}
                      >
                        {editedAgent.apis?.abc_apis?.includes(api) && (
                          <div className="size-2  bg-white rounded-full" />
                        )}
                      </div>
                    </div>
                    <span className="truncate">{api}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="space-y-2">
            <h3 className="text-sm font-medium">eMarking APIs:</h3>
            <div className="grid grid-cols-2 gap-2">
              {availableApis.emarking_apis?.map((api) => (
                <div key={api} className="flex items-center">
                  <div
                    className={`w-full p-2 border rounded-md flex items-center gap-2 cursor-pointer text-sm ${
                      editedAgent.apis?.emarking_apis?.includes(api)
                        ? "bg-blue-50 border-blue-500"
                        : "bg-white"
                    }`}
                    onClick={() => handleApiToggle(api, "emarking_apis")}
                  >
                    <div className="shrink-0">
                      <div
                        className={`size-5 rounded-full border flex items-center justify-center ${
                          editedAgent.apis?.emarking_apis?.includes(api)
                            ? "border-blue-500 bg-blue-500"
                            : "border-gray-300"
                        }`}
                      >
                        {editedAgent.apis?.emarking_apis?.includes(api) && (
                          <div className="size-2 bg-white rounded-full" />
                        )}
                      </div>
                    </div>
                    <span className="truncate">{api}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <SheetFooter className="pt-2 flex justify-between">
          <div>
            <Button
              variant="destructive"
              onClick={handleDelete}
              disabled={isSubmitting || isDeleting}
            >
              {isDeleting ? "Deleting..." : "Delete Agent"}
            </Button>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setIsOpen(false)}
              disabled={isSubmitting || isDeleting}
            >
              Cancel
            </Button>
            <Button
              onClick={handleSaveChanges}
              disabled={isSubmitting || isDeleting}
            >
              {isSubmitting ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};

export default AgentEditSheet;
