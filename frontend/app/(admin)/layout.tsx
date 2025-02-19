import { cookies } from "next/headers";
import { getUser } from "@/lib/supabase/user";
import { AdminSidebar } from "@/components/admin-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { AdminHeader } from "@/components/admin-header";
import { redirect } from "next/navigation";

export default async function AdminLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, cookieStore] = await Promise.all([getUser(), cookies()]);
  const isCollapsed = cookieStore.get("sidebar:state")?.value !== "true";

  if (!user) {
    redirect("/sign-in");
  } else if (user.role !== "admin") {
    redirect("/");
  }

  return (
    <SidebarProvider defaultOpen={!isCollapsed}>
      <AdminSidebar user={user} />
      <SidebarInset>
        <AdminHeader />
        {children}
      </SidebarInset>
    </SidebarProvider>
  );
}
