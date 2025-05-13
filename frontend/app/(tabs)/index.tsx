// App.tsx
import React from 'react';
import { SafeAreaView, StyleSheet } from 'react-native';
import { LocationButton } from './../../components/LocationButton';

export default function App() {
  return (
    <SafeAreaView style={styles.container}>
      <LocationButton />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 16
  }
});