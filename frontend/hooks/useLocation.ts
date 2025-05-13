import * as Location from 'expo-location';
import { useState } from 'react';
import { sendLocationData } from '../services/api';

export interface LocationState {
  coords: {
    latitude: number;
    longitude: number;
  };
  timestamp: number;
}

export const useLocation = () => {
  const [location, setLocation] = useState<LocationState | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const getLocation = async () => {
    setIsLoading(true);
    try {
      const { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') throw new Error('Permission denied!');

      const location = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.High,
      });

      setLocation({
        coords: {
          latitude: location.coords.latitude,
          longitude: location.coords.longitude,
        },
        timestamp: location.timestamp,
      });

      await sendLocationData({
        latitude: location.coords.latitude,
        longitude: location.coords.longitude,
        timestamp: location.timestamp
      });

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  return { location, error, isLoading, getLocation };
};