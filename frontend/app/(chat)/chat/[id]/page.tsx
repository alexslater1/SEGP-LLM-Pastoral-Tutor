"use client";

import { notFound, useParams } from 'next/navigation';
import { Chat } from '@/components/chat';

export default function Page() {
  const { id } = useParams();

  if (!id) {
    notFound();
  }

  return (
    <>
      <Chat
        id={id as string}
        isReadonly={false}
      />
    </>
  );
}
