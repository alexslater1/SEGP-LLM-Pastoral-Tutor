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
import { toast } from "sonner";

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
  mode?: "create" | "edit";
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
  mode = "edit",
}) => {
  const [editedAgent, setEditedAgent] = useState<AgentsResponse | null>(null);
  const [validationErrors, setValidationErrors] = useState<{
    name?: string;
    description?: string;
    prompt?: string;
  }>({});

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

  const validateForm = () => {
    const errors: {
      name?: string;
      description?: string;
      prompt?: string;
    } = {};

    if (!editedAgent?.name?.trim()) {
      errors.name = "* Agent Name cannot be blank";
    }
    if (!editedAgent?.description?.trim()) {
      errors.description = "* Agent Description cannot be blank";
    }
    if (!editedAgent?.prompt?.trim()) {
      errors.prompt = "* Agent Prompt cannot be blank";
    }

    setValidationErrors(errors);
    return Object.keys(errors).length === 0;
  };

  const handleSaveChanges = () => {
    if (editedAgent) {
      if (!validateForm()) {
        toast.error("Missing Required Fields", {
          description: "Please fill in all required fields marked with *",
        });
        return;
      }
      onSave(editedAgent);
      toast.success("Changes Saved", {
        description: "Agent configuration has been updated successfully",
      });
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
          <SheetTitle className="text-2xl font-bold text-primary flex items-center justify-between">
            {mode === "create" ? "Create Agent" : "Edit Agent"}
            {mode === "edit" && (
              <Button
                variant="destructive"
                onClick={handleDelete}
                disabled={isSubmitting || isDeleting}
                size="sm"
                className="ml-8"
              >
                {isDeleting ? "Deleting..." : "Delete Agent"}
              </Button>
            )}
          </SheetTitle>
          <SheetDescription className="text-foreground">
            {mode === "create"
              ? "Set up the configuration for a new agent."
              : "Modify the agent's configuration."}
            <br />
            <span className="text-error">*</span> Indicates a required field
          </SheetDescription>
        </SheetHeader>

        <div className="py-4">
          <div className="space-y-4">
            <div className="space-y-2">
              <label
                htmlFor="agentName"
                className="text-md font-bold text-primary flex items-center gap-1"
              >
                <span className="text-error">*</span>
                Agent Name
              </label>
              <input
                id="agentName"
                className={`w-full p-2 border rounded-md bg-muted/25 dark:bg-muted ${
                  validationErrors.name ? "border-error" : "border-input"
                }`}
                value={editedAgent.name || ""}
                onChange={(e) => handleInputChange("name", e.target.value)}
              />
              {validationErrors.name && (
                <p className="text-sm text-error">{validationErrors.name}</p>
              )}
            </div>

            <div className="space-y-2">
              <label
                htmlFor="routerDescription"
                className="text-md font-bold text-primary flex items-center gap-1"
              >
                <span className="text-error">*</span>
                Router Description
                <TooltipProvider delayDuration={200}>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Info className="size-4 text-muted-foreground hover:text-primary" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p className="max-w-xs font-normal text-foreground">
                        This is the description of the agent which the router
                        LLM uses to decide who to route the request to
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </label>
              <input
                id="routerDescription"
                className={`w-full p-2 border rounded-md bg-muted/25 dark:bg-muted ${
                  validationErrors.description ? "border-error" : "border-input"
                }`}
                value={editedAgent.description || ""}
                onChange={(e) =>
                  handleInputChange("description", e.target.value)
                }
              />
              {validationErrors.description && (
                <p className="text-sm text-error">
                  {validationErrors.description}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <label
                htmlFor="agentPrompt"
                className="text-md font-bold text-primary flex items-center gap-1"
              >
                <span className="text-error">*</span>
                Agent Prompt
              </label>
              <textarea
                id="agentPrompt"
                className={`w-full p-2 border rounded-md min-h-[150px] bg-muted/25 dark:bg-muted ${
                  validationErrors.prompt ? "border-error" : "border-input"
                }`}
                value={editedAgent.prompt || ""}
                onChange={(e) => handleInputChange("prompt", e.target.value)}
              />
              {validationErrors.prompt && (
                <p className="text-sm text-error">{validationErrors.prompt}</p>
              )}
            </div>

            <div className="space-y-2">
              <h3 className="text-md font-bold text-primary">Tools</h3>
              <div className="space-y-2">
                {availableTools.map((tool) => (
                  <div key={tool} className="flex items-center">
                    <div
                      className={`w-full p-3 border rounded-md flex items-center gap-2 cursor-pointer ${
                        editedAgent.tools?.includes(tool)
                          ? "bg-primary/10 border-primary"
                          : "bg-background border-input hover:bg-accent hover:dark:bg-accent/50"
                      }`}
                      onClick={() => handleToolToggle(tool)}
                    >
                      <div className="shrink-0">
                        <div
                          className={`size-5 rounded-full border flex items-center justify-center ${
                            editedAgent.tools?.includes(tool)
                              ? "border-primary bg-primary"
                              : "border-input"
                          }`}
                        >
                          {editedAgent.tools?.includes(tool) && (
                            <div className="size-2 bg-white rounded-full" />
                          )}
                        </div>
                      </div>
                      <span className="text-foreground">{tool}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <h3 className="text-md font-bold text-primary">ABC APIs</h3>
              <div className="grid grid-cols-2 gap-2">
                {availableApis.abc_apis?.map((api) => (
                  <div key={api} className="flex items-center">
                    <div
                      className={`w-full p-2 border rounded-md flex items-center gap-2 cursor-pointer text-sm ${
                        editedAgent.apis?.abc_apis?.includes(api)
                          ? "bg-primary/10 border-primary"
                          : "bg-background border-input hover:bg-accent hover:dark:bg-accent/50"
                      }`}
                      onClick={() => handleApiToggle(api, "abc_apis")}
                    >
                      <div className="shrink-0">
                        <div
                          className={`size-5 rounded-full border flex items-center justify-center ${
                            editedAgent.apis?.abc_apis?.includes(api)
                              ? "border-primary bg-primary"
                              : "border-input"
                          }`}
                        >
                          {editedAgent.apis?.abc_apis?.includes(api) && (
                            <div className="size-2 bg-white rounded-full" />
                          )}
                        </div>
                      </div>
                      <span className="text-foreground truncate">{api}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="space-y-2">
              <h3 className="text-md font-bold text-primary">eMarking APIs</h3>
              <div className="grid grid-cols-2 gap-2">
                {availableApis.emarking_apis?.map((api) => (
                  <div key={api} className="flex items-center">
                    <div
                      className={`w-full p-2 border rounded-md flex items-center gap-2 cursor-pointer text-sm ${
                        editedAgent.apis?.emarking_apis?.includes(api)
                          ? "bg-primary/10 border-primary"
                          : "bg-background border-input hover:bg-accent hover:dark:bg-accent/50"
                      }`}
                      onClick={() => handleApiToggle(api, "emarking_apis")}
                    >
                      <div className="shrink-0">
                        <div
                          className={`size-5 rounded-full border flex items-center justify-center ${
                            editedAgent.apis?.emarking_apis?.includes(api)
                              ? "border-primary bg-primary"
                              : "border-input"
                          }`}
                        >
                          {editedAgent.apis?.emarking_apis?.includes(api) && (
                            <div className="size-2 bg-white rounded-full" />
                          )}
                        </div>
                      </div>
                      <span className="text-foreground truncate">{api}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="flex w-full gap-2 pt-10">
            {mode === "create" ? (
              <>
                <div className="flex-1">
                  <Button
                    variant="outline"
                    onClick={() => setIsOpen(false)}
                    disabled={isSubmitting}
                    className="w-full bg-muted hover:bg-muted/50"
                  >
                    Cancel
                  </Button>
                </div>
                <div className="flex-1">
                  <Button
                    onClick={handleSaveChanges}
                    disabled={isSubmitting}
                    className="w-full bg-button text-button-foreground hover:bg-button/90"
                  >
                    {isSubmitting ? "Creating..." : "Create"}
                  </Button>
                </div>
              </>
            ) : (
              <>
                <div className="flex-1">
                  <Button
                    variant="outline"
                    onClick={() => setIsOpen(false)}
                    disabled={isSubmitting || isDeleting}
                    className="w-full bg-muted hover:bg-muted/50"
                  >
                    Cancel
                  </Button>
                </div>
                <div className="flex-1">
                  <Button
                    onClick={handleSaveChanges}
                    disabled={isSubmitting || isDeleting}
                    className="w-full bg-button text-button-foreground hover:bg-button/90"
                  >
                    {isSubmitting ? "Saving..." : "Save Changes"}
                  </Button>
                </div>
                {/* <div className="flex-1">
                  <Button
                    variant="destructive"
                    onClick={handleDelete}
                    disabled={isSubmitting || isDeleting}
                    className="w-full"
                  >
                    {isDeleting ? "Deleting..." : "Delete Agent"}
                  </Button>
                </div> */}
              </>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default AgentEditSheet;
