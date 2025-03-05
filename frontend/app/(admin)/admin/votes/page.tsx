"use client";

import { useAllDownvotes, useAllUpvotes, useAllVotes } from "@/hooks/use-all-votes";
import { 
  PagedObjectTable, 
  VotesTableContents, 
  VotesExpandedTableContents, 
  VoteHeadingType,
  VotesTableHeadings 
} from "@/components/paged-object-table";
import { useState } from "react";
import { VoteAndMessage } from "@/lib/supabase/vote";
import { FilterBar } from "@/components/filter-bar";
import { ThumbsDown, ThumbsUp, Vote } from "lucide-react";

type VoteFilter = "downvotes" | "upvotes" | "all";

const voteFilterOptions = [
  {
    id: "upvotes",
    label: "Upvotes",
    icon: ThumbsUp,
    theme: "success" as const
  },
  {
    id: "downvotes",
    label: "Downvotes",
    icon: ThumbsDown,
    theme: "destructive" as const
  },
  {
    id: "all",
    label: "All Votes",
    icon: Vote,
    theme: "default" as const
  }
];

export default function VotesPage() {
  const [filter, setFilter] = useState<VoteFilter>("downvotes");

  return (
    <>
      <FilterBar
        filter={filter}
        setFilter={setFilter}
        options={voteFilterOptions}
      />
      {filter === "all" ? (
        <PagedVotesTable />
      ) : filter === "downvotes" ? (
        <PagedDownvotesTable />
      ) : (
        <PagedUpvotesTable />
      )}
    </>
  );
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
