"use server";

import { z } from "zod";
import { redirect } from "next/navigation";
import { createClient } from "@/lib/supabase/server";

const signInSchema = z.object({
  email: z.string().email().min(3).max(255),
  password: z.string().min(8).max(100),
  searchParams: z.string().optional(),
});

export async function signIn(formData: FormData) {
  const supabase = await createClient();
  const validation = signInSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
    searchParams: formData.get("searchParams"),
  });
  if (!validation.success) {
    return { error: validation.error.message };
  }

  const { email, password, searchParams } = validation.data;

  const { data: signInData, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  });
  if (error) {
    return { error: "Invalid credentials. Please try again." };
  }

  const { data: userData, error: userDataError } = await supabase
    .from("user_data")
    .select("*")
    .eq("id", signInData.user?.id)
    .single();

  if (userDataError && userDataError.code === "PGRST116") {
    // No user_data entry found, create one
    const { error: insertError } = await supabase
      .from("user_data")
      .insert({ id: signInData.user?.id, role: "student", name: email.split("@")[0] });
    if (insertError) {
      return { error: insertError.message };
    }
  } else if (userDataError) {
    return { error: userDataError.message };
  }

  // If sign-in is successful, redirect to dashboard
  searchParams ? redirect(`/?${searchParams}`) : redirect("/");
  return { error: null };
};

const signUpSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
  searchParams: z.string().optional(),
});

export async function signUp(formData: FormData) {
  const supabase = await createClient();
  const validation = signUpSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
    searchParams: formData.get("searchParams"),
  });
  if (!validation.success) {
    return { error: validation.error.message };
  }

  const { email, password, searchParams } = validation.data;

  const { data: signUpData, error: signUpError } = await supabase.auth.signUp({
    email,
    password,
  });
  if (signUpError) {
    return { error: signUpError.message };
  }

  // Create a new user_data entry
  const { error: insertError } = await supabase
    .from("user_data")
    .insert({ id: signUpData?.user?.id, role: "student", name: email.split("@")[0] });

  if (insertError) {
    return { error: insertError.message };
  }

  searchParams ? redirect(`/?${searchParams}`) : redirect("/");
  return { error: null };
}

export const signOut = async () => {
  const supabase = await createClient();
  await supabase.auth.signOut();
  redirect("/");
};