import { Tabs } from 'expo-router';
import { FontAwesome } from '@expo/vector-icons';
import { StatusBar } from 'expo-status-bar';
import { View, Image } from 'react-native';

export default function TabLayout() {
  return (
    <View style={{ flex: 1, backgroundColor: '#121212' }}>
      <StatusBar style="light" backgroundColor="#1E1E1E" translucent={false} />
      <Tabs
        screenOptions={{
          tabBarActiveTintColor: '#FFD600',
          tabBarInactiveTintColor: '#888',
          tabBarLabelStyle: {
            fontSize: 12,
            fontWeight: '600',
          },
          tabBarStyle: {
            backgroundColor: '#1E1E1E',
            borderTopColor: '#333',
            paddingBottom: 6,
            paddingTop: 4,
            height: 60,
          },
          headerStyle: {
            backgroundColor: '#1E1E1E',
          },
          headerTitleStyle: {
            color: '#FFD600',
            fontWeight: 'bold',
            fontSize: 20,
          },
          headerTintColor: '#FFD600',

          headerRight: () => (
            <Image
              source={require('../../assets/images/stride_wars.png')}
              style={{
                width: 53,
                height: 53,
                marginRight: 16,
                resizeMode: 'contain',
              }}
            />
          ),
        }}
      >
        <Tabs.Screen
          name="index"
          options={{
            title: 'Home',
            tabBarIcon: ({ color, focused }) => (
              <FontAwesome
                name="home"
                size={24}
                color={color}
                style={{
                  textShadowColor: focused ? '#FFD600AA' : 'transparent',
                  textShadowOffset: { width: 0, height: 0 },
                  textShadowRadius: 8,
                }}
              />
            ),
          }}
        />
        <Tabs.Screen
          name="map"
          options={{
            title: 'Map',
            tabBarIcon: ({ color, focused }) => (
              <FontAwesome
                name="map-marker"
                size={24}
                color={color}
                style={{
                  textShadowColor: focused ? '#FFD600AA' : 'transparent',
                  textShadowOffset: { width: 0, height: 0 },
                  textShadowRadius: 8,
                }}
              />
            ),
          }}
        />
        <Tabs.Screen
          name="profile"
          options={{
            title: 'Profile',
            tabBarIcon: ({ color, focused }) => (
              <FontAwesome
                name="user"
                size={24}
                color={color}
                style={{
                  textShadowColor: focused ? '#FFD600AA' : 'transparent',
                  textShadowOffset: { width: 0, height: 0 },
                  textShadowRadius: 8,
                }}
              />
            ),
          }}
        />
      </Tabs>
    </View>
  );
}
