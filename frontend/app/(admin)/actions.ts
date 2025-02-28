"use server";

import { createClient } from "@/lib/supabase/server";
import { getAllVotes, VoteType } from "@/lib/supabase/vote";

const ITEMS_PER_PAGE = 10;

export interface AgentEvent {
  id?: string;
  created_at?: string;
  type: string;
  request_id?: string;
  metadata: {
    answer?: string;
    reason?: string;
    toolCallChoice?: {
      name: string;
      arguments: string;
    };
    toolCallResult?: string;
  };
}

export interface AgentRequest {
  id?: string;
  created_at?: string;
  endpoint: string;
  metadata: {
    query: string;
  };
}

export async function fetchAgentEvents(page: number) {
  const supabase = await createClient();

  // Calculate the range for pagination
  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  // First, get total count
  const { count } = await supabase
    .from("agent_events")
    .select("*", { count: "exact", head: true });

  // Then get paginated data
  const { data, error } = await supabase
    .from("agent_events")
    .select("*")
    .order("created_at", { ascending: false })
    .range(start, end);

  if (error) throw error;

  return {
    data: data.map((event: any) => ({
      id: event.id,
      created_at: event.created_at,
      type: event.type,
      request_id: event.request_id,
      metadata: {
        answer: event.metadata?.answer,
        reason: event.metadata?.reason,
        toolCallChoice: event.metadata?.toolCallChoice,
        toolCallResult: event.metadata?.toolCallResult,
      },
    })) as AgentEvent[],
    totalPages: Math.ceil((count || 0) / ITEMS_PER_PAGE),
  };
}

export async function fetchAgentRequests(page: number) {
  const supabase = await createClient();

  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  // Get total count
  const { count } = await supabase
    .from("agent_requests")
    .select("*", { count: "exact", head: true });

  // Get paginated data
  const { data, error } = await supabase
    .from("agent_requests")
    .select("*")
    .order("created_at", { ascending: false })
    .range(start, end);

  if (error) throw error;

  return {
    data: data.map((request: any) => ({
      id: request.id,
      created_at: request.created_at,
      endpoint: request.endpoint,
      metadata: request.metadata,
    })) as AgentRequest[],
    totalPages: Math.ceil((count || 0) / ITEMS_PER_PAGE),
  };
}

export async function fetchVotes(type: VoteType, page: number) {
  const votes = await getAllVotes(type);

  // Sort downvotes by creation date in descending order
  const sortedVotes = votes.sort((a, b) => 
    new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  );

  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  const slicedVotes = sortedVotes.slice(start, end);

  return {
    data: slicedVotes,
    totalPages: Math.ceil((votes.length || 0) / ITEMS_PER_PAGE),
  }
}

