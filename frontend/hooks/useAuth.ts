import { useState } from 'react';
import { Alert } from 'react-native';
import { useRouter } from 'expo-router';
import { handleError } from '@/utils/handleError';
import { api } from '@/api';
import AsyncStorage from '@react-native-async-storage/async-storage';

export function useAuth() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!email) return 'Email is required';
    if (!emailRegex.test(email)) return 'Invalid email format';
    return '';
  };

  const validatePassword = (password: string) => {
    if (!password) return 'Password is required';
    if (password.length < 6) return 'Password must be at least 6 characters';
    return '';
  };

  const validateUsername = (username: string) => {
    if (!username) return 'Username is required';
    if (username.length < 3) return 'Username must be at least 3 characters';
    return '';
  };

  const validateForm = (isLogin: boolean = false) => {
    const newErrors: Record<string, string> = {};
    
    const emailError = validateEmail(email);
    if (emailError) newErrors.email = emailError;

    const passwordError = validatePassword(password);
    if (passwordError) newErrors.password = passwordError;

    if (!isLogin) {
      const usernameError = validateUsername(username);
      if (usernameError) newErrors.username = usernameError;

      if (password !== confirmPassword) {
        newErrors.confirmPassword = 'Passwords do not match';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleRegister = async () => {
    if (!validateForm(false)) return;

    try {
      const response = await api.signUp(username, email, password);
      if (response.data) {
        // Store auth data
        await AsyncStorage.setItem('access_token', response.data.session.access_token);
        await AsyncStorage.setItem('refresh_token', response.data.session.refresh_token);
        await AsyncStorage.setItem('user', JSON.stringify({
          id: response.data.user_id,
          username: response.data.username,
          email: response.data.email,
          external_user: response.data.external_user
        }));
        router.replace('/');
      }
    } catch (error) {
      handleError(error);
    }
  };

  const handleLogin = async () => {
    if (!validateForm(true)) return;

    try {
      const response = await api.signIn(email, password);
      console.log('%c Login Response:', 'color: #4CAF50; font-weight: bold', response);
      
      if (response.error) {
        throw new Error(response.error);
      }
      
      if (!response.data) {
        console.log('%c No data in response', 'color: #FF9800; font-weight: bold');
        throw new Error('Invalid email or password');
      }
      
      // Store auth data
      await AsyncStorage.setItem('access_token', response.data.session.access_token);
      await AsyncStorage.setItem('refresh_token', response.data.session.refresh_token);
      await AsyncStorage.setItem('user', JSON.stringify({
        id: response.data.user_id,
        username: response.data.username,
        email: response.data.email,
        external_user: response.data.external_user
      }));
      router.replace('/');
    } catch (error) {
      console.log('%c Login Error:', 'color: #F44336; font-weight: bold', error);
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('An error occurred during login');
    }
  };

  return {
    email,
    setEmail,
    username,
    setUsername,
    password,
    setPassword,
    confirmPassword,
    setConfirmPassword,
    errors,
    handleRegister,
    handleLogin,
  };
}