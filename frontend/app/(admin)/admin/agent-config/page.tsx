"use client";

import React, { useCallback, useMemo } from "react";
import ReactFlow, {
  Node,
  Edge,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  Position,
  MarkerType,
} from "reactflow";
import "reactflow/dist/style.css";
import { Split } from "lucide-react";

// Define the agent type
interface Agent {
  id: string;
  name: string;
}

// Sample agents data - this could come from an API or props
const agentsData: Agent[] = [
  { id: "agent1", name: "Agent 1" },
  { id: "agent2", name: "Agent 2" },
  { id: "agent3", name: "Agent 3" },
  // Add more agents as needed
];

const AgentConfigPage = () => {
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

  return (
    <div style={{ width: "100vw", height: "100vh" }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        fitView
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  );
};

export default AgentConfigPage;
