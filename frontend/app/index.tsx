import React, { useEffect, useState } from 'react'
import {
  View,
  Text,
  StyleSheet,
  SafeAreaView,
  TouchableOpacity,
  ActivityIndicator,
  Image,
  StatusBar,
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
        <StatusBar barStyle="light-content" backgroundColor="#121212" />
        <ActivityIndicator size="large" color="#FFD600" />
      </View>
    )
  }

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="light-content" backgroundColor="#121212" />
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
            {/* Add logo */}
            <Image
              source={require('../assets/images/stride_wars.png')} // Replace with actual logo path
              style={styles.logo}
              resizeMode="contain"
            />
            <Text style={styles.title}>Stride Wars</Text>
            <Text style={styles.subtitle}>Welcome to the game!</Text>
            <View style={styles.buttonContainer}>
              <TouchableOpacity
                style={[styles.button, styles.loginButton]}
                onPress={() => router.push('/login')}
              >
                <Text style={styles.buttonText}>Login</Text>
              </TouchableOpacity>

              <TouchableOpacity
                style={[styles.button, styles.registerButton]}
                onPress={() => router.push('/register')}
              >
                <Text style={styles.registerButtonText}>Register</Text>
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
    borderRadius: 20,
    padding: 32,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.5,
    shadowRadius: 16,
    elevation: 10,
  },
logo: {
    width: 300,       
    height: 300,      
    marginBottom: 24,  
    resizeMode: 'contain',
},
  title: {
    fontSize: 36,
    fontWeight: 'bold',
    color: '#FFD600',
    marginBottom: 16,
    textShadowColor: 'rgba(250, 204, 21, 0.8)',
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 10,
  },
  subtitle: {
    fontSize: 16,
    color: '#A1A1AA',
    textAlign: 'center',
    marginBottom: 40,
  },
  buttonContainer: {
    width: '100%',
    gap: 16,
  },
  button: {
    paddingVertical: 16,
    borderRadius: 14,
    width: '100%',
    alignItems: 'center',
    elevation: 6,
  },
  loginButton: {
    backgroundColor: '#FFD600',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.5,
    shadowRadius: 8,
  },
  registerButton: {
    backgroundColor: '#2C3E50',
    shadowColor: '#00FFFF',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.9,
    shadowRadius: 12,
    borderWidth: 1,
    borderColor: '#00FFFF',
  },
  buttonText: {
    color: '#121212',
    fontSize: 18,
    fontWeight: '700',
  },
  registerButtonText: {
    color: '#00FFFF',
    fontSize: 18,
    fontWeight: '700',
  },
  welcomeText: {
    fontSize: 28,
    fontWeight: '700',
    color: '#FFD600',
    marginLeft: 12,
    textShadowColor: 'rgba(250, 204, 21, 0.8)',
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 8,
  },
  logoutButton: {
    marginTop: 28,
    backgroundColor: '#FF3B30',
    paddingVertical: 14,
    paddingHorizontal: 36,
    borderRadius: 14,
    shadowColor: '#FF3B30',
    shadowOffset: { width: 0, height: 3 },
    shadowOpacity: 0.5,
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
    marginBottom: 12,
  },
  colorCircle: {
    width: 18,
    height: 18,
    borderRadius: 9,
  },
})
