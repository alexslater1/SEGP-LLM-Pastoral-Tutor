import React, { useState } from "react";
import { AgentRequest, AgentEvent } from "@/app/(admin)/actions";
import { ChevronRight } from "lucide-react";
import { cn, getRelativeTimeString } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";
import { useAgentEvents, useAgentRequests } from "@/hooks/use-agent-data";
import { useAllDownvotes, useAllUpvotes } from "@/hooks/use-all-downvotes";
import { UseQueryResult } from "@tanstack/react-query";
import { VoteAndMessage, VoteType } from "@/lib/supabase/vote";
import { Markdown } from "@/components/markdown";

export function PagedAgentRequestsTable() {
  return (
    <PagedObjectTable
      dataHook={useAgentRequests}
      idField="id"
      TableHeadings={AgentRequestTableHeadings}
      RowContents={AgentRequestsTableContents}
      ExpandedRowContents={AgentRequestsExpandedTableContents}
    />
  )
}

export function PagedAgentEventsTable() {
  return (
    <PagedObjectTable
      dataHook={useAgentEvents}
      idField="id"
      TableHeadings={AgentEventsTableHeadings}
      RowContents={AgentEventsTableContents}
      ExpandedRowContents={AgentEventsExpandedTableContents}
    />
  )
}

export function PagedDownvotesTable() {
  return (
    <PagedObjectTable
      dataHook={useAllDownvotes}
      idField="id"
      TableHeadings={DownvotesTableHeadings}
      RowContents={VotesTableContents}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}

export function PagedUpvotesTable() {
  return (
    <PagedObjectTable
      dataHook={useAllUpvotes}
      idField="id"
      TableHeadings={UpvotesTableHeadings}
      RowContents={VotesTableContents}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}

type props<T> = {
  dataHook: (page: number) => UseQueryResult<{ data: T[]; totalPages: number; }, Error>;
  idField: string;
  TableHeadings: () => React.ReactNode;
  RowContents: ({data}: {data: T}) => React.ReactNode;
  ExpandedRowContents: ({data, index}: {data: T, index: number}) => React.ReactNode;
}

function PagedObjectTable<T>({
  dataHook,
  idField,
  TableHeadings,
  RowContents,
  ExpandedRowContents,
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
        <h2 className="text-2xl text-primary font-bold tracking-tight">Agent Requests</h2>
        <p className="text-foreground">View incoming agent requests.</p>
      </div>
      <div>
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    <th className="h-12 w-[48px] px-4 first:rounded-tl-xl"></th>
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
                          onClick={() => data[idField] && toggleRow(data[idField])}
                        >
                          <td className="p-4 align-middle text-center">
                            <ChevronRight 
                              className={`size-4 text-muted-foreground transition-transform ${
                                expandedRows.has(data[idField]!) ? 'rotate-90' : ''
                              }`}
                            />
                          </td>
                          <RowContents data={data as T} />
                        </tr>
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

function getRequestsEndpointStyles (endpoint: string) {
  switch (endpoint) {
    case "/completion":
      return "bg-blue-100 text-blue-800";
    // Add more endpoint styles here if needed
    default:
      return "bg-gray-100 text-gray-800";
  }
}

function formatRequestsEndpoint (endpoint: string) {
  switch (endpoint) {
    case "/completion":
      return "Completion";
    default:
      return endpoint;
  }
};

function AgentRequestTableHeadings() {
  return (
    <>
      <th className="h-12 w-[500px] px-4 text-left align-middle font-semibold text-primary">
        Endpoint
      </th>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        Query
      </th>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
        Time
      </th>
    </>
  )
}

function AgentRequestsTableContents<T>({data}: {data: AgentRequest}) {
  return (
    <>
      <td className="p-4 align-middle">
        <span
          className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getRequestsEndpointStyles(
            data.endpoint
          )}`}
        >
          {formatRequestsEndpoint(data.endpoint)}
        </span>
      </td>
      <td className="p-4 align-middle">
        <div className="truncate max-w-[500px]">
          {data.metadata.query}
        </div>
      </td>
      <td className="p-4 align-middle text-sm text-foreground">
        {data.created_at &&
          getRelativeTimeString(new Date(data.created_at))}
      </td>
    </>
  )
}

function AgentRequestsExpandedTableContents({data}: {data: AgentRequest, index: number}) {
  return (
    <>
      <h4 className="font-semibold text-primary">Full Query:</h4>
      <p className="text-foreground mt-1 whitespace-pre-wrap">
        {data.metadata.query}
      </p>
    </>
  )
}

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

function AgentEventsTableHeadings() {
  return (
    <>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary">
        Type
      </th>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        Request ID
      </th>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
        Time
      </th>
    </>
  )
}

function AgentEventsTableContents<T>({data}: {data: AgentEvent}) {
  return (
    <>
      <td className="p-4 align-middle">
        <span
          className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${getEventTypeStyles(
            data.type
          )}`}
        >
          {data.type}
        </span>
      </td>
      <td className="p-4 align-middle font-mono text-sm text-foreground">
        {data.request_id}
      </td>
      <td className="p-4 align-middle text-sm text-foreground">
        {data.created_at &&
          getRelativeTimeString(
            new Date(data.created_at)
          )}
      </td>
    </>
  )
}

function AgentEventsExpandedTableContents({data, index}: {data: AgentEvent, index: number}) {
  return (
    <>
      {data.metadata.answer && (
        <div>
          <h4 className="font-semibold text-primary">Answer:</h4>
          <p className="text-foreground mt-1 whitespace-pre-wrap">
            {data.metadata.answer}
          </p>
        </div>
      )}
      {data.metadata.reason && (
        <div>
          <h4 className="font-semibold text-primary">Reason:</h4>
          <p className="text-foreground mt-1 whitespace-pre-wrap">
            {data.metadata.reason}
          </p>
        </div>
      )}
      {data.metadata.toolCallChoice && (
        <div>
          <h4 className="font-semibold text-primary">Tool Call:</h4>
          <p className="text-foreground mt-1">
            Name: {data.metadata.toolCallChoice.name}
          </p>
          <pre className={cn(
            "mt-2 p-4 rounded-md overflow-x-auto whitespace-pre-wrap break-words font-mono",
            index % 2 === 0 
              ? "bg-table-row-even" 
              : "bg-table-row-odd"
          )}>
            {JSON.stringify(
              JSON.parse(data.metadata.toolCallChoice.arguments),
              null,
              2
            )}
          </pre>
        </div>
      )}
      {data.metadata.toolCallResult && (
        <div>
          <h4 className="font-semibold text-primary">Tool Result:</h4>
          <pre className={cn(
            "mt-2 p-4 rounded-md whitespace-pre-wrap break-all font-mono",
            index % 2 === 0 
              ? "bg-table-row-even" 
              : "bg-table-row-odd"
          )}>
            {data.metadata.toolCallResult}
          </pre>
        </div>
      )}
    </>
  )
}

function DownvotesTableHeadings() {
  return (
    <>
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
    </>
  )
}

function UpvotesTableHeadings() {
  return (
    <>
      <th className="h-12 w-[350px] px-4 text-left align-middle font-semibold text-primary">
        Upvote ID
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
    </>
  )
}

function VotesTableContents<T>({data}: {data: VoteAndMessage}) {
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

function VotesExpandedTableContents({data}: {data: VoteAndMessage}) {
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
        <Markdown>{data.query}</Markdown>
      </div>
      <div>
        <h4 className="font-semibold text-primary">Agent Answer:</h4>
        <p className="text-foreground mt-1 whitespace-pre-wrap">
          {data.answer}
        </p>
      </div>
      <div>
        <h4 className="font-semibold text-primary">User Downvote Reason:</h4>
        <p className={cn(
          "mt-1 whitespace-pre-wrap",
          data.reason === "" 
            ? "italic text-muted-foreground" 
            : (data.type === 'downvote' ? "text-red-500" : "text-green-500")
        )}>
          {data.reason === "" ? "User did not provide a reason" : data.reason}
        </p>
      </div>
    </>
  )
}
