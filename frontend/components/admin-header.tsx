'use client';

import { SidebarToggle } from '@/components/sidebar-toggle';
import { Button } from './ui/button';
import { usePathname, useRouter } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';

export function AdminHeader() {
  const pathname = usePathname();
  const router = useRouter();

  return (
    <header className="flex sticky top-0 bg-background py-1.5 items-center px-2 md:px-2 gap-2">
      <SidebarToggle />
      {pathname.includes('/tutees/') && (
        <Button
          onClick={() => router.push('/admin/tutees')}
          variant="outline"
          className="md:px-2 md:h-fit hover:text-primary"
        >
          <ArrowLeft size={30} />
          <div className="flex flex-row gap-2 items-center">
            Back to Tutees
          </div>
        </Button>
      )}
    </header>
  );
} 