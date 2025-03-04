"use server";

import { createClient } from "@/lib/supabase/server";
import { User } from "@/lib/supabase/user";
import { SupabaseClient } from "@supabase/supabase-js";

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

  try {
    await Promise.all(votes.map(async (vote) => {
      const additionalVoteData = await getAdditionalVoteData(vote, supabase);
      votesArray.push({ 
        type: type,
        id: vote.id,
        createdAt: vote.created_at,
        requestID: vote.request_id,
        userEmail: additionalVoteData.userEmail,
        agent: additionalVoteData.agent,
        reason: vote.reason,
        query: (vote.agent_requests.metadata).query,
        answer: additionalVoteData.answer
    })}));
  } catch (error) {
    console.error('Error fetching additional vote data:', error);
    throw error;
  }

  return votesArray;
}

type AdditionalVoteData = {
  userEmail: string;
  agent: string;
  answer: string;
}

async function getAdditionalVoteData(vote: Vote, supabase: SupabaseClient): Promise<AdditionalVoteData> {
  const [eventsData, userData] = await Promise.all([
    supabase
      .from('agent_events')
      .select('type, metadata')
      .eq('request_id', vote.request_id),
    supabase
      .from('user_data')
      .select('email')
      .eq('id', vote.user_id)
      .single()
  ])

  if (eventsData.error) {
    console.error('Error fetching events:', eventsData.error);
    throw eventsData.error;
  }

  const answerEvent = eventsData.data.find(event => event.type === 'answer_success');

  if (userData.error) {
    console.error('Error fetching user:', userData.error);
    throw userData.error;
  }

  return {
    userEmail: userData.data.email,
    agent: (answerEvent?.metadata).agentID,
    answer: (answerEvent?.metadata).answer
  };
}