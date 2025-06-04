import React, { createContext, useContext, useEffect, useState } from 'react';
import { api } from '../api';
import { GlobalLeaderboardEntry } from '../consts/types';

type LeaderboardContextType = {
  leaderboard: GlobalLeaderboardEntry[] | null;
  loading: boolean;
  error: string | null;
  refresh: () => void;
};

const LeaderboardContext = createContext<LeaderboardContextType | undefined>(undefined);

export const LeaderboardProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [leaderboard, setLeaderboard] = useState<GlobalLeaderboardEntry[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLeaderboard = async () => {
    setLoading(true);
    setError(null);
    const res = await api.getGlobalLeaderboard();
    if (res.data) setLeaderboard(res.data);
    else setError(res.error || 'Failed to load leaderboard');
    setLoading(false);
  };

  useEffect(() => {
    fetchLeaderboard();
  }, []);

  return (
    <LeaderboardContext.Provider value={{
      leaderboard,
      loading,
      error,
      refresh: fetchLeaderboard,
    }}>
      {children}
    </LeaderboardContext.Provider>
  );
};

export const useGlobalLeaderboard = () => {
  const ctx = useContext(LeaderboardContext);
  if (!ctx) throw new Error('useGlobalLeaderboard must be used within a LeaderboardProvider');
  return ctx;
};