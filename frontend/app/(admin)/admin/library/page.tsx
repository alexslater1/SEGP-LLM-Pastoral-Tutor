'use client';

import { useState } from "react";
import { PagedRagDocumentsTable } from "./paged-rag-documents-table";
import { PagedRagWebpagesTable } from "./paged-rag-webpages-table";
import { PageSelectionBar } from "@/components/page-selection-bar";

export default function LibraryPage() {
  return <PagedRagSourcesTable />;
} 

function PagedRagSourcesTable() {
  const [page, setPage] = useState(0);

  return (
    <>
      <PageSelectionBar pageNumber={page} setPage={setPage} pageNames={["Documents", "Webpages"]} />
      {page === 0 ? <PagedRagDocumentsTable /> : <PagedRagWebpagesTable />}
    </>
  );
}