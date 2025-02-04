import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { supabase } from "@/lib/supabase";
import { RagDocument } from "@/lib/db/schema";

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



export function useRagUploadDocs() {
  // const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (files: File[]) => {
      const uploads = await Promise.all(files.map(uploadRagDocument));
      return uploads;
    },

    onError: (error: Error) => {
      console.error("Upload failed", error);
    },
  });

  return mutation;
}

