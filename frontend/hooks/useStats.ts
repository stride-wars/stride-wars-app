import { useState, useEffect, useCallback } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { api } from '@/api';
import type { ApiResponse, GetActivityStatsResponse } from '../consts/types';

export function useStats() {
  const [username, setUsername] = useState<string>('');
  const [hexesVisited, setHexesVisited] = useState<number>(0);
  const [activitiesRecorded, setActivitiesRecorded] = useState<number>(0);
  const [distanceCovered, setDistanceCovered] = useState<number>(0);
  const [weeklyActivities, setWeeklyActivities] = useState<number[]>([0, 0, 0, 0, 0, 0, 0]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  // fetchStats is memoized with useCallback so that refresh has a stable reference
  const fetchStats = useCallback(async () => {
    setLoading(true);
    setError('');

    try {
      const storedUser = await AsyncStorage.getItem('user');
      if (!storedUser) {
        setError('User not found in storage.');
        setLoading(false);
        return;
      }

      const parsedUser: { id: string; username: string } = JSON.parse(storedUser);
      setUsername(parsedUser.username);

      const result: ApiResponse<GetActivityStatsResponse> = await api.getUserActivityStats(parsedUser.id);

      if (result.error) {
        setError(result.error);
      } else if (result.data) {
        setHexesVisited(result.data.hexes_visited);
        setActivitiesRecorded(result.data.activities_recorded);
        setDistanceCovered(result.data.distance_covered);
        if (result.data.weekly_activities.length === 7) {
          setWeeklyActivities(result.data.weekly_activities);
        } else {
          setWeeklyActivities([0, 0, 0, 0, 0, 0, 0]);
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch activity stats.');
    } finally {
      setLoading(false);
    }
  }, []);


  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  return {
    username,
    hexesVisited,
    activitiesRecorded,
    distanceCovered,
    weeklyActivities,
    loading,
    error,
    refresh: fetchStats, 
  };
}
