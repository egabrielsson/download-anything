import { useState } from "react";
import { useDownload } from "./hooks/useDownload";
import { Platform, DownloadType } from "./services/api";

interface DownloadSectionProps {
  title: string;
  placeholder: string;
  platform: Platform;
  type: DownloadType;
  enabled: boolean;
}

function DownloadSection({
  title,
  placeholder,
  platform,
  type,
  enabled,
}: DownloadSectionProps) {
  const [url, setUrl] = useState("");
  const { download, loading, error, clearError } = useDownload();

  const handleDownload = async () => {
    if (!enabled) {
      alert("Coming soon!");
      return;
    }
    await download(url, platform, type);
    if (!error) {
      setUrl("");
    }
  };

  return (
    <div className="border-2 border-[#808080] bg-[#c0c0c0] p-2">
      <div className="bg-gradient-to-r from-[#000080] to-[#1084d0] px-2 py-1 mb-2">
        <span className="text-white text-sm font-bold">{title}</span>
      </div>
      <input
        type="text"
        value={url}
        onChange={(e) => {
          setUrl(e.target.value);
          clearError();
        }}
        placeholder={placeholder}
        className="w-full px-2 py-1 border-2 border-t-[#808080] border-l-[#808080] border-b-white border-r-white bg-white text-black text-sm font-mono mb-2"
        onKeyDown={(e) => e.key === "Enter" && handleDownload()}
        disabled={loading}
      />
      {error && (
        <div className="bg-[#ffcccc] border border-[#cc0000] px-2 py-1 mb-2 text-xs text-[#cc0000]">
          {error}
        </div>
      )}
      <button
        onClick={handleDownload}
        disabled={loading}
        className={`w-full px-4 py-1 bg-[#c0c0c0] border-2 border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-black text-sm font-bold active:border-t-[#808080] active:border-l-[#808080] active:border-b-white active:border-r-white ${
          !enabled ? "opacity-50" : ""
        } ${loading ? "cursor-wait" : ""}`}
      >
        {loading ? "Downloading..." : "Download"}
      </button>
    </div>
  );
}

interface PlatformColumnProps {
  platform: Platform;
  displayName: string;
  color: string;
  enabled: boolean;
}

function PlatformColumn({
  platform,
  displayName,
  color,
  enabled,
}: PlatformColumnProps) {
  return (
    <div className="border-2 border-t-white border-l-white border-b-[#808080] border-r-[#808080] bg-[#c0c0c0]">
      {/* Title Bar */}
      <div className={`${color} px-2 py-1 flex items-center justify-between`}>
        <span className="text-white text-sm font-bold">
          {displayName} {!enabled && "(Coming Soon)"}
        </span>
        <div className="flex gap-1">
          <button className="w-4 h-4 bg-[#c0c0c0] border border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-[10px] leading-none">
            _
          </button>
          <button className="w-4 h-4 bg-[#c0c0c0] border border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-[10px] leading-none">
            x
          </button>
        </div>
      </div>

      {/* Content */}
      <div className="p-3 space-y-3">
        <DownloadSection
          title="Download Audio"
          placeholder="Paste URL here..."
          platform={platform}
          type="audio"
          enabled={enabled}
        />
        <DownloadSection
          title="Download Video"
          placeholder="Paste URL here..."
          platform={platform}
          type="video"
          enabled={enabled}
        />
      </div>
    </div>
  );
}

function App() {
  return (
    <div className="min-h-screen bg-[#008080] p-4">
      {/* Main Window */}
      <div className="max-w-6xl mx-auto border-2 border-t-white border-l-white border-b-[#808080] border-r-[#808080] bg-[#c0c0c0]">
        {/* Window Title Bar */}
        <div className="bg-gradient-to-r from-[#000080] to-[#1084d0] px-2 py-1 flex items-center justify-between">
          <span className="text-white font-bold">Download Anything</span>
          <div className="flex gap-1">
            <button className="w-5 h-5 bg-[#c0c0c0] border border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-xs leading-none">
              _
            </button>
            <button className="w-5 h-5 bg-[#c0c0c0] border border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-xs leading-none">
              []
            </button>
            <button className="w-5 h-5 bg-[#c0c0c0] border border-t-white border-l-white border-b-[#808080] border-r-[#808080] text-xs leading-none font-bold">
              X
            </button>
          </div>
        </div>

        {/* Menu Bar */}
        <div className="bg-[#c0c0c0] border-b border-[#808080] px-1 py-1">
          <span className="px-2 text-sm">File</span>
          <span className="px-2 text-sm">Edit</span>
          <span className="px-2 text-sm">View</span>
          <span className="px-2 text-sm">Help</span>
        </div>

        {/* Content Area */}
        <div className="p-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <PlatformColumn
              platform="youtube"
              displayName="YouTube"
              color="bg-[#cc0000]"
              enabled={true}
            />
            <PlatformColumn
              platform="instagram"
              displayName="Instagram"
              color="bg-gradient-to-r from-[#833ab4] via-[#fd1d1d] to-[#fcb045]"
              enabled={true}
            />
            <PlatformColumn
              platform="tiktok"
              displayName="TikTok"
              color="bg-black"
              enabled={true}
            />
          </div>
        </div>

        {/* Status Bar */}
        <div className="border-t-2 border-[#808080] bg-[#c0c0c0] px-2 py-1 flex">
          <div className="border border-t-[#808080] border-l-[#808080] border-b-white border-r-white px-2 flex-1">
            <span className="text-sm">Ready</span>
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
