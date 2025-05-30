import * as h3 from 'h3-js';

export const getHexagonsInRadius = (lat: number, lng: number, radiusInMeters: number, resolution: number = 9) => {
  const radiusInDegrees = radiusInMeters / 111000; // 111km per degree

  const centerHex = h3.latLngToCell(lat, lng, resolution);

  const hexagons = h3.gridDisk(centerHex, Math.ceil(radiusInDegrees * 10));

  return hexagons.map(hex => {
    const [hexLat, hexLng] = h3.cellToLatLng(hex);
    return {
      hexId: hex,
      coordinates: h3.cellToBoundary(hex).map(([lat, lng]) => ({
        latitude: lat,
        longitude: lng,
      })),
    };
  });
};

export const getHexagonColor = (hexId: string) => {
  let hash = 0;
  for (let i = 0; i < hexId.length; i++) {
    hash = hexId.charCodeAt(i) + ((hash << 5) - hash);
  }
  
  // Convert to hex color
  const hue = Math.abs(hash % 360);
  return `hsl(${hue}, 70%, 50%)`;
}; 