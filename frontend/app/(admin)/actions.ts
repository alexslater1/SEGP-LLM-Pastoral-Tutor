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

export interface RagDocument {
  id: string;
  url: string;
  name: string;
  type: string;
  date_uploaded: string;
  document_size: number;
  document_type: string;
  user_id: string;
  backend_source_id: string;
}

export interface RagWebpage {
  id: string;
  url: string;
  date_uploaded: string;
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

export async function fetchRagSources(page: number, type: string) {
  const supabase = await createClient();

  const start = page * ITEMS_PER_PAGE;
  const end = start + ITEMS_PER_PAGE - 1;

  const { count } = await supabase
    .from("rag_sources")
    .select("*", { count: "exact", head: true })
    .eq("type", type);

  const { data, error } = await supabase
    .from("rag_sources")
    .select("*")
    .eq("type", type)
    .order("date_uploaded", { ascending: false })
    .range(start, end);

  if (error) throw error;

  return type === "DOCUMENT" ? {
    data: data.map((source: any) => ({
      id: source.id,
      url: source.url,
      name: source.name,
      date_uploaded: source.date_uploaded,
      document_size: source.document_size,
      document_type: source.document_type,
      user_id: source.user_id,
      backend_source_id: source.backend_source_id,
    })) as RagDocument[],
    totalPages: Math.ceil((count || 0) / ITEMS_PER_PAGE),
  } : {
    data: data.map((source: any) => ({
      id: source.id,
      url: source.url,
      date_uploaded: source.date_uploaded,
    })) as RagWebpage[],
    totalPages: Math.ceil((count || 0) / ITEMS_PER_PAGE),
  };
}

export async function fetchRagDocuments(page: number){
  return fetchRagSources(page, "DOCUMENT") as Promise<{ data: RagDocument[], totalPages: number }>;
}

export async function fetchRagWebpages(page: number) {
  return fetchRagSources(page, "WEBSITE") as Promise<{ data: RagWebpage[], totalPages: number }>;
}

export async function downloadRagDocument(name: string) {
  try { 
    console.log("Starting download for:", name);
    
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_RAG_URL}/rag-doc/download?name=${encodeURIComponent(name)}`, {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(`Download failed with status: ${response.status}`);
    }

    const arrayBuffer = await response.arrayBuffer();
    const contentType = response.headers.get("content-type") || "application/octet-stream";

    return { 
      data: arrayBuffer, 
      contentType,
      fileName: name 
    };
  } catch (error) {
    console.error("Download error details:", error);
    throw error;
  }
}

export async function uploadRagDocument(file: File) {
  try {
    console.log("Starting upload for:", file.name);
    const formData = new FormData();
    formData.append("file", file);
    
    const response = await fetch(process.env.NEXT_PUBLIC_BACKEND_RAG_URL + "/rag-doc", {
      method: "POST",
      body: formData,
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(`Upload failed with status: ${response.status}`);
    }

    const data = await response.json();
    console.log("Response data:", data);
    return data;
  } catch (error) {
    console.error("Upload error details:", error);
    throw error;
  }
};

export async function deleteRagDocument(name: string) {
  try {
    console.log("Starting delete for:", name);
    
    const response = await fetch(process.env.NEXT_PUBLIC_BACKEND_RAG_URL + `/rag-doc?name=${encodeURIComponent(name)}`, {
      method: "DELETE",
      credentials: "include",
    });


    if (!response.ok) {
      throw new Error(`Delete failed with status: ${response.status}`);
    }
  } catch (error) {
    console.error("Delete error details:", error);
    throw error;
  }
};

export async function uploadRagUrl(url: string) {
  try {
    console.log("Starting URL upload for:", url);
    
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_RAG_URL}/rag-url?url=${encodeURIComponent(url)}`, {
      method: "POST",
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(`Upload failed with status: ${response.status}`);
    }

    const data = await response.json();
    console.log("Response data:", data);
    return data;
  } catch (error) {
    console.error("Upload error details:", error);
    throw error;
  }
}

export async function deleteRagUrl(url: string) {
  try {
    console.log("Starting URL delete for:", url);
    
    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_RAG_URL}/rag-url?url=${encodeURIComponent(url)}`, {
      method: "DELETE",
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error(`Delete failed with status: ${response.status}`);
    }

    const data = await response.json();
    console.log("Response data:", data);
    return data;
  } catch (error) {
    console.error("Delete error details:", error);
    throw error;
  }
}