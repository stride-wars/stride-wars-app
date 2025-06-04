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

const res = 9; // the size of hexes
const API_BASE = 'https://e39c-188-146-191-28.ngrok-free.app/api/v1';

type Coordinate = { latitude: number; longitude: number };
type LeaderboardEntry = { name: string; points: number };

// [utils]
function interpolatePolygon(from: Coordinate[], to: Coordinate[], t: number): Coordinate[] {
  return from.map((point, i) => ({
    latitude: point.latitude + (to[i].latitude - point.latitude) * t,
    longitude: point.longitude + (to[i].longitude - point.longitude) * t,
  }));
}

function scalePolygon(coordinates: Coordinate[], scale: number): Coordinate[] {
  const latAvg = coordinates.reduce((sum, p) => sum + p.latitude, 0) / coordinates.length;
  const lngAvg = coordinates.reduce((sum, p) => sum + p.longitude, 0) / coordinates.length;
  return coordinates.map(p => ({
    latitude: latAvg + (p.latitude - latAvg) * scale,
    longitude: lngAvg + (p.longitude - lngAvg) * scale,
  }));
}

function calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371; // Radius of the earth in km
  const dLat = deg2rad(lat2 - lat1);
  const dLon = deg2rad(lon2 - lon1);
  const a = 
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(deg2rad(lat1)) * Math.cos(deg2rad(lat2)) * 
    Math.sin(dLon/2) * Math.sin(dLon/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a)); 
  const d = R * c; // Distance in km
  return d * 1000; // Convert to meters
}

function deg2rad(deg: number): number {
  return deg * (Math.PI/180);
}

const SCREEN_WIDTH = Dimensions.get('window').width;

  export const getCurrentUser = async () => {
    const userString = await AsyncStorage.getItem('user');
    if (!userString) return null;
    return JSON.parse(userString);
  };

export default function MapScreen() {
  const { location, error, isLoading, getLocation } = useLocation();
  const [selectedHexId, setSelectedHexId] = useState<string | null>(null);
  const [hexagons, setHexagons] = useState<
    Array<{ hexId: string; coordinates: Coordinate[]; animatedCoordinates: Coordinate[] }>
  >([]);
  const [leaderboardData, setLeaderboardData] = useState<Record<string, { user: string; score: number }[]>>({});

  const mapRef = useRef<MapView>(null);
  const [isRecording, setIsRecording] = useState(false);
  const [elapsedTime, setElapsedTime] = useState(0);
  const [distanceTraveled, setDistanceTraveled] = useState(0);
  const [previousLocation, setPreviousLocation] = useState<Location.LocationObject | null>(null);
  const timerRef = useRef<NodeJS.Timer | null>(null);
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

    // Calculate distance if we have a previous location
    if (previousLocation) {
      const newDistance = calculateDistance(
        previousLocation.coords.latitude,
        previousLocation.coords.longitude,
        location.coords.latitude,
        location.coords.longitude
      );
      setDistanceTraveled(prev => prev + newDistance);
    }
    
    // Update previous location
    setPreviousLocation(location);
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

  // Handle timer for recording
  useEffect(() => {
    if (isRecording) {
      timerRef.current = setInterval(() => {
        setElapsedTime(prev => prev + 1);
      }, 1000);
    }

    return () => {
      if (timerRef.current) {
        clearInterval(timerRef.current);
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

        const boundary = h3.cellToBoundary(hexId, false);
        const coordinates = boundary.map(([lat, lng]: [number, number]) => ({
          latitude: lat,
          longitude: lng,
        }));

        newHexagons.push({
          hexId,
          coordinates,
          animatedCoordinates: coordinates,
        });

        hexMap[hexId] = leaderboard.top_users.map((user: any) => ({
          user: user.username ?? user.name ?? user.user_id,
          score: user.score,
        }));
      }
      console.log('hexy', newHexagons.map(v => v.hexId))
      setHexagons(newHexagons);
      setLeaderboardData(hexMap);
    } catch (err) {
      console.error('Failed to fetch leaderboard data:', err);
    }
  };



  // Fetch leaderboard data on region change
  const fetchLeaderboards = async (region: Region) => {
    // Only update if region has changed significantly
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

  // Initialize map when location is available
  useEffect(() => {
    if (!location) return;
    
    const initialRegion = {
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
      latitudeDelta: 0.01,
      longitudeDelta: 0.01,
    };
    
    mapRef.current?.animateToRegion(initialRegion);
    
    // Fetch initial leaderboard data
    const minLat = initialRegion.latitude - initialRegion.latitudeDelta / 2;
    const maxLat = initialRegion.latitude + initialRegion.latitudeDelta / 2;
    const minLng = initialRegion.longitude - initialRegion.longitudeDelta / 2;
    const maxLng = initialRegion.longitude + initialRegion.longitudeDelta / 2;
    
    fetchLeaderboardDataForBounds(minLat, minLng, maxLat, maxLng);
  }, [location]);

  // Animate polygon scaling
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
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = null;
    }

    const hexesToSend = Array.from(visitedHexIds);
    const durationToSend = elapsedTime;
    const distanceToSend = distanceTraveled;

    if (hexesToSend.length > 0) {
      try {
        const activityData = {
          user_id: userId,
          h3_indexes: hexesToSend,
          duration: durationToSend,
          distance: distanceToSend + 1, // for demo purposes we want this to be positive
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
    setVisitedHexIds(new Set());
    setElapsedTime(0);
    setDistanceTraveled(0);
    setPreviousLocation(null);
  };

  // Handle hexagon press
  const handleHexPress = (hexId: string) => {
    if (selectedHexId === hexId) {
      animateHexScaling(hexId, 1.0);
      Animated.timing(leaderboardAnim, {
        toValue: 0,
        duration: 200,
        useNativeDriver: true,
      }).start(() => setSelectedHexId(null));
    } else {
      if (selectedHexId) {
        animateHexScaling(selectedHexId, 1.0);
      }
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
              {index + 1}. {entry.user} - {entry.score} pts
            </Text>
          ))}
          <TouchableOpacity onPress={() => handleHexPress(selectedHexId)} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>Hide</Text>
          </TouchableOpacity>
        </Animated.View>
      )}

      <View style={styles.activityControls}>
        {!isRecording ? (
          <TouchableOpacity onPress={() => setIsRecording(true)} style={styles.startButton}>
            <Text style={styles.startButtonText}>Start activity</Text>
          </TouchableOpacity>
        ) : (
          <View style={styles.timerContainer}>
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
  },
  errorText: {
    color: '#FFD600',
    fontSize: 16,
    textAlign: 'center',
    marginTop: 20,
  },
  locateButton: {
    position: 'absolute',
    bottom: 20,
    right: 20,
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
    gap: 10,
  },
  statsContainer: {
    backgroundColor: 'rgba(0,0,0,0.6)',
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: 20,
    alignItems: 'center',
  },
  timerText: {
    color: '#fff',
    fontSize: 16,
  },
  distanceText: {
    color: '#fff',
    fontSize: 12,
    marginTop: 4,
  },
  stopButton: {
    backgroundColor: '#ff5252',
    paddingVertical: 10,
    paddingHorizontal: 20,
    borderRadius: 30,
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
  closeButtonText: {
    color: '#fff',
    fontSize: 14,
  },
});