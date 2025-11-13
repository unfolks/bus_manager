import axios, { AxiosInstance, AxiosResponse } from 'axios';
import {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  Company,
  Depot,
  Bus,
  Route,
  Trip,
  CreateCompanyRequest,
  CreateDepotRequest,
  CreateBusRequest,
  CreateTripRequest,
} from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

class ApiService {
  private api: AxiosInstance;

  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add auth token to requests
    this.api.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    // Handle auth errors
    this.api.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          localStorage.removeItem('user');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Authentication
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/login', credentials);
    return response.data;
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/register', userData);
    return response.data;
  }

  async logout(): Promise<void> {
    await this.api.post('/auth/logout');
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  async refreshToken(): Promise<AuthResponse> {
    const response: AxiosResponse<AuthResponse> = await this.api.post('/auth/refresh');
    return response.data;
  }

  // Company Management
  async getCompany(): Promise<Company> {
    const response: AxiosResponse<Company> = await this.api.get('/company');
    return response.data;
  }

  async createCompany(companyData: CreateCompanyRequest): Promise<Company> {
    const response: AxiosResponse<Company> = await this.api.post('/company', companyData);
    return response.data;
  }

  // Depot Management
  async getDepots(): Promise<Depot[]> {
    const response: AxiosResponse<Depot[]> = await this.api.get('/depots');
    return response.data;
  }

  async createDepot(depotData: CreateDepotRequest): Promise<Depot> {
    const response: AxiosResponse<Depot> = await this.api.post('/depots', depotData);
    return response.data;
  }

  // Bus Management
  async getBuses(): Promise<Bus[]> {
    const response: AxiosResponse<Bus[]> = await this.api.get('/buses');
    return response.data;
  }

  async createBus(busData: CreateBusRequest): Promise<Bus> {
    const response: AxiosResponse<Bus> = await this.api.post('/buses', busData);
    return response.data;
  }

  // Route Management
  async getRoutes(): Promise<Route[]> {
    const response: AxiosResponse<Route[]> = await this.api.get('/routes');
    return response.data;
  }

  // Trip Management
  async getActiveTrips(): Promise<Trip[]> {
    const response: AxiosResponse<Trip[]> = await this.api.get('/trips/active');
    return response.data;
  }

  async createTrip(tripData: CreateTripRequest): Promise<Trip> {
    const response: AxiosResponse<Trip> = await this.api.post('/trips', tripData);
    return response.data;
  }

  // Utility methods
  setAuthToken(token: string): void {
    localStorage.setItem('token', token);
  }

  getAuthToken(): string | null {
    return localStorage.getItem('token');
  }

  setUser(user: User): void {
    localStorage.setItem('user', JSON.stringify(user));
  }

  getUser(): User | null {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  }

  isAuthenticated(): boolean {
    return !!this.getAuthToken();
  }
}

export const apiService = new ApiService();
export default apiService;
