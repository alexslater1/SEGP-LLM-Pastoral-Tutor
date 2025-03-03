import React, { useState } from "react";
import { ChevronRight, Download, Trash2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { useAgentEvents, useAgentRequests } from "@/hooks/use-agent-data";
import { useAllDownvotes, useAllUpvotes } from "@/hooks/use-all-downvotes";
import { UseQueryResult } from "@tanstack/react-query";
import { VoteAndMessage, VoteType } from "@/lib/supabase/vote";
import { Markdown } from "@/components/markdown";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  AgentRequest,
  AgentEvent,
  RagDocument,
  RagWebpage
} from "@/app/(admin)/actions";
import {
  cn,
  getRelativeTimeString,
  humanReadableSize,
  truncateUrl
} from "@/lib/utils";
import {
  useRagDocuments,
  useDownloadRagDoc,
  useDeleteRagDoc,
  useRagWebpages
} from "@/hooks/use-rag";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogOverlay,
  DialogPortal,
  DialogTitle,
} from "@/components/ui/dialog";

export function PagedAgentRequestsTable() {
  return (
    <PagedObjectTable
      title="Agent Requests"
      description="View incoming agent requests."
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
      title="Agent Events"
      description="View agent events."
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
      title="Downvotes"
      description="View all downvotes."
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
      title="Upvotes"
      description="View all upvotes."
      dataHook={useAllUpvotes}
      idField="id"
      TableHeadings={UpvotesTableHeadings}
      RowContents={VotesTableContents}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}

export function PagedRagDocumentsTable() {
  return (
    <PagedObjectTable
      title="Document Library"
      description="View all documents uploaded to the RAG library."
      dataHook={useRagDocuments}
      idField="id"
      TableHeadings={RagDocumentsTableHeadings}
      RowContents={RagDocumentsTableContents}
      ExpandedRowContents={RagDocumentsExpandedTableContents}
    />
  )
}

export function PagedRagWebpagesTable() {
  return (
    <PagedObjectTable
      title="Webpage Library"
      description="View all webpages uploaded to the RAG library."
      dataHook={useRagWebpages}
      idField="id"
      TableHeadings={RagWebpagesTableHeadings}
      RowContents={RagWebpagesTableContents}
      ExpandedRowContents={RagWebpagesExpandedTableContents}
    />
  )
}

type props<T> = {
  title: string;
  description: string;
  dataHook: (page: number) => UseQueryResult<{ data: T[]; totalPages: number; }, Error>;
  idField: string;
  TableHeadings: () => React.ReactNode;
  RowContents: ({data}: {data: T}) => React.ReactNode;
  ExpandedRowContents: ({data, index}: {data: T, index: number}) => React.ReactNode;
}

function PagedObjectTable<T>({
  title,
  description,
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
                                expandedRows.has(data[idField]!) ? "rotate-90" : ""
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
      <div>
        <h4 className="font-semibold text-primary">Full Query:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.metadata.query}</Markdown>
        </div>
      </div>
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

export function AgentEventsExpandedTableContents({data, index}: {data: AgentEvent, index: number}) {
  return (
    <>
      {data.metadata.answer && (
        <div>
          <h4 className="font-semibold text-primary">Answer:</h4>
          <div className="text-foreground mt-1">
            <Markdown>{data.metadata.answer}</Markdown>
          </div>
        </div>
      )}
      {data.metadata.reason && (
        <div>
          <h4 className="font-semibold text-primary">Reason:</h4>
          <div className="text-foreground mt-1">
            <Markdown>{data.metadata.reason}</Markdown>
          </div>
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

function RagDocumentsTableHeadings() {
  return (
    <>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        Document Name
      </th>
      <th className="h-12 w-[150px] px-4 text-left align-middle font-semibold text-primary">
        Size
      </th>
      <th className="h-12 w-[180px] px-4 text-left align-middle font-semibold text-primary">
        Upload Date
      </th>
      <th className="h-12 w-[100px] px-4 text-right align-middle font-semibold text-primary last:rounded-tr-xl">
        Actions
      </th>
    </>
  )
}

function RagDocumentsTableContents({data}: {data: RagDocument}) {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const downloadMutation = useDownloadRagDoc();
  const deleteMutation = useDeleteRagDoc();

  const handleDownload = async (doc: RagDocument) => {
    try {
      toast.promise(downloadMutation.mutateAsync(doc.name), {
        loading: `Downloading ${doc.name}...`,
        success: () => {
          return `Downloaded ${doc.name}`;
        },
        error: () => {
          return "Failed to download file";
        }
      });
    } catch (error) {
      console.error("Download error:", error);
    }
  };

  const handleDelete = async(doc: RagDocument) => {
    try {
      setIsDialogOpen(false);
      toast.promise(deleteMutation.mutateAsync(doc.name), {
        loading: `Deleting ${doc.name}...`,
        success: () => {
          return `Deleted ${doc.name}`;
        },
        error: () => {
          return "Failed to delete file";
        }
      });
    } catch (error) {
      console.error("Delete error:", error);
    }
  };

  return (
    <>
      <td className="p-4 align-middle">
        <div className="flex items-center">
          <span>{data.name}</span>
        </div>
      </td>
      <td className="p-4 align-middle">
        {data.document_size && humanReadableSize(data.document_size)}
      </td>
      <td className="p-4 align-middle">
        {data.date_uploaded && getRelativeTimeString(new Date(data.date_uploaded))}
      </td>
      <td className="p-4 align-middle text-right">
        <div className="flex justify-end gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => {
              e.stopPropagation();
              handleDownload(data);
            }}
            className="size-8 hover:text-primary"
          >
            <Download className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => {
              e.stopPropagation();
              setIsDialogOpen(true);
            }}
            className="size-8 text-destructive/80 hover:bg-destructive/30 hover:text-destructive"
          >
            <Trash2 className="size-4" />
          </Button>
        </div>
      </td>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogTitle className="text-destructive">Delete Document</DialogTitle>
            <DialogDescription className="text-foreground">
              Are you sure you want to delete &quot;{data.name}&quot;? This action cannot be undone.
            </DialogDescription>
            <div className="flex justify-end gap-2">
              <Button
                className="bg-muted hover:bg-muted/50"
                variant="outline"
                onClick={(e) => {
                  e.stopPropagation();
                  setIsDialogOpen(false)
                }}
              >
                Cancel
              </Button>
              <Button
                className="bg-destructive hover:bg-destructive/50"
                variant="destructive"
                onClick={(e) => {
                  e.stopPropagation();
                  handleDelete(data);
                }}
              >
                Delete
              </Button>
            </div>
          </DialogContent>
        </DialogPortal>
      </Dialog>
    </>
  )
}

function RagDocumentsExpandedTableContents({data}: {data: RagDocument}) {
  return (
    <>
      <div>
        <h4 className="font-semibold text-primary">Document Type:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.document_type}</Markdown>
        </div>
      </div>
      <div>
        <h4 className="font-semibold text-primary">URL:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.url}</Markdown>
        </div>
      </div>
    </>
  )
}

function RagWebpagesTableHeadings() {
  return (
    <>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        URL
      </th>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
        Upload Date
      </th>
    </>
  )
}

function RagWebpagesTableContents({data}: {data: RagWebpage}) {
  return (
    <>
      <td className="p-4 align-middle">
        <div className="flex items-center">
          <span className="font-mono">{truncateUrl(data.url)}</span>
        </div>
      </td>
      <td className="p-4 align-middle">
        {data.date_uploaded && getRelativeTimeString(new Date(data.date_uploaded))}
      </td>
    </>
  )
}

function RagWebpagesExpandedTableContents({data}: {data: RagWebpage}) {
  return (
    <>
      <div>
        <h4 className="font-semibold text-primary">Full URL:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.url}</Markdown>
        </div>
      </div>
    </>
  )
}