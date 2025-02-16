import { Chat } from '@/components/chat';
import { generateUUID } from '@/lib/utils';

export default async function Page() {
  const key = generateUUID();
    return (
    <>
      <Chat
        key={key}
        id={null}
        selectedVisibilityType="private"
        isReadonly={false}
      />
    </>
  );
}
