"use client";

import type { User } from "@/lib/supabase/user";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { SidebarUserNav } from "@/components/sidebar-user-nav";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
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
import { cn } from "@/lib/utils";

export function AdminSidebar({ user }: { user: User | null }) {
  const { setOpenMobile } = useSidebar();
  const pathname = usePathname();

  const getSelectedStyles = (path: string) => {
    const isSelected =
      path === "/admin" ? pathname === "/admin" : pathname.includes(path);

    return (
      isSelected &&
      "bg-[hsl(var(--sidebar-selected-bg)_/_0.15)] text-[hsl(var(--sidebar-selected))]"
    );
  };

  return (
    <Sidebar className="group-data-[side=left]:border-r-0">
      <SidebarHeader>
        <SidebarMenu>
          <div className="flex flex-row justify-between items-center">
            <span className="text-lg font-semibold px-2">Admin Dashboard</span>
          </div>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarMenu className="px-2 space-y-2">
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin")
              )}
            >
              <Link
                href="/admin"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <Home size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Dashboard
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/upload")
              )}
            >
              <Link
                href="/admin/upload"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <FileUp size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/upload")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Document Upload
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/library")
              )}
            >
              <Link
                href="/admin/library"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <Library size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/library")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Document Library
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/agent-events")
              )}
            >
              <Link
                href="/admin/agent-events"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <Activity size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/agent-events")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Agent Events
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/agent-requests")
              )}
            >
              <Link
                href="/admin/agent-requests"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <MessageSquare size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/agent-requests")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Agent Requests
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/agent-requests")
              )}
            >
              <Link
                href="/admin/agent-config"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <Settings size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/agent-config")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Agent Config
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles("/admin/downvotes")
              )}
            >
              <Link
                href="/admin/downvotes"
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <ThumbsDown size={24} />
                <span
                  className={cn(
                    "text-base",
                    getSelectedStyles("/admin/downvotes")
                      ? "font-semibold"
                      : "font-normal"
                  )}
                >
                  Downvoted Responses
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarContent>
      <SidebarFooter>
        {user && <SidebarUserNav user={user} adminPage />}
      </SidebarFooter>
    </Sidebar>
  );
}
