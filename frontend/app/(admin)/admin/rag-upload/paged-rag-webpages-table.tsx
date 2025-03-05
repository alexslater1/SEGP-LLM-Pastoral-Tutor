import { useDeleteRagUrl, useRagWebpages } from "@/hooks/use-rag";
import { PagedObjectTable, formatDate } from "@/components/paged-object-table";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { truncateUrl } from "@/lib/utils";
import { Link, Trash2 } from "lucide-react";
import { Dialog, DialogContent, DialogOverlay, DialogPortal, DialogTitle } from "@/components/ui/dialog";
import { Markdown } from "@/components/markdown";
import { RagWebpage } from "../../actions";

export function PagedRagWebpagesTable() {
  return (
    <PagedObjectTable
      title=""
      description=""
      dataHook={useRagWebpages}
      idField="id"
      TableHeadings={RagWebpagesTableHeadings}
      RowContents={RagWebpagesTableContents}
      ExpandedRowContents={RagWebpagesExpandedTableContents}
    />
  )
}

function RagWebpagesTableHeadings() {
  return (
    <>
      <th className="h-12 px-4 text-left align-middle font-semibold text-primary">
        URL
      </th>
      <th className="h-12 w-[200px] px-4 text-left align-middle font-semibold text-primary">
        Upload Date
      </th>
      <th className="h-12 w-[100px] px-4 text-right align-middle font-semibold text-primary last:rounded-tr-xl">
        Actions
      </th>
    </>
  );
}

function RagWebpagesTableContents({data}: {data: RagWebpage}) {
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const deleteMutation = useDeleteRagUrl();

  const handleDelete = async(webpage: RagWebpage) => {
    try {
      setIsDialogOpen(false);
      toast.promise(deleteMutation.mutateAsync(webpage.url), {
        loading: `Deleting ${truncateUrl(webpage.url)}...`,
        success: () => {
          return `Deleted ${truncateUrl(webpage.url)}`;
        },
        error: () => {
          return "Failed to delete URL";
        }
      });
    } catch (error) {
      console.error("Delete error:", error);
    }
  };
    
  return (
    <>
      <td className="p-4 align-middle">
        <div className="flex items-center gap-3">
          <span className="font-mono">{truncateUrl(data.url)}</span>
        </div>
      </td>
      <td className="p-4 align-middle">
        {data.date_uploaded && formatDate(data.date_uploaded)}
      </td>
      <td className="p-4 align-middle text-right">
        <div className="flex justify-end gap-2">
          <a
            href={data.url}
            target="_blank"
            rel="noopener noreferrer"
            onClick={(e) => e.stopPropagation()}
          >
            <Button
              variant="ghost"
              size="icon"
              className="size-8 hover:text-primary"
            >
              <Link className="size-4" />
            </Button>
          </a>
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
            <DialogTitle className="text-destructive">Delete URL</DialogTitle>
            <div className="text-foreground">
              Are you sure you want to delete:
              <div className="mt-2 font-mono break-all bg-muted/25 p-2 rounded-md">
                {data.url}
              </div>
              <div className="mt-2">
                This action cannot be undone.
              </div>
            </div>
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
  );
}

function RagWebpagesExpandedTableContents({data}: {data: RagWebpage}) {
  return (
    <>
      <div>
        <h4 className="font-semibold text-primary">Full URL:</h4>
        <div className="text-foreground mt-1">
          <Markdown>{data.url}</Markdown>
        </div>
      </div>
    </>
  )
}