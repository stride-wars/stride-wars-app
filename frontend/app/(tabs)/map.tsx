import React, { useEffect, useState } from 'react';
import { StyleSheet, View, Text } from 'react-native';
import MapView, { Polygon, PROVIDER_DEFAULT } from 'react-native-maps';
import * as Location from 'expo-location';
import { getHexagonsInRadius, getHexagonColor } from '../../utils/h3Utils';

export default function MapScreen() {
  const [location, setLocation] = useState<Location.LocationObject | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [hexagons, setHexagons] = useState<Array<{
    hexId: string;
    coordinates: Array<{ latitude: number; longitude: number }>;
  }>>([]);

  useEffect(() => {
    (async () => {
      let { status } = await Location.requestForegroundPermissionsAsync();
      if (status !== 'granted') {
        setErrorMsg('Permission to access location was denied');
        return;
      }

      let location = await Location.getCurrentPositionAsync({
        accuracy: Location.Accuracy.Balanced,
      });
      setLocation(location);

      // Generate hexagons around the user's location
      const hexes = getHexagonsInRadius(
        location.coords.latitude,
        location.coords.longitude,
        1000, // 1km radius
        9 // resolution (9 is a good balance between detail and performance)
      );
      setHexagons(hexes);
    })();
  }, []);

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
        <Text style={styles.loadingText}>Loading map...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <MapView
        style={styles.map}
        provider={PROVIDER_DEFAULT}
        initialRegion={{
          latitude: location.coords.latitude,
          longitude: location.coords.longitude,
          latitudeDelta: 0.01,
          longitudeDelta: 0.01,
        }}
        showsUserLocation={true}
        showsMyLocationButton={true}
        showsCompass={true}
        mapType="standard"
      >
        {hexagons.map((hexagon) => (
          <Polygon
            key={hexagon.hexId}
            coordinates={hexagon.coordinates}
            fillColor={getHexagonColor(hexagon.hexId)}
            strokeColor="rgba(255, 255, 255, 0.5)"
            strokeWidth={1}
            tappable={true}
            onPress={() => {
              // Handle hexagon press if needed
              console.log('Hexagon pressed:', hexagon.hexId);
            }}
          />
        ))}
      </MapView>
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
  loadingText: {
    color: '#FFD600',
    fontSize: 16,
    textAlign: 'center',
    marginTop: 20,
  },
}); 