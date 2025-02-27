"use client";

import { getRelativeTimeString } from "@/lib/utils";
import { useEffect, useState } from "react";
import { useAllDownvotes } from "@/hooks/use-all-downvotes";
import { cn } from "@/lib/utils";
import { Markdown } from "@/components/markdown";

export default function AgentRequestsPage() {
  const [page, setPage] = useState(0);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const { data, error, isLoading } = useAllDownvotes(page);

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
        <h2 className="text-2xl font-bold tracking-tight">Downvotes</h2>
        <p className="text-muted-foreground">View all downvotes</p>
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
                          Downvote ID
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          Request ID
                        </th>
                        <th className="h-12 px-4 text-left align-middle font-medium">
                          User Email
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
                      {data?.downvotes.map((downvote) => (
                        <>
                          <tr
                            key={downvote.downvoteID}
                            className="hover:bg-muted/50 cursor-pointer"
                            onClick={() => downvote.downvoteID && toggleRow(downvote.downvoteID)}
                          >
                            <td className="p-4 align-middle">
                              <span
                                className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium`}
                              >
                                {downvote.downvoteID}
                              </span>
                            </td>
                            <td className="p-4 align-middle">
                              <span
                                className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium`}
                              >
                                {downvote.requestID}
                              </span>
                            </td>
                            <td className="p-4 align-middle">
                              <div className="truncate max-w-[400px]">
                                {downvote.userEmail}
                              </div>
                            </td>
                            <td className="p-4 align-middle text-muted-foreground">
                              {downvote.createdAt &&
                                getRelativeTimeString(
                                  new Date(downvote.createdAt)
                                )}
                            </td>
                            <td className="p-4 align-middle">
                              <button className="text-blue-600 hover:text-blue-800">
                                {expandedRows.has(downvote.downvoteID!)
                                  ? "Hide details"
                                  : "Show details"}
                              </button>
                            </td>
                          </tr>
                          {expandedRows.has(downvote.downvoteID!) && (
                            <tr className="bg-muted/50">
                              <td colSpan={4} className="p-4">
                                <div className="rounded-lg bg-background p-4">
                                  <h4 className="font-semibold">Agent Name:</h4>
                                  <p className="text-muted-foreground mt-1 whitespace-pre-wrap">
                                    {downvote.agent}
                                  </p>
                                </div>
                                <div className="rounded-lg bg-background p-4">
                                  <h4 className="font-semibold">Full Query:</h4>
                                  <Markdown>{downvote.query}</Markdown>
                                </div>
                                <div className="rounded-lg bg-background p-4">
                                  <h4 className="font-semibold">Agent Answer:</h4>
                                  <p className="text-muted-foreground mt-1 whitespace-pre-wrap">
                                    {downvote.answer}
                                  </p>
                                </div>
                                <div className="rounded-lg bg-background p-4">
                                  <h4 className="font-semibold">User Downvote Reason:</h4>
                                  <p className={cn("text-muted-foreground mt-1 whitespace-pre-wrap", 
                                    downvote.reason ? "text-red-500" : "italic")}>
                                    {downvote.reason ?? "User did not provide a reason"}
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
