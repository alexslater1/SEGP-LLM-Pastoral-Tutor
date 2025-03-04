'use client';

import { SidebarToggle } from '@/components/sidebar-toggle';
import { Button } from './ui/button';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';
import { useUser } from '@/providers/user-provider';
import { UserRoleEnum } from '@/app/(admin)/role-authorization';

export function AdminHeader() {
  const pathname = usePathname();
  const router = useRouter();

  return (
    <header className="flex sticky top-0 bg-background py-1.5 items-center gap-2">
      <div className="px-0"></div>
      <SidebarToggle className="z-20" />
      <AdditionalAdminHeaderElements pathname={pathname} />
    </header>
  );
} 

function AdditionalAdminHeaderElements({ pathname }: { pathname: string }) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const userContext = useUser();
  const subordinatesName = userContext.user?.role === UserRoleEnum.TUTOR ? "Tutee" : "User";

  if (pathname.includes('/users/')) {
    return (
      <BackButton 
        text={"Back to " + subordinatesName + "s"}
        path="/admin/users" 
      />
    )
  } else if (pathname.includes('/user-chat/')) {
    const title = searchParams.get('title');
    const user_id = searchParams.get('user_id');
    return (
      <>
        <div className="absolute z-10 bg-muted size-full">
          <div className="flex flex-col justify-center h-full">
            <div className="items-center text-center text-primary text-lg">
              {title ? title + " (Read-Only)" : "Read-Only Chat"}
            </div>
          </div>
        </div>
        {user_id ? (
          <BackButton 
            className="z-20" 
            text={"Back to " + subordinatesName + " Chat History"} 
            path={`/admin/users/${user_id}`} 
          />
        ) : (
          <BackButton 
            className="z-20" 
            text={"Back to " + subordinatesName + "s"} 
            path="/admin/users" 
          />
        )}
      </>
    )
  }
}

function BackButton({ text, path, className }: { text: string, path: string, className?: string }) {
  const router = useRouter();

  return (
    <Button
      onClick={() => router.push(path)}
      variant="outline"
      className={`md:px-2 md:h-fit hover:text-primary ${className}`}
    >
      <ArrowLeft size={30} />
      <div className="flex flex-row gap-2 items-center">
        {text}
      </div>
    </Button>
  )
} 
