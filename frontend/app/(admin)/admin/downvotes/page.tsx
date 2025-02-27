"use client";

import React, { useState } from "react";
import { useAllDownvotes } from "@/hooks/use-all-downvotes";
import { cn, getRelativeTimeString } from "@/lib/utils";
import { Markdown } from "@/components/markdown";
import { ChevronRight } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

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
        <h2 className="text-2xl text-primary font-bold tracking-tight">Downvotes</h2>
        <p className="text-foreground">View all downvotes</p>
      </div>
      <div>
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    <th className="h-12 w-[48px] px-4 first:rounded-tl-xl"></th>
                    <th className="h-12 w-[350px] px-4 text-left align-middle font-semibold text-primary">
                      Downvote ID
                    </th>
                    <th className="h-12 w-[350px] px-4 text-left align-middle font-semibold text-primary">
                      Request ID
                    </th>
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
                      User Email
                    </th>
                    <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
                      Time
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {isLoading ? (
                    <tr>
                      <td colSpan={5} className="p-0">
                        <div className="h-[569px] flex items-center justify-center bg-secondary/50 dark:bg-muted/90 font-bold text-5xl text-primary">
                          Loading...
                        </div>
                      </td>
                    </tr>
                  ) : (
                    data?.downvotes.map((downvote, index) => (
                      <React.Fragment key={downvote.downvoteID}>
                        <tr
                          className={cn(
                            "transition-colors hover:bg-muted/50 cursor-pointer",
                            index % 2 === 0 
                              ? "bg-table-row-odd" 
                              : "bg-table-row-even"
                          )}
                          onClick={() => downvote.downvoteID && toggleRow(downvote.downvoteID)}
                        >
                          <td className="p-4 align-middle text-center">
                            <ChevronRight 
                              className={`size-4 text-muted-foreground transition-transform ${
                                expandedRows.has(downvote.downvoteID!) ? 'rotate-90' : ''
                              }`}
                            />
                          </td>
                          <td className="p-4 align-middle font-mono text-sm text-foreground">
                            {downvote.downvoteID}
                          </td>
                          <td className="p-4 align-middle font-mono text-sm text-foreground">
                            {downvote.requestID}
                          </td>
                          <td className="p-4 align-middle">
                            <div className="truncate max-w-[400px] text-foreground">
                              {downvote.userEmail}
                            </div>
                          </td>
                          <td className="p-4 align-middle text-sm text-foreground">
                            {downvote.createdAt &&
                              getRelativeTimeString(new Date(downvote.createdAt))}
                          </td>
                        </tr>
                        <AnimatePresence>
                          {expandedRows.has(downvote.downvoteID!) && (
                            <motion.tr
                              initial={{ opacity: 1 }}
                              exit={{ opacity: 1 }}
                              className={cn(
                                index % 2 === 0 
                                  ? "bg-table-row-odd" 
                                  : "bg-table-row-even"
                              )}
                            >
                              <td colSpan={5} className="p-0">
                                <motion.div
                                  initial={{ height: 0 }}
                                  animate={{ height: "auto" }}
                                  exit={{ height: 0 }}
                                  transition={{ duration: 0.2, ease: "easeInOut" }}
                                  className="overflow-hidden"
                                >
                                  <div className="p-4 space-y-4">
                                    <div className="rounded-lg bg-background/50 p-4">
                                      <h4 className="font-semibold text-primary">Agent Name:</h4>
                                      <p className="text-foreground mt-1 whitespace-pre-wrap">
                                        {downvote.agent}
                                      </p>
                                    </div>
                                    <div className="rounded-lg bg-background/50 p-4">
                                      <h4 className="font-semibold text-primary">Full Query:</h4>
                                      <Markdown>{downvote.query}</Markdown>
                                    </div>
                                    <div className="rounded-lg bg-background/50 p-4">
                                      <h4 className="font-semibold text-primary">Agent Answer:</h4>
                                      <p className="text-foreground mt-1 whitespace-pre-wrap">
                                        {downvote.answer}
                                      </p>
                                    </div>
                                    <div className="rounded-lg bg-background/50 p-4">
                                      <h4 className="font-semibold text-primary">User Downvote Reason:</h4>
                                      <p className={cn(
                                        "mt-1 whitespace-pre-wrap",
                                        downvote.reason === "" 
                                          ? "italic text-muted-foreground" 
                                          : "text-red-500"
                                      )}>
                                        {downvote.reason === "" ? "User did not provide a reason" : downvote.reason}
                                      </p>
                                    </div>
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
                    disabled={!data?.totalPages || page >= data.totalPages - 1}
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
