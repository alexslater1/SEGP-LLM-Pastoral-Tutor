"use server";

import { createClient } from "@/lib/supabase/server";
import { User } from "@/lib/supabase/user";

export type VoteType = 'downvote' | 'upvote';

export type Vote = {
  type: VoteType;
  id: string;
  created_at: string;
  request_id: string;
  user_id: string;
  reason: string;
}

export type VoteAndMessage = {
  type: VoteType;
  id: string;
  createdAt: string;
  requestID: string;
  userEmail: string;
  agent: string;
  reason: string;
  query: string;
  answer: string;
}

export async function vote(type: VoteType, messageId: string, reason: string, user: User): Promise<Vote> {
  const supabase = await createClient();

  const { data: existingVote, error: fetchError } = await supabase
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .select('*')
    .eq('request_id', messageId)

  if (existingVote?.length && existingVote.length > 0) {
    throw new Error('Message already ' + (type === 'downvote' ? 'downvoted' : 'upvoted'));
  } else if (fetchError) {
    console.error('Error fetching ' + (type === 'downvote' ? 'downvoted' : 'upvoted') + ' message:', fetchError);
    throw fetchError;
  }

  const { data, error } = await supabase
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .insert({
      request_id: messageId,
      reason: reason,
      user_id: user.id
    })

  if (error) {
    console.error('Error ' + (type === 'downvote' ? 'downvoting' : 'upvoting') + ' message:', error);
    throw error;
  } else {
    const { data: voteData, error: voteError } = await supabase
      .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
      .select('*')
      .eq('request_id', messageId)
      .single();

    if (voteError) {
      console.error('Error fetching ' + (type === 'downvote' ? 'downvoted' : 'upvoted') + ' message:', voteError);
      throw voteError;
    }

    return { ...voteData, type: type };
  }
}

export async function removeVote(type: VoteType, messageId: string): Promise<Vote> {
  const supabase = await createClient();

  const { data: existingVote, error: fetchError } = await supabase
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .select('*')
    .eq('request_id', messageId)
    .single();

  if (!existingVote) {
    throw new Error('Message not ' + (type === 'downvote' ? 'downvoted' : 'upvoted'));
  }

  const { data, error } = await supabase
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .delete()
    .eq('request_id', messageId)

  if (error) {
    console.error('Error removing ' + (type === 'downvote' ? 'downvote' : 'upvote') + ':', error);
    throw error;
  }

  return { ...existingVote, type: type };
}

export async function getVotesByChatID(type: VoteType, chatId: string): Promise<Vote[]> {
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
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .select('*')
    .in('request_id', sessionData.map(row => row.request_id));

  if (error) {
    console.error('Error fetching ' + (type === 'downvote' ? 'downvotes' : 'upvotes') + ':', error);
    throw error;
  }

  return data.map((vote) => ({ ...vote, type: type }));
}

export async function getAllVotes(type: VoteType): Promise<VoteAndMessage[]> {
  const supabase = await createClient();
  let votesArray: VoteAndMessage[] = [];

  const { data: votes, error: votesError } = await supabase
    .from(type === 'downvote' ? 'downvoted_responses' : 'upvoted_responses')
    .select(`
      *,
      agent_requests (
        metadata
      )
    `);

  if (votesError) {
    console.error('Error fetching ' + (type === 'downvote' ? 'downvotes' : 'upvotes') + ':', votesError);
    throw votesError;
  }

  for (const vote of votes) {
    const { data: events, error: eventsError } = await supabase
      .from('agent_events')
      .select('type, metadata')
      .eq('request_id', vote.request_id)

    if (eventsError) {
      console.error('Error fetching events:', eventsError);
      throw eventsError;
    }

    const answerEvent = events.find(event => event.type === 'answer_success');

    const { data: user, error: userError } = await supabase
      .from('user_data')
      .select('email')
      .eq('id', vote.user_id)
      .single();
      
    if (userError) {
      console.error('Error fetching user:', userError);
      throw userError;
    }

    votesArray.push({
      type: type,
      id: vote.id,
      createdAt: vote.created_at,
      requestID: vote.request_id,
      userEmail: user.email,
      agent: (answerEvent?.metadata).agentID,
      reason: vote.reason,
      query: (vote.agent_requests.metadata).query,
      answer: (answerEvent?.metadata).answer
    });
  }

  return votesArray;
}