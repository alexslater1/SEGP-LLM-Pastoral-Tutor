"use server";

import { createClient } from "@/lib/supabase/server";
import { getAllDownvotes } from "@/lib/supabase/vote";

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

interface AgentRequest {
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
    events: data.map((event: any) => ({
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
    requests: data.map((request: any) => ({
      id: request.id,
      created_at: request.created_at,
      endpoint: request.endpoint,
      metadata: request.metadata,
    })) as AgentRequest[],
    totalPages: Math.ceil((count || 0) / ITEMS_PER_PAGE),
  };
}

export async function fetchDownvotes(page: number) {
  const downvotes = await getAllDownvotes();

  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  const slicedDownvotes = downvotes.slice(start, end);

  return {
    downvotes: slicedDownvotes,
    totalPages: Math.ceil((downvotes.length || 0) / ITEMS_PER_PAGE),
  }
}

