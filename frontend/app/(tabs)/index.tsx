import React from 'react';
import { View, Text, StyleSheet, SafeAreaView, TouchableOpacity } from 'react-native';
import { useAuth } from '../../hooks/useAuth';

export default function HomeScreen() {
  const { handleLogout } = useAuth();

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.mainCard}>
        <Text style={styles.title}>Stride Wars</Text>
        <Text style={styles.subtitle}>Welcome to the game!</Text>

        {/* Mock small leaderboard panel */}
        <View style={styles.leaderboardPanel}>
          <Text style={styles.leaderboardTitle}>Leaderboard</Text>
          <Text style={styles.leaderboardEntry}>StrideRunner42 - 120 pts</Text>
          <Text style={styles.leaderboardEntry}>WalkMaster - 85 pts</Text>
          <Text style={styles.leaderboardEntry}>HikerPro - 95 pts</Text>
        </View>

        {/* Logout Button */}
        <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
          <Text style={styles.logoutButtonText}>Logout</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#121212',
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
});