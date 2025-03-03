"use client";

import { PagedObjectTable, VotesExpandedTableContents, VotesTableContents } from "@/components/paged-object-table";
import { useAllUpvotes } from "@/hooks/use-all-votes";

export default function UpvotesPage() {
  return (
    <PagedUpvotesTable />
  )
}

function PagedUpvotesTable() {
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