import AsyncStorage from '@react-native-async-storage/async-storage';
import { 
  SignUpRequest, 
  SignInRequest, 
  SignUpResponse, 
  SignInResponse, 
  ApiResponse,
  Session,
  GlobalLeaderboardEntry
} from '../consts/types';
//const API_BASE = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const API_BASE = 'https://4d85-188-146-191-2.ngrok-free.app/api/v1';

class ApiClient {
  private async refreshToken(): Promise<Session | null> {
    try {
      const refreshToken = await AsyncStorage.getItem('refresh_token');
      if (!refreshToken) return null;

      const response = await fetch(`${API_BASE}/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!response.ok) {
        await this.signOut();
        return null;
      }

      const data = await response.json();
      const session = data.data.session;

      // Store new tokens
      await AsyncStorage.setItem('access_token', session.access_token);
      await AsyncStorage.setItem('refresh_token', session.refresh_token);

      return session;
    } catch (error) {
      console.error('Failed to refresh token:', {
        message: (error as Error).message,
        stack: (error as Error).stack,
        name: (error as Error).name,
      });
      await this.signOut();
      return null;
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      let token = await AsyncStorage.getItem('access_token');
      
      // If no token, try to refresh
      if (!token) {
        const session = await this.refreshToken();
        if (session) {
          token = session.access_token;
        }
      }

      const headers: HeadersInit = {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...options.headers,
      };

      const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
      });

      const data = await response.json();
      console.log('%c API Response:', 'color: #2196F3; font-weight: bold', data);

      // Handle 401 errors differently for auth endpoints
      if (response.status === 401) {
        // For auth endpoints, pass through the error message
        if (endpoint.includes('/auth/')) {
          if (data.error) {
            throw new Error(data.error);
          }
          throw new Error('Invalid email or password');
        }

        // For other endpoints, handle session expiration
        const session = await this.refreshToken();
        if (session) {
          // Retry the request with new token
          const retryResponse = await fetch(`${API_BASE}${endpoint}`, {
            ...options,
            headers: {
              ...headers,
              Authorization: `Bearer ${session.access_token}`,
            },
          });

          if (!retryResponse.ok) {
            await this.signOut();
            return { error: 'Session expired. Please sign in again.' };
          }

          const retryData = await retryResponse.json();
          return { data: retryData };
        }

        await this.signOut();
        return { error: 'Session expired. Please sign in again.' };
      }

      if (!response.ok) {
        console.log('%c API Error:', 'color: #F44336; font-weight: bold', data.error);
        if (data.error) {
          throw new Error(data.error);
        }
        throw new Error('Invalid email or password');
      }

      return { data };
    } catch (error) {
      return { error: error instanceof Error ? error.message : 'An error occurred' };
    }
  }

  async signUp(username: string, email: string, password: string): Promise<ApiResponse<SignUpResponse>> {
    const request: SignUpRequest = { username, email, password };
    return this.request<SignUpResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async signIn(email: string, password: string): Promise<ApiResponse<SignInResponse>> {
    const request: SignInRequest = { email, password }; 
    return this.request<SignInResponse>('/auth/signin', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async signOut(): Promise<void> {
    await AsyncStorage.multiRemove(['access_token', 'refresh_token', 'user']);
  }

  async getGlobalLeaderboard(): Promise<ApiResponse<GlobalLeaderboardEntry[]>> {
    return this.request<GlobalLeaderboardEntry[]>('/leaderboard/global');
  }
}

export const api = new ApiClient(); 