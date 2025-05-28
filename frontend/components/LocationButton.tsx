import React from 'react';
import { View, Button, Text, ActivityIndicator } from 'react-native';
import { useLocation } from '../hooks/useLocation';

export const LocationButton = () => {
  const { location, error, isLoading, getLocation } = useLocation();

  return (
    <View>
      <Button
        title={isLoading ? 'Fetching...' : 'Get Location'}
        onPress={getLocation}
        disabled={isLoading}
      />

      {isLoading && <ActivityIndicator size="small" />}

      {error && <Text style={{ color: 'red' }}>Error: {error}</Text>}

      {location && (
        <View >
          <Text style={{ color: 'grey' }}>Latitude: {location.coords.latitude}</Text>
          <Text style={{ color: 'grey' }}>Longitude: {location.coords.longitude}</Text>
          <Text style={{ color: 'grey' }}>Last Updated: {new Date(location.timestamp).toLocaleString()}</Text>
        </View>
      )}
    </View>
  );
};