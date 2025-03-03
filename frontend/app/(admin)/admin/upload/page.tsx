"use client";

import { Button } from "@/components/ui/button";
import { File, CloudUpload } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { cn, getRelativeTimeString, humanReadableSize } from "@/lib/utils";
import { useRagDocuments, useRagUploadDocs } from "@/hooks/use-rag";

const getFileExtension = (filename: string) => {
  return filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2);
};

export default function UploadPage() {
  const [isDragging, setIsDragging] = useState(false);
  const { data: documents } = useRagDocuments(0);
  const uploadMutation = useRagUploadDocs();

  const getRecentDocuments = () => {
    if (!documents?.data) return [];
    
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
    
    return documents.data
      .filter(doc => new Date(doc.date_uploaded) > sevenDaysAgo)
      .sort((a, b) => new Date(b.date_uploaded).getTime() - new Date(a.date_uploaded).getTime())
      .slice(0, 10); // Cap at 10 documents
  };

  const recentDocuments = getRecentDocuments();

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

      <div className="border rounded-xl bg-card">
        <div className="overflow-hidden rounded-xl">
          <div className="border-b bg-table-header">
            <h3 className="p-4 font-semibold text-primary">
              Recently Uploaded Documents (last 7 days)
            </h3>
          </div>

          <div className="divide-y divide-border">
            {recentDocuments.length === 0 ? (
              <div className="p-4 text-center bg-table-row-odd text-foreground">
                No documents uploaded in the last 7 days.
              </div>
            ) : (
              <div className="divide-y divide-border">
                {recentDocuments.map((doc, index) => (
                  <div
                    key={doc.id}
                    className={cn(
                      "p-4 flex items-center transition-colors",
                      index % 2 === 0 
                        ? "bg-table-row-odd" 
                        : "bg-table-row-even",
                      "hover:bg-muted/50"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <File className="size-5 text-muted-foreground" />
                      <div>
                        <p className="font-medium text-foreground">{doc.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {humanReadableSize(doc.document_size)} • {getRelativeTimeString(new Date(doc.date_uploaded))}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
} 