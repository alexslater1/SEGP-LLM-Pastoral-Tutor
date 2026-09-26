"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { UserContextProvider } from "@/providers/user-provider";
import { User } from "@/lib/supabase/user";

export function Providers({ children, user, subordinates }: { children: React.ReactNode, user: User | null, subordinates: User[] | null }) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <UserContextProvider user={user} subordinates={subordinates}>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </UserContextProvider>
  );
}
