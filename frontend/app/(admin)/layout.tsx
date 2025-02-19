import { cookies } from "next/headers";
import { getUser } from "@/lib/supabase/user";
import { redirect } from "next/navigation";
import { AdminSidebar } from "@/components/admin-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AdminHeader } from "@/components/admin-header";

export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, cookieStore] = await Promise.all([getUser(), cookies()]);
  const isCollapsed = cookieStore.get("sidebar:state")?.value !== "true";

  // if (!session?.user || session.user.role !== 'admin') {
  //   redirect('/login');
  // }

  const fakeUser = {
    id: "1",
    role: "admin",
    email: user?.email,
  };

  return (
    <SidebarProvider defaultOpen={!isCollapsed}>
      <AdminSidebar user={fakeUser as any} />
      <SidebarInset>
        <AdminHeader />
        {children}
      </SidebarInset>
    </SidebarProvider>
  );
}
