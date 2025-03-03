'use client';

import { PagedObjectTable } from "@/components/paged-object-table";
import { useDeleteRagDoc, useDownloadRagDoc, useRagDocuments } from "@/hooks/use-rag";
import { RagDocument } from "../../actions";
import { humanReadableSize, getRelativeTimeString } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogOverlay, DialogPortal, DialogTitle } from "@/components/ui/dialog";
import { toast } from "sonner";
import { useState } from "react";
import { Trash2, Download } from "lucide-react";
import { Markdown } from "@/components/markdown";

export default function LibraryPage() {
  return <PagedRagDocumentsTable />
}

function PagedRagDocumentsTable() {
  return (
    <PagedObjectTable
      title="Document Library"
      description="View all RAG documents uploaded to the library."
      dataHook={useRagDocuments}
      idField="id"
      TableHeadings={RagDocumentsTableHeadings}
      RowContents={RagDocumentsTableContents}
      ExpandedRowContents={RagDocumentsExpandedTableContents}
    />
  )
}

function RagDocumentsTableHeadings() {
  return (
    <>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        Document Name
      </th>
      <th className="h-12 w-[150px] px-4 text-left align-middle font-semibold text-primary">
        Size
      </th>
      <th className="h-12 w-[180px] px-4 text-left align-middle font-semibold text-primary">
        Upload Date
      </th>
      <th className="h-12 w-[100px] px-4 text-right align-middle font-semibold text-primary last:rounded-tr-xl">
        Actions
      </th>
    </>
  )
}

function RagDocumentsTableContents({data}: {data: RagDocument}) {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const downloadMutation = useDownloadRagDoc();
  const deleteMutation = useDeleteRagDoc();

  const handleDownload = async (doc: RagDocument) => {
    try {
      toast.promise(downloadMutation.mutateAsync(doc.name), {
        loading: `Downloading ${doc.name}...`,
        success: () => {
          toast.success(`Downloaded ${doc.name}`);
          return `Downloaded ${doc.name}`;
        },
        error: () => {
          toast.error("Failed to download file");
          return "Failed to download file";
        }
      });
    } catch (error) {
      console.error("Download error:", error);
    }
  };

  const handleDelete = async(doc: RagDocument) => {
    try {
      setIsDialogOpen(false);
      toast.promise(deleteMutation.mutateAsync(doc.name), {
        loading: `Deleting ${doc.name}...`,
        success: () => {
          toast.success(`Deleted ${doc.name}`);
          return `Deleted ${doc.name}`;
        },
        error: () => {
          toast.error("Failed to delete file");
          return "Failed to delete file";
        }
      });
    } catch (error) {
      console.error("Delete error:", error);
    }
  };

  return (
    <>
      <td className="p-4 align-middle">
        <div className="flex items-center">
          <span>{data.name}</span>
        </div>
      </td>
      <td className="p-4 align-middle">
        {data.document_size && humanReadableSize(data.document_size)}
      </td>
      <td className="p-4 align-middle">
        {data.date_uploaded && getRelativeTimeString(new Date(data.date_uploaded))}
      </td>
      <td className="p-4 align-middle text-right">
        <div className="flex justify-end gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => {
              e.stopPropagation();
              handleDownload(data);
            }}
            className="size-8 hover:text-primary"
          >
            <Download className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={(e) => {
              e.stopPropagation();
              setIsDialogOpen(true);
            }}
            className="size-8 text-destructive/80 hover:bg-destructive/30 hover:text-destructive"
          >
            <Trash2 className="size-4" />
          </Button>
        </div>
      </td>

      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogTitle className="text-destructive">Delete Document</DialogTitle>
            <DialogDescription className="text-foreground">
              Are you sure you want to delete &quot;{data.name}&quot;? This action cannot be undone.
            </DialogDescription>
            <div className="flex justify-end gap-2">
              <Button
                className="bg-muted hover:bg-muted/50"
                variant="outline"
                onClick={(e) => {
                  e.stopPropagation();
                  setIsDialogOpen(false)
                }}
              >
                Cancel
              </Button>
              <Button
                className="bg-destructive hover:bg-destructive/50"
                variant="destructive"
                onClick={(e) => {
                  e.stopPropagation();
                  handleDelete(data);
                }}
              >
                Delete
              </Button>
            </div>
          </DialogContent>
        </DialogPortal>
      </Dialog>
    </>
  )
}

function RagDocumentsExpandedTableContents({data}: {data: RagDocument}) {
  return (
    <>
      <div>
        <h4 className="font-semibold text-primary">Document Type:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.document_type}</Markdown>
        </div>
      </div>
      <div>
        <h4 className="font-semibold text-primary">URL:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.url}</Markdown>
        </div>
      </div>
    </>
  )
}


