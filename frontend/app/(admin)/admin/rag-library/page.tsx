"use client";

import { useState } from "react";
import { FileText, Globe } from "lucide-react";
import { FilterBar } from "@/components/filter-bar";
import { PagedRagDocumentsTable } from "./paged-rag-documents-table";
import { PagedRagWebpagesTable } from "./paged-rag-webpages-table";

type RagFilter = "documents" | "webpages";

const ragFilterOptions = [
  {
    id: "documents",
    label: "Documents",
    icon: FileText,
    theme: "default" as const
  },
  {
    id: "webpages",
    label: "Webpages",
    icon: Globe,
    theme: "default" as const
  }
];

export default function RagLibraryPage() {
  const [filter, setFilter] = useState<RagFilter>("documents");

  return (
    <>
      <FilterBar
        filter={filter}
        setFilter={setFilter}
        options={ragFilterOptions}
      />
      {filter === "documents" ? (
        <PagedRagDocumentsTable />
      ) : (
        <PagedRagWebpagesTable />
      )}
    </>
  );
}
