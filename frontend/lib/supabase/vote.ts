"use server";

import { createClient } from "@/lib/supabase/server";

export async function downvote(messageId: string) {
  const supabase = await createClient();
  
}
