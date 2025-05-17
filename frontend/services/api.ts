import { h3IndexToSplitLong,  latLngToCell} from "h3-js";

interface LocationData {
  latitude: number;
  longitude: number;
  timestamp: number;
}

const res = 10;

export const sendLocationData = async (location: LocationData): Promise<void> => {
  try {
    // TODO: zrobic komunikacje z backendem
    const response = await fetch('/api/data', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(latLngToCell(location.latitude, location.longitude, res)),
    });

    if (!response.ok) {
      throw new Error('Failed to send location data');
    }
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
};