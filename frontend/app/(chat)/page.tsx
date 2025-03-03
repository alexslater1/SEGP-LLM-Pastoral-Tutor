import { Chat } from '@/components/chat';
import { generateUUID } from '@/lib/utils';

export default function Page() {
  const key = generateUUID();
    return (
    <>
      <Chat
        key={key}
        id={null}
        isReadonly={false}
      />
    </>
  );
}
