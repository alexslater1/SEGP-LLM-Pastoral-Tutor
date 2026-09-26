"use client";

import {
  useAllDownvotes,
  useAllUpvotes,
  useAllVotes,
} from "@/hooks/use-all-votes";
import {
  PagedObjectTable,
  VotesTableContents,
  VotesExpandedTableContents,
  VoteHeadingType,
  VotesTableHeadings,
} from "@/components/paged-object-table";
import { useState } from "react";
import { VoteAndMessage } from "@/lib/supabase/vote";
import { FilterBar } from "@/components/filter-bar";
import { ThumbsDown, ThumbsUp, Vote } from "lucide-react";
import { SwitchCard } from "@/components/switch-card";
import { useAutoUpdatePrompts } from "@/hooks/use-auto-update-prompts";
import { Separator } from "@/components/ui/separator";

type VoteFilter = "downvotes" | "upvotes" | "all";

const voteFilterOptions = [
  {
    id: "upvotes",
    label: "Likes",
    icon: ThumbsUp,
    theme: "success" as const,
  },
  {
    id: "downvotes",
    label: "Dislikes",
    icon: ThumbsDown,
    theme: "destructive" as const,
  },
  {
    id: "all",
    label: "All Feedback",
    icon: Vote,
    theme: "default" as const,
  },
];

export default function VotesPage() {
  const [filter, setFilter] = useState<VoteFilter>("all");
  const { enabled, setEnabled, loading, error } = useAutoUpdatePrompts();

  return (
    <>
      <div className="absolute right-[24.5rem] top-20 p-1 shadow-sm">
        <div className="flex flex-row items-center gap-2">
          <SwitchCard 
            titleText="Toggle Auto-Feedback Processing" 
            descriptionText={"When enabled, " +
              "all feedback will be automatically processed in regular internvals and will be used to adjust the agent prompts."} 
            enabled={enabled} 
            setEnabled={setEnabled} 
            loading={loading} 
            error={error} 
          />
          <div className="relative h-[30px] px-5">
            <Separator orientation="vertical" />
          </div>
        </div>
      </div>
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

function UpAndDownvotesTableContentsInstance({
  data,
}: {
  data: VoteAndMessage;
}) {
  return <VotesTableContents includeVoteType={false} data={data} />;
}

function BothVotesTableContentsInstance({ data }: { data: VoteAndMessage }) {
  return <VotesTableContents includeVoteType={true} data={data} />;
}

function DownvotesTableHeadings() {
  return <VotesTableHeadings headingType={VoteHeadingType.DOWNVOTE} />;
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
  );
}

function UpvotesTableHeadings() {
  return <VotesTableHeadings headingType={VoteHeadingType.UPVOTE} />;
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
  );
}

function VotesTableHeadingsInstance() {
  return <VotesTableHeadings headingType={VoteHeadingType.BOTH} />;
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
  );
}
