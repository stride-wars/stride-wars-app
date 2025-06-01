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

export function getHexagonColor(hexId: string): string {
  // Simple hash function from hexId → 0–360
  let hash = 0;
  for (let i = 0; i < hexId.length; i++) {
    hash = (hash << 5) - hash + hexId.charCodeAt(i);
    hash |= 0; // Convert to 32bit int
  }

  const hue = Math.abs(hash) % 360; // Map to hue value
  return `hsla(${hue}, 70%, 50%, 0.3)`; // 30% transparent HSL color
}