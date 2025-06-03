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
import MapView, { Polygon, PROVIDER_DEFAULT } from 'react-native-maps';
import { MaterialIcons } from '@expo/vector-icons';
import { getHexagonsInRadius, getHexagonColor } from '../../utils/h3Utils';
import { useLocation } from '../../hooks/useLocation';

type Coordinate = { latitude: number; longitude: number };

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

const SCREEN_WIDTH = Dimensions.get('window').width;

export default function MapScreen() {
  const { location, error, isLoading, getLocation } = useLocation();
  const [selectedHexId, setSelectedHexId] = useState<string | null>(null);
  const [hexagons, setHexagons] = useState<Array<{
    hexId: string;
    coordinates: Coordinate[];
    animatedCoordinates: Coordinate[];
  }>>([]);
  const mapRef = useRef<MapView>(null);

  const [isRecording, setIsRecording] = useState(false);
  const [elapsedTime, setElapsedTime] = useState(0);
  const timerRef = useRef<NodeJS.Timer | null>(null);
  const leaderboardAnim = useRef(new Animated.Value(0)).current;

  useEffect(() => {
    getLocation();
  }, []);

  useEffect(() => {
    if (isRecording) {
      timerRef.current = setInterval(() => {
        setElapsedTime(prev => prev + 1);
      }, 1000);
    } else {
      if (timerRef.current) {
        clearInterval(timerRef.current);
        timerRef.current = null;
      }
      setElapsedTime(0);
    }

    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [isRecording]);

  useEffect(() => {
    if (location) {
      const rawHexes = getHexagonsInRadius(
        location.coords.latitude,
        location.coords.longitude,
        1000,
        9
      );

      const enrichedHexes = rawHexes.map(h => ({
        ...h,
        animatedCoordinates: h.coordinates,
      }));

      setHexagons(enrichedHexes);

      mapRef.current?.animateToRegion({
        latitude: location.coords.latitude,
        longitude: location.coords.longitude,
        latitudeDelta: 0.01,
        longitudeDelta: 0.01,
      });
    }
  }, [location]);

  const animateHexScaling = (hexId: string, toScale: number, duration: number = 150) => {
    const hex = hexagons.find(h => h.hexId === hexId);
    if (!hex) return;

    const from = hex.animatedCoordinates;
    const to = scalePolygon(hex.coordinates, toScale);
    const steps = 10;
    let currentStep = 0;
    const updatedFrames: Coordinate[][] = [];

    const animateStep = () => {
      currentStep++;
      const t = currentStep / steps;
      const intermediate = interpolatePolygon(from, to, t);
      updatedFrames.push(intermediate);

      if (currentStep < steps) {
        requestAnimationFrame(animateStep);
      } else {
        // Apply final update once to reduce state churn
        setHexagons(prev =>
          prev.map(h =>
            h.hexId === hexId
              ? { ...h, animatedCoordinates: to }
              : h
          )
        );
      }
    };

    animateStep();
  };

  const handleHexPress = (hexId: string) => {
    if (selectedHexId === hexId) {
      animateHexScaling(hexId, 1.0);
      Animated.timing(leaderboardAnim, {
        toValue: 0,
        duration: 200,
        useNativeDriver: true, // changed
      }).start(() => setSelectedHexId(null));
    } else {
      if (selectedHexId) {
        animateHexScaling(selectedHexId, 1.0);
      }
      animateHexScaling(hexId, 1.05);
      setSelectedHexId(hexId); // moved up
      Animated.timing(leaderboardAnim, {
        toValue: 1,
        duration: 200,
        useNativeDriver: true, // changed
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
          <Text style={styles.leaderboardEntry}>1. Alice - 120 pts</Text>
          <Text style={styles.leaderboardEntry}>2. Bob - 95 pts</Text>
          <Text style={styles.leaderboardEntry}>3. Charlie - 80 pts</Text>
          <TouchableOpacity onPress={() => handleHexPress(selectedHexId)} style={styles.closeButton}>
            <Text style={styles.closeButtonText}>Dismiss</Text>
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
            <Text style={styles.timerText}>
              {Math.floor(elapsedTime / 60).toString().padStart(2, '0')}:
              {(elapsedTime % 60).toString().padStart(2, '0')}
            </Text>
            <TouchableOpacity onPress={() => setIsRecording(false)} style={styles.stopButton}>
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
  timerText: {
    color: '#fff',
    fontSize: 16,
    backgroundColor: 'rgba(0,0,0,0.6)',
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: 20,
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
