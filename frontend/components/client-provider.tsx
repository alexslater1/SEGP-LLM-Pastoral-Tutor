'use client';

import { UserContext } from "@/lib/userContext";
import { User } from "@/lib/supabase/user";

export function UserContextProvider({
  children,
  user
}: {
  children: React.ReactNode;
  user: User | null;
}) {
  return (
    <UserContext.Provider value={user}>
      {children}
    </UserContext.Provider>
  );
}
