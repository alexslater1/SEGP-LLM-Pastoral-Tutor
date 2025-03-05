"use client";

import { User } from "@/lib/supabase/user";
import { useUser } from "@/providers/user-provider";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";
import { UserRoleEnum } from "@/app/(admin)/role-authorization";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";

export default function UsersPage() {
  const userContext = useUser();
  const router = useRouter();

  if (!userContext.user) {
    return <div>Loading...</div>;
  }

  const capitalize = (name: string) => name.charAt(0).toUpperCase() + name.slice(1);

  const redirectToChatHistory = (user: User) => router.push(`/admin/users/${user.id}`)

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl text-primary font-bold tracking-tight">{userContext.user.role === UserRoleEnum.ADMIN ? "Users" : "Tutees"}</h2>
        <p className="text-foreground">{userContext.user.role === UserRoleEnum.ADMIN ? 
          "View all users. Click on a user to view their chat history." : 
          "View all tutees. Click on a tutee to view their chat history."}
        </p>
      </div>

      <div className="h-px bg-border" />

      <div>
        <div className="bg-background">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">Name</th>
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">Email</th>
                    {userContext.user.role === UserRoleEnum.ADMIN && (
                      <>
                        <th className="h-12 px-4 text-left align-middle font-semibold text-primary">Role</th>
                        <th className="h-12 w-24 px-4 text-left align-middle font-semibold text-primary">Superior</th>
                      </>
                    )}
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary"></th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {!!!userContext.subordinates ? (
                    <tr>
                      <td colSpan={5} className="p-0">
                        <div className="h-[569px] flex items-center justify-center bg-secondary/50 dark:bg-muted/90 font-bold text-5xl text-primary">
                          <div className="animate-pulse">
                            Loading...
                          </div>
                        </div>
                      </td>
                    </tr>
                  ) : (
                    userContext.subordinates.map((user: User, index) => (
                      <tr
                        key={user.id}
                        className={cn(
                          "transition-colors hover:bg-muted/50 cursor-pointer",
                          index % 2 === 0 
                            ? "bg-table-row-odd" 
                            : "bg-table-row-even"
                        )}
                        onClick={() => redirectToChatHistory(user)}
                      >
                        <td className="p-4 align-middle text-foreground">
                          {capitalize(user.firstName)} {capitalize(user.lastName || "")} {user.id === userContext.user?.id && "(You)"}
                        </td>
                        <td className="p-4 align-middle text-foreground">
                          {user.email}
                        </td>
                        {userContext.user?.role === UserRoleEnum.ADMIN && (
                          <>
                            <td className="p-4 align-middle text-foreground">
                              {user.role}
                            </td>
                            <td className="p-4 align-middle text-foreground">
                              {userContext.subordinates?.find((subordinate: User) => subordinate.id === user.superior)?.email || "N/A"}
                            </td>
                          </>
                        )}
                        <td className="p-4 align-middle text-right">
                          <div className="flex justify-end gap-2">
                            <Button
                              variant="ghost"
                              size="icon"
                              onClick={() => redirectToChatHistory(user)}
                              className="size-8 hover:text-primary"
                            >
                              <ArrowRight className="size-4" />
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
} 