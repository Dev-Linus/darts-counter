export interface Player {
  id: string;
  name: string;
  matches: number;
  throws: number;
  totalScore: number;
}

export interface Match {
  id: string;
  players: string[];
  currentThrow: number;
  currentPlayer: string;
  startAt: number;
  startMode: number; 
  endMode: number;
  scores: Record<string, number>;
}

export interface ApiLog {
  time: string;
  request: { method: string; url: string; body?: any };
  response?: { status: number; body?: any };
  error?: string;
}