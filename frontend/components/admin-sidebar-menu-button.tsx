import Link from "next/link";
import { LucideIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import {
  SidebarMenuItem,
  SidebarMenuButton,
} from "@/components/ui/sidebar";

interface AdminSidebarMenuButtonProps {
  href: string;
  icon: LucideIcon;
  label: string;
  isSelected: boolean;
  onMobileClose: () => void;
}

export function AdminSidebarMenuButton({
  href,
  icon: Icon,
  label,
  isSelected,
  onMobileClose,
}: AdminSidebarMenuButtonProps) {
  return (
    <SidebarMenuItem>
      <SidebarMenuButton
        asChild
        className={cn(
          "w-full px-5 py-5 rounded-lg transition-colors",
          "!hover:bg-muted/50 !hover:text-foreground",
          isSelected && "bg-[hsl(var(--sidebar-selected-bg)_/_0.15)] text-[hsl(var(--sidebar-selected))]"
        )}
      >
        <Link
          href={href}
          onClick={onMobileClose}
          className="flex items-center gap-3"
        >
          <Icon size={24} />
          <span
            className={cn(
              "text-base",
              isSelected ? "font-semibold" : "font-normal"
            )}
          >
            {label}
          </span>
        </Link>
      </SidebarMenuButton>
    </SidebarMenuItem>
  );
}
