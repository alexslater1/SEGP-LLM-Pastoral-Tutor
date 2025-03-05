"use client";

import { usePathname } from "next/navigation";
import { SidebarUserNav } from "@/components/sidebar-user-nav";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  useSidebar,
} from "@/components/ui/sidebar";
import {
  Home,
  Upload,
  Library,
  Activity,
  MessageSquare,
  Settings,
  Vote,
  UserRound,
} from "lucide-react";
import { useUser } from "@/providers/user-provider";

import { AdminSidebarMenuButton } from "@/components/admin-sidebar-menu-button";
import { UserRoleEnum } from "@/app/(admin)/role-authorization";

export function AdminSidebar() {
  const { setOpenMobile } = useSidebar();
  const pathname = usePathname();
  const userContext = useUser();

  const isSelected = (path: string) =>
    path === "/admin" ? pathname === "/admin" : pathname.includes(path);

  const adminMenuItems = [
    { href: "/admin", icon: Home, label: "Dashboard" },
    { href: "/admin/users", icon: UserRound, label: "Users" },
    { href: "/admin/rag-upload", icon: Upload, label: "RAG Upload" },
    { href: "/admin/rag-library", icon: Library, label: "RAG Library" },
    { href: "/admin/agent-events", icon: Activity, label: "Agent Events" },
    // { href: "/admin/agent-requests", icon: MessageSquare, label: "Agent Requests" },
    { href: "/admin/agent-config", icon: Settings, label: "Agent Config" },
    { href: "/admin/feedback", icon: Vote, label: "Feedback" },
  ];

  const tutorMenuItems = [
    { href: "/admin", icon: Home, label: "Dashboard" },
    { href: "/admin/users", icon: UserRound, label: "Tutees" },
  ];

  if (!userContext.user) {
    return (
      <Sidebar className="group-data-[side=left]:border-r-0">
        <SidebarHeader>
          <SidebarMenu>
            <div className="flex flex-row justify-between items-center">
              <span className="text-lg text-primary font-semibold px-2">
                Dashboard
              </span>
            </div>
          </SidebarMenu>
        </SidebarHeader>
        <SidebarContent>
          <div className="flex flex-col gap-2 text-error text-center py-4">
            Couldn&apos;t load user
          </div>
        </SidebarContent>
      </Sidebar>
    );
  }

  const menuItems =
    userContext.user.role === UserRoleEnum.ADMIN
      ? adminMenuItems
      : tutorMenuItems;

  return (
    <Sidebar className="group-data-[side=left]:border-r-0">
      <SidebarHeader>
        <SidebarMenu>
          <div className="flex flex-row justify-between items-center">
            <span className="text-lg text-primary font-semibold px-2">
              {userContext.user.role === UserRoleEnum.ADMIN
                ? "Admin Dashboard"
                : "Tutor Dashboard"}
            </span>
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
        {userContext.user && (
          <SidebarUserNav user={userContext.user} adminPage />
        )}
      </SidebarFooter>
    </Sidebar>
  );
}
