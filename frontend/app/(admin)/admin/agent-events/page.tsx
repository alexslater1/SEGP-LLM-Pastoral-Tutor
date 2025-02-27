"use client";

import React, { useState } from "react";
import { useAgentEvents } from "@/hooks/use-agent-data";
import { ChevronRight } from "lucide-react";
import { cn, getRelativeTimeString } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";

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
        <h2 className="text-2xl text-primary font-bold tracking-tight">Agent Events</h2>
        <p className="text-foreground">View and manage agent events.</p>
      </div>
      <div>
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    <th className="h-12 w-[48px] px-4 first:rounded-tl-xl"></th>
                    <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary">
                      Type
                    </th>
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
                      Request ID
                    </th>
                    <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
                      Time
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {isLoading ? (
                    <tr>
                      <td colSpan={4} className="p-0">
                        <div className="h-[569px] flex items-center justify-center bg-secondary/50 dark:bg-muted/90 font-bold text-5xl text-primary">
                          Loading...
                        </div>
                      </td>
                    </tr>
                  ) : (
                    data?.events.map((event, index) => (
                      <React.Fragment key={event.id}>
                        <tr
                          key={event.id}
                          className={cn(
                            "transition-colors hover:bg-muted/50 cursor-pointer",
                            index % 2 === 0 
                              ? "bg-table-row-odd" 
                              : "bg-table-row-even"
                          )}
                          onClick={() => event.id && toggleRow(event.id)}
                        >
                          <td className="p-4 align-middle text-center">
                            <ChevronRight 
                              className={`size-4 text-muted-foreground transition-transform ${
                                expandedRows.has(event.id!) ? 'rotate-90' : ''
                              }`}
                            />
                          </td>
                          <td className="p-4 align-middle">
                            <span
                              className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getEventTypeStyles(
                                event.type
                              )}`}
                            >
                              {event.type}
                            </span>
                          </td>
                          <td className="p-4 align-middle font-mono text-sm text-foreground">
                            {event.request_id}
                          </td>
                          <td className="p-4 align-middle text-sm text-foreground">
                            {event.created_at &&
                              getRelativeTimeString(
                                new Date(event.created_at)
                              )}
                          </td>
                        </tr>
                        <AnimatePresence>
                          {expandedRows.has(event.id!) && (
                            <motion.tr
                              initial={{ opacity: 1 }}
                              exit={{ opacity: 1 }}
                              className={cn(
                                index % 2 === 0 
                                  ? "bg-table-row-odd" 
                                  : "bg-table-row-even"
                              )}
                            >
                              <td colSpan={4} className="p-0">
                                <motion.div
                                  initial={{ height: 0 }}
                                  animate={{ height: "auto" }}
                                  exit={{ height: 0 }}
                                  transition={{ duration: 0.2, ease: "easeInOut" }}
                                  className="overflow-hidden"
                                >
                                  <div className="p-4 space-y-4">
                                    {event.metadata.answer && (
                                      <div>
                                        <h4 className="font-semibold text-primary">Answer:</h4>
                                        <p className="text-foreground mt-1 whitespace-pre-wrap">
                                          {event.metadata.answer}
                                        </p>
                                      </div>
                                    )}
                                    {event.metadata.reason && (
                                      <div>
                                        <h4 className="font-semibold text-primary">Reason:</h4>
                                        <p className="text-foreground mt-1 whitespace-pre-wrap">
                                          {event.metadata.reason}
                                        </p>
                                      </div>
                                    )}
                                    {event.metadata.toolCallChoice && (
                                      <div>
                                        <h4 className="font-semibold text-primary">Tool Call:</h4>
                                        <p className="text-foreground mt-1">
                                          Name: {event.metadata.toolCallChoice.name}
                                        </p>
                                        <pre className={cn(
                                          "mt-2 p-4 rounded-md overflow-x-auto whitespace-pre-wrap break-words font-mono",
                                          index % 2 === 0 
                                            ? "bg-table-row-even" 
                                            : "bg-table-row-odd"
                                        )}>
                                          {JSON.stringify(
                                            JSON.parse(event.metadata.toolCallChoice.arguments),
                                            null,
                                            2
                                          )}
                                        </pre>
                                      </div>
                                    )}
                                    {event.metadata.toolCallResult && (
                                      <div>
                                        <h4 className="font-semibold text-primary">Tool Result:</h4>
                                        <pre className={cn(
                                          "mt-2 p-4 rounded-md whitespace-pre-wrap break-all font-mono",
                                          index % 2 === 0 
                                            ? "bg-table-row-even" 
                                            : "bg-table-row-odd"
                                        )}>
                                          {event.metadata.toolCallResult}
                                        </pre>
                                      </div>
                                    )}
                                  </div>
                                </motion.div>
                              </td>
                            </motion.tr>
                          )}
                        </AnimatePresence>
                      </React.Fragment>
                    ))
                  )}
                </tbody>
              </table>

              {/* Pagination Controls */}
              <div className="flex items-center justify-between p-4 border-t bg-table-footer rounded-b-xl">
                <div className="flex-1 text-sm text-foreground">
                  Page {page + 1} of {data?.totalPages || '...'}
                </div>
                <div className="flex space-x-2">
                  <button
                    onClick={() => setPage((p) => Math.max(0, p - 1))}
                    disabled={page === 0}
                    className="px-3 py-2 text-sm font-medium rounded-md border bg-button text-button-foreground hover:bg-button/50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>
                  <button
                    onClick={() => setPage((p) => p + 1)}
                    disabled={data?.totalPages ? page >= data.totalPages - 1 : false}
                    className="px-3 py-2 text-sm font-medium rounded-md border bg-button text-button-foreground hover:bg-button/50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
