import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  RagDocument,
  fetchRagDocuments,
  downloadRagDocument,
  uploadRagDocument,
  deleteRagDocument,
  fetchRagWebpages,
  uploadRagUrl,
  deleteRagUrl,
  fetchRagDocumentsCount,
  fetchRagWebpagesCount,
  fetchAgentRequestsCount,
  fetchAgentEventsCount,
  fetchDownvotesCount,
  fetchUpvotesCount,
  deleteRagSource,
} from "@/app/(admin)/actions";

type CountType =
  | "rag-documents"
  | "rag-webpages"
  | "agent-requests"
  | "agent-events";

const countFunctions = {
  "rag-documents": fetchRagDocumentsCount,
  "rag-webpages": fetchRagWebpagesCount,
  "agent-requests": fetchAgentRequestsCount,
  "agent-events": fetchAgentEventsCount,
} as const;

export function useCount(type: CountType, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: [`${type}-count`],
    queryFn: countFunctions[type],
    ...options,
  });
}

export function useRagDocuments(page: number) {
  return useQuery({
    queryKey: ["rag-documents", page],
    queryFn: () => fetchRagDocuments(page),
  });
}

export function useRagWebpages(page: number) {
  return useQuery({
    queryKey: ["rag-webpages", page],
    queryFn: () => fetchRagWebpages(page),
  });
}

export function useRagDocumentsCount(options?: { enabled?: boolean }) {
  return useCount("rag-documents", options);
}

export function useRagWebpagesCount(options?: { enabled?: boolean }) {
  return useCount("rag-webpages", options);
}

export function useAgentRequestsCount(options?: { enabled?: boolean }) {
  return useCount("agent-requests", options);
}

export function useAgentEventsCount(options?: { enabled?: boolean }) {
  return useCount("agent-events", options);
}

export function useDownvotesCount(options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["downvotes-count"],
    queryFn: () => fetchDownvotesCount(),
    ...options,
  });
}

export function useUpvotesCount(options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ["upvotes-count"],
    queryFn: () => fetchUpvotesCount(),
    ...options,
  });
}

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
}

export function useRagUploadDocs() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (files: File[]) => {
      const uploads = await Promise.all(files.map(uploadRagDocument));
      return uploads;
    },
    onMutate: async (files: File[]) => {
      await queryClient.cancelQueries({ queryKey: ["rag-documents"] });
      const previousData = queryClient.getQueryData<{
        data: RagDocument[];
        totalPages: number;
      }>(["rag-documents", 0]) || { data: [], totalPages: 1 };

      const optimisticDocs: RagDocument[] = files.map((file) => ({
        id: `-1`,
        name: file.name,
        url: "",
        type: "DOCUMENT",
        date_uploaded: new Date().toISOString(),
        document_size: file.size,
        document_type: file.name.split(".").pop()?.toUpperCase() || "UNKNOWN",
        user_id: "-1",
        backend_source_id: "-1",
      }));

      queryClient.setQueryData<{ data: RagDocument[]; totalPages: number }>(
        ["rag-documents", 0],
        (old) => ({
          data: [...(old?.data ?? []), ...optimisticDocs],
          totalPages: old?.totalPages ?? 1,
        })
      );

      return { previousData };
    },
    onError: (
      err: Error,
      variables: File[],
      context?: { previousData: { data: RagDocument[]; totalPages: number } }
    ) => {
      if (context?.previousData) {
        queryClient.setQueryData(["rag-documents", 0], context.previousData);
      }
      console.error("Upload failed", err);
    },
    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-documents"] });
    },
  });

  return mutation;
}

export function useDeleteRagSource() {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: deleteRagSource,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-documents"] });
      queryClient.invalidateQueries({ queryKey: ["rag-webpages"] });
    },
    onError: (error: Error) => {
      console.error("Delete failed", error);
    },
  });

  return mutation;
}

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
}

export function useRagUploadUrl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: uploadRagUrl,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-webpages"] });
    },
    onError: (error: Error) => {
      console.error("URL upload failed", error);
    },
  });
}

export function useDeleteRagUrl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteRagUrl,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["rag-webpages"] });
    },
    onError: (error: Error) => {
      console.error("URL delete failed", error);
    },
  });
}
