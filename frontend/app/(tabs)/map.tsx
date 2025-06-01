import React, { useEffect, useState, useRef } from 'react';
import { StyleSheet, View, Text, TouchableOpacity, Platform } from 'react-native';
import MapView, { Polygon, PROVIDER_DEFAULT, Region } from 'react-native-maps';
import * as Location from 'expo-location';
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

export default function MapScreen() {
  const [location, setLocation] = useState<Location.LocationObject | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [selectedHexId, setSelectedHexId] = useState<string | null>(null);
  const [hexagons, setHexagons] = useState<Array<{
    hexId: string;
    coordinates: Coordinate[];
    animatedCoordinates: Coordinate[];
  }>>([]);

  const mapRef = useRef<MapView>(null);

  useEffect(() => {
    (async () => {
      const { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setErrorMsg('Permission to access location was denied');
        return;
      }

      const loc = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.High,
      });
      setLocation(loc);

      const rawHexes = getHexagonsInRadius(
        loc.coords.latitude,
        loc.coords.longitude,
        1000,
        9
      );

      const enrichedHexes = rawHexes.map(h => ({
        ...h,
        animatedCoordinates: h.coordinates,
      }));

      setHexagons(enrichedHexes);
    })();
  }, []);

  const animateHexScaling = (hexId: string, toScale: number, duration: number = 150) => {
    const hex = hexagons.find(h => h.hexId === hexId);
    if (!hex) return;

    const from = hex.animatedCoordinates;
    const to = scalePolygon(hex.coordinates, toScale);
    const steps = 10;
    let currentStep = 0;

    const interval = setInterval(() => {
      currentStep++;
      const t = currentStep / steps;
      const intermediate = interpolatePolygon(from, to, t);

      setHexagons(prev =>
        prev.map(h =>
          h.hexId === hexId
            ? { ...h, animatedCoordinates: intermediate }
            : h
        )
      );

      if (currentStep >= steps) clearInterval(interval);
    }, duration / steps);
  };

  const handleHexPress = (hexId: string) => {
    if (selectedHexId === hexId) {
      animateHexScaling(hexId, 1.0);
      setSelectedHexId(null);
    } else {
      if (selectedHexId) {
        animateHexScaling(selectedHexId, 1.0);
      }
      animateHexScaling(hexId, 1.05);
      setSelectedHexId(hexId);
    }
  };

  const defaultRegion: Region = {
    latitude: location?.coords.latitude || 0,
    longitude: location?.coords.longitude || 0,
    latitudeDelta: 0.01,
    longitudeDelta: 0.01,
  };

  if (errorMsg) {
    return (
      <View style={styles.container}>
        <Text style={styles.errorText}>{errorMsg}</Text>
      </View>
    );
  }

  if (!location) {
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
        initialRegion={defaultRegion}
        showsUserLocation
        showsMyLocationButton={Platform.OS === 'android'}
        showsCompass
        mapType="standard"
      >
        {hexagons.map(hex => (
          <Polygon
            key={hex.hexId}
            coordinates={hex.animatedCoordinates}
            fillColor={getHexagonColor(hex.hexId)} // translucent and varied
            strokeColor="rgba(255, 255, 255, 0.4)"
            strokeWidth={1}
            tappable
            onPress={() => handleHexPress(hex.hexId)}
          />
        ))}
      </MapView>

      {/* Custom "Locate Me" Button (iOS + Android fallback) */}
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
});
