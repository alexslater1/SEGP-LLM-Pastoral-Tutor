"use client";

import { notFound, useParams } from 'next/navigation';
import { Chat } from '@/components/chat';

export default function Page() {
  const { chat_id } = useParams();
  if (!chat_id) {
    notFound();
  }

  return (
    <Chat
      id={chat_id as string}
      isReadonly={true}      
    />
  );
}
