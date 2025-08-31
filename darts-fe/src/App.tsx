import { useEffect, useState } from "react";
import { createApi } from "./lib/api";
import { usePlayers } from "./hooks/usePlayers";
import { useMatches } from "./hooks/useMatches";
import ErrorBanner from "./components/common/ErrorBanner";
import BottomNav from "./components/common/BottomNav";
import DebugDrawer from "./components/common/DebugDrawer";

import StartScreen from "./components/pages/StartScreen";
import PlayersScreen from "./components/pages/PlayersScreen";
import StatsScreen from "./components/pages/StatsScreen";
import MatchesScreen from "./components/pages/MatchesScreen";
import SettingsScreen from "./components/pages/SettingsScreen";

const api = createApi();

export default function App() {
  const apiState = (api as any).use(); // subscribe to error/logs
  const playersState = usePlayers(api);
  const matchesState = useMatches(api);

  const [tab, setTab] = useState("start");

  useEffect(() => {
    document.body.classList.add("bg-black", "text-white");
  }, []);

  return (
    <div className="max-w-md mx-auto min-h-screen pb-24">
      <ErrorBanner message={apiState.lastError} onClose={apiState.clearError} />

      {tab === "start" && (
        <StartScreen
          api={api}
          playersState={playersState}
          matchesState={matchesState}
        />
      )}
      {tab === "players" && (
        <PlayersScreen api={api} playersState={playersState} />
      )}
      {tab === "stats" && <StatsScreen players={playersState.players} />}
      {tab === "matches" && (
        <MatchesScreen
          api={api}
          players={playersState.players}
          matchesState={matchesState}
        />
      )}
      {tab === "settings" && <SettingsScreen />}

      <BottomNav tab={tab} setTab={setTab} />
      <DebugDrawer logs={apiState.logs} />
    </div>
  );
}
