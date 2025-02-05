import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { RagDocument } from "@/lib/db/schema";

function humanReadableSize(sizeInBytes: number): string {
  const units = ["B", "KB", "MB", "GB", "TB", "PB"];
  let size = sizeInBytes;
  let unitIndex = 0;
  while (size >= 1024 && unitIndex < units.length - 1) {
      size /= 1024;
      unitIndex++;
  }
  return `${size.toFixed(2)} ${units[unitIndex]}`;
}

const uploadRagDocument = async (file: File) => {
  try {
    console.log('Starting upload for:', file.name);
    const formData = new FormData();
    formData.append('file', file);
    
    // Use relative path to Next.js API route
    const response = await fetch('http://127.0.0.1:8000/rag-doc', {
      method: 'POST',
      body: formData,
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Upload failed with status: ${response.status}`);
    }

    const data = await response.json();
    console.log('Response data:', data);
    return data;
  } catch (error) {
    console.error('Upload error details:', error);
    throw error;
  }
};

type RagDocumentResponse = {
    documents: Document[];
}

type Document = {
  id: number;
  name: string;
  date_uploaded: string;
  document_size: number;
  document_type: string;
  user_id?: number;
  backend_source_id?: number;
};

const fetchRagDocuments = async (): Promise<RagDocument[]> => {
  console.log("fetching docs");
  const response = await fetch('http://127.0.0.1:8000/rag-doc', {
    method: 'GET',
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error(`Failed to fetch RAG documents: ${response.status}`);
  }

  const resp = await response.json() as RagDocumentResponse;
 
  console.log("RESP:", resp);
  // Map over the array of documents from the API
 
  const ragDocs: RagDocument[] = resp.documents.map((doc: Document) => ({
      id: doc.id.toString(),
      name: doc.name,
      uploadedAt: new Date(doc.date_uploaded),
      size: humanReadableSize(doc.document_size),
      type: doc.document_type as "PDF" | "DOCX" | "TXT" | "PPTX",
      userId: doc.user_id?.toString() ?? '-1',
      backendSourceId: doc.backend_source_id ?? -1
  }));

  console.log("RagDocs:", ragDocs);
  return ragDocs;
};

export function useRagUploadDocs() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (files: File[]) => {
      const uploads = await Promise.all(files.map(uploadRagDocument));
      return uploads;
    },
    onSuccess: () => {
      // Invalidate and refetch
      queryClient.invalidateQueries({ queryKey: ['rag-documents'] });
    },
    onError: (error: Error) => {
      console.error("Upload failed", error);
    },
  });

  return mutation;
}

export function useRagDocuments() {
  return useQuery({
    queryKey: ['rag-documents'],
    queryFn: fetchRagDocuments,
    staleTime: 1000 * 60,
  });
}

const deleteRagDocument = async (name: string) => {
  try {
    console.log('Starting delete for:', name);
    
    const response = await fetch(`http://127.0.0.1:8000/rag-doc?name=${encodeURIComponent(name)}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error(`Delete failed with status: ${response.status}`);
    }
  } catch (error) {
    console.error('Delete error details:', error);
    throw error;
  }
};


export function useDeleteRagDoc() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: deleteRagDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['rag-documents'] });
    },
    onError: (error: Error) => {
      console.error("Delete failed", error);
    },
  });

  return mutation;
}