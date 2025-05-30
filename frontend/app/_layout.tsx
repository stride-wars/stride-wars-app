import { Stack, useRouter, useSegments } from 'expo-router';
import React, { useEffect } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';

export default function RootLayout() {
  const segments = useSegments();
  const router = useRouter();

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const userData = await AsyncStorage.getItem('user');
      const inAuthGroup = segments[0] === '(tabs)';

      if (userData && !inAuthGroup) {
        // Redirect to tabs if user is logged in but not in tabs
        router.replace('/(tabs)');
      } else if (!userData && inAuthGroup) {
        // Redirect to login if user is not logged in but trying to access tabs
        router.replace('/login');
      }
    } catch (error) {
      console.error('Error checking auth:', error);
    }
  };

  return (
    <Stack
      screenOptions={{
        headerShown: false,
      }}
    >
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      <Stack.Screen name="login" options={{ headerShown: false }} />
      <Stack.Screen name="register" options={{ headerShown: false }} />
    </Stack>
  );
}
