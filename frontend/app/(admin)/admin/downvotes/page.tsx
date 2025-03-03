"use client";

import { useAllDownvotes } from "@/hooks/use-all-votes";
import { PagedObjectTable, VotesTableContents, VotesExpandedTableContents } from "@/components/paged-object-table";

export default function DownvotesPage() {
  return (
    <PagedDownvotesTable />
  )
}

function PagedDownvotesTable() {
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