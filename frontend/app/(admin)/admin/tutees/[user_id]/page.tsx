"use client";

import { PagedObjectTable, TableRowType } from "@/components/paged-object-table";
import { useParams, useRouter } from "next/navigation";
import { useUser } from "@/providers/user-provider";
import { User } from "@/lib/supabase/user";
import { useUserChatHistory } from "@/hooks/use-user-chat-history-adapter";
import { Chat } from "@/hooks/use-chat-history";
import { getRelativeTimeString } from "@/lib/utils";

export default function TuteeChatHistoryPage() {
  return (
    <PagedChatHistoryTable />
  )
}

function UsePartiallyApplliedUserChatHistory(page: number) {
  const { user_id } = useParams();
  const userContext = useUser();
  const tutee = userContext.subordinates?.find((user: User) => user.id === user_id) || null;

  return useUserChatHistory(tutee ? tutee.id : null, page);
}

function PagedChatHistoryTable() {
  const { user_id } = useParams();
  const userContext = useUser();
  const tutee = userContext.subordinates?.find((user: User) => user.id === user_id);
  const router = useRouter();

  const capitalize = (name: string) => name.charAt(0).toUpperCase() + name.slice(1);

  return (
    <PagedObjectTable
      title="Chat History"
      description={tutee ? "View all chats for " + capitalize(tutee.firstName) + "." : ""}
      dataHook={UsePartiallyApplliedUserChatHistory}
      idField="id"
      TableHeadings={ChatHistoryTableHeadings}
      RowContents={ChatHistoryTableContents}
      ExpandedRowContents={() => null}
      rowType={TableRowType.CLICKABLE}
      rowClickHandler={(data: Chat) => {
        router.push(`/admin/tutee-chat/${data.id}?title=${data.title}&user_id=${user_id}`);
      }}
    />
  )
}

function ChatHistoryTableHeadings() {
  return (
    <>
      <th className="h-12 w-[350px] px-4 text-left align-middle font-semibold text-primary">
        Chat ID
      </th>
      <th className="h-12 w-[350px] px-4 text-left align-middle font-semibold text-primary">
        Chat Title
      </th>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary last:rounded-tr-xl">
        Time
      </th>
    </>
  )
}

function ChatHistoryTableContents<T>({data}: {data: Chat}) {
  return (
    <>
      <td className="p-4 align-middle font-mono text-sm text-foreground">
        {data.id}
      </td>
      <td className="p-4 align-middle text-sm text-foreground">
        {data.title}
      </td>
      <td className="p-4 align-middle text-sm text-foreground">
        {data.createdAt &&
          getRelativeTimeString(new Date(data.createdAt))}
      </td>
    </>
  )
}