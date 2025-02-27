"use server";

import { createClient } from "@/lib/supabase/server";
import { User } from "@/lib/supabase/user";

export async function downvote(messageId: string, reason: string, user: User) {
  const supabase = await createClient();

  console.log({messageId, reason, user});
  
  const { error } = await supabase
    .from('downvoted_responses')
    .insert({
      message_id: messageId,
      reason: reason,
      user_id: user.id
    })

  if (error) {
    console.error('Error downvoting message:', error);
    throw error;
  }
}

export async function removeDownvote(messageId: string, user: User) {
  const supabase = await createClient();

  const { error } = await supabase
    .from('downvoted_responses')
    .delete()
    .eq('message_id', messageId)
    .eq('user_id', user.id);

  if (error) {
    console.error('Error removing downvote:', error);
    throw error;
  }
}
