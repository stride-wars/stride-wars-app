import React, { useState } from 'react';
import {
  View,
  StyleSheet,
  SafeAreaView,
  KeyboardAvoidingView,
  Platform,
  TouchableWithoutFeedback,
  Keyboard,
  Image,
  Text,
  TouchableOpacity,
} from 'react-native';
import { useRouter } from 'expo-router';
import { useAuth } from '../hooks/useAuth';
import { Snackbar } from '../components/Snackbar';
import { Input } from '../components/Input';
import { Button } from '../components/Button';
import { Feather } from '@expo/vector-icons';

export default function LoginScreen() {
  const router = useRouter();
  const { email, setEmail, password, setPassword, errors, handleLogin } = useAuth();
  const [snackbarVisible, setSnackbarVisible] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');

  const handleLoginWithSnackbar = async () => {
    try {
      await handleLogin();
    } catch (error) {
      console.log('%c Login Screen Error:', 'color: #F44336; font-weight: bold', error);
      if (error instanceof Error) {
        setSnackbarMessage(error.message);
        setSnackbarVisible(true);
      }
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardAvoid}
      >
        <TouchableWithoutFeedback >
          <View style={styles.inner}>
            <View style={styles.logoContainer}>
              <Image
                source={require('../assets/images/stride_wars.png')} 
                style={styles.logo}
                resizeMode="contain"
              />
              <Text style={styles.appTitle}>Stride Wars</Text>
              <Text style={styles.subtitle}>Conquer territory with every step</Text>
            </View>

            <View style={styles.formContainer}>
              <Input
                label="Email"
                placeholder="Enter your email"
                value={email}
                onChangeText={setEmail}
                error={errors.email}
                autoCapitalize="none"
                autoComplete="email"
                keyboardType="email-address"
                leftIcon={<Feather name="user" size={20} color="#888" />}
              />

              <Input
                label="Password"
                placeholder="Enter your password"
                value={password}
                onChangeText={setPassword}
                error={errors.password}
                secureTextEntry
                autoCapitalize="none"
                autoComplete="password"
                leftIcon={<Feather name="lock" size={20} color="#888" />}
              />

              <Button onPress={handleLoginWithSnackbar} style={styles.loginButton}>
                <Feather name="log-in" size={20} color="#fff" style={{ marginRight: 8 }} />
                Login
              </Button>

              <TouchableOpacity onPress={() => router.push('/register')} style={styles.registerLink}>
                <Text style={styles.registerText}>
                  Don't have an account? <Text style={styles.registerTextBold}>Register</Text>
                </Text>
              </TouchableOpacity>
            </View>
          </View>
        </TouchableWithoutFeedback>
      </KeyboardAvoidingView>

      <Snackbar
        visible={snackbarVisible}
        message={snackbarMessage}
        onDismiss={() => setSnackbarVisible(false)}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#111827', // dark background
  },
  keyboardAvoid: {
    flex: 1,
  },
  inner: {
    flex: 1,
    padding: 20,
    justifyContent: 'center',
  },
  logoContainer: {
    alignItems: 'center',
    marginBottom: 32,
  },
  logo: {
    width: 120,
    height: 120,
    marginBottom: 12,
  },
  appTitle: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#FACC15', // yellow-500
  },
  subtitle: {
    fontSize: 14,
    color: '#9CA3AF', // gray-400
    marginTop: 4,
  },
  formContainer: {
    gap: 20,
  },
  loginButton: {
    marginTop: 10,
    backgroundColor: '#2563EB', // blue-600
  },
  registerLink: {
    marginTop: 16,
    alignItems: 'center',
  },
  registerText: {
    fontSize: 14,
    color: '#FACC15', // yellow-500
  },
  registerTextBold: {
    fontWeight: '600',
  },
});
