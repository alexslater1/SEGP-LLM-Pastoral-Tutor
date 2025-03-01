import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  RagDocument,
  fetchRagDocuments,
  downloadRagDocument,
  uploadRagDocument,
  deleteRagDocument
} from "@/app/(admin)/actions";

export function useRagDocuments(page: number) {
  return useQuery({
    queryKey: ["rag-documents", page],
    queryFn: () => fetchRagDocuments(page)
  });
};

export function useDownloadRagDoc() {
  return useMutation({
    mutationFn: async (name: string) => {
      const result = await downloadRagDocument(name);

      // Create blob from array buffer
      const blob = new Blob([result.data], { type: result.contentType });

      // Create download link
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = result.fileName;
      document.body.appendChild(a);
      a.click();

      // Cleanup
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
      return result;
    },
    onError: (error: Error) => {
      console.error("Download failed", error);
    },
  });
};

export function useRagUploadDocs() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (files: File[]) => {
      const uploads = await Promise.all(files.map(uploadRagDocument));
      return uploads;
    },
    onMutate: async (files: File[]) => {
      await queryClient.cancelQueries({ queryKey: ["rag-documents"] });
      const previousData = queryClient.getQueryData<{ data: RagDocument[]; totalPages: number; }>(["rag-documents", 0]) || 
        { data: [], totalPages: 1 };
      
      const optimisticDocs: RagDocument[] = files.map(file => ({
        id: `-1`,
        name: file.name,
        url: "",
        type: "DOCUMENT",
        date_uploaded: new Date().toISOString(),
        document_size: file.size,
        document_type: file.name.split(".").pop()?.toUpperCase() || "UNKNOWN",
        user_id: "-1",
        backend_source_id: "-1"
      }));

      queryClient.setQueryData<{ data: RagDocument[]; totalPages: number; }>(
        ["rag-documents", 0],
        old => ({
          data: [...(old?.data ?? []), ...optimisticDocs],
          totalPages: old?.totalPages ?? 1
        })
      );

      return { previousData }
    },
    onError: (err: Error, variables: File[], context?: { previousData: { data: RagDocument[]; totalPages: number; } }) => {
      if (context?.previousData) {
        queryClient.setQueryData(["rag-documents", 0], context.previousData);
      }
      console.error("Upload failed", err);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-documents"] });
    }
  });

  return mutation;
};

export function useDeleteRagDoc() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: deleteRagDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-documents"] });
    },
    onError: (error: Error) => {
      console.error("Delete failed", error);
    },
  });

  return mutation;
};
