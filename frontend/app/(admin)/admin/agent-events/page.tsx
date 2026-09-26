"use client";

import { useAgentEvents } from "@/hooks/use-agent-data";
import { PagedObjectTable, formatDate } from "@/components/paged-object-table";
import { AgentEvent } from "../../actions";
import { cn } from "@/lib/utils";
import { Markdown } from "@/components/markdown";

export default function AgentEventsPage() {
  return <PagedAgentEventsTable />
}

function PagedAgentEventsTable() {
  return (
    <PagedObjectTable
      title="Agent Events"
      description="View all agent events."
      dataHook={useAgentEvents}
      idField="id"
      TableHeadings={AgentEventsTableHeadings}
      RowContents={AgentEventsTableContents}
      ExpandedRowContents={AgentEventsExpandedTableContents}
    />
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
        {data.created_at && formatDate(data.created_at)}
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
      {data.metadata.query && (
        <div>
          <h4 className="font-semibold text-primary">Query:</h4>
          <p className="text-foreground mt-1">
            {data.metadata.query}
          </p>
        </div>
      )}
      {data.metadata.agentID && (
        <div>
          <h4 className="font-semibold text-primary">Agent ID:</h4>
          <p className="text-foreground mt-1">
            {data.metadata.agentID}
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



