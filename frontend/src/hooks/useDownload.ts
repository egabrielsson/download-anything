import { useState, useCallback } from "react";
import { downloadMedia, Platform, DownloadType } from "../services/api";

interface UseDownloadReturn {
  download: (
    url: string,
    platform: Platform,
    type: DownloadType,
  ) => Promise<void>;
  loading: boolean;
  error: string | null;
  clearError: () => void;
}

export function useDownload(): UseDownloadReturn {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const download = useCallback(
    async (url: string, platform: Platform, type: DownloadType) => {
      if (!url.trim()) {
        setError("Please enter a URL");
        return;
      }

      setLoading(true);
      setError(null);

      try {
        await downloadMedia({ url: url.trim(), platform, type });
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      } finally {
        setLoading(false);
      }
    },
    [],
  );

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return { download, loading, error, clearError };
}
