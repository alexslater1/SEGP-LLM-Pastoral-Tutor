import { AppSidebar } from '@/components/app-sidebar';
import { cookies } from 'next/headers';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';
import Script from 'next/script';
import { SessionHistoryProvider } from '@/providers/session-history-provider';
import { ChatUrlProvider } from '@/providers/chat-url-provider';

export const experimental_ppr = true;

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const cookieStore = await cookies();
  const isCollapsed = cookieStore.get('sidebar:state')?.value !== 'true';

  return (
    <>
      <Script
        src="https://cdn.jsdelivr.net/pyodide/v0.23.4/full/pyodide.js"
        strategy="beforeInteractive"
      />
      <SidebarProvider defaultOpen={!isCollapsed}>
        <SessionHistoryProvider>
          <ChatUrlProvider>
            <AppSidebar />
            <SidebarInset>{children}</SidebarInset>
          </ChatUrlProvider>
        </SessionHistoryProvider>
      </SidebarProvider>
    </>
  );
}
