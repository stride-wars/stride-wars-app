import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  SafeAreaView,
  TouchableOpacity,
  ScrollView,
} from 'react-native';
import { useAuth } from '../../hooks/useAuth';
import { FontAwesome5 } from '@expo/vector-icons';
import { useGlobalLeaderboard } from '../../hooks/useGlobalLeaderboards';

const leaderboardData = [
  { name: 'StrideRunner42', score: 120 },
  { name: 'HikerPro', score: 95 },
  { name: 'WalkMaster', score: 85 },
  { name: 'TrailBlazer99', score: 70 },
  { name: 'SpeedWalker', score: 65 },
  { name: 'UrbanStroller', score: 60 },
  { name: 'FootFrenzy', score: 58 },
  { name: 'JoggerX', score: 55 },
  { name: 'ExplorerElle', score: 53 },
  { name: 'MarathonMike', score: 50 },
];

export default function HomeScreen() {
  const { handleLogout } = useAuth();
  const { leaderboard, loading, error, refresh } = useGlobalLeaderboard();

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

  if (loading) return <Text>Loading leaderboard...</Text>;
  if (error) return <Text>Error: {error}</Text>;

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.mainCard}>
        <Text style={styles.title}>Stride Wars</Text>
        <Text style={styles.subtitle}>Welcome to the game!</Text>

        {/* Leaderboard */}
        <View style={styles.leaderboardPanel}>
          <Text style={styles.leaderboardTitle}>🏆 Leaderboard</Text>
          <ScrollView style={{ maxHeight: 240 }}>
            {leaderboard && leaderboard.map((entry, index) => (
              <View style={styles.leaderboardRow} key={entry.user_id}>
                <Text style={styles.rank}>{index + 1}.</Text>
                {index < 3 && getTrophyIcon(index)}
                <Text style={styles.name}>{entry.name}</Text>
                <Text style={styles.score}>{entry.score} pts</Text>
              </View>
            ))}
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
