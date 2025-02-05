"use client";

import { useQuery } from "@tanstack/react-query";
import { getRelativeTimeString } from "@/lib/utils";
import { createClient } from "@/lib/supabase/client";
import { useState } from "react";

interface AgentRequest {
  id?: string;
  created_at?: string;
  endpoint: string;
  metadata: {
    query: string;
  };
}

const ITEMS_PER_PAGE = 10;

function useAgentRequests(page: number) {
  const supabase = createClient();

  return useQuery({
    queryKey: ["agent-requests", page],
    queryFn: async () => {
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
    },
  });
}

const getEndpointStyles = (endpoint: string) => {
  switch (endpoint) {
    case "/completion":
      return "bg-blue-100 text-blue-800";
    // Add more endpoint styles here if needed
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const formatEndpoint = (endpoint: string) => {
  switch (endpoint) {
    case "/completion":
      return "Completion";
    default:
      return endpoint;
  }
};

export default function AgentRequestsPage() {
  const [page, setPage] = useState(0);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const { data, error, isLoading } = useAgentRequests(page);

  const toggleRow = (id: string) => {
    const newExpandedRows = new Set(expandedRows);
    if (expandedRows.has(id)) {
      newExpandedRows.delete(id);
    } else {
      newExpandedRows.add(id);
    }
    setExpandedRows(newExpandedRows);
  };

  if (error) {
    return <div>Error loading requests: {error.message}</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Agent Requests</h2>
        <p className="text-muted-foreground">View incoming agent requests</p>
      </div>
      <div className="border-t">
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-x-auto">
              {isLoading ? (
                <div className="p-8 text-center">Loading...</div>
              ) : (
                <>
                  <table className="w-full">
                    <thead>
                      <tr className="border-b">
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Endpoint
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Query
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Time
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Actions
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      {data?.requests.map((request) => (
                        <>
                          <tr
                            key={request.id}
                            className="hover:bg-muted/50 cursor-pointer"
                            onClick={() => request.id && toggleRow(request.id)}
                          >
                            <td className="p-4 align-middle">
                              <span
                                className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getEndpointStyles(
                                  request.endpoint
                                )}`}
                              >
                                {formatEndpoint(request.endpoint)}
                              </span>
                            </td>
                            <td className="p-4 align-middle">
                              <div className="truncate max-w-[400px]">
                                {request.metadata.query}
                              </div>
                            </td>
                            <td className="p-4 align-middle text-muted-foreground">
                              {request.created_at &&
                                getRelativeTimeString(
                                  new Date(request.created_at)
                                )}
                            </td>
                            <td className="p-4 align-middle">
                              <button className="text-blue-600 hover:text-blue-800">
                                {expandedRows.has(request.id!)
                                  ? "Hide details"
                                  : "Show details"}
                              </button>
                            </td>
                          </tr>
                          {expandedRows.has(request.id!) && (
                            <tr className="bg-muted/50">
                              <td colSpan={4} className="p-4">
                                <div className="rounded-lg bg-background p-4">
                                  <h4 className="font-semibold">Full Query:</h4>
                                  <p className="text-muted-foreground mt-1 whitespace-pre-wrap">
                                    {request.metadata.query}
                                  </p>
                                </div>
                              </td>
                            </tr>
                          )}
                        </>
                      ))}
                    </tbody>
                  </table>

                  {/* Pagination Controls */}
                  <div className="flex items-center justify-between px-4 py-4 border-t">
                    <div className="flex-1 text-sm text-muted-foreground">
                      Page {page + 1} of {data?.totalPages}
                    </div>
                    <div className="flex space-x-2">
                      <button
                        onClick={() => setPage((p) => Math.max(0, p - 1))}
                        disabled={page === 0}
                        className="px-3 py-2 text-sm font-medium rounded-md border hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Previous
                      </button>
                      <button
                        onClick={() => setPage((p) => p + 1)}
                        disabled={
                          !data?.totalPages || page >= data.totalPages - 1
                        }
                        className="px-3 py-2 text-sm font-medium rounded-md border hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        Next
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
