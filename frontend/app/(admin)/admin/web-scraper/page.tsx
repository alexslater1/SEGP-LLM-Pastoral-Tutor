"use client";

import { Button } from "@/components/ui/button";
import { Link, CloudUpload, Globe } from "lucide-react";
import { useState } from "react";
import { cn, getRelativeTimeString, truncateUrl } from "@/lib/utils";
import { useRagWebpages, useRagUploadUrl } from "@/hooks/use-rag";
import { toast } from "sonner";

export default function WebScraperPage() {
  const [isDragging, setIsDragging] = useState(false);
  const [url, setUrl] = useState("");
  const { data: webpages } = useRagWebpages(0);
  const uploadUrlMutation = useRagUploadUrl();

  const getRecentWebpages = () => {
    if (!webpages?.data) return [];
    
    const sevenDaysAgo = new Date();
    sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7);
    
    return webpages.data
      .filter(webpage => new Date(webpage.date_uploaded) > sevenDaysAgo)
      .sort((a, b) => new Date(b.date_uploaded).getTime() - new Date(a.date_uploaded).getTime())
      .slice(0, 10); // Cap at 10 webpages
  };

  const recentWebpages = getRecentWebpages();

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
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    // TODO: yeah
  };

  const handleUrlSubmit = async () => {
    const trimmedUrl = url.trim();
    
    if (!trimmedUrl) {
      toast.error("Please enter a URL");
      return;
    }

    try {
      await toast.promise(uploadUrlMutation.mutateAsync(trimmedUrl), {
        loading: "Adding URL. This may take a while...",
        success: "URL added successfully",
        error: "Failed to add URL"
      });
      
      setUrl(""); // Clear input on success
    } catch (error) {
      console.error("URL upload error:", error);
    }
  };

  return (
    <div className="space-y-6 p-6">
      <div>
        <h2 className="text-2xl font-bold text-primary tracking-tight">Web Scraper</h2>
        <p className="text-foreground">
          Scrape webpages for the AI tutor to use as extra knowledge.
        </p>
      </div>

      <div className="flex gap-2">
        <div className="flex-1 relative">
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
            <Link className="size-4" />
          </div>
          <input
            type="url"
            placeholder="Enter webpage URL"
            className="w-full pl-9 p-2 bg-muted/25 dark:bg-muted border rounded-md"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
          />
        </div>
        <Button 
          className="bg-button text-button-foreground hover:bg-button/50"
          onClick={handleUrlSubmit}
          disabled={uploadUrlMutation.isPending}
        >
          {uploadUrlMutation.isPending ? "Adding..." : "Add URL"}
        </Button>
      </div>

      <div>
        <h3 className="text-lg font-semibold text-primary mb-2">Bulk URL Upload</h3>
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
            <p className="text-lg">Drag and drop CSV file here</p>
            <p className="text-muted-foreground">or</p>
            <Button 
              variant="outline"
              className="bg-button text-button-foreground hover:bg-button/50"
            >
              <label className="cursor-pointer">
                <input
                  type="file"
                  className="hidden"
                  onChange={handleFileInput}
                  accept=".csv"
                />
                Browse Files
              </label>
            </Button>
            <p className="text-sm text-muted-foreground">
              Upload a CSV file containing webpage URLs
            </p>
          </div>
        </div>
      </div>

      <div className="border rounded-xl bg-card">
        <div className="overflow-hidden rounded-xl">
          <div className="border-b bg-table-header">
            <h3 className="p-4 font-semibold text-primary">
              Recently Added URLs (last 7 days)
            </h3>
          </div>

          <div className="divide-y divide-border">
            {recentWebpages.length === 0 ? (
              <div className="p-4 text-center bg-table-row-odd text-foreground">
                No URLs added in the last 7 days.
              </div>
            ) : (
              <div className="divide-y divide-border">
                {recentWebpages.map((webpage, index) => (
                  <div
                    key={webpage.id}
                    className={cn(
                      "p-4 flex items-center transition-colors",
                      index % 2 === 0 
                        ? "bg-table-row-odd" 
                        : "bg-table-row-even",
                      "hover:bg-muted/50"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <Globe className="size-5 text-muted-foreground" />
                      <div>
                        <p className="font-medium text-foreground font-mono">{truncateUrl(webpage.url)}</p>
                        <p className="text-sm text-muted-foreground">
                          {getRelativeTimeString(new Date(webpage.date_uploaded))}
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
