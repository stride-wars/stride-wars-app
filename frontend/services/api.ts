import {latLngToCell} from "h3-js";

interface LocationData {
  latitude: number;
  longitude: number;
  timestamp: number;
}

const res = 9; // the size of hexes
const API_BASE = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const sendLocationData = async (location: LocationData): Promise<void> => {
  try {
    const cellId = String(latLngToCell(location.latitude, location.longitude, res));
    console.log('cellId:', cellId, typeof cellId);
    const response = await fetch(`${API_BASE}/data/get`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ id: cellId }),
    });

    if (!response.ok) {
      console.log(response)
      throw new Error('Failed to send location data');
    }
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
};