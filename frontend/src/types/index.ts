export interface User {
  id: number;
  email: string;
  username: string;
  created_at: string;
  updated_at: string;
}

export interface Company {
  id: number;
  user_id: number;
  name: string;
  money: number;
  reputation: number;
  level: number;
  experience: number;
  created_at: string;
  updated_at: string;
  depots?: Depot[];
  buses?: Bus[];
}

export interface Depot {
  id: number;
  company_id: number;
  name: string;
  latitude: number;
  longitude: number;
  capacity: number;
  current_buses: number;
  level: number;
  created_at: string;
  updated_at: string;
}

export interface Bus {
  id: number;
  company_id: number;
  depot_id: number;
  name: string;
  type: string;
  capacity: number;
  fuel_capacity: number;
  current_fuel: number;
  range: number;
  service_type: string;
  status: 'available' | 'on_trip' | 'maintenance';
  condition: number;
  purchase_price: number;
  operating_cost: number;
  created_at: string;
  updated_at: string;
  depot?: Depot;
}

export interface Route {
  id: number;
  origin: string;
  destination: string;
  origin_lat: number;
  origin_lng: number;
  dest_lat: number;
  dest_lng: number;
  distance: number;
  duration: number;
  base_fare: number;
  popularity: number;
  type: 'intra_province' | 'inter_province';
}

export interface Driver {
  id: number;
  company_id: number;
  name: string;
  experience: number;
  salary: number;
  status: 'available' | 'on_trip' | 'rest';
  created_at: string;
  updated_at: string;
}

export interface Trip {
  id: number;
  bus_id: number;
  route_id: number;
  driver_id: number;
  status: 'planned' | 'active' | 'completed' | 'cancelled';
  passengers: number;
  revenue: number;
  cost: number;
  profit: number;
  progress: number;
  current_lat?: number;
  current_lng?: number;
  planned_start: string;
  actual_start?: string;
  actual_end?: string;
  created_at: string;
  updated_at: string;
  bus?: Bus;
  route?: Route;
  driver?: Driver;
}

export interface Transaction {
  id: number;
  company_id: number;
  type: 'income' | 'expense';
  description: string;
  amount: number;
  balance: number;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
  expires_in: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  username: string;
  password: string;
}

export interface CreateCompanyRequest {
  name: string;
}

export interface CreateDepotRequest {
  name: string;
  latitude: number;
  longitude: number;
}

export interface CreateBusRequest {
  name: string;
  type: string;
  capacity: number;
  service_type: string;
  purchase_price: number;
}

export interface CreateTripRequest {
  bus_id: number;
  route_id: number;
  driver_id?: number;
}

export interface WebSocketMessage {
  type: string;
  data: any;
  trip_id?: number;
  bus_id?: number;
}
