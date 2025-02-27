"use client";

import type { User } from "@/lib/supabase/user";
import { usePathname } from "next/navigation";
import { SidebarUserNav } from "@/components/sidebar-user-nav";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  Home,
  FileUp,
  Library,
  Activity,
  MessageSquare,
  Settings,
  ThumbsDown,
} from "lucide-react";
import { AdminSidebarMenuButton } from "@/components/admin-sidebar-menu-button";

export function AdminSidebar({ user }: { user: User | null }) {
  const { setOpenMobile } = useSidebar();
  const pathname = usePathname();

  const isSelected = (path: string) =>
    path === "/admin" ? pathname === "/admin" : pathname.includes(path);

  const menuItems = [
    { href: "/admin", icon: Home, label: "Dashboard" },
    { href: "/admin/upload", icon: FileUp, label: "Document Upload" },
    { href: "/admin/library", icon: Library, label: "Document Library" },
    { href: "/admin/agent-events", icon: Activity, label: "Agent Events" },
    { href: "/admin/agent-requests", icon: MessageSquare, label: "Agent Requests" },
    { href: "/admin/agent-config", icon: Settings, label: "Agent Config" },
    { href: "/admin/downvotes", icon: ThumbsDown, label: "Downvoted Responses" },
  ];

  return (
    <Sidebar className="group-data-[side=left]:border-r-0">
      <SidebarHeader>
        <SidebarMenu>
          <div className="flex flex-row justify-between items-center">
            <span className="text-lg text-primary font-semibold px-2">Admin Dashboard</span>
          </div>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarMenu className="px-2 space-y-2">
          {menuItems.map((item) => (
            <AdminSidebarMenuButton
              key={item.href}
              href={item.href}
              icon={item.icon}
              label={item.label}
              isSelected={isSelected(item.href)}
              onMobileClose={() => setOpenMobile(false)}
            />
          ))}
        </SidebarMenu>
      </SidebarContent>
      <SidebarFooter>
        {user && <SidebarUserNav user={user} adminPage />}
      </SidebarFooter>
    </Sidebar>
  );
}
