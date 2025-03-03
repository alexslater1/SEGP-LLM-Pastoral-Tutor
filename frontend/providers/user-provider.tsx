'use client';

import { User } from "@/lib/supabase/user";
import { useContext } from "react";
import { createContext } from "react";

const UserContext = createContext<User | null>(null);

export function UserContextProvider({children, user}: {children: React.ReactNode; user: User | null}) {

  return (
    <UserContext.Provider value={user}>
      {children}
    </UserContext.Provider>
  );
}

export const useUser = () => {
  const user = useContext(UserContext);
  return user;
};
