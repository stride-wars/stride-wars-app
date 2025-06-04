import React, { useEffect, useState, useRef } from 'react';
import {
  StyleSheet,
  View,
  Text,
  TouchableOpacity,
  Platform,
  Animated,
  Dimensions,
} from 'react-native';
import MapView, { Polygon, PROVIDER_DEFAULT, Region } from 'react-native-maps';
import { MaterialIcons } from '@expo/vector-icons';
import { getHexagonColor } from '../../utils/h3Utils';
import { useLocation } from '../../hooks/useLocation';
import * as h3 from 'h3-js';
import * as Location from 'expo-location';
import AsyncStorage from '@react-native-async-storage/async-storage';

// Hex resolution
const res = 9;

// Fallback API base (Expo will replace EXPO_PUBLIC_API_URL at build time)
const API_BASE = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

type Coordinate = { latitude: number; longitude: number };
type LeaderboardEntry = { name: string; points: number };

// Utility: interpolate between two polygons (arrays of coords), based on t ∈ [0..1]
function interpolatePolygon(from: Coordinate[], to: Coordinate[], t: number): Coordinate[] {
  return from.map((point, i) => ({
    latitude: point.latitude + (to[i].latitude - point.latitude) * t,
    longitude: point.longitude + (to[i].longitude - point.longitude) * t,
  }));
}

// Utility: scale a polygon (all points) around its centroid by `scale` factor
function scalePolygon(coordinates: Coordinate[], scale: number): Coordinate[] {
  const latAvg = coordinates.reduce((sum, p) => sum + p.latitude, 0) / coordinates.length;
  const lngAvg = coordinates.reduce((sum, p) => sum + p.longitude, 0) / coordinates.length;
  return coordinates.map(p => ({
    latitude: latAvg + (p.latitude - latAvg) * scale,
    longitude: lngAvg + (p.longitude - lngAvg) * scale,
  }));
}

function calculateDistance(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number {
  const R = 6371; // km
  const dLat = deg2rad(lat2 - lat1);
  const dLon = deg2rad(lon2 - lon1);
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(deg2rad(lat1)) * Math.cos(deg2rad(lat2)) *
    Math.sin(dLon / 2) * Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  const d = R * c; // km
  return d * 1000; // meters
}

function deg2rad(deg: number): number {
  return deg * (Math.PI / 180);
}

export const getCurrentUser = async () => {
  const userString = await AsyncStorage.getItem('user');
  if (!userString) return null;
  return JSON.parse(userString);
};

const SCREEN_WIDTH = Dimensions.get('window').width;

export default function MapScreen() {
  const { location, error, isLoading, getLocation } = useLocation();

  const [selectedHexId, setSelectedHexId] = useState<string | null>(null);

  const [hexagons, setHexagons] = useState<
    Array<{ hexId: string; coordinates: Coordinate[]; animatedCoordinates: Coordinate[] }>
  >([]);

  const [leaderboardData, setLeaderboardData] = useState<Record<string, { user: string; score: number }[]>>({});

  const mapRef = useRef<MapView>(null);

  const [isRecording, setIsRecording] = useState(false);
  const [elapsedTime, setElapsedTime] = useState(0); // in seconds
  const [distanceTraveled, setDistanceTraveled] = useState(0); // in meters

  const [previousLocation, setPreviousLocation] = useState<Coordinate | null>(null);

 const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  const leaderboardAnim = useRef(new Animated.Value(0)).current;

  const [visitedHexIds, setVisitedHexIds] = useState<Set<string>>(new Set());

  const lastRegionRef = useRef<Region | null>(null);

  const [userId, setUserId] = useState<string | null>(null);

  useEffect(() => {
    if (!isRecording || !location) return;

    const nearestHex = String(h3.latLngToCell(location.coords.latitude, location.coords.longitude, res));
    setVisitedHexIds(prev => {
      const updated = new Set(prev);
      updated.add(nearestHex);
      return updated;
    });

    if (previousLocation) {
      const newDistance = calculateDistance(
        previousLocation.latitude,
        previousLocation.longitude,
        location.coords.latitude,
        location.coords.longitude
      );
      setDistanceTraveled(prev => prev + newDistance);
    }

    setPreviousLocation({
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
    });
  }, [location, isRecording]);

  useEffect(() => {
    getLocation();
  }, []);

  useEffect(() => {
    const loadUser = async () => {
      const user = await getCurrentUser();
      if (user?.id) {
        setUserId(user.id);
      }
    };
    loadUser();
  }, []);

  useEffect(() => {
    if (isRecording) {
      const id = setInterval(() => {
        setElapsedTime(prev => prev + 1);
      }, 1000);
      timerRef.current = id;
    }

    return () => {
      if (timerRef.current !== null) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
    };
  }, [isRecording]);

  const fetchLeaderboardDataForBounds = async (
    minLat: number,
    minLng: number,
    maxLat: number,
    maxLng: number
  ) => {
    try {
      const url = `${API_BASE}/leaderboard/bbox?min_lat=${minLat}&min_lng=${minLng}&max_lat=${maxLat}&max_lng=${maxLng}`;
      console.log('Requesting leaderboard:', url);
      const res = await fetch(url);
      if (!res.ok) throw new Error(`HTTP error ${res.status}`);

      const data = await res.json();
      console.log('[DEBUG] leaderboard API raw response:', data);

      const leaderboards = data?.data?.leaderboards ?? [];
      const hexMap: Record<string, { user: string; score: number }[]> = {};
      const newHexagons: typeof hexagons = [];

      for (const leaderboard of leaderboards) {
        const hexId = leaderboard.h3_index;
        if (!hexId) continue;

        // Convert H3 cell → polygon boundary
        const boundary = h3.cellToBoundary(hexId, false);
        const coordinates: Coordinate[] = boundary.map(([lat, lng]: [number, number]) => ({
          latitude: lat,
          longitude: lng,
        }));

        newHexagons.push({
          hexId,
          coordinates,
          animatedCoordinates: coordinates,
        });

        hexMap[hexId] = leaderboard.top_users.map((user: any) => ({
          user: user.user_name ?? user.user_id,
          score: user.score,
        }));
      }

      console.log('Loaded hexagons:', newHexagons.map(v => v.hexId));
      setHexagons(newHexagons);
      setLeaderboardData(hexMap);
    } catch (err) {
      console.error('Failed to fetch leaderboard data:', err);
    }
  };

  // Called whenever the map region changes
  const fetchLeaderboards = async (region: Region) => {
    if (
      lastRegionRef.current &&
      Math.abs(lastRegionRef.current.latitude - region.latitude) < region.latitudeDelta / 10 &&
      Math.abs(lastRegionRef.current.longitude - region.longitude) < region.longitudeDelta / 10 &&
      Math.abs(lastRegionRef.current.latitudeDelta - region.latitudeDelta) < region.latitudeDelta / 5 &&
      Math.abs(lastRegionRef.current.longitudeDelta - region.longitudeDelta) < region.longitudeDelta / 5
    ) {
      return;
    }

    const minLat = region.latitude - region.latitudeDelta / 2;
    const maxLat = region.latitude + region.latitudeDelta / 2;
    const minLng = region.longitude - region.longitudeDelta / 2;
    const maxLng = region.longitude + region.longitudeDelta / 2;

    fetchLeaderboardDataForBounds(minLat, minLng, maxLat, maxLng);
    lastRegionRef.current = region;
  };

  useEffect(() => {
    if (!location) return;

    const initialRegion: Region = {
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
      latitudeDelta: 0.01,
      longitudeDelta: 0.01,
    };

    mapRef.current?.animateToRegion(initialRegion);

    const minLat = initialRegion.latitude - initialRegion.latitudeDelta / 2;
    const maxLat = initialRegion.latitude + initialRegion.latitudeDelta / 2;
    const minLng = initialRegion.longitude - initialRegion.longitudeDelta / 2;
    const maxLng = initialRegion.longitude + initialRegion.longitudeDelta / 2;

    fetchLeaderboardDataForBounds(minLat, minLng, maxLat, maxLng);
  }, [location]);

  const animateHexScaling = (hexId: string, toScale: number, duration: number = 150) => {
    const hex = hexagons.find(h => h.hexId === hexId);
    if (!hex) return;

    const from = hex.animatedCoordinates;
    const to = scalePolygon(hex.coordinates, toScale);
    const steps = 10;
    let currentStep = 0;

    const animateStep = () => {
      currentStep++;
      const t = currentStep / steps;
      const intermediate = interpolatePolygon(from, to, t);

      if (currentStep < steps) {
        setHexagons(prev =>
          prev.map(h =>
            h.hexId === hexId ? { ...h, animatedCoordinates: intermediate } : h
          )
        );
        requestAnimationFrame(animateStep);
      } else {
        setHexagons(prev =>
          prev.map(h =>
            h.hexId === hexId ? { ...h, animatedCoordinates: to } : h
          )
        );
      }
    };
    animateStep();
  };

  const handleStopRecording = async () => {
    setIsRecording(false);

    if (timerRef.current !== null) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }

    const hexesToSend = Array.from(visitedHexIds);
    const durationToSend = elapsedTime;
    const distanceToSend = distanceTraveled;

    if (hexesToSend.length > 0 && userId) {
      try {
        const activityData = {
          user_id: userId,
          h3_indexes: hexesToSend,
          duration: durationToSend,
          distance: distanceToSend + 1, // ensure > 0 for demo
        };

        console.log('Submitting activity with data:', activityData);

        const res = await fetch(`${API_BASE}/activity/create`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(activityData),
        });

        if (!res.ok) {
          const text = await res.text();
          console.error('Failed to submit activity:', text);
          throw new Error(`HTTP error ${res.status}`);
        }

        const json = await res.json();
        console.log('Activity saved:', json);
      } catch (err) {
        console.error('Error submitting activity:', err);
      }
    }

    // Reset for next session
    setVisitedHexIds(new Set());
    setElapsedTime(0);
    setDistanceTraveled(0);
    setPreviousLocation(null);
  };

  // Called when a hexagon on the map is tapped
  const handleHexPress = (hexId: string) => {
    if (selectedHexId === hexId) {
      // Un‐select: scale back down & hide leaderboard
      animateHexScaling(hexId, 1.0);
      Animated.timing(leaderboardAnim, {
        toValue: 0,
        duration: 200,
        useNativeDriver: true,
      }).start(() => setSelectedHexId(null));
    } else {
      // If another hex was selected, shrink it back first
      if (selectedHexId) {
        animateHexScaling(selectedHexId, 1.0);
      }
      // Enlarge the newly tapped hex
      animateHexScaling(hexId, 1.05);
      setSelectedHexId(hexId);
      Animated.timing(leaderboardAnim, {
        toValue: 1,
        duration: 200,
        useNativeDriver: true,
      }).start();
    }
  };

  const leaderboardTranslateY = leaderboardAnim.interpolate({
    inputRange: [0, 1],
    outputRange: [200, 0],
  });

  if (error) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>{error}</Text>
      </View>
    );
  }

  if (isLoading || !location) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>Loading location...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <MapView
        ref={mapRef}
        style={styles.map}
        provider={PROVIDER_DEFAULT}
        initialRegion={{
          latitude: location.coords.latitude,
          longitude: location.coords.longitude,
          latitudeDelta: 0.01,
          longitudeDelta: 0.01,
        }}
        onRegionChangeComplete={fetchLeaderboards}
        showsUserLocation
        showsMyLocationButton={Platform.OS === 'android'}
        showsCompass
        mapType="standard"
      >
        {hexagons.map(hex => (
          <Polygon
            key={hex.hexId}
            coordinates={hex.animatedCoordinates}
            fillColor={getHexagonColor(hex.hexId)}
            strokeColor="rgba(255, 255, 255, 0.4)"
            strokeWidth={1}
            tappable
            onPress={() => handleHexPress(hex.hexId)}
          />
        ))}
      </MapView>

      {selectedHexId && (
        <Animated.View
          style={[
            styles.leaderboard,
            {
              transform: [{ translateY: leaderboardTranslateY }],
              opacity: leaderboardAnim,
            },
          ]}
          pointerEvents="box-none"
        >
          <Text style={styles.leaderboardTitle}>🏆 Leaderboard</Text>
          {(leaderboardData[selectedHexId] || []).map((entry, index) => (
            <Text key={index} style={styles.leaderboardEntry}>
              {index + 1}. {entry.user} – {entry.score} pts
            </Text>
          ))}
          <TouchableOpacity onPress={() => handleHexPress(selectedHexId)} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>Hide</Text>
          </TouchableOpacity>
        </Animated.View>
      )}

      <View style={styles.bottomButtonContainer}>
        {!isRecording ? (
          <TouchableOpacity onPress={() => setIsRecording(true)} style={styles.fullWidthButton}>
            <Text style={{ fontWeight: 'bold', fontSize: 16 }}>Start striding!</Text>
          </TouchableOpacity>
        ) : (
          <View style={styles.timerContainer}>
            {/* Translucent stats container */}
            <View style={styles.statsContainer}>
              <Text style={styles.timerText}>
                {Math.floor(elapsedTime / 60).toString().padStart(2, '0')}:
                {(elapsedTime % 60).toString().padStart(2, '0')}
              </Text>
              <Text style={styles.distanceText}>
                {(distanceTraveled / 1000).toFixed(2)} km
              </Text>
            </View>

            <TouchableOpacity onPress={handleStopRecording} style={styles.stopButton}>
              <Text style={styles.stopButtonText}>Stop</Text>
            </TouchableOpacity>
          </View>
        )}
      </View>

      <TouchableOpacity
        style={styles.locateButton}
        onPress={() => {
          if (location) {
            mapRef.current?.animateToRegion({
              latitude: location.coords.latitude,
              longitude: location.coords.longitude,
              latitudeDelta: 0.01,
              longitudeDelta: 0.01,
            });
          }
        }}
      >
        <MaterialIcons name="my-location" size={24} color="white" />
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#121212',
  },
  map: {
    width: '100%',
    height: '100%',
    flex: 1,
    margin: 0,
    padding: 0,
  },
  errorText: {
    color: '#FFD600',
    fontSize: 16,
    textAlign: 'center',
    marginTop: 20,
  },
  locateButton: {
    position: 'absolute',
    bottom: 80,
    right: 30,
    backgroundColor: 'rgba(0,0,0,0.6)',
    padding: 12,
    borderRadius: 24,
    zIndex: 999,
  },
  activityControls: {
    position: 'absolute',
    bottom: 20,
    left: '59%',
    transform: [{ translateX: -100 }],
    alignItems: 'center',
    gap: 10,
    zIndex: 998,
  },
  startButton: {
    backgroundColor: '#FFD600',
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 30,
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.9,
    shadowRadius: 10,
    elevation: 10,
  },
  startButtonText: {
    color: '#000',
    fontWeight: 'bold',
    fontSize: 16,
  },
  timerContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    width: '100%',
    paddingHorizontal: 20,
  },
  statsContainer: {
    height: '80%',
    backgroundColor: 'rgba(0, 0, 0, 0.5)', // Semi‐transparent black
    borderRadius: 30,
    paddingHorizontal: 15,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 15,
  },
  timerText: {
    color: '#fff',
    fontSize: 12,
    backgroundColor: 'rgba(0,0,0,0.6)',
    paddingHorizontal: 14,
    borderRadius: 20,
    width: '40%',
    textAlign: 'center',
  },
  distanceText: {
    color: '#fff',
    fontSize: 12,
  },
  stopButton: {
    backgroundColor: '#ff5252',
    paddingHorizontal: 20,
    borderRadius: 30,
    height: '80%',
    aspectRatio: 2.5,
    alignItems: 'center',
    justifyContent: 'center',
    marginLeft: 10,
    width: '25%',
    textAlign: 'center',
  },
  fullWidthButton: {
    width: '100%',
    height: '100%',
    backgroundColor: 'transparent',
    alignItems: 'center',
    justifyContent: 'center',
    padding: 0,
    margin: 0,
    borderRadius: 0,
  },
  stopButtonText: {
    color: '#fff',
    fontWeight: 'bold',
    fontSize: 16,
  },
  leaderboard: {
    position: 'absolute',
    top: 80,
    alignSelf: 'center',
    backgroundColor: '#222',
    padding: 20,
    borderRadius: 20,
    width: SCREEN_WIDTH * 0.8,
    zIndex: 1000,
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.8,
    shadowRadius: 10,
    elevation: 10,
  },
  leaderboardTitle: {
    color: '#FFD600',
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 10,
    textAlign: 'center',
  },
  leaderboardEntry: {
    color: '#fff',
    fontSize: 16,
    marginBottom: 4,
  },
  closeButton: {
    marginTop: 10,
    alignSelf: 'center',
    backgroundColor: '#444',
    paddingVertical: 6,
    paddingHorizontal: 16,
    borderRadius: 16,
  },
  bottomButtonContainer: {
    width: '100%',
    backgroundColor: '#FFD600',
    alignItems: 'center',
    justifyContent: 'center',
    height: 56,
    padding: 0,
    margin: 0,
    borderRadius: 0,
    borderTopWidth: 0,
  },
  closeButtonText: {
    color: '#fff',
    fontSize: 14,
  },
});
