'use client';

import { Button } from "@/components/ui/button";
import { File, CloudUpload } from 'lucide-react';
import { useState, useEffect } from "react";
import { toast } from "sonner";
import { getRelativeTimeString } from "@/lib/utils";
import { RagDocument } from "@/lib/db/schema";
import { useRagDocuments, useRagUploadDocs } from "@/hooks/use-rag";

const getFileExtension = (filename: string) => {
  return filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2);
};

export default function UploadPage() {
  const { data: documents, error } = useRagDocuments();
  const [isDragging, setIsDragging] = useState(false);
  
  useEffect(() => {
    if (error) {
      toast.error("Error fetching documents", {
        description: error.message
      });
    }
  }, [error]);

  const uploadMutation = useRagUploadDocs();

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
      'application/pdf',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      'application/vnd.openxmlformats-officedocument.presentationml.presentation',
      'text/plain'
    ];
    return supportedTypes.includes(file.type);
  };

  const handleFiles = async (files: File[]) => {
    console.log("Uploading files:", files);
    console.log("File type:", files[0].type);

    const invalidFiles = files.filter(file => !isValidFileType(file));
    
    if (invalidFiles.length > 0) {
      toast.error(
        <div>
          <p>
            Unsupported file format{invalidFiles.length > 1 ? 's' : ''}:   
            {invalidFiles.map(f => ` ${getFileExtension(f.name).toUpperCase()}`).join(', ')}
          </p>
          <p>Only PDF, DOCX, PPTX, and TXT files are supported.</p>
        </div>
      );
      return;
    }
  
    try {
      await uploadMutation.mutateAsync(files);
      toast.success(`Uploaded ${files.length} file(s)`);
    } catch (error) {
      toast.error('Failed to upload files');
      console.error('Upload error:', error);
    }
  };

  const getRecentDocuments = () => {
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
    
    return documents ? documents
      .filter(doc => doc.uploadedAt > sevenDaysAgo)
      .sort((a, b) => b.uploadedAt.getTime() - a.uploadedAt.getTime())
      : [];
  };

  const recentDocuments = getRecentDocuments();

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Document Upload</h2>
        <p className="text-muted-foreground">
          Upload documents for the AI pastoral tutor to learn from.
        </p>
      </div>

      <div
        className={`border-2 border-dashed rounded-lg p-8 transition-colors ${
          isDragging 
            ? 'border-primary bg-primary/5' 
            : 'border-muted-foreground/25 hover:border-muted-foreground/50'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <div className="flex flex-col items-center justify-center gap-4">
          <CloudUpload className={`size-12 ${isDragging ? 'text-primary' : 'text-muted-foreground'}`} />
          <p className="text-lg">Drag and drop files here</p>
          <p className="text-muted-foreground">or</p>
          <Button>
            <label className="cursor-pointer">
              <input
                type="file"
                className="hidden"
                multiple
                onChange={handleFileInput}
                accept=".pdf,.docx,.txt,.pptx"
              />
              Browse Files
            </label>
          </Button>
          <p className="text-sm text-muted-foreground">
            Supported formats: PDF, DOCX, TXT, PPTX
          </p>
        </div>
      </div>

      <div className="border rounded-lg">
        <div className="p-4">
          <h3 className="font-semibold">Recent Documents (last 7 days)</h3>
        </div>

        <div className="border-t">
          {recentDocuments.length === 0 ? (
            <div className="p-4 text-center text-muted-foreground">
              No documents uploaded in the last 7 days
            </div>
          ) : (
            <div className="divide-y">
              {recentDocuments.map((doc) => (
                <div
                  key={doc.id}
                  className="p-4 flex items-center hover:bg-muted/50"
                >
                  <div className="flex items-center gap-3">
                    <File className="size-5 text-muted-foreground" />
                    <div>
                      <p className="font-medium">{doc.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {doc.size} • {getRelativeTimeString(doc.uploadedAt)}
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
  );
} 