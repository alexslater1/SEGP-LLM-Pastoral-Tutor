"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";
import { UserContextProvider } from "@/providers/client-provider";
import { User } from "@/lib/supabase/user";

export function Providers({ children, user }: { children: React.ReactNode, user: User | null }) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <UserContextProvider user={user}>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </UserContextProvider>
  );
}
