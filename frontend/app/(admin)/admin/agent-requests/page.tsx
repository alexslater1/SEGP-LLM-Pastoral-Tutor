"use client";

import { PagedObjectTable } from "@/components/paged-object-table";
import { useAgentRequests } from "@/hooks/use-agent-data";
import { AgentRequest } from "../../actions";
import { getRelativeTimeString, cn } from "@/lib/utils";
import { Markdown } from "@/components/markdown";

export default function AgentRequestsPage() {
  return <PagedAgentRequestsTable />
}

function PagedAgentRequestsTable() {
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
        <div className={cn(
          "truncate max-w-[500px]",
          !data.metadata.query && "text-gray-500"
        )}>
          {data.metadata.query || "N/A"}
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
      {data.metadata.query && (
        <div>
          <h4 className="font-semibold text-primary">Full Query:</h4>
          <div className="text-foreground mt-1">
            <Markdown>{data.metadata.query}</Markdown>
          </div>
        </div>
      )}
    </>
  )
}
