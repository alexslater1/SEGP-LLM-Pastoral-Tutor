"use client";

import { useState } from "react";
import { PageSelectionBar } from "@/components/page-selection-bar";
import { PagedRagDocumentsTable } from "./paged-rag-documents-table";
import { PagedRagWebpagesTable } from "./paged-rag-webpages-table";

export default function RagLibraryPage() {
  const [page, setPage] = useState(0);

  return (
    <>
      <PageSelectionBar pageNumber={page} setPage={setPage} pageNames={["Documents", "Webpages"]} />
      {page === 0 ? <PagedRagDocumentsTable /> : <PagedRagWebpagesTable />}
    </>
  );
}
