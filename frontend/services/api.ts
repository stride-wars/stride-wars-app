import { h3IndexToSplitLong,  latLngToCell} from "h3-js";

interface LocationData {
  latitude: number;
  longitude: number;
  timestamp: number;
}

const res = 9; // the size of hexes

export const sendLocationData = async (location: LocationData): Promise<void> => {
  try {
    const cellId = String(latLngToCell(location.latitude, location.longitude, res));
    const response = await fetch('http://localhost:8080/api/data', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ hex: cellId }),
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