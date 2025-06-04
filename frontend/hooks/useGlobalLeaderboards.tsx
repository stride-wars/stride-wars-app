import { useEffect, useState, useCallback } from 'react';
import { api } from '../api';
import { GlobalLeaderboardEntry } from '../consts/types';

export function useGlobalLeaderboard() {
  const [leaderboard, setLeaderboard] = useState<GlobalLeaderboardEntry[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLeaderboard = useCallback(async () => {
    setLoading(true);
    setError(null);
    const res = await api.getGlobalLeaderboard();
    if (res.data) setLeaderboard(res.data);
    else setError(res.error || 'Failed to load leaderboard');
    setLoading(false);
  }, []);

  useEffect(() => {
    fetchLeaderboard();
  }, [fetchLeaderboard]);

  return { leaderboard, loading, error, refresh: fetchLeaderboard };
} 