import { useState } from "react";
import { Session } from "@supabase/supabase-js";
import { getUserSession } from "@/lib/supabase/client";
import { useMutation, useQuery } from "@tanstack/react-query";

type AutoUpdateItems = {
  enabled: boolean;
  setEnabled: (enabled: boolean) => void;
  loading: boolean;
  error: Error | null;
};

export function useAutoUpdatePrompts(): AutoUpdateItems {
  const [autoUpdatePrompts, setAutoUpdatePrompts] = useState(false);

  const {isLoading: isAutoUpdateStatusLoading, error: autoUpdateStatusError } = useQuery({
    queryKey: ["autoUpdateStatus"],
    queryFn: async () => {
      const autoUpdateStatus = await getAutoUpdateStatus();
      setAutoUpdatePrompts(autoUpdateStatus);
      return autoUpdateStatus;
    },
    staleTime: Infinity,
  });

  const { mutate: mutateAutoUpdateStatus, isPending: isUpdatingAutoUpdateStatus, error: updateAutoUpdateStatusError } = useMutation({
    mutationFn: async (enabled: boolean) => updateAutoUpdateStatus(enabled),
    onSuccess: (_, enabled) => {
      setAutoUpdatePrompts(enabled);
    },
  });

  const getAutoUpdateStatus = async () => {
    const login_session = await getUserSession();
    if (!login_session) {
      throw new Error("User not logged in");
    } else {
      const { enabled, error: loadError } = await fetchAutoUpdateStatus(
        login_session
      );
      if (loadError) {
        throw new Error(loadError);
      } else {
        return enabled;
      }
    }
  };

  const updateAutoUpdateStatus = async (enabled: boolean) => {
    const login_session = await getUserSession();
    if (!login_session) {
      throw new Error("User not logged in");
    } else {
      const { error: updateError } = await postAutoUpdateStatus(login_session, enabled);
      if (updateError) {
        throw new Error(updateError);
      }
    }
  };

  return { 
    enabled: autoUpdatePrompts, 
    setEnabled: mutateAutoUpdateStatus,
    loading: (isAutoUpdateStatusLoading || isUpdatingAutoUpdateStatus),
    error: (autoUpdateStatusError || updateAutoUpdateStatusError)
  };
}

type AutoUpdateStatus = {
  enabled: boolean;
  error: string | null;
};

async function fetchAutoUpdateStatus(
  session: Session
): Promise<AutoUpdateStatus> {
  try {
    const completionEndpoint = "/agents/auto-update-prompts";
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + completionEndpoint,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${session.access_token}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    let backendResponse =
      (await response.json()) as { enabled: boolean };

    return { enabled: backendResponse.enabled, error: null };
  } catch (currentError) {
    console.error("Error fetching auto update status:", (currentError as Error).message);
    return { enabled: false, error: (currentError as Error).message };
  }
}

async function postAutoUpdateStatus(
  session: Session,
  enabled: boolean
): Promise<{ error: string | null }> {
  try {
    const completionEndpoint = "/agents/auto-update-prompts";
    const response = await fetch(
      process.env.NEXT_PUBLIC_BACKEND_AGENT_URL + completionEndpoint,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({ enabled }),
      }
    );

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    let backendResponse =
      (await response.json()) as { enabled: boolean };

    return { error: null };
  } catch (currentError) {
    console.error("Error updating auto update status:", (currentError as Error).message);
    return { error: (currentError as Error).message };
  }
}

