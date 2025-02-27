"use server";

import { createClient } from "@/lib/supabase/server";
import { User } from "@/lib/supabase/user";

export type Downvote = {
  id: string;
  created_at: string;
  request_id: string;
  user_id: string;
  reason: string;
}

export type DownvoteAndMessage = {
  downvoteID: string;
  createdAt: string;
  requestID: string;
  userEmail: string;
  agent: string;
  reason: string;
  query: string;
  answer: string;
}

export async function downvote(messageId: string, reason: string, user: User): Promise<Downvote> {
  const supabase = await createClient();

  const { data: existingDownvote, error: fetchError } = await supabase
    .from('downvoted_responses')
    .select('*')
    .eq('request_id', messageId)

  if (existingDownvote?.length && existingDownvote.length > 0) {
    throw new Error('Message already downvoted');
  } else if (fetchError) {
    console.error('Error fetching downvoted message:', fetchError);
    throw fetchError;
  }

  const { data, error } = await supabase
    .from('downvoted_responses')
    .insert({
      request_id: messageId,
      reason: reason,
      user_id: user.id
    })

  if (error) {
    console.error('Error downvoting message:', error);
    throw error;
  } else {
    const { data: downvoteData, error: downvoteError } = await supabase
      .from('downvoted_responses')
      .select('*')
      .eq('request_id', messageId)
      .single();

    if (downvoteError) {
      console.error('Error fetching downvoted message:', downvoteError);
      throw downvoteError;
    }

    return downvoteData;
  }
}

export async function removeDownvote(messageId: string): Promise<Downvote> {
  const supabase = await createClient();

  const { data: existingDownvote, error: fetchError } = await supabase
    .from('downvoted_responses')
    .select('*')
    .eq('request_id', messageId)
    .single();

  if (!existingDownvote) {
    throw new Error('Message not downvoted');
  }

  const { data, error } = await supabase
    .from('downvoted_responses')
    .delete()
    .eq('request_id', messageId)

  if (error) {
    console.error('Error removing downvote:', error);
    throw error;
  }

  return existingDownvote;
}

export async function getDownvotesByChatID(chatId: string): Promise<Downvote[]> {
  const supabase = await createClient();

  const { data: sessionData, error: sessionError } = await supabase
    .from('request_sessions')
    .select('request_id')
    .eq('session_id', chatId);

  if (sessionError) {
    console.error('Error fetching session data:', sessionError);
    throw sessionError;
  }

  const { data, error } = await supabase
    .from('downvoted_responses')
    .select('*')
    .in('request_id', sessionData.map(row => row.request_id));

  if (error) {
    console.error('Error fetching downvotes:', error);
    throw error;
  }

  console.log(data);

  return data;
}

export async function getAllDownvotes(): Promise<DownvoteAndMessage[]> {
  const supabase = await createClient();
  let downvotesArray: DownvoteAndMessage[] = [];

  const { data: downvotes, error: downvotesError } = await supabase
    .from('downvoted_responses')
    .select(`
      *,
      agent_requests (
        metadata
      )
    `);

  if (downvotesError) {
    console.error('Error fetching downvotes:', downvotesError);
    throw downvotesError;
  }

  for (const downvote of downvotes) {
    const { data: events, error: eventsError } = await supabase
      .from('agent_events')
      .select('type, metadata')
      .eq('request_id', downvote.request_id)

    if (eventsError) {
      console.error('Error fetching events:', eventsError);
      throw eventsError;
    }

    const answerEvent = events.find(event => event.type === 'answer_success');

    const { data: user, error: userError } = await supabase
      .from('user_data')
      .select('email')
      .eq('id', downvote.user_id)
      .single();
      
    if (userError) {
      console.error('Error fetching user:', userError);
      throw userError;
    }

    downvotesArray.push({
      downvoteID: downvote.id,
      createdAt: downvote.created_at,
      requestID: downvote.request_id,
      userEmail: user.email,
      agent: (answerEvent?.metadata).agentID,
      reason: downvote.reason,
      query: (downvote.agent_requests.metadata).query,
      answer: (answerEvent?.metadata).answer
    });
  }

  return downvotesArray;
}