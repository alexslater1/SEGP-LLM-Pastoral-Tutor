"use client";

import React, { useCallback, useMemo, useState } from "react";
import ReactFlow, {
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
import { Split, Info } from "lucide-react";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetFooter,
  SheetClose,
} from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

// Define the agent type
interface Agent {
  id: string;
  name: string;
  description?: string;
  routerDescription?: string;
  prompt?: string;
  tools?: string[];
  apis?: string[];
}

// Sample agents data - this could come from an API or props
const agentsData: Agent[] = [
  {
    id: "agent1",
    name: "Agent 1",
    description: "A general-purpose assistant that can handle various tasks.",
    routerDescription: "General purpose assistant",
    prompt: "You are a helpful assistant that provides general information.",
    tools: ["Google Search"],
    apis: ["Weather API"],
  },
  {
    id: "agent2",
    name: "Agent 2",
    description: "Specialized in technical support and troubleshooting.",
  },
  {
    id: "agent3",
    name: "Agent 3",
    description: "Creative assistant focused on content generation.",
  },
  // Add more agents as needed
];

const AgentConfigPage = () => {
  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [editedAgent, setEditedAgent] = useState<Agent | null>(null);

  // Generate nodes and edges dynamically based on agents
  const { initialNodes, initialEdges } = useMemo(() => {
    // Create the fixed nodes (user and router)
    const nodes: Node[] = [
      {
        id: "user",
        data: { label: "User" },
        position: { x: 100, y: 50 },
        style: {
          background: "#f5f5f5",
          border: "1px solid #ddd",
          borderRadius: "5px",
          padding: "10px",
          width: 150,
          height: 80,
          textAlign: "center",
          justifyContent: "center",
          alignItems: "center",
          display: "flex",
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
          background: "#ffffff",
          border: "1px solid #ddd",
          borderRadius: "50%",
          padding: "10px",
          width: 80,
          height: 80,
          textAlign: "center",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        },
        sourcePosition: Position.Right,
        targetPosition: Position.Top,
      },
    ];

    // Create agent nodes dynamically
    agentsData.forEach((agent, index) => {
      nodes.push({
        id: agent.id,
        data: { label: agent.name },
        position: { x: 200, y: 300 + index * 100 }, // Position agents vertically with spacing
        style: {
          background: "#e6f7ff",
          border: "1px solid #91d5ff",
          borderRadius: "5px",
          padding: "10px",
          width: 150,
          height: 60,
          textAlign: "center",
          display: "flex",
          justifyContent: "center",
          alignItems: "center",
        },
        sourcePosition: Position.Right,
        targetPosition: Position.Left,
      });
    });

    // Create the fixed edge from user to router
    const edges: Edge[] = [
      {
        id: "user-to-router",
        source: "user",
        target: "router",
        animated: true,
        style: { stroke: "#555" },
        type: "smoothstep",
      },
    ];

    // Create edges from router to each agent and from each agent back to user
    agentsData.forEach((agent) => {
      edges.push({
        id: `router-to-${agent.id}`,
        source: "router",
        target: agent.id,
        animated: true,
        style: { stroke: "#555" },
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
        style: { stroke: "#555" },
        type: "smoothstep",
        markerEnd: {
          type: MarkerType.Arrow,
        },
      });
    });

    return { initialNodes: nodes, initialEdges: edges };
  }, [agentsData]);

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes);
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

  const onNodeClick: NodeMouseHandler = useCallback(
    (_, node) => {
      // Find if the clicked node is an agent
      const agent = agentsData.find((a) => a.id === node.id);
      if (agent) {
        setSelectedAgent(agent);
        setEditedAgent({ ...agent }); // Create a copy for editing
        setIsSheetOpen(true);
      }
    },
    [agentsData]
  );

  const handleSaveChanges = () => {
    if (editedAgent) {
      // In a real app, you would save these changes to your backend
      // For now, we'll just update the local state
      console.log("Saving changes:", editedAgent);
      setSelectedAgent(editedAgent);
      setIsSheetOpen(false);
    }
  };

  const handleInputChange = (field: keyof Agent, value: string) => {
    if (editedAgent) {
      setEditedAgent({
        ...editedAgent,
        [field]: value,
      });
    }
  };

  // Available tools options
  const availableTools = ["Google Search", "Google Maps"];

  // Available API options
  const availableApis = [
    "Weather API",
    "Stripe Payments",
    "Twilio SMS",
    "SendGrid Email",
    "GitHub API",
    "Slack API",
    "Twitter API",
    "Spotify API",
    "Google Calendar",
    "Salesforce CRM",
    "Shopify API",
    "Zoom API",
    "LinkedIn API",
    "OpenAI API",
    "HubSpot API",
    "Dropbox API",
    "PayPal API",
    "AWS S3",
    "Azure Cognitive Services",
    "Google Cloud Vision",
  ];

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

  const handleApiToggle = (api: string) => {
    if (editedAgent) {
      const updatedApis = editedAgent.apis || [];
      if (updatedApis.includes(api)) {
        // Remove API if already selected
        setEditedAgent({
          ...editedAgent,
          apis: updatedApis.filter((a) => a !== api),
        });
      } else {
        // Add API if not selected
        setEditedAgent({
          ...editedAgent,
          apis: [...updatedApis, api],
        });
      }
    }
  };

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        fitView
      >
        <Background />
        <Controls />
      </ReactFlow>

      <Sheet open={isSheetOpen} onOpenChange={setIsSheetOpen}>
        <SheetContent className="overflow-y-auto">
          <SheetHeader>
            <SheetTitle>Edit Agent</SheetTitle>
            <SheetDescription>
              Modify the agent's configuration
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
                value={editedAgent?.name || ""}
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
                      <Info className="h-4 w-4 text-muted-foreground" />
                    </TooltipTrigger>
                    <TooltipContent>
                      <p className="max-w-xs">
                        This is the description of the agent which the router
                        LLM uses to decide who to route the request to
                      </p>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              </label>
              <input
                id="routerDescription"
                className="w-full p-2 border rounded-md"
                value={editedAgent?.routerDescription || ""}
                onChange={(e) =>
                  handleInputChange("routerDescription", e.target.value)
                }
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="agentPrompt" className="text-sm font-medium">
                Agent Prompt
              </label>
              <textarea
                id="agentPrompt"
                className="w-full p-2 border rounded-md min-h-[150px]"
                value={editedAgent?.prompt || ""}
                onChange={(e) => handleInputChange("prompt", e.target.value)}
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="agentDescription" className="text-sm font-medium">
                Description
              </label>
              <textarea
                id="agentDescription"
                className="w-full p-2 border rounded-md"
                value={editedAgent?.description || ""}
                onChange={(e) =>
                  handleInputChange("description", e.target.value)
                }
              />
            </div>

            <div className="space-y-2">
              <h3 className="text-sm font-medium">Tools:</h3>
              <div className="space-y-2">
                {availableTools.map((tool) => (
                  <div key={tool} className="flex items-center">
                    <div
                      className={`w-full p-3 border rounded-md flex items-center gap-2 cursor-pointer ${
                        editedAgent?.tools?.includes(tool)
                          ? "bg-blue-50 border-blue-500"
                          : "bg-white"
                      }`}
                      onClick={() => handleToolToggle(tool)}
                    >
                      <div className="flex-shrink-0">
                        <div
                          className={`w-5 h-5 rounded-full border flex items-center justify-center ${
                            editedAgent?.tools?.includes(tool)
                              ? "border-blue-500 bg-blue-500"
                              : "border-gray-300"
                          }`}
                        >
                          {editedAgent?.tools?.includes(tool) && (
                            <div className="w-2 h-2 bg-white rounded-full" />
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
              <h3 className="text-sm font-medium">APIs:</h3>
              <div className="grid grid-cols-2 gap-2">
                {availableApis.map((api) => (
                  <div key={api} className="flex items-center">
                    <div
                      className={`w-full p-2 border rounded-md flex items-center gap-2 cursor-pointer text-sm ${
                        editedAgent?.apis?.includes(api)
                          ? "bg-blue-50 border-blue-500"
                          : "bg-white"
                      }`}
                      onClick={() => handleApiToggle(api)}
                    >
                      <div className="flex-shrink-0">
                        <div
                          className={`w-4 h-4 rounded-full border flex items-center justify-center ${
                            editedAgent?.apis?.includes(api)
                              ? "border-blue-500 bg-blue-500"
                              : "border-gray-300"
                          }`}
                        >
                          {editedAgent?.apis?.includes(api) && (
                            <div className="w-2 h-2 bg-white rounded-full" />
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

          <SheetFooter className="pt-2">
            <Button variant="outline" onClick={() => setIsSheetOpen(false)}>
              Cancel
            </Button>
            <Button onClick={handleSaveChanges}>Save Changes</Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>
    </div>
  );
};

export default AgentConfigPage;
