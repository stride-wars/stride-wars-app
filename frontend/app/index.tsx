import React, { useEffect, useState } from 'react'
import {
  View,
  Text,
  StyleSheet,
  SafeAreaView,
  TouchableOpacity,
  ActivityIndicator,
} from 'react-native'
import { useRouter } from 'expo-router'
import AsyncStorage from '@react-native-async-storage/async-storage'
import { useAuth } from '../hooks/useAuth'

interface User {
  id: string
  username: string
  email: string
}

export default function IndexScreen() {
  const router = useRouter()
  const [isLoading, setIsLoading] = useState(true)
  const [user, setUser] = useState<User | null>(null)
  const { handleLogout } = useAuth()

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    try {
      const userData = await AsyncStorage.getItem('user')
      if (userData) {
        setUser(JSON.parse(userData))
        router.replace('/(tabs)')
      }
    } catch (error) {
      console.error('Error checking auth:', error)
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <View style={[styles.container, styles.centerContent]}>
        <ActivityIndicator size="large" color="#FFD600" />
      </View>
    )
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.mainCard}>
        {user ? (
          <>
            <View style={styles.userBadge}>
              <View style={[styles.colorCircle, { backgroundColor: '#FFD600' }]} />
              <Text style={styles.welcomeText}>Welcome, {user.username}!</Text>
            </View>

            <Text style={styles.subtitle}>You are now logged in to Stride Wars</Text>

            <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
              <Text style={styles.logoutButtonText}>Logout</Text>
            </TouchableOpacity>
          </>
        ) : (
          <>
            <Text style={styles.title}>Stride Wars</Text>
            <Text style={styles.subtitle}>Welcome to the game!</Text>

            <View style={styles.buttonContainer}>
              <TouchableOpacity
                style={[styles.button, { backgroundColor: '#FFD600' }]}
                onPress={() => router.push('/login')}
              >
                <Text style={styles.buttonText}>Login</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.button, { backgroundColor: '#34C759' }]}
                onPress={() => router.push('/register')}
              >
                <Text style={styles.buttonText}>Register</Text>
              </TouchableOpacity>
            </View>
          </>
        )}
      </View>
    </SafeAreaView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#121212',
  },
  centerContent: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  mainCard: {
    flex: 1,
    margin: 24,
    backgroundColor: '#1E1E1E',
    borderRadius: 16,
    padding: 28,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.6,
    shadowRadius: 12,
    elevation: 10,
  },
  title: {
    fontSize: 36,
    fontWeight: 'bold',
    color: '#FFD600',
    marginBottom: 12,
  },
  subtitle: {
    fontSize: 18,
    color: '#ccc',
    textAlign: 'center',
    marginBottom: 40,
  },
  buttonContainer: {
    width: '100%',
    gap: 16,
  },
  button: {
    paddingVertical: 16,
    borderRadius: 12,
    width: '100%',
    alignItems: 'center',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.5,
    shadowRadius: 8,
    elevation: 6,
  },
  buttonText: {
    color: '#121212',
    fontSize: 18,
    fontWeight: '700',
  },
  welcomeText: {
    fontSize: 28,
    fontWeight: '700',
    color: '#FFD600',
    marginLeft: 12,
  },
  logoutButton: {
    marginTop: 28,
    backgroundColor: '#FF3B30',
    paddingVertical: 14,
    paddingHorizontal: 32,
    borderRadius: 12,
    shadowColor: '#FF3B30',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.6,
    shadowRadius: 10,
    elevation: 8,
  },
  logoutButtonText: {
    color: '#fff',
    fontSize: 18,
    fontWeight: '700',
  },
  userBadge: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 8,
  },
  colorCircle: {
    width: 18,
    height: 18,
    borderRadius: 9,
  },
  leaderboardPanel: {
    marginTop: 48,
    backgroundColor: '#272727',
    padding: 16,
    borderRadius: 12,
    width: '100%',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.4,
    shadowRadius: 10,
  },
  leaderboardTitle: {
    fontSize: 20,
    fontWeight: '700',
    color: '#FFD600',
    marginBottom: 8,
    textAlign: 'center',
  },
  leaderboardEntry: {
    fontSize: 16,
    color: '#ddd',
    marginBottom: 4,
  },
})
