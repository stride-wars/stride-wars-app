interface LocationData {
  latitude: number;
  longitude: number;
  timestamp: number;
}

export const sendLocationData = async (location: LocationData): Promise<void> => {
  try {
    // TODO: zrobic komunikacje z backendem
    const response = await fetch('JAKIS_ENDPOINT', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(location),
    });

    if (!response.ok) {
      throw new Error('Failed to send location data');
    }
  } catch (error) {
    console.error('API Error:', error);
    throw error;
  }
};