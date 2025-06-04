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

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}
