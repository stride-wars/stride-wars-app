import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  Dimensions,
  ScrollView,
} from 'react-native';
import { FontAwesome5 } from '@expo/vector-icons';
import { LineChart } from 'react-native-chart-kit';

const screenWidth = Dimensions.get('window').width;

export default function ExploreTab() {
  const stats = {
    hexesVisited: 128,
    activitiesRecorded: 37,
  };

  const chartData = {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datasets: [
      {
        data: [2, 4, 3, 6, 4, 5, 7],
        strokeWidth: 2,
        color: (opacity = 1) => `rgba(255, 214, 0, ${opacity})`,
      },
    ],
    legend: ['Activities per Day'],
  };

  return (
    <ScrollView style={styles.container} contentContainerStyle={styles.content}>
      {/* User Avatar */}
      <View style={styles.avatarWrapper}>
        <View style={styles.avatarGlow}>
          <FontAwesome5 name="user-astronaut" size={48} color="#FFD600" />
        </View>
        <Text style={styles.username}>@YourStride</Text>
      </View>

      {/* Stats */}
      <View style={styles.statsRow}>
        <View style={styles.statCard}>
          <Text style={styles.statNumber}>{stats.hexesVisited}</Text>
          <Text style={styles.statLabel}>Hexes Visited</Text>
        </View>
        <View style={styles.statCard}>
          <Text style={styles.statNumber}>{stats.activitiesRecorded}</Text>
          <Text style={styles.statLabel}>Activities</Text>
        </View>
      </View>

      {/* Chart */}
      <Text style={styles.chartTitle}>Weekly Activity</Text>
      <LineChart
        data={chartData}
        width={screenWidth - 40}
        height={220}
        chartConfig={{
          backgroundGradientFrom: '#1E1E1E',
          backgroundGradientTo: '#1E1E1E',
          decimalPlaces: 0,
          color: (opacity = 1) => `rgba(255, 214, 0, ${opacity})`,
          labelColor: () => '#ccc',
          propsForDots: {
            r: '5',
            strokeWidth: '2',
            stroke: '#FFD600',
          },
        }}
        bezier
        style={styles.chart}
      />
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#121212',
  },
  content: {
    alignItems: 'center',
    paddingVertical: 24,
  },
  avatarWrapper: {
    alignItems: 'center',
    marginBottom: 20,
  },
  avatarGlow: {
    backgroundColor: '#272727',
    borderRadius: 64,
    padding: 24,
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.8,
    shadowRadius: 20,
    elevation: 10,
  },
  username: {
    color: '#FFD600',
    fontSize: 18,
    fontWeight: '600',
    marginTop: 12,
  },
  statsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    width: '90%',
    marginVertical: 20,
  },
  statCard: {
    backgroundColor: '#1E1E1E',
    padding: 20,
    borderRadius: 12,
    alignItems: 'center',
    width: '45%',
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.4,
    shadowRadius: 10,
    elevation: 8,
  },
  statNumber: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#FFD600',
  },
  statLabel: {
    fontSize: 14,
    color: '#aaa',
    marginTop: 4,
  },
  chartTitle: {
    fontSize: 20,
    color: '#FFD600',
    fontWeight: '600',
    marginBottom: 8,
  },
  chart: {
    marginVertical: 8,
    borderRadius: 16,
  },
});
