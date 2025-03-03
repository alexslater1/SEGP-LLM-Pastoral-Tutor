"use client";

import { redirect } from "next/navigation";
import { Login } from "@/components/login";
import { useUser } from "@/providers/user-provider";

export default function SignInPage() {
  const user = useUser();
  if (user) {
    return redirect("/");
  }

  return <Login mode="signin" />;
}