"use client";

import { Button } from "@/components/ui/button";
import { Link, CloudUpload, Globe, Download } from "lucide-react";
import { useState } from "react";
import { cn, getRelativeTimeString, truncateUrl } from "@/lib/utils";
import { useRagWebpages, useRagUploadUrl } from "@/hooks/use-rag";
import { toast } from "sonner";

export function WebScraper() {
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

  const handleFiles = async (files: File[]) => {
    const file = files[0];
    if (!file || !file.name.toLowerCase().endsWith('.csv')) {
      toast.error("Please upload a CSV file");
      return;
    }

    try {
      const text = await file.text();
      const lines = text.split(/\r?\n/).filter(line => line.trim() && !line.startsWith('#'));
      
      const urls = lines.map(line => {
        if (line.startsWith('"') && line.endsWith('"')) {
          return line.slice(1, -1);
        }
        return line.split(',')[0];
      }).filter(url => url.trim());

      if (urls.length === 0) {
        toast.error("No valid URLs found in the file");
        return;
      }

      await toast.promise(
        (async () => {
          for (const url of urls) {
            await uploadUrlMutation.mutateAsync(url);
          }
          return urls.length;
        })(),
        {
          loading: "Adding URLs. This may take a while...",
          success: (count) => `Added ${count} URL${count !== 1 ? 's' : ''} successfully`,
          error: "Failed to add URLs"
        }
      );

    } catch (error) {
      console.error("CSV processing error:", error);
      toast.error("Failed to process CSV file");
    }
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
      e.target.value = '';
    }
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
        success: "Added URL successfully",
        error: "Failed to add URL"
      });
      
      setUrl(""); // Clear input on success
    } catch (error) {
      console.error("URL upload error:", error);
    }
  };

  const downloadTemplate = () => {
    const template = `# Instructions:
# 1. Each URL should be on a new line
# 2. Enclose URLs in double quotes if they contain commas
# 3. Empty lines and lines starting with # are ignored

"https://example.com/page1"
"https://example.com/path,with,commas"
"https://example.com/another/page"`;

    const blob = new Blob([template], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'url_template.csv';
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  };

  return (
    <>
      <div className="space-y-6 p-6">
        <div>
          <h2 className="text-2xl font-bold text-primary tracking-tight">Web Scraper</h2>
          <p className="text-foreground">
            Scrape webpages for the AI tutor to use as extra knowledge.
          </p>
        </div>

        <div className="flex gap-2">
          <div className="flex-1 relative group">
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground transition-colors group-focus-within:text-primary">
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
      </div>

      <div className="p-6">
        <div className="flex justify-between items-center mb-1">
          <h3 className="text-lg font-semibold text-primary">Bulk URL Upload</h3>
          <Button
            variant="ghost"
            size="sm"
            className="text-foreground text-lg !hover:bg-muted/50 hover:text-primary flex items-center gap-1.5 [&_svg]:!size-5"
            onClick={downloadTemplate}
          >
            Download Template
            <Download className="size-8" />
          </Button>
        </div>

        <div className="mt-2">
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
      </div>

      <div className="px-6">
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
    </>
  );
}
