'use client';

import { Download, Trash2 } from 'lucide-react';
import { useState } from "react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { getRelativeTimeString } from "@/lib/utils";
import { RagDocument } from "@/lib/db/schema";
import { useDeleteRagDoc, useRagDocuments } from '@/hooks/use-rag';

export default function LibraryPage() {
  const { data: documents, error } = useRagDocuments();
  const deleteMutation = useDeleteRagDoc();

  const handleDownload = async (doc: RagDocument) => {
    // TODO: Implement download logic here
    
    toast.success(`Downloading ${doc.name}`);
  };

  const handleDelete = async(doc: RagDocument) => {
    // TODO: Implement delete logic here
    try {
      await deleteMutation.mutateAsync(doc.name);
      toast.success(`Deleted ${doc.name}`);
    } catch (error) {
      toast.error('Failed to delete files');
      console.error('Delete error:', error);
    }
    toast.success(`Deleted ${doc.name}`);
  };

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Document Library</h2>
        <p className="text-muted-foreground">
          View and manage uploaded documents.
        </p>
      </div>

      <div className="border rounded-lg">
        <div className="min-w-full">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead>
              <tr className="bg-muted/50">
                <th scope="col" className="px-6 py-3 text-left text-sm font-semibold">
                  Document Name
                </th>
                <th scope="col" className="px-6 py-3 text-left text-sm font-semibold">
                  Size
                </th>
                <th scope="col" className="px-6 py-3 text-left text-sm font-semibold">
                  Upload Date
                </th>
                <th scope="col" className="px-6 py-3 text-right text-sm font-semibold">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {documents?.map((doc: RagDocument) => (
                <tr key={doc.id} className="hover:bg-muted/50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    <div className="flex items-center">
                      <span>{doc.name}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {doc.size}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm">
                    {getRelativeTimeString(doc.uploadedAt)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-right">
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
          {documents?.length === 0 && (
            <div className="text-center py-4 text-muted-foreground">
              No documents found
            </div>
          )}
        </div>
      </div>
    </div>
  );
} 