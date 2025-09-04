import React, { useState } from "react";
import MediaPopup from "./MediaPopup";

const Demo: React.FC = () => {
  const [showVideo, setShowVideo] = useState(false);

  return (
    <div>
      <button
        onClick={() => setShowVideo(true)}
        className="p-2 bg-blue-500 text-white rounded-lg"
      >
        Play Video Popup
      </button>

      {showVideo && (
        <MediaPopup
          src="http://localhost:8080/streamFile?file=video/myvideo.mp4"
          type="video"
          onClose={() => setShowVideo(false)}
        />
      )}
    </div>
  );
};

export default Demo;
