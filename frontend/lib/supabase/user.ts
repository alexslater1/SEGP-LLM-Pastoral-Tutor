"use server";

import { createClient } from "./server";

export type User = {
  id: string;
  role: string;
  name: string;
  email: string;
};

export async function getUser(): Promise<User | null> {
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

export async function getUserData(id: string): Promise<User | null> {
  const supabase = await createClient();
  const { data, error } = await supabase.from("user_data").select("*").eq("id", id).single()

  if (error) {
    console.error("Error fetching user_data:", error.message);
    return null;
  }

  return {
    id: data.id,
    role: data.role,
    name: data.name,
    email: data.email,
  };
};

