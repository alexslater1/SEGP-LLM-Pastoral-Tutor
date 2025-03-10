"use client";

import { useState } from "react";
import { FileText, Globe, RefreshCw } from "lucide-react";
import { FilterBar } from "@/components/filter-bar";
import { DocumentUpload } from "./document-upload";
import { WebScraper } from "./web-scraper";
import { PagedRagDocumentsTable } from "./paged-rag-documents-table";
import { PagedRagWebpagesTable } from "./paged-rag-webpages-table";
import { SwitchCard } from "@/components/switch-card";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";

type RagUploadFilter = "documents" | "webpages";

const ragUploadFilterOptions = [
  {
    id: "documents",
    label: "Upload Documents",
    icon: FileText,
    theme: "default" as const
  },
  {
    id: "webpages",
    label: "Upload Webpages",
    icon: Globe,
    theme: "default" as const
  }
];

export default function RagUploadPage() {
  const [filter, setFilter] = useState<RagUploadFilter>("documents");

  return (
    <>
      
      <FilterBar
        filter={filter}
        setFilter={setFilter}
        options={ragUploadFilterOptions}
      />
      {filter === "documents" ? (
        <>
          <DocumentUpload />
          <PagedRagDocumentsTable />
        </>
      ) : (
        <>
          <WebScraper />
          <PagedRagWebpagesTable />
        </>
      )}
    </>
  );
} 