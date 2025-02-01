'use client';

import { Button } from "@/components/ui/button";
import { File, CloudUpload } from 'lucide-react';
import { useState } from "react";
import { toast } from "sonner";

interface Document {
  id: string;
  name: string;
  uploadedAt: Date;
  size: string;
}

const getRelativeTimeString = (date: Date) => {
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);
  
  // A few seconds ago
  if (diffInSeconds < 60) {
    return 'Uploaded a few seconds ago';
  }
  
  // Minutes ago
  const diffInMinutes = Math.floor(diffInSeconds / 60);
  if (diffInMinutes < 60) {
    return `Uploaded ${diffInMinutes} ${diffInMinutes === 1 ? 'minute' : 'minutes'} ago`;
  }
  
  // Hours ago
  const diffInHours = Math.floor(diffInMinutes / 60);
  if (diffInHours < 24) {
    return `Uploaded ${diffInHours} ${diffInHours === 1 ? 'hour' : 'hours'} ago`;
  }
  
  // Days ago
  const diffInDays = Math.floor(diffInHours / 24);
  return `Uploaded ${diffInDays} ${diffInDays === 1 ? 'day' : 'days'} ago`;
};

const getFileExtension = (filename: string) => {
  return filename.slice((filename.lastIndexOf(".") - 1 >>> 0) + 2);
};

export default function DocumentsPage() {
  const [documents, setDocuments] = useState<Document[]>([
    {
      id: '1',
      name: 'dummy-doc.pdf',
      uploadedAt: new Date(),
      size: '2.4 MB'
    }
  ]);
  const [isDragging, setIsDragging] = useState(false);

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
      'text/plain'
    ];
    return supportedTypes.includes(file.type);
  };

  const handleFiles = (files: File[]) => {
    const invalidFiles = files.filter(file => !isValidFileType(file));
    
    if (invalidFiles.length > 0) {
      toast.error(
        <div>
          <p>
            Unsupported file format{invalidFiles.length > 1 ? 's' : ''}:   
            {invalidFiles.map(f => ` ${getFileExtension(f.name).toUpperCase()}`).join(', ')}
          </p>
          <p>Only PDF, DOCX, and TXT files are supported.</p>
        </div>
      );
      return;
    }

    // For now, just log the files and add them to the documents list
    console.log('Files received:', files);
    
    const newDocuments = files.map(file => ({
      id: Math.random().toString(36).slice(2, 11),
      name: file.name,
      uploadedAt: new Date(),
      size: `${(file.size / (1024 * 1024)).toFixed(2)} MB`
    }));

    setDocuments(prev => [...prev, ...newDocuments]);
  };

  const getRecentDocuments = () => {
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
    
    return documents
      .filter(doc => doc.uploadedAt > sevenDaysAgo)
      .sort((a, b) => b.uploadedAt.getTime() - a.uploadedAt.getTime());
  };

  const recentDocuments = getRecentDocuments();

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Upload Documents</h2>
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
          <CloudUpload className={`h-12 w-12 ${isDragging ? 'text-primary' : 'text-muted-foreground'}`} />
          <p className="text-lg">Drag and drop files here</p>
          <p className="text-muted-foreground">or</p>
          <Button>
            <label className="cursor-pointer">
              <input
                type="file"
                className="hidden"
                multiple
                onChange={handleFileInput}
                accept=".pdf,.docx,.txt"
              />
              Browse Files
            </label>
          </Button>
          <p className="text-sm text-muted-foreground">
            Supported formats: PDF, DOCX, TXT
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
                    <File className="h-5 w-5 text-muted-foreground" />
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