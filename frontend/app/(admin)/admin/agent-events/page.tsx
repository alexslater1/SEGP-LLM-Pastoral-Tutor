"use client";

import { getRelativeTimeString } from "@/lib/utils";
import { useState } from "react";
import { useAgentEvents } from "@/hooks/use-agent-data";

const getEventTypeStyles = (type: string) => {
  switch (type) {
    case "answer_success":
      return "bg-green-100 text-green-800";
    case "tool_call_choice":
      return "bg-blue-100 text-blue-800";
    case "tool_call_result":
      return "bg-purple-100 text-purple-800";
    case "error":
      return "bg-red-100 text-red-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

export default function AgentEventsPage() {
  const [page, setPage] = useState(0);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const { data, error, isLoading } = useAgentEvents(page);

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
    return <div>Error loading events: {error.message}</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Agent Events</h2>
        <p className="text-muted-foreground">View and manage agent events</p>
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
                          Type
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Request ID
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
                      {data?.events.map((event) => (
                        <>
                          <tr
                            key={event.id}
                            className="hover:bg-muted/50 cursor-pointer"
                            onClick={() => event.id && toggleRow(event.id)}
                          >
                            <td className="p-4 align-middle">
                              <span
                                className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getEventTypeStyles(
                                  event.type
                                )}`}
                              >
                                {event.type}
                              </span>
                            </td>
                            <td className="p-4 align-middle font-mono text-sm">
                              {event.request_id}
                            </td>
                            <td className="p-4 align-middle text-muted-foreground">
                              {event.created_at &&
                                getRelativeTimeString(
                                  new Date(event.created_at)
                                )}
                            </td>
                            <td className="p-4 align-middle">
                              <button className="text-blue-600 hover:text-blue-800">
                                {expandedRows.has(event.id!)
                                  ? "Hide details"
                                  : "Show details"}
                              </button>
                            </td>
                          </tr>
                          {expandedRows.has(event.id!) && (
                            <tr className="bg-muted/50">
                              <td colSpan={3} className="p-4">
                                <div className="space-y-4">
                                  {event.metadata.answer && (
                                    <div className="rounded-lg bg-background p-4">
                                      <h4 className="font-semibold">Answer:</h4>
                                      <p className="text-muted-foreground mt-1 whitespace-pre-wrap">
                                        {event.metadata.answer}
                                      </p>
                                    </div>
                                  )}
                                  {event.metadata.reason && (
                                    <div className="rounded-lg bg-background p-4">
                                      <h4 className="font-semibold">Reason:</h4>
                                      <p className="text-muted-foreground mt-1 whitespace-pre-wrap">
                                        {event.metadata.reason}
                                      </p>
                                    </div>
                                  )}
                                  {event.metadata.toolCallChoice && (
                                    <div className="rounded-lg bg-background p-4">
                                      <h4 className="font-semibold">
                                        Tool Call:
                                      </h4>
                                      <p className="text-muted-foreground mt-1">
                                        Name:{" "}
                                        {event.metadata.toolCallChoice.name}
                                      </p>
                                      <pre className="mt-2 rounded-md bg-muted p-4 overflow-x-auto">
                                        {JSON.stringify(
                                          JSON.parse(
                                            event.metadata.toolCallChoice
                                              .arguments
                                          ),
                                          null,
                                          2
                                        )}
                                      </pre>
                                    </div>
                                  )}
                                  {event.metadata.toolCallResult && (
                                    <div className="rounded-lg bg-background p-4">
                                      <h4 className="font-semibold">
                                        Tool Result:
                                      </h4>
                                      <div className="mt-2 rounded-md bg-muted p-4 max-w-full overflow-x-auto">
                                        <pre className="whitespace-pre-wrap break-all">
                                          {event.metadata.toolCallResult}
                                        </pre>
                                      </div>
                                    </div>
                                  )}
                                </div>
                              </td>
                            </tr>
                          )}
                        </>
                      ))}
                    </tbody>
                  </table>

                  {/* Pagination Controls */}
                  <div className="flex items-center justify-between p-4 border-t">
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
