"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { signIn, signUp } from "../app/(auth)/actions";
import { useState } from "react";
import { getUserSession } from "@/lib/supabase/client";

export function Login({ mode = "signin" }: { mode?: "signin" | "signup" }) {
  const [error, setError] = useState("");

  const formSignIn = async (formData: FormData) => {
    const { error } = await signIn(formData);
    if (error) {
      setError(error);
    }

    const login_session = await getUserSession();
    if (login_session) {
      console.log("login_session: ", login_session.access_token);
    }
  }

  const formSignUp = async (formData: FormData) => {
    const { error } = await signUp(formData);
    if (error) {
      setError(error);
    }
  }

  return (
    <div className="min-h-dvh bg-gradient-to-b from-white to-gray-50 flex items-center justify-center px-4 py-16 sm:px-6 lg:px-8">
      <div className="w-full max-w-md">
        <h1 className="mt-10 text-2xl font-semibold tracking-tight text-center text-gray-900">
          {mode === "signin" ? "Welcome back" : "Create your account"}
        </h1>
        <p className="mt-2 text-sm text-center text-gray-600">
          {mode === "signin"
            ? "Sign in to continue to your account"
            : "Get started with your new account"}
        </p>

        <div className="mt-10">
          <div className="space-y-6">
            <form className="space-y-4">
              <Input
                name="email"
                type="email"
                placeholder="name@example.com"
                required
                className="px-4 h-12 bg-white text-zinc-950 rounded-lg border-gray-200 shadow-sm transition-colors focus:border-blue-500 focus:ring-blue-500"
              />
              <Input
                name="password"
                type="password"
                placeholder="password"
                required
                className="px-4 h-12 bg-white text-zinc-950 rounded-lg border-gray-200 shadow-sm transition-colors focus:border-blue-500 focus:ring-blue-500"
              />

              <Button
                formAction={mode === "signin" ? formSignIn : formSignUp}
                type="submit"
                className="w-full h-12 font-medium text-white bg-blue-600 rounded-lg transition-colors hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              >
                Continue with Email
              </Button>
            </form>
          </div>

          {error && (
            <div className="mt-4 text-sm text-red-600 text-center">
              {error}
            </div>
          )}

          <p className="mt-8 text-sm text-center text-gray-600">
            {mode === "signin"
              ? "New to our platform? "
              : "Already have an account? "}
            <Link
              href={mode === "signin" ? "/sign-up" : "/sign-in"}
              className="font-medium text-blue-600 hover:text-blue-500"
            >
              {mode === "signin" ? "Create an account" : "Sign in"}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}