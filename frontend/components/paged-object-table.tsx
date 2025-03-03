"use client";

import React, { useState } from "react";
import { cn, getRelativeTimeString } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";
import { UseQueryResult } from "@tanstack/react-query";
import { VoteAndMessage } from "@/lib/supabase/vote";
import { Markdown } from "@/components/markdown";
import { ChevronRight } from "lucide-react";

export const ITEMS_PER_PAGE = 10;

export enum TableRowType {
  EXPANDABLE = "expandable",
  CLICKABLE = "clickable",
}

type hookResultData<T> = { data: T[]; totalPages: number; }
type hookResult<T> = { data: hookResultData<T>, error: Error | null, isLoading: boolean }

type props<T> = {
  title: string;
  description: string;
  dataHook: (page: number) => (UseQueryResult<hookResultData<T>, Error> | hookResult<T>);
  idField: string;
  TableHeadings: () => React.ReactNode;
  RowContents: ({data}: {data: T}) => React.ReactNode;
  ExpandedRowContents: ({data, index}: {data: T, index: number}) => React.ReactNode;
  rowType?: TableRowType;
  rowClickHandler?: (data: T) => void;
}

export function PagedObjectTable<T>({
  title,
  description,
  dataHook,
  idField,
  TableHeadings,
  RowContents,
  ExpandedRowContents,
  rowType = TableRowType.EXPANDABLE,
  rowClickHandler = () => {},
}: props<T>) {
  const [page, setPage] = useState(0);
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());
  const { data, error, isLoading } = dataHook(page);

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
        <h2 className="text-2xl text-primary font-bold tracking-tight">{title}</h2>
        <p className="text-foreground">{description}</p>
      </div>
      <div>
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    {rowType === TableRowType.EXPANDABLE && (
                      <th className="h-12 w-[48px] px-4 first:rounded-tl-xl"></th>
                    )}
                    <TableHeadings />
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
                  ) : !data || data.data.length === 0 ? 
                  (
                    <tr>
                      <td colSpan={5} className="p-0">
                        <div className="h-[569px] flex items-center justify-center bg-secondary/50 dark:bg-muted/90 font-bold text-5xl text-primary">
                          No table entries found
                        </div>
                      </td>
                    </tr>
                  ) : (
                    data?.data.map((data: any, index) => (
                      <React.Fragment key={data[idField]}>
                        <tr
                          key={data[idField]}
                          className={cn(
                            "transition-colors hover:bg-muted/50 cursor-pointer",
                            index % 2 === 0 
                              ? "bg-table-row-odd" 
                              : "bg-table-row-even"
                          )}
                          onClick={rowType === TableRowType.CLICKABLE ? 
                            () => rowClickHandler(data as T) : 
                            () => data[idField] && toggleRow(data[idField])}
                        >
                          {rowType === TableRowType.EXPANDABLE && (
                            <td className="p-4 align-middle text-center">
                              <ChevronRight 
                                className={`size-4 text-muted-foreground transition-transform ${
                                  expandedRows.has(data[idField]!) ? "rotate-90" : ""
                                }`}
                              />
                            </td>
                          )}
                          <RowContents data={data as T} />
                        </tr>
                        {rowType === TableRowType.EXPANDABLE && (
                          <AnimatePresence>
                            {expandedRows.has(data[idField]!) && (
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
                                      <ExpandedRowContents data={data as T} index={index} />
                                    </div>
                                  </motion.div>
                                </td>
                              </motion.tr>
                            )}
                          </AnimatePresence>
                        )}
                      </React.Fragment>
                    ))
                  )}
                </tbody>
              </table>

              <div className="flex items-center justify-between p-4 border-t bg-table-footer rounded-b-xl">
                <div className="flex-1 text-sm text-foreground">
                  Page {page + 1} of {data?.totalPages || "..."}
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
                    disabled={!isLoading && data?.totalPages ? page >= data.totalPages - 1 : false}
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
  )
}

export function VotesTableContents<T>({data}: {data: VoteAndMessage}) {
  return (
    <>
      <td className="p-4 align-middle font-mono text-sm text-foreground">
        {data.id}
      </td>
      <td className="p-4 align-middle font-mono text-sm text-foreground">
        {data.requestID}
      </td>
      <td className="p-4 align-middle">
        <div className="truncate max-w-[400px] text-foreground">
          {data.userEmail}
        </div>
      </td>
      <td className="p-4 align-middle text-sm text-foreground">
        {data.createdAt &&
          getRelativeTimeString(new Date(data.createdAt))}
      </td>
    </>
  )
}

export function VotesExpandedTableContents({data}: {data: VoteAndMessage}) {
  return (
    <>
      <div>
        <h4 className="font-semibold text-primary">Agent Name:</h4>
        <p className="text-foreground mt-1 whitespace-pre-wrap">
          {data.agent}
        </p>
      </div>
      <div>
        <h4 className="font-semibold text-primary">Full Query:</h4>
        {data.query}
      </div>
      <div>
        <h4 className="font-semibold text-primary">Agent Answer:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.answer}</Markdown>
        </div>
      </div>
      <div>
        <h4 className="font-semibold text-primary">{data.type === "downvote" ? "User Downvote Reason:" : "User Upvote Reason:"}</h4>
        <p className={cn(
          "mt-1 whitespace-pre-wrap",
          data.reason === "" 
            ? "italic text-muted-foreground" 
            : (data.type === "downvote" ? "text-error" : "text-success")
        )}>
          {data.reason === "" ? "User did not provide a reason" : data.reason}
        </p>
      </div>
    </>
  )
}