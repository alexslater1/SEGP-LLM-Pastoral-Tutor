'use client';

import { SidebarToggle } from '@/components/sidebar-toggle';
import { Button } from './ui/button';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';

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

  if (pathname.includes('/tutees/')) {
    return (
      <BackButton text="Back to Tutees" path="/admin/tutees" />
    )
  } else if (pathname.includes('/tutee-chat/')) {
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
          <BackButton className="z-20" text="Back to Tutee Chat History" path={`/admin/tutees/${user_id}`} />
        ) : (
          <BackButton className="z-20" text="Back to Tutees" path="/admin/tutees" />
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
