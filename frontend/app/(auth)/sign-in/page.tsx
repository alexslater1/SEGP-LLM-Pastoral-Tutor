import { redirect } from "next/navigation";
import { Login } from "@/components/login";
import { getUser } from "@/lib/supabase/user";

export default async function SignInPage() {
  const user = await getUser();
  if (user) {
    return redirect("/");
  }

  return <Login mode="signin" />;
}