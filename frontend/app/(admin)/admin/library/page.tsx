'use client';

import { Download, Trash2 } from 'lucide-react';
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { RagDocument } from "@/lib/db/schema";
import { cn, getRelativeTimeString } from "@/lib/utils";

export default function LibraryPage() {
  const [documents] = useState<RagDocument[]>([
    {
      id: '1',
      name: 'Computing-UG-Handbook-2425-v1b.pdf',
      uploadedAt: new Date('2024-03-15T10:00:00'),
      size: '2.4 MB',
      type: 'PDF',
      status: 'ready',
      userId: 'dummy-user-id',
      backendSourceId: 1
    },
    {
      id: '2',
      name: 'student-guidelines.docx',
      uploadedAt: new Date('2024-03-14T15:30:00'),
      size: '1.2 MB',
      type: 'DOCX',
      status: 'ready',
      userId: 'dummy-user-id',
      backendSourceId: 2
    },
    {
      id: '3',
      name: 'course-outline.txt',
      uploadedAt: new Date('2024-03-13T09:15:00'),
      size: '156 KB',
      type: 'TXT',
      status: 'ready',
      userId: 'dummy-user-id',
      backendSourceId: 3
    }
  ]);

  const handleDownload = (doc: RagDocument) => {
    // TODO: Implement download logic here
    toast.success(`Downloading ${doc.name}`);
  };

  const handleDelete = (doc: RagDocument) => {
    // TODO: Implement delete logic here
    toast.success(`Deleted ${doc.name}`);
  };

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl text-primary font-bold tracking-tight">Document Library</h2>
        <p className="text-foreground">
          View and manage uploaded documents.
        </p>
      </div>

      <div className="border rounded-xl bg-card">
        <div className="overflow-hidden rounded-xl">
          <table className="w-full">
            <thead>
              <tr className="border-b bg-table-header">
                <th scope="col" className="h-12 px-4 text-left align-middle font-semibold text-primary first:rounded-tl-xl">
                  Document Name
                </th>
                <th scope="col" className="h-12 px-4 text-left align-middle font-semibold text-primary">
                  Size
                </th>
                <th scope="col" className="h-12 px-4 text-left align-middle font-semibold text-primary">
                  Upload Date
                </th>
                <th scope="col" className="h-12 px-4 text-right align-middle font-semibold text-primary last:rounded-tr-xl">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {documents.map((doc, index) => (
                <tr
                  key={doc.id}
                  className={cn(
                    "transition-colors hover:bg-muted/50",
                    index % 2 === 0 
                      ? "bg-table-row-odd" 
                      : "bg-table-row-even"
                  )}
                >
                  <td className="p-4 align-middle">
                    <div className="flex items-center">
                      <span>{doc.name}</span>
                    </div>
                  </td>
                  <td className="p-4 align-middle">
                    {doc.size}
                  </td>
                  <td className="p-4 align-middle">
                    {getRelativeTimeString(doc.uploadedAt)}
                  </td>
                  <td className="p-4 align-middle text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDownload(doc)}
                        className="size-8 hover:text-blue-500"
                      >
                        <Download className="size-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDelete(doc)}
                        className="size-8 text-destructive hover:text-red-500"
                      >
                        <Trash2 className="size-4" />
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {documents.length === 0 && (
            <div className="text-center py-4 text-muted-foreground">
              No documents found
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 