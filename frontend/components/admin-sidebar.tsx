'use client';

import type { User } from 'next-auth';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { SidebarUserNav } from '@/components/sidebar-user-nav';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  useSidebar,
} from '@/components/ui/sidebar';
import { Home, FileUp } from 'lucide-react';
import { cn } from '@/lib/utils';

export function AdminSidebar({ user }: { user: User | undefined }) {
  const { setOpenMobile } = useSidebar();
  const pathname = usePathname();

  const getSelectedStyles = (path: string) => {
    const isSelected = path === '/admin' 
      ? pathname === '/admin'
      : pathname.includes(path);
    
    return isSelected && "bg-[hsl(var(--sidebar-selected-bg)_/_0.15)] text-[hsl(var(--sidebar-selected))]";
  };

  return (
    <Sidebar className="group-data-[side=left]:border-r-0">
      <SidebarHeader>
        <SidebarMenu>
          <div className="flex flex-row justify-between items-center">
            <span className="text-lg font-semibold px-2">
              Admin Dashboard
            </span>
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
                getSelectedStyles('/admin')
              )}
            >
              <Link 
                href="/admin" 
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <Home size={24} />
                <span className="text-base font-semibold">Dashboard</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton 
              asChild
              className={cn(
                "w-full px-5 py-5 rounded-lg transition-colors",
                "hover:bg-muted/50",
                getSelectedStyles('/admin/documents')
              )}
            >
              <Link 
                href="/admin/documents" 
                onClick={() => setOpenMobile(false)}
                className="flex items-center gap-3"
              >
                <FileUp size={24} />
                <span className="text-base font-semibold">Document Upload</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarContent>
      <SidebarFooter>{user && <SidebarUserNav user={user} />}</SidebarFooter>
    </Sidebar>
  );
} 