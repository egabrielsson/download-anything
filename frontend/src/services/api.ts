const API_BASE = "/api";

export type Platform = "youtube" | "instagram" | "tiktok";
export type DownloadType = "audio" | "video";

interface DownloadOptions {
  url: string;
  platform: Platform;
  type: DownloadType;
}

export async function downloadMedia(options: DownloadOptions): Promise<void> {
  const { url, platform, type } = options;

  // Build endpoint based on platform
  const endpoint = `${API_BASE}/${platform}`;

  const response = await fetch(endpoint, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ url, type }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Download failed");
  }

  // Get the blob and trigger download
  const blob = await response.blob();
  const downloadUrl = window.URL.createObjectURL(blob);

  // Extract filename from Content-Disposition header
  const contentDisposition = response.headers.get("Content-Disposition");
  let filename = `${platform}_${type}.${type === "audio" ? "mp3" : "mp4"}`;
  if (contentDisposition) {
    const match = contentDisposition.match(/filename="(.+)"/);
    if (match) filename = match[1];
  }

  // Create and click download link
  const a = document.createElement("a");
  a.href = downloadUrl;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  window.URL.revokeObjectURL(downloadUrl);
}
