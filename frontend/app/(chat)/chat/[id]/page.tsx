import { notFound } from 'next/navigation';

import { auth } from '@/app/(auth)/auth';
import { Chat } from '@/components/chat';
import { VisibilityType } from '@/components/visibility-selector';

export default async function Page(props: { params: Promise<{ id: string }> }) {
  const params = await props.params;
  const { id } = params;
  //const chat = await getChatById({ id });
  const chat = {
    id: id,
    visibility: 'private' as VisibilityType,
    userId: '123',
  };

  if (!chat) {
    notFound();
  }

  /*const session = await auth();

  if (chat.visibility === 'private') {
    if (!session || !session.user) {
      return notFound();
    }

    if (session.user.id !== chat.userId) {
      return notFound();
    }
  }
  */

  return (
    <>
      <Chat
        id={chat.id}
        selectedVisibilityType={chat.visibility}
        isReadonly={chat.visibility === 'public'} // session?.user?.id !== chat.userId}
      />
    </>
  );
}
