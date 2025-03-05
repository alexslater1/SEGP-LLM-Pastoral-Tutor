import { cn } from "@/lib/utils";
import { Button } from "./ui/button";
import { ThumbsDown, ThumbsUp, Vote } from "lucide-react";

type VoteFilter = "downvotes" | "upvotes" | "all";

export function VoteFilterBar({filter, setFilter}: {
  filter: VoteFilter;
  setFilter: (filter: VoteFilter) => void;
}) {
  return (
    <div className="absolute right-6 top-20">
      <div className="relative flex items-center gap-2 bg-background rounded-lg p-1 shadow-sm">
        <Button
          variant="outline"
          className={cn(
            "flex items-center gap-2",
            filter === "upvotes" && "bg-success/10 border-success text-success hover:bg-success/20"
          )}
          onClick={() => setFilter("upvotes")}
        >
          <ThumbsUp size={16} />
          Upvotes
        </Button>
        <Button
          variant="outline"
          className={cn(
            "flex items-center gap-2",
            filter === "downvotes" && "bg-destructive/10 border-destructive text-destructive hover:bg-destructive/20"
          )}
          onClick={() => setFilter("downvotes")}
        >
          <ThumbsDown size={16} />
          Downvotes
        </Button>
        <Button
          variant="outline"
          className={cn(
            "flex items-center gap-2",
            filter === "all" && "bg-primary/10 border-primary text-primary hover:bg-primary/20"
          )}
          onClick={() => setFilter("all")}
        >
          <Vote size={16} />
          All Votes
        </Button>
      </div>
    </div>
  );
} 