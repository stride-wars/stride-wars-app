import React, { useState } from 'react';
import {
  Alert,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
  SafeAreaView,
} from 'react-native';
import { router } from 'expo-router';
import { useAuth } from '@/hooks/useAuth';

export default function RegisterScreen() {
  const {
    email,
    setEmail,
    password,
    setPassword,
    confirmPassword,
    setConfirmPassword,
    username,
    setUsername,
    handleRegister,
  } = useAuth();
  const [loading, setLoading] = useState(false);

  const onRegister = async () => {
    setLoading(true);
    await handleRegister();
    setLoading(false);
  };

  return (
    <SafeAreaView style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Create Account</Text>

        <TextInput
          style={styles.input}
          placeholder="Email"
          placeholderTextColor="#999"
          value={email}
          onChangeText={setEmail}
          autoCapitalize="none"
          keyboardType="email-address"
          autoComplete="email"
        />

        <TextInput
          style={styles.input}
          placeholder="Username"
          placeholderTextColor="#999"
          value={username}
          onChangeText={setUsername}
          autoCapitalize="none"
          autoComplete="username"
        />

        <TextInput
          style={styles.input}
          placeholder="Password"
          placeholderTextColor="#999"
          value={password}
          onChangeText={setPassword}
          secureTextEntry
          autoComplete="password-new"
        />

        <TextInput
          style={styles.input}
          placeholder="Confirm Password"
          placeholderTextColor="#999"
          value={confirmPassword}
          onChangeText={setConfirmPassword}
          secureTextEntry
          autoComplete="password-new"
        />

        <TouchableOpacity
          style={styles.registerButton}
          onPress={onRegister}
          disabled={loading}
        >
          <Text style={styles.buttonText}>
            {loading ? 'Creating Account...' : 'Register'}
          </Text>
        </TouchableOpacity>

        <TouchableOpacity
          style={styles.linkButton}
          onPress={() => router.push('/login')}
        >
          <Text style={styles.linkText}>Already have an account? Sign in</Text>
        </TouchableOpacity>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#0D0D0D', // Rich black
  },
  card: {
    flex: 1,
    margin: 24,
    backgroundColor: '#1C1F2A',
    borderRadius: 20,
    padding: 32,
    justifyContent: 'center',
    shadowColor: '#00FFFF',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.6,
    shadowRadius: 20,
    elevation: 12,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#FFD700',
    textAlign: 'center',
    marginBottom: 28,
    textShadowColor: 'rgba(255, 215, 0, 0.8)',
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 12,
  },
  input: {
    backgroundColor: '#2D2F38',
    borderRadius: 14,
    paddingVertical: 14,
    paddingHorizontal: 16,
    color: '#FFF',
    fontSize: 16,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: '#00BFFF',
    shadowColor: '#00BFFF',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.3,
    shadowRadius: 6,
  },
  registerButton: {
    backgroundColor: '#00BFFF',
    paddingVertical: 16,
    borderRadius: 16,
    alignItems: 'center',
    marginTop: 10,
    shadowColor: '#00FFFF',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.8,
    shadowRadius: 10,
    elevation: 10,
  },
  buttonText: {
    color: '#0D0D0D',
    fontSize: 18,
    fontWeight: 'bold',
  },
  linkButton: {
    marginTop: 28,
    alignItems: 'center',
  },
  linkText: {
    color: '#A0A0FF',
    fontSize: 16,
    textDecorationLine: 'underline',
  },
});
