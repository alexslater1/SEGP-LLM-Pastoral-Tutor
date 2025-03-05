'use client';

import { useState } from "react";
import { PageSelectionBar } from "@/components/page-selection-bar";
import { WebScraper } from "./web-scraper";
import { DocumentUpload } from "./document-upload";

export default function RagUploadPage() {
  const [page, setPage] = useState(0);

  return (
    <>
      <PageSelectionBar pageNumber={page} setPage={setPage} pageNames={["Documents", "Webpages"]} />
      {page === 0 ? <DocumentUpload /> : <WebScraper />}
    </>
  );
} 