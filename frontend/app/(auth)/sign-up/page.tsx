"use client";

import { Login } from "@/components/login";
import { redirect } from "next/navigation";
import { useUser } from "@/providers/user-provider";

export default function SignUpPage() {
  const user = useUser();

  if (user) {
    return redirect("/");
  }

  return <Login mode="signup" />;
}