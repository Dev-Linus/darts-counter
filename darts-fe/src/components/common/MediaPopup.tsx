import React, { useEffect, useRef, useState } from "react";

type MediaPopupProps = {
  src: string; // API endpoint to stream (e.g., http://localhost:8080/streamFile?file=video/myvideo.mp4)
  type: "audio" | "video"; // Type of media
  onClose?: () => void; // Optional callback when finished
};

const MediaPopup: React.FC<MediaPopupProps> = ({ src, type, onClose }) => {
  const [visible, setVisible] = useState(true);
  const mediaRef = useRef<HTMLVideoElement & HTMLAudioElement>(null);

  useEffect(() => {
    const media = mediaRef.current;
    if (!media) return;

    // Auto play
    media.play().catch((err) => {
      console.error("Autoplay failed:", err);
    });

    // Cleanup on end
    const handleEnd = () => {
      setVisible(false);
      if (onClose) onClose();
    };

    media.addEventListener("ended", handleEnd);
    return () => {
      media.removeEventListener("ended", handleEnd);
    };
  }, [onClose]);

  if (!visible) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-60 z-50">
      <div className="bg-white p-4 rounded-2xl shadow-lg">
        {type === "video" ? (
          <video
            ref={mediaRef}
            src={src}
            autoPlay
            playsInline
            className="max-w-[80vw] max-h-[80vh] rounded-xl"
          />
        ) : (
          <audio ref={mediaRef} src={src} autoPlay />
        )}
      </div>
    </div>
  );
};

export default MediaPopup;
