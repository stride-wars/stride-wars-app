import { useState } from 'react';
import { Alert } from 'react-native';

const API_BASE = 'http://localhost:8080';

export function useAuth() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleRegister = async () => {
    try {
      const response = await fetch(`${API_BASE}/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Registration failed');
      }

      const data = await response.json();
      Alert.alert('Success', `User registered: ${data.email}`);
    } catch (err) {
      // I will add some better error later
      Alert.alert('Error');
    }
  };

  const handleLogin = async () => {
    try {
      const response = await fetch(`${API_BASE}/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || 'Login failed');
      }

      const data = await response.json();
      Alert.alert('Success', `Logged in as: ${data.email}`);
    } catch (err) {
      Alert.alert('Error');
    }
  };

  return {
    email,
    setEmail,
    username, 
    setUsername,
    password,
    setPassword,
    handleRegister,
    handleLogin,
  };
}