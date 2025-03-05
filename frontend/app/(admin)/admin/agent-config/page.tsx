"use client";

import React, { useCallback, useMemo, useState, useEffect } from "react";
import {
  ReactFlow,
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Position,
  MarkerType,
  NodeMouseHandler,
} from "reactflow";
import "reactflow/dist/style.css";
import { Split } from "lucide-react";
import AgentEditSheet from "./_components/AgentEditSheet";
import {
  useAgentConfigs,
  useAgentConfigOptions,
  AgentsResponse,
  useUpdateAgentConfig,
  AgentProviderAgentConfig,
} from "@/hooks/use-agent-config";
import { Button } from "@/components/ui/button";
import { Dialog,
  DialogContent,
  DialogTitle,
  DialogDescription,
  DialogPortal,
  DialogOverlay,
} from "@/components/ui/dialog";
import { toast } from "sonner";

// Define the agent type
interface Agent {
  id: string;
  name: string;
  description?: string;
  routerDescription?: string;
  prompt?: string;
  tools?: string[];
  apis?: {
    abc_apis?: string[];
    emarking_apis?: string[];
  };
  created_at?: string;
}

const AgentConfigPage = () => {
  const {
    agents,
    isLoading: agentsLoading,
    error: agentsError,
  } = useAgentConfigs();
  const {
    configOptions,
    isLoading: optionsLoading,
    error: optionsError,
  } = useAgentConfigOptions();

  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [editedAgent, setEditedAgent] = useState<Agent | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  const { mutate: updateAgents, isPending: isUpdating } =
    useUpdateAgentConfig();

  // Generate nodes and edges dynamically based on agents
  const { initialNodes, initialEdges } = useMemo(() => {
    // Create the fixed nodes (user and router)
    const nodes: Node[] = [
      {
        id: "user",
        data: { label: "User" },
        position: { x: 100, y: 50 },
        style: {
          background: "hsl(var(--flow-node-user))",
          border: "1px solid hsl(var(--flow-node-user-border))",
          borderRadius: "5px",
          padding: "10px",
          width: 150,
          height: 80,
          textAlign: "center",
          justifyContent: "center",
          alignItems: "center",
          display: "flex",
          color: "hsl(var(--foreground))",
        },
        sourcePosition: Position.Left,
        targetPosition: Position.Right,
      },
      {
        id: "router",
        data: {
          label: (
            <div className="flex justify-center items-center h-full">
              <Split className="rotate-90" size={24} />
            </div>
          ),
        },
        position: { x: 0, y: 350 },
        style: {
          background: "hsl(var(--flow-node-router))",
          border: "1px solid hsl(var(--flow-node-router-border))",
          borderRadius: "50%",
          padding: "10px",
          width: 80,
          height: 80,
          textAlign: "center",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
          color: "hsl(var(--foreground))",
        },
        sourcePosition: Position.Right,
        targetPosition: Position.Top,
      },
    ];

    // Create agent nodes dynamically from the fetched agents
    if (agents && agents.length > 0) {
      console.log("Creating nodes for agents:", agents);
      agents.forEach((agent, index) => {
        nodes.push({
          id: agent.id,
          data: { label: agent.name },
          position: { x: 200, y: 300 + index * 100 }, // Position agents vertically with spacing
          style: {
            background: "hsl(var(--flow-node-agent))",
            border: "1px solid hsl(var(--flow-node-agent-border))",
            borderRadius: "5px",
            padding: "10px",
            width: 150,
            height: 60,
            textAlign: "center",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            color: "hsl(var(--foreground))",
          },
          sourcePosition: Position.Right,
          targetPosition: Position.Left,
        });
      });
    }

    // Create the fixed edge from user to router
    const edges: Edge[] = [
      {
        id: "user-to-router",
        source: "user",
        target: "router",
        animated: true,
        style: { stroke: "hsl(var(--flow-edge))" },
        type: "smoothstep",
      },
    ];

    // Create edges from router to each agent and from each agent back to user
    if (agents && agents.length > 0) {
      agents.forEach((agent) => {
        edges.push({
          id: `router-to-${agent.id}`,
          source: "router",
          target: agent.id,
          animated: true,
          style: { stroke: "hsl(var(--flow-edge))" },
          type: "smoothstep",
          markerEnd: {
            type: MarkerType.Arrow,
          },
        });

        edges.push({
          id: `${agent.id}-to-user`,
          source: agent.id,
          target: "user",
          animated: true,
          style: { stroke: "hsl(var(--flow-edge))" },
          type: "smoothstep",
          markerEnd: {
            type: MarkerType.Arrow,
          },
        });
      });
    }

    return { initialNodes: nodes, initialEdges: edges };
  }, [agents]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  // Update nodes and edges when agents change
  useEffect(() => {
    if (agents && agents.length > 0) {
      setNodes(initialNodes);
      setEdges(initialEdges);
    }
  }, [agents, initialNodes, initialEdges, setNodes, setEdges]);

  const onNodeClick: NodeMouseHandler = useCallback(
    (_, node) => {
      // Find if the clicked node is an agent
      if (agents) {
        const agent = agents.find((a) => a.id === node.id);
        if (agent) {
          setSelectedAgent(agent);
          setEditedAgent({ ...agent }); // Create a copy for editing
          setIsSheetOpen(true);
        }
      }
    },
    [agents]
  );

  const handleSaveChanges = (updatedAgent: AgentsResponse) => {
    // Create a new array with all agents, either updating an existing one or adding a new one
    let updatedAgents: AgentsResponse[] = [];

    if (agents) {
      // Check if this is an existing agent (update) or a new one (add)
      const existingAgentIndex = agents.findIndex(
        (agent) => agent.id === updatedAgent.id
      );

      if (existingAgentIndex >= 0) {
        // Update existing agent
        updatedAgents = agents.map((agent) =>
          agent.id === updatedAgent.id ? updatedAgent : agent
        );
      } else {
        // Add new agent
        updatedAgents = [...agents, updatedAgent];
      }
    } else {
      // If agents is null/undefined, just use the updated agent
      updatedAgents = [updatedAgent];
    }

    // Convert each agent to AgentProviderAgentConfig format
    const configsToUpdate: AgentProviderAgentConfig[] = updatedAgents.map(
      (agent) => ({
        name: agent.name,
        prompt: agent.prompt,
        description: agent.description,
        tool_names: agent.tools || [],
        abc_apis: agent.apis?.abc_apis || [],
        emarking_apis: agent.apis?.emarking_apis || [],
      })
    );

    // Call the mutation function with the entire updated list
    updateAgents(configsToUpdate, {
      onSuccess: () => {
        console.log("Agent updated successfully");
        setIsSheetOpen(false);
      },
      onError: (error) => {
        console.error("Failed to update agent:", error);
      },
    });
  };

  // Use available tools from the API response
  const availableTools = configOptions?.all_tool_names || [];

  // Use available APIs from the API response
  const availableApis = {
    abc_apis: configOptions?.abc_endpoints || [],
    emarking_apis: configOptions?.emarking_endpoints || [],
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

  // Add this new function to handle creating a new agent
  const handleAddAgent = () => {
    // Create a new empty agent with a temporary ID
    const newAgent: AgentsResponse = {
      id: `temp-${Date.now()}`, // Temporary ID that will be replaced by the backend
      created_at: new Date().toISOString(),
      name: "New Agent",
      prompt: "",
      description: "",
      tools: [],
      apis: {
        abc_apis: [],
        emarking_apis: [],
      },
    };

    setSelectedAgent(newAgent);
    setEditedAgent(newAgent);
    setIsSheetOpen(true);
  };

  // Add this new function to handle deleting an agent
  const handleDeleteAgent = (agentId: string) => {
    if (!agents) return;

    setIsDeleting(true);

    // Filter out the agent to delete
    const updatedAgents = agents.filter((agent) => agent.id !== agentId);

    // Convert each agent to AgentProviderAgentConfig format
    const configsToUpdate: AgentProviderAgentConfig[] = updatedAgents.map(
      (agent) => ({
        name: agent.name,
        prompt: agent.prompt,
        description: agent.description,
        tool_names: agent.tools || [],
        abc_apis: agent.apis?.abc_apis || [],
        emarking_apis: agent.apis?.emarking_apis || [],
      })
    );

    // Call the mutation function with the updated list (minus the deleted agent)
    updateAgents(configsToUpdate, {
      onSuccess: () => {
        console.log("Agent deleted successfully");
        toast.success("Agent Deleted", {
          description: "The agent has been successfully deleted"
        });
        setIsSheetOpen(false);
        setIsDeleting(false);
      },
      onError: (error) => {
        console.error("Failed to delete agent:", error);
        toast.error("Delete Failed", {
          description: "Failed to delete the agent. Please try again."
        });
        setIsDeleting(false);
      },
    });
  };

  // Show loading state
  if (agentsLoading || optionsLoading) {
    return (
      <div className="flex items-center justify-center h-screen font-bold text-5xl text-primary animate-pulse">
        Loading agent configuration...
      </div>
    );
  }

  // Show error state
  if (agentsError || optionsError) {
    return (
      <div className="flex items-center justify-center h-screen text-red-500">
        Error loading agent configuration:{" "}
        {(agentsError || optionsError)?.toString()}
      </div>
    );
  }

  return (
    <div style={{ flex: 1 }}>
      {/* Add button for creating a new agent */}
      <div className="absolute top-4 right-4 z-10">
        <Button
          className="bg-button text-button-foreground hover:bg-button/50 disabled:opacity-50 disabled:cursor-not-allowed"
          onClick={handleAddAgent}>
            Add New Agent
        </Button>
      </div>

      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        fitView
        style={{ background: "hsl(var(--flow-background))" }}
      >
        <Background />
        <Controls />
      </ReactFlow>

      <AgentEditSheet
        isOpen={isSheetOpen}
        setIsOpen={setIsSheetOpen}
        agent={selectedAgent as AgentsResponse}
        onSave={handleSaveChanges}
        setIsDeleteDialogOpen={setIsDialogOpen}
        availableTools={availableTools}
        availableApis={{
          abc_apis: availableApis.abc_apis,
          emarking_apis: availableApis.emarking_apis,
        }}
        isSubmitting={isUpdating}
        isDeleting={isDeleting}
        mode={selectedAgent?.id?.startsWith('temp-') ? 'create' : 'edit'}
      />

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogTitle className="text-destructive">Delete Agent</DialogTitle>
            <DialogDescription className="text-foreground">
              Are you sure you want to delete &quot;{selectedAgent?.name}&quot;? This action cannot be undone.
            </DialogDescription>
            <div className="flex justify-end gap-2">
              <Button
                className="bg-muted hover:bg-muted/50"
                variant="outline"
                onClick={() => setIsDialogOpen(false)}>
                  Cancel
              </Button>
              <Button
                className="bg-destructive hover:bg-destructive/50"
                variant="destructive"
                onClick={() => {
                  handleDeleteAgent(selectedAgent?.id || "");
                  setIsDialogOpen(false);
                }}>
                Delete
              </Button>
            </div>
          </DialogContent>
        </DialogPortal>
      </Dialog>
    </div>
  );
};

export default AgentConfigPage;
