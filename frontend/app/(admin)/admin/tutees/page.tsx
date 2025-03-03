"use client";

import { User } from "@/lib/supabase/user";
import { useUser } from "@/providers/user-provider";
import { cn } from "@/lib/utils";
import { useRouter } from "next/navigation";

export default function TuteesPage() {
  const userContext = useUser();
  const router = useRouter();

  if (!userContext.user) {
    return <div>Loading...</div>;
  }

  const capitalize = (name: string) => name.charAt(0).toUpperCase() + name.slice(1);

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl text-primary font-bold tracking-tight">Tutees</h2>
        <p className="text-foreground">View all tutees. Click on a tutee to view their chat history.</p>
      </div>
      <div>
        <div className="bg-background md:w-1/2 w-full">
          <div className="rounded-xl border bg-card">
            <div className="overflow-hidden rounded-xl">
              <table className="w-full">
                <thead>
                  <tr className="border-b bg-table-header">
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">Name</th>
                    <th className="h-12 px-4 text-left align-middle font-semibold text-primary">Email</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {!!!userContext.subordinates ? (
                    <tr>
                      <td colSpan={5} className="p-0">
                        <div className="h-[569px] flex items-center justify-center bg-secondary/50 dark:bg-muted/90 font-bold text-5xl text-primary">
                          Loading...
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
                        onClick={() => router.push(`/admin/tutees/${user.id}`)}
                      >
                        <td className="p-4 align-middle">
                          <span
                            className={`inline-flex items-center rounded-full px-2.5 py-0.5 font-medium`}
                          >
                            {capitalize(user.firstName)} {capitalize(user.lastName)}
                          </span>
                        </td>
                        <td className="p-4 align-middle">
                          <span
                            className={`inline-flex items-center rounded-full px-2.5 py-0.5 font-medium`}
                          >
                            {user.email}
                          </span>
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