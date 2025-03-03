'use client';

import { User } from "@/lib/supabase/user";
import { useContext } from "react";
import { createContext } from "react";

export type UserContextType = {
  user: User | null;
  subordinates: User[] | null;
}

const UserContext = createContext<UserContextType>({
  user: null,
  subordinates: null,
});

export function UserContextProvider({children, user, subordinates}: {children: React.ReactNode; user: User | null; subordinates: User[] | null}) {

  return (
    <UserContext.Provider value={{user, subordinates}}>
      {children}
    </UserContext.Provider>
  );
}

export function useUser () {
  return useContext(UserContext);
};
