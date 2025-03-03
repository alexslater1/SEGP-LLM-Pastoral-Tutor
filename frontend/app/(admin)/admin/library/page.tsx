'use client';

import { useState } from "react";
import { PagedRagDocumentsTable } from "./paged-rag-documents-table";
import { PagedRagWebpagesTable } from "./paged-rag-webpages-table";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export default function LibraryPage() {
  return <PagedRagSourcesTable />;
} 

function PagedRagSourcesTable() {
  const [showDocuments, setShowDocuments] = useState(true);

  return (
    <div className="relative">
      <div className="absolute right-6 top-8">
        <div className="relative flex items-center bg-card rounded-lg p-1 border shadow-sm">
          <div 
            className="absolute h-[85%] top-[7.5%] bg-primary/10 rounded-md transition-all duration-300 ease-out"
            style={{
              left: showDocuments ? '3px' : '50%',
              width: 'calc(50% - 6px)',
            }}
          />
          
          <Button
            onClick={() => setShowDocuments(true)}
            variant="ghost"
            className={cn(
              "relative px-2 z-10 transition-colors duration-300",
              showDocuments ? "font-bold text-primary hover:bg-transparent" : "text-muted-foreground"
            )}
          >
            Documents
          </Button>
          <Button
            onClick={() => setShowDocuments(false)}
            variant="ghost"
            className={cn(
              "relative px-3 z-10 transition-colors duration-300",
              !showDocuments ? "font-bold text-primary hover:bg-transparent" : "text-muted-foreground"
            )}
          >
            Webpages
          </Button>
        </div>
      </div>

      {showDocuments ? <PagedRagDocumentsTable /> : <PagedRagWebpagesTable />}
    </div>
  );
}