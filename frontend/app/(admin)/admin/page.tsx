"use client";

import { useUser } from "@/providers/user-provider";
import {
  FileText,
  Globe,
  MessageSquare,
  Activity,
  Users,
  ThumbsUp,
} from "lucide-react";
import {
  useRagDocumentsCount,
  useRagWebpagesCount,
  useAgentEventsCount,
  useAgentRequestsCount,
  useUpvotesCount,
  useDownvotesCount,
} from "@/hooks/use-rag";
import { UserRoleEnum } from "../role-authorization";

export default function AdminPage() {
  const userContext = useUser();
  const isAdmin = userContext.user?.role === UserRoleEnum.ADMIN;

  const { data: documents, isLoading: isDocumentsLoading } =
    useRagDocumentsCount({ enabled: isAdmin });
  const { data: webpages, isLoading: isWebpagesLoading } = useRagWebpagesCount({
    enabled: isAdmin,
  });
  const { data: agentEvents, isLoading: isEventsLoading } = useAgentEventsCount(
    { enabled: isAdmin }
  );
  // const { data: agentRequests, isLoading: isRequestsLoading } = useAgentRequestsCount({ enabled: isAdmin });
  const { data: upvotes, isLoading: isUpvotesLoading } = useUpvotesCount({
    enabled: isAdmin,
  });
  const { data: downvotes, isLoading: isDownvotesLoading } = useDownvotesCount({
    enabled: isAdmin,
  });

  const totalUpvotes = upvotes ? upvotes : 0;
  const totalDownvotes = downvotes ? downvotes : 0;
  const ratio =
    totalDownvotes === 0
      ? totalUpvotes
      : (totalUpvotes / totalDownvotes).toFixed(2);

  const allStats = [
    {
      name: `Active ${isAdmin ? "Users" : "Tutees"}`,
      value: userContext.subordinates?.length || 0,
      icon: Users,
      description: `Total number of ${
        isAdmin ? "users" : "tutees"
      } being managed`,
      isLoading: false,
    },
    {
      name: "RAG Documents Uploaded",
      value: documents,
      icon: FileText,
      description: "Documents in the RAG library",
      isLoading: isDocumentsLoading,
    },
    {
      name: "Webpages Scraped",
      value: webpages,
      icon: Globe,
      description: "Total webpages in the RAG library",
      isLoading: isWebpagesLoading,
    },
    {
      name: "Agent Events",
      value: agentEvents,
      icon: Activity,
      description: "Total agent events recorded",
      isLoading: isEventsLoading,
    },
    // {
    //   name: "Agent Requests",
    //   value: agentRequests,
    //   icon: MessageSquare,
    //   description: "Total agent requests processed",
    //   isLoading: isRequestsLoading
    // },
    {
      name: "Upvote/Downvote Ratio",
      value: `${ratio}:1`,
      icon: ThumbsUp,
      description: `${totalUpvotes} upvotes to ${totalDownvotes} downvotes`,
      isLoading: isUpvotesLoading || isDownvotesLoading,
    },
  ];

  const stats = isAdmin ? allStats : allStats.slice(0, 1); // Only show Active Tutees for tutors

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold text-primary tracking-tight">
          {isAdmin ? "Admin Dashboard" : "Tutor Dashboard"}
        </h2>
        <p className="text-foreground">
          Welcome to the {isAdmin ? "admin" : "tutor"} dashboard.
          {isAdmin ? "View key metrics and system status below." : ""}
        </p>
      </div>
      <div className="border-t">
        <div className="bg-background">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 p-6">
            {stats.map((stat) => (
              <div
                key={stat.name}
                className="rounded-xl border bg-muted p-6 hover:border-primary transition-all duration-300"
              >
                <div className="flex items-center gap-2">
                  <stat.icon className="size-5 text-muted-foreground" />
                  <h3 className="font-semibold text-primary">{stat.name}</h3>
                </div>
                <div className="mt-4 space-y-2">
                  {stat.isLoading ? (
                    <p className="text-2xl font-bold text-muted-foreground animate-pulse">
                      Loading...
                    </p>
                  ) : (
                    <p className="text-2xl font-bold text-foreground">
                      {stat.value}
                    </p>
                  )}
                  <p className="text-sm text-muted-foreground">
                    {stat.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
