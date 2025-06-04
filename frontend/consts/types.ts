import { NumberArray } from "react-native-svg";
import { Float, Int32 } from "react-native/Libraries/Types/CodegenTypes";

export interface SignUpRequest {
  username: string;
  email: string;
  password: string;
}

export interface SignInRequest {
  email: string;
  password: string;
}

export interface Data{
  session: Session
  email: string;
  external_user: string;
  user_id: string;
  username: string;
}


export interface Session {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  external_user: string;
}

export interface SignUpResponse {
  session: Session;
  user_id: string;
  external_user: string;
  username: string;
  email: string;
}

export interface SignInResponse {
  data: Data;
  session: Session;
  user_id: string;
  external_user: string;
  username: string;
  email: string;
}

export interface GetActivityStatsResponse {
  hexes_visited: number;
  activities_recorded: number;
  distance_covered: number;
  weekly_activities: number[];
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  success?: boolean;
}

export type GlobalLeaderboardEntry = {
  user_id: string;
  username?: string;
  top_count: number;
};

