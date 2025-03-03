"use server";

import { createClient } from "./server";
import { UserRoleEnum } from "@/app/(admin)/role-authorization";

export type User = {
  id: string;
  role: UserRoleEnum;
  firstName: string;
  lastName: string;
  email: string;
  superior: string;
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

  return await getUserData(user.id);
};

export async function getUserData(id: string): Promise<User> {
  const supabase = await createClient();
  const { data, error } = await supabase.from("user_data").select("*").eq("id", id).single()

  if (error) {
    console.error("Error fetching user_data:", error.message);
    throw new Error("Error fetching user_data: " + error.message);
  }

  return {
    id: data.id,
    role: data.role as UserRoleEnum,
    firstName: data.first_name,
    lastName: data.last_name,
    email: data.email,
    superior: data.superior,
  };
};

export async function getUserSubordinates(user: User): Promise<User[]> {
  const supabase = await createClient();
  const { data, error } = await supabase.from("user_data").select("*").eq("superior", user.id);

  if (error) {
    console.error("Error fetching user_data:", error.message);
    throw new Error("Error fetching user_data: " + error.message);
  }

  return data.map((subordinate) => ({
    id: subordinate.id,
    role: subordinate.role as UserRoleEnum,
    firstName: subordinate.first_name,
    lastName: subordinate.last_name,
    email: subordinate.email,
    superior: subordinate.superior,
  }));
};

