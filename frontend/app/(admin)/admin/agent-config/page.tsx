"use client";

import { useState } from "react";
import { Plus, Trash2, Save } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

// Mock data for tools and APIs
const TOOLS = [
  { id: "search", name: "Web Search" },
  { id: "calculator", name: "Calculator" },
  { id: "code_interpreter", name: "Code Interpreter" },
];

// Mock API data - in a real app you'd have many more
const APIS = Array.from({ length: 60 }, (_, i) => ({
  id: `api-${i + 1}`,
  name: `API ${i + 1}`,
  category: [`Finance`, `Weather`, `E-commerce`, `Social Media`, `Maps`][
    Math.floor(Math.random() * 5)
  ],
}));

type Agent = {
  id: string;
  name: string;
  prompt: string;
  tools: string[];
  apis: string[];
};

export default function AgentConfigPage() {
  const [agents, setAgents] = useState<Agent[]>([
    {
      id: "1",
      name: "Customer Support Agent",
      prompt:
        "You are a helpful customer support agent. Answer questions politely and concisely.",
      tools: ["search"],
      apis: ["api-1", "api-5"],
    },
  ]);

  const [selectedAgentId, setSelectedAgentId] = useState<string | null>("1");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [newAgentName, setNewAgentName] = useState("");

  const selectedAgent = agents.find((agent) => agent.id === selectedAgentId);

  const handleCreateAgent = () => {
    if (!newAgentName.trim()) return;

    const newAgent: Agent = {
      id: Date.now().toString(),
      name: newAgentName,
      prompt: "",
      tools: [],
      apis: [],
    };

    setAgents([...agents, newAgent]);
    setSelectedAgentId(newAgent.id);
    setNewAgentName("");
    setIsCreateDialogOpen(false);
  };

  const handleDeleteAgent = (id: string) => {
    setAgents(agents.filter((agent) => agent.id !== id));
    if (selectedAgentId === id) {
      setSelectedAgentId(agents.length > 1 ? agents[0].id : null);
    }
  };

  const updateSelectedAgent = (updates: Partial<Agent>) => {
    if (!selectedAgentId) return;

    setAgents(
      agents.map((agent) =>
        agent.id === selectedAgentId ? { ...agent, ...updates } : agent
      )
    );
  };

  const toggleTool = (toolId: string) => {
    if (!selectedAgent) return;

    const tools = selectedAgent.tools.includes(toolId)
      ? selectedAgent.tools.filter((id) => id !== toolId)
      : [...selectedAgent.tools, toolId];

    updateSelectedAgent({ tools });
  };

  const toggleApi = (apiId: string) => {
    if (!selectedAgent) return;

    const apis = selectedAgent.apis.includes(apiId)
      ? selectedAgent.apis.filter((id) => id !== apiId)
      : [...selectedAgent.apis, apiId];

    updateSelectedAgent({ apis });
  };

  // Group APIs by category for better organization
  const apisByCategory = APIS.reduce((acc, api) => {
    if (!acc[api.category]) {
      acc[api.category] = [];
    }
    acc[api.category].push(api);
    return acc;
  }, {} as Record<string, typeof APIS>);

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Agent Configuration</h1>
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="mr-2 h-4 w-4" />
              Create Agent
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create New Agent</DialogTitle>
              <DialogDescription>
                Give your new AI agent a name to get started.
              </DialogDescription>
            </DialogHeader>
            <Input
              placeholder="Agent Name"
              value={newAgentName}
              onChange={(e) => setNewAgentName(e.target.value)}
            />
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsCreateDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button onClick={handleCreateAgent}>Create</Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      <div className="grid grid-cols-12 gap-6">
        {/* Agent List Sidebar */}
        <div className="col-span-3">
          <Card>
            <CardHeader>
              <CardTitle>Your Agents</CardTitle>
              <CardDescription>Select an agent to configure</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                {agents.map((agent) => (
                  <div
                    key={agent.id}
                    className={`flex justify-between items-center p-3 rounded-md cursor-pointer ${
                      selectedAgentId === agent.id
                        ? "bg-secondary"
                        : "hover:bg-secondary/50"
                    }`}
                    onClick={() => setSelectedAgentId(agent.id)}
                  >
                    <span>{agent.name}</span>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDeleteAgent(agent.id);
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Agent Configuration Area */}
        <div className="col-span-9">
          {selectedAgent ? (
            <Card>
              <CardHeader>
                <CardTitle>Configure {selectedAgent.name}</CardTitle>
                <CardDescription>
                  Customize your agent's behavior, tools, and API access
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Tabs defaultValue="basic">
                  <TabsList className="mb-4">
                    <TabsTrigger value="basic">Basic Settings</TabsTrigger>
                    <TabsTrigger value="tools">Tools</TabsTrigger>
                    <TabsTrigger value="apis">APIs</TabsTrigger>
                  </TabsList>

                  <TabsContent value="basic" className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="agent-name">Agent Name</Label>
                      <Input
                        id="agent-name"
                        value={selectedAgent.name}
                        onChange={(e) =>
                          updateSelectedAgent({ name: e.target.value })
                        }
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="agent-prompt">Agent Prompt</Label>
                      <Textarea
                        id="agent-prompt"
                        rows={10}
                        value={selectedAgent.prompt}
                        onChange={(e) =>
                          updateSelectedAgent({ prompt: e.target.value })
                        }
                        placeholder="Enter the system prompt for your agent..."
                      />
                    </div>
                  </TabsContent>

                  <TabsContent value="tools">
                    <div className="space-y-4">
                      <p className="text-sm text-muted-foreground">
                        Select the tools this agent can use
                      </p>
                      {TOOLS.map((tool) => (
                        <div
                          key={tool.id}
                          className="flex items-center space-x-2"
                        >
                          <Checkbox
                            id={`tool-${tool.id}`}
                            checked={selectedAgent.tools.includes(tool.id)}
                            onCheckedChange={() => toggleTool(tool.id)}
                          />
                          <Label htmlFor={`tool-${tool.id}`}>{tool.name}</Label>
                        </div>
                      ))}
                    </div>
                  </TabsContent>

                  <TabsContent value="apis">
                    <div className="space-y-4">
                      <p className="text-sm text-muted-foreground mb-4">
                        Select the APIs this agent can access
                      </p>

                      <div className="mb-4">
                        <p className="font-medium mb-2">Selected APIs:</p>
                        <div className="flex flex-wrap gap-2">
                          {selectedAgent.apis.length === 0 ? (
                            <p className="text-sm text-muted-foreground">
                              No APIs selected
                            </p>
                          ) : (
                            selectedAgent.apis.map((apiId) => {
                              const api = APIS.find((a) => a.id === apiId);
                              return api ? (
                                <Badge
                                  key={api.id}
                                  variant="secondary"
                                  className="cursor-pointer"
                                  onClick={() => toggleApi(api.id)}
                                >
                                  {api.name} ✕
                                </Badge>
                              ) : null;
                            })
                          )}
                        </div>
                      </div>

                      <ScrollArea className="h-[300px] pr-4">
                        {Object.entries(apisByCategory).map(
                          ([category, apis]) => (
                            <div key={category} className="mb-6">
                              <h3 className="font-medium mb-2">{category}</h3>
                              <div className="grid grid-cols-2 gap-2">
                                {apis.map((api) => (
                                  <div
                                    key={api.id}
                                    className="flex items-center space-x-2"
                                  >
                                    <Checkbox
                                      id={`api-${api.id}`}
                                      checked={selectedAgent.apis.includes(
                                        api.id
                                      )}
                                      onCheckedChange={() => toggleApi(api.id)}
                                    />
                                    <Label htmlFor={`api-${api.id}`}>
                                      {api.name}
                                    </Label>
                                  </div>
                                ))}
                              </div>
                            </div>
                          )
                        )}
                      </ScrollArea>
                    </div>
                  </TabsContent>
                </Tabs>
              </CardContent>
              <CardFooter className="flex justify-end">
                <Button>
                  <Save className="mr-2 h-4 w-4" />
                  Save Changes
                </Button>
              </CardFooter>
            </Card>
          ) : (
            <Card>
              <CardContent className="flex items-center justify-center h-[400px]">
                <div className="text-center">
                  <p className="text-muted-foreground mb-4">
                    No agent selected or create a new agent to get started
                  </p>
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button>
                        <Plus className="mr-2 h-4 w-4" />
                        Create Agent
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Create New Agent</DialogTitle>
                        <DialogDescription>
                          Give your new AI agent a name to get started.
                        </DialogDescription>
                      </DialogHeader>
                      <Input
                        placeholder="Agent Name"
                        value={newAgentName}
                        onChange={(e) => setNewAgentName(e.target.value)}
                      />
                      <DialogFooter>
                        <Button
                          variant="outline"
                          onClick={() => setIsCreateDialogOpen(false)}
                        >
                          Cancel
                        </Button>
                        <Button onClick={handleCreateAgent}>Create</Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
