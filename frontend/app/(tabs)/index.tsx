import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  SafeAreaView,
  TouchableOpacity,
  ScrollView,
  RefreshControl,
} from 'react-native';
import { useAuth } from '../../hooks/useAuth';
import { FontAwesome5 } from '@expo/vector-icons';
import { useGlobalLeaderboard } from '../../hooks/useGlobalLeaderboards';

export default function HomeScreen() {
  const { handleLogout } = useAuth();
  const {
    leaderboard,    // typed as Array<…> | null
    loading,
    error,
    refresh,
  } = useGlobalLeaderboard();

  const getTrophyIcon = (index: number) => {
    const colors = ['#FFD700', '#C0C0C0', '#CD7F32']; // gold, silver, bronze
    return (
      <FontAwesome5
        name="trophy"
        size={16}
        color={colors[index]}
        style={{ marginRight: 8 }}
      />
    );
  };

  // If still loading and there's no leaderboard data at all yet, show a placeholder
  if (loading && !leaderboard) {
    return (
      <SafeAreaView style={styles.centeredContainer}>
        <Text style={styles.loadingText}>Loading leaderboard...</Text>
      </SafeAreaView>
    );
  }

  // If an error occurred, show it (and let user retry)
  if (error) {
    return (
      <SafeAreaView style={styles.centeredContainer}>
        <Text style={styles.errorText}>Error: {error}</Text>
        <TouchableOpacity style={styles.retryButton} onPress={refresh}>
          <Text style={styles.retryButtonText}>Try Again</Text>
        </TouchableOpacity>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.mainCard}>
        <Text style={styles.title}>Stride Wars</Text>
        <Text style={styles.subtitle}>Welcome to the game!</Text>

        {/* Leaderboard Panel */}
        <View style={styles.leaderboardPanel}>
          <Text style={styles.leaderboardTitle}>🏆 Leaderboard</Text>
          <ScrollView
            style={styles.scrollArea}
            refreshControl={
              <RefreshControl
                refreshing={loading}
                onRefresh={refresh}
                colors={['#FFD600']}     // Android indicator color
                tintColor="#FFD600"      // iOS indicator color
              />
            }
          >
            {/* Safely map over `leaderboard` only if non-null; else show “No entries yet.” */}
            {leaderboard && leaderboard.length > 0 ? (
              leaderboard.map((entry, index) => (
                <View style={styles.leaderboardRow} key={entry.user_id}>
                  <Text style={styles.rank}>{index + 1}.</Text>
                  {index < 3 && getTrophyIcon(index)}
                  <Text style={styles.name}>{entry.username}</Text>
                  <Text style={styles.score}>{entry.top_count} pts</Text>
                </View>
              ))
            ) : (
              <Text style={styles.noDataText}>No leaderboard entries yet.</Text>
            )}
          </ScrollView>
        </View>
      </View>

      {/* Logout Button */}
      <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
        <Text style={styles.logoutButtonText}>Logout</Text>
      </TouchableOpacity>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  centeredContainer: {
    flex: 1,
    backgroundColor: '#121212',
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    color: '#FFD600',
    fontSize: 18,
    textAlign: 'center',
  },
  errorText: {
    color: '#FF4444',
    fontSize: 18,
    textAlign: 'center',
  },
  retryButton: {
    marginTop: 16,
    alignSelf: 'center',
    backgroundColor: '#FFD600',
    paddingVertical: 8,
    paddingHorizontal: 16,
    borderRadius: 8,
  },
  retryButtonText: {
    color: '#000',
    fontWeight: '600',
  },
  container: {
    flex: 1,
    backgroundColor: '#121212',
    justifyContent: 'space-between',
  },
  mainCard: {
    margin: 24,
    backgroundColor: '#1E1E1E',
    borderRadius: 16,
    padding: 28,
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
    marginBottom: 12,
    textAlign: 'center',
  },
  scrollArea: {
    maxHeight: 240,
  },
  leaderboardRow: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 6,
    borderBottomColor: '#444',
    borderBottomWidth: 1,
    paddingHorizontal: 6,
  },
  rank: {
    color: '#fff',
    fontWeight: 'bold',
    width: 24,
    textAlign: 'center',
  },
  name: {
    color: '#ddd',
    flex: 1,
    marginLeft: 4,
  },
  score: {
    color: '#FFD600',
    fontWeight: '600',
  },
  noDataText: {
    color: '#aaa',
    textAlign: 'center',
    marginVertical: 20,
  },
  logoutButton: {
    backgroundColor: '#2C2C2C',
    paddingVertical: 12,
    margin: 20,
    borderRadius: 10,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.4,
    shadowRadius: 4,
    elevation: 4,
  },
  logoutButtonText: {
    color: '#FF3B30',
    fontSize: 16,
    fontWeight: '600',
  },
});
