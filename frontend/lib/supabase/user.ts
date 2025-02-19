"use server";

import { createClient } from "./server";

export type User = {
  id: string;
  role: string;
  name: string;
  email: string;
};

export const getUser = async (): Promise<User | null> => {
  const supabase = await createClient();
  const {
    data: { user },
    error,
  } = await supabase.auth.getUser();

  if (error || !user) {
    return null;
  }

  return await getUserData(user?.id);
};

async function getUserData(id: string): Promise<User | null> {
  const supabase = await createClient();
  const { data, error } = await supabase.from("user_data").select("*").eq("id", id);

  if (error) {
    console.error("Error fetching user_data:", error.message);
    return null;
  }

  return {
    id: data[0].id,
    role: data[0].role,
    name: data[0].name,
    email: data[0].email,
  };
};


