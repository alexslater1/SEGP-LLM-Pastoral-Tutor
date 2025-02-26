"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { signIn, signUp } from "../app/(auth)/actions";
import { useState } from "react";
import { getUserSession } from "@/lib/supabase/client";
import { Card, CardHeader, CardContent, CardFooter, CardTitle, CardDescription } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { useSearchParams } from "next/navigation";

export function Login({ mode = "signin" }: { mode?: "signin" | "signup" }) {
  const [error, setError] = useState("");
  const searchParams = useSearchParams().toString();

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
    <div className="min-h-screen flex items-center justify-center p-4 bg-background">
      <Card className="w-full max-w-md bg-secondary">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold tracking-tight text-center text-primary">
            {mode === "signin" ? "Welcome back" : "Create account"}
          </CardTitle>
          <CardDescription className="text-center text-foreground">
            {mode === "signin"
              ? "Sign in to continue to your account"
              : "Enter your details to create your account"}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email" className="text-primary">Email</Label>
              <Input
                id="email"
                name="email"
                type="email"
                placeholder="Enter your email"
                required
                className="w-full bg-muted"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password" className="text-primary">Password</Label>
              <Input
                id="password"
                name="password"
                type="password"
                placeholder="Enter your password"
                required
                className="w-full bg-muted"
              />
            </div>
            <input type="hidden" name="searchParams" value={searchParams} />

            {error && (
              <div className="text-sm text-destructive text-center">
                {error}
              </div>
            )}

            <Button
              formAction={mode === "signin" ? formSignIn : formSignUp}
              type="submit"
              className="w-full"
            >
              {mode === "signin" ? "Sign In" : "Create Account"}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="flex flex-col space-y-4">
          <div className="text-sm text-foreground text-center">
            {mode === "signin" ? "New to our platform? " : "Already have an account? "}
            <Link
              href={searchParams ? mode === "signin" ? `/sign-in?${searchParams}` : `/sign-up?${searchParams}` : 
                                   mode === "signin" ? "/sign-up" : "/sign-in"}
              className="font-medium text-primary hover:underline"
            >
              {mode === "signin" ? "Create an account" : "Sign in"}
            </Link>
          </div>
        </CardFooter>
      </Card>
    </div>
  );
}