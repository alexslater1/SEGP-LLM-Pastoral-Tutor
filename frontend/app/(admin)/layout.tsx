import { cookies } from 'next/headers';
import { auth } from '../(auth)/auth';
import { redirect } from 'next/navigation';
import { AdminSidebar } from '@/components/admin-sidebar';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';

export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [session, cookieStore] = await Promise.all([auth(), cookies()]);
  const isCollapsed = cookieStore.get('sidebar:state')?.value !== 'true';

  if (!session?.user || session.user.role !== 'admin') {
    redirect('/login');
  }

  return (
    <SidebarProvider defaultOpen={!isCollapsed}>
      <AdminSidebar user={session.user} />
      <SidebarInset>{children}</SidebarInset>
    </SidebarProvider>
  );
} 