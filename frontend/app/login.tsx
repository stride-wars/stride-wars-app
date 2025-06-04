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
        <TouchableWithoutFeedback onPress={Keyboard.dismiss}>
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
    backgroundColor: '#121212', // dark background like MapScreen
  },
  keyboardAvoid: {
    flex: 1,
  },
  inner: {
    flex: 1,
    padding: 24,
    justifyContent: 'center',
  },
  logoContainer: {
    alignItems: 'center',
    marginBottom: 40,
  },
  logo: {
    width: 180,
    height: 180,
    marginBottom: 16,
    shadowColor: '#FACC15',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.8,
    shadowRadius: 15,
    elevation: 15,
  },
  appTitle: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#FACC15', // bright yellow
    textShadowColor: 'rgba(250, 204, 21, 0.8)',
    textShadowOffset: { width: 0, height: 0 },
    textShadowRadius: 10,
  },
  subtitle: {
    fontSize: 16,
    color: '#9CA3AF',
    marginTop: 6,
  },
  formContainer: {
    gap: 24,
  },
  loginButton: {
    marginTop: 12,
    backgroundColor: '#FFD600', // bright yellow
    shadowColor: '#FFD600',
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.9,
    shadowRadius: 10,
    elevation: 10,
  },
  registerLink: {
    marginTop: 20,
    alignItems: 'center',
  },
  registerText: {
    fontSize: 15,
    color: '#FACC15',
  },
  registerTextBold: {
    fontWeight: '700',
  },
});
