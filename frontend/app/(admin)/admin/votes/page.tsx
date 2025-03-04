"use client";

import { useAllDownvotes, useAllUpvotes, useAllVotes } from "@/hooks/use-all-votes";
import { 
  PagedObjectTable, 
  VotesTableContents, 
  VotesExpandedTableContents, 
  VoteHeadingType,
  VotesTableHeadings 
} from "@/components/paged-object-table";
import { PageSelectionBarMultiple } from "@/components/page-selection-bar";
import { useState } from "react";
import { VoteAndMessage } from "@/lib/supabase/vote";

export default function VotesPage() {
  const [pages, setPages] = useState([0]);

  return (
    <>
      <PageSelectionBarMultiple pageNumbers={pages} setPages={setPages} pageNames={["Downvotes", "Upvotes"]} allowNoneSelected={false} />
      {pages.includes(0) && pages.includes(1) ? <PagedVotesTable /> : pages.includes(0) ? <PagedDownvotesTable /> : <PagedUpvotesTable />}
    </>
  )
}

function UpAndDownvotesTableContentsInstance({data}: {data: VoteAndMessage}) {
  return (
    <VotesTableContents includeVoteType={false} data={data} />
  )
}

function BothVotesTableContentsInstance({data}: {data: VoteAndMessage}) {
  return (
    <VotesTableContents includeVoteType={true} data={data} />
  )
}

function DownvotesTableHeadings() {
  return (
    <VotesTableHeadings headingType={VoteHeadingType.DOWNVOTE} />
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
      RowContents={UpAndDownvotesTableContentsInstance}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}

function UpvotesTableHeadings() {
  return (
    <VotesTableHeadings headingType={VoteHeadingType.UPVOTE} />
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
      RowContents={UpAndDownvotesTableContentsInstance}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}

function VotesTableHeadingsInstance() {
  return (
    <VotesTableHeadings headingType={VoteHeadingType.BOTH} />
  )
}

function PagedVotesTable() {
  return (
    <PagedObjectTable
      title="Votes"
      description="View all votes."
      dataHook={useAllVotes}
      idField="id"
      TableHeadings={VotesTableHeadingsInstance}
      RowContents={BothVotesTableContentsInstance}
      ExpandedRowContents={VotesExpandedTableContents}
    />
  )
}
