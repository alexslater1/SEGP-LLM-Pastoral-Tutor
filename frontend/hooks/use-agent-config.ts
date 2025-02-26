"use client";

import { getUserSession } from "@/lib/supabase/client";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export interface AgentsResponse {
  id: string;
  created_at: string;
  name: string;
  prompt: string;
  description: string;
  tools?: string[];
  apis: {
    abc_apis?: string[];
    emarking_apis?: string[];
  };
}

interface AgentConfigOptionsResponse {
  abc_endpoints: string[];
  emarking_endpoints: string[];
  all_tool_names: string[];
}

export interface AgentProviderAgentConfig {
  abc_apis: string[];
  emarking_apis: string[];
  tool_names: string[];
  name: string;
  prompt: string;
  description: string;
}

export function useAgentConfigs() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["agent-config"],
    queryFn: async () => {
      const data = await fetchAgents();
      console.log(data);
      return data;
    },
  });

  return { agents: data, isLoading, error };
}

export function useAgentConfigOptions() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["agent-config-options"],
    queryFn: () => fetchAgentConfigOptions(),
  });

  return { configOptions: data, isLoading, error };
}

export function useUpdateAgentConfig() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (configs: AgentProviderAgentConfig[]) =>
      updateAgentConfig(configs),
    onSuccess: () => {
      // Invalidate and refetch the agent configs after a successful update
      queryClient.invalidateQueries({ queryKey: ["agent-config"] });
    },
  });
}

async function fetchAgents(): Promise<AgentsResponse[]> {
  const login_session = await getUserSession();
  if (!login_session) {
    return [];
  }

  const response = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_AGENT_URL}/agents`,
    {
      headers: {
        Authorization: `Bearer ${login_session.access_token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch agent configs: ${response.statusText}`);
  }

  return response.json();
}

async function fetchAgentConfigOptions(): Promise<AgentConfigOptionsResponse> {
  const login_session = await getUserSession();
  if (!login_session) {
    throw new Error("User session not found");
  }

  const response = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_AGENT_URL}/agents/config`,
    {
      headers: {
        Authorization: `Bearer ${login_session.access_token}`,
      },
    }
  );

  if (!response.ok) {
    throw new Error(
      `Failed to fetch agent config options: ${response.statusText}`
    );
  }

  return response.json();
}

async function updateAgentConfig(
  configs: AgentProviderAgentConfig[]
): Promise<void> {
  const login_session = await getUserSession();
  if (!login_session) {
    throw new Error("User session not found");
  }

  const response = await fetch(
    `${process.env.NEXT_PUBLIC_BACKEND_AGENT_URL}/agents`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${login_session.access_token}`,
      },
      body: JSON.stringify({ configs }),
    }
  );

  if (!response.ok) {
    throw new Error(`Failed to update agent config: ${response.statusText}`);
  }
}
