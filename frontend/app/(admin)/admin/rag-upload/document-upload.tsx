"use client";

import { Button } from "@/components/ui/button";
import { CloudUpload } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { cn } from "@/lib/utils";
import { useRagDocuments, useRagUploadDocs } from "@/hooks/use-rag";
import { RecentItemsTable, getRecentItems } from "@/components/recent-items-table";

const getFileExtension = (filename: string) => {
  return filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2);
};

export function DocumentUpload() {
  const [isDragging, setIsDragging] = useState(false);
  const { data: documents } = useRagDocuments(0);
  const uploadMutation = useRagUploadDocs();

  const recentDocuments = getRecentItems(documents?.data).map(doc => ({
    ...doc,
    type: 'document' as const
  }));

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    handleFiles(files);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      handleFiles(files);
    }
  };

  const isValidFileType = (file: File) => {
    const supportedTypes = [
      "application/pdf",
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "application/vnd.openxmlformats-officedocument.presentationml.presentation",
      "text/plain"
    ];
    return supportedTypes.includes(file.type);
  };

  const handleFiles = async (files: File[]) => {
    console.log("Uploading files:", files);
    
    const invalidFiles = files.filter(file => !isValidFileType(file));
    
    if (invalidFiles.length > 0) {
      toast.error(
        <div>
          <p>
            Unsupported file format{invalidFiles.length > 1 ? "s" : ""}:   
            {invalidFiles.map(f => ` ${getFileExtension(f.name).toUpperCase()}`).join(", ")}
          </p>
          <p>Only PDF, DOCX, PPTX, and TXT files are supported.</p>
        </div>
      );
      return;
    }

    try {
      toast.promise(uploadMutation.mutateAsync(files), {
        loading: `Uploading ${files.length} file${files.length > 1 ? "s" : ""}...`,
        success: () => {
          return `Uploaded ${files.length} file${files.length > 1 ? "s" : ""}`;
        },
        error: () => {
          return "Failed to upload files";
        }
      });
    } catch (error) {
      console.error("Upload error:", error);
    }
  };

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold text-primary tracking-tight">Document Upload</h2>
        <p className="text-foreground">
          Upload documents for the AI tutor to use as extra knowledge.
        </p>
      </div>

      <div className="h-px bg-border" />

      <div
        className={cn(
          "border-2 border-dashed rounded-lg p-8 transition-colors",
          isDragging 
            ? "border-primary bg-primary/5" 
            : "border-muted-foreground/25 hover:border-muted-foreground/50"
        )}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <div className="flex flex-col items-center justify-center gap-4">
          <CloudUpload className={cn(
            "size-12",
            isDragging ? "text-primary" : "text-muted-foreground"
          )} />
          <p className="text-lg">Drag and drop files here</p>
          <p className="text-muted-foreground">or</p>
          <Button 
            variant="outline"
            className="bg-button text-button-foreground hover:bg-button/50 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={uploadMutation.isPending}
          >
            <label className="cursor-pointer">
              <input
                type="file"
                className="hidden"
                multiple
                onChange={handleFileInput}
                accept=".pdf,.docx,.txt,.pptx"
                disabled={uploadMutation.isPending}
              />
              {uploadMutation.isPending ? "Uploading..." : "Browse Files"}
            </label>
          </Button>
          <p className="text-sm text-muted-foreground">
            Supported formats: PDF, DOCX, TXT, PPTX
          </p>
        </div>
      </div>

      {/*<div className="h-px bg-border" />

      <RecentItemsTable
        items={recentDocuments}
        emptyMessage="No documents uploaded in the last 7 days."
        title="Recently Uploaded Documents (last 7 days)"
      />*/}
    </div>
  );
} 