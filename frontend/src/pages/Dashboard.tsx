import React, { useState, useEffect } from 'react';
import Map from '../components/Map';
import apiService from '../services/api';
import { Company, Depot, Bus, Route, Trip } from '../types';

const Dashboard: React.FC = () => {
  const [company, setCompany] = useState<Company | null>(null);
  const [depots, setDepots] = useState<Depot[]>([]);
  const [buses, setBuses] = useState<Bus[]>([]);
  const [routes, setRoutes] = useState<Route[]>([]);
  const [activeTrips, setActiveTrips] = useState<Trip[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Map center (Java, Indonesia)
  const mapCenter: [number, number] = [-7.2575, 112.7521];
  const mapZoom = 7;

  useEffect(() => {
    fetchGameData();
  }, []);

  const fetchGameData = async () => {
    try {
      setLoading(true);
      const [companyData, depotsData, busesData, routesData, tripsData] = await Promise.all([
        apiService.getCompany(),
        apiService.getDepots(),
        apiService.getBuses(),
        apiService.getRoutes(),
        apiService.getActiveTrips(),
      ]);

      setCompany(companyData);
      setDepots(depotsData);
      setBuses(busesData);
      setRoutes(routesData);
      setActiveTrips(tripsData);
    } catch (err: any) {
      if (err.response?.status === 404) {
        // Company doesn't exist yet, show setup
        setCompany(null);
      } else {
        setError(err.message || 'Failed to load game data');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCreateCompany = async (companyName: string) => {
    try {
      const newCompany = await apiService.createCompany({ name: companyName });
      setCompany(newCompany);
      fetchGameData(); // Refresh all data
    } catch (err: any) {
      setError(err.message || 'Failed to create company');
    }
  };

  const handleCreateDepot = async (name: string, lat: number, lng: number) => {
    try {
      await apiService.createDepot({ name, latitude: lat, longitude: lng });
      fetchGameData(); // Refresh data
    } catch (err: any) {
      setError(err.message || 'Failed to create depot');
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl">Loading Bus Manager...</div>
      </div>
    );
  }

  if (!company) {
    return <CompanySetup onCreateCompany={handleCreateCompany} />;
  }

  if (depots.length === 0) {
    return <DepotSetup onCreateDepot={handleCreateDepot} />;
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-blue-600 text-white p-4 shadow-lg">
        <div className="container mx-auto flex justify-between items-center">
          <div>
            <h1 className="text-2xl font-bold">{company.name}</h1>
            <p className="text-sm opacity-90">Level {company.level} • {company.experience} XP</p>
          </div>
          <div className="flex gap-6">
            <div className="text-right">
              <p className="text-sm opacity-90">Balance</p>
              <p className="text-xl font-bold">{formatCurrency(company.money)}</p>
            </div>
            <div className="text-right">
              <p className="text-sm opacity-90">Reputation</p>
              <p className="text-xl font-bold">{company.reputation}</p>
            </div>
          </div>
        </div>
      </header>

      <div className="container mx-auto p-4">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Map */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-lg p-4">
              <h2 className="text-xl font-bold mb-4">Company Map</h2>
              <div className="h-96 lg:h-full min-h-96">
                <Map
                  center={mapCenter}
                  zoom={mapZoom}
                  depots={depots}
                  buses={buses.map(bus => ({
                    ...bus,
                    latitude: bus.depot?.latitude || mapCenter[0],
                    longitude: bus.depot?.longitude || mapCenter[1],
                  }))}
                  routes={routes}
                />
              </div>
            </div>
          </div>

          {/* Side Panel */}
          <div className="space-y-6">
            {/* Quick Stats */}
            <div className="bg-white rounded-lg shadow-lg p-4">
              <h3 className="text-lg font-bold mb-3">Quick Stats</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Depots:</span>
                  <span className="font-semibold">{depots.length}</span>
                </div>
                <div className="flex justify-between">
                  <span>Buses:</span>
                  <span className="font-semibold">{buses.length}</span>
                </div>
                <div className="flex justify-between">
                  <span>Active Trips:</span>
                  <span className="font-semibold">{activeTrips.length}</span>
                </div>
                <div className="flex justify-between">
                  <span>Available Buses:</span>
                  <span className="font-semibold text-green-600">
                    {buses.filter(b => b.status === 'available').length}
                  </span>
                </div>
              </div>
            </div>

            {/* Active Trips */}
            <div className="bg-white rounded-lg shadow-lg p-4">
              <h3 className="text-lg font-bold mb-3">Active Trips</h3>
              {activeTrips.length === 0 ? (
                <p className="text-gray-500">No active trips</p>
              ) : (
                <div className="space-y-2 max-h-64 overflow-y-auto">
                  {activeTrips.map(trip => (
                    <div key={trip.id} className="border-l-4 border-blue-500 pl-3 py-2">
                      <div className="font-semibold text-sm">
                        {trip.bus?.name} • {trip.route?.origin} → {trip.route?.destination}
                      </div>
                      <div className="text-xs text-gray-600">
                        Progress: {Math.round(trip.progress)}% • {trip.passengers} passengers
                      </div>
                      <div className="w-full bg-gray-200 rounded-full h-2 mt-1">
                        <div
                          className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                          style={{ width: `${trip.progress}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Quick Actions */}
            <div className="bg-white rounded-lg shadow-lg p-4">
              <h3 className="text-lg font-bold mb-3">Quick Actions</h3>
              <div className="space-y-2">
                <button className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition">
                  Dispatch Bus
                </button>
                <button className="w-full bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition">
                  Buy Bus
                </button>
                <button className="w-full bg-purple-600 text-white px-4 py-2 rounded hover:bg-purple-700 transition">
                  Upgrade Depot
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {error && (
        <div className="fixed bottom-4 right-4 bg-red-500 text-white px-6 py-3 rounded-lg shadow-lg">
          {error}
          <button
            onClick={() => setError(null)}
            className="ml-4 text-white hover:text-gray-200"
          >
            ×
          </button>
        </div>
      )}
    </div>
  );
};

// Company Setup Component
const CompanySetup: React.FC<{ onCreateCompany: (name: string) => void }> = ({ onCreateCompany }) => {
  const [companyName, setCompanyName] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!companyName.trim()) return;

    setIsSubmitting(true);
    await onCreateCompany(companyName);
    setIsSubmitting(false);
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
        <h1 className="text-3xl font-bold text-center mb-6 text-blue-600">
          Welcome to Bus Manager
        </h1>
        <p className="text-gray-600 text-center mb-8">
          Start your bus company empire in Indonesia!
        </p>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Company Name
            </label>
            <input
              type="text"
              value={companyName}
              onChange={(e) => setCompanyName(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter your company name"
              required
              minLength={3}
              maxLength={100}
            />
          </div>
          <button
            type="submit"
            disabled={isSubmitting || !companyName.trim()}
            className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition"
          >
            {isSubmitting ? 'Creating...' : 'Create Company'}
          </button>
        </form>
      </div>
    </div>
  );
};

// Depot Setup Component
const DepotSetup: React.FC<{ onCreateDepot: (name: string, lat: number, lng: number) => void }> = ({ onCreateDepot }) => {
  const [depotName, setDepotName] = useState('');
  const [selectedLocation, setSelectedLocation] = useState<{ lat: number; lng: number } | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleMapClick = (lat: number, lng: number) => {
    setSelectedLocation({ lat, lng });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!depotName.trim() || !selectedLocation) return;

    setIsSubmitting(true);
    await onCreateDepot(depotName, selectedLocation.lat, selectedLocation.lng);
    setIsSubmitting(false);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto p-4">
        <div className="bg-white rounded-lg shadow-xl p-6 max-w-4xl mx-auto">
          <h1 className="text-3xl font-bold text-center mb-6 text-blue-600">
            Place Your First Depot
          </h1>
          <p className="text-gray-600 text-center mb-8">
            Click on the map to select a location for your first depot
          </p>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <Map
                center={[-7.2575, 112.7521]}
                zoom={7}
                onMapClick={handleMapClick}
              />
              {selectedLocation && (
                <div className="mt-4 p-3 bg-green-100 border border-green-400 rounded-lg">
                  <p className="text-sm text-green-800">
                    Selected: {selectedLocation.lat.toFixed(4)}, {selectedLocation.lng.toFixed(4)}
                  </p>
                </div>
              )}
            </div>

            <div>
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Depot Name
                  </label>
                  <input
                    type="text"
                    value={depotName}
                    onChange={(e) => setDepotName(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Enter depot name"
                    required
                    minLength={3}
                    maxLength={100}
                  />
                </div>

                <div className="bg-blue-50 p-4 rounded-lg">
                  <h3 className="font-semibold text-blue-900 mb-2">Starting Bonus:</h3>
                  <ul className="text-sm text-blue-800 space-y-1">
                    <li>• Starting capital: Rp 1,000,000</li>
                    <li>• Free small bus (25 seats)</li>
                    <li>• Basic depot (10 bus capacity)</li>
                  </ul>
                </div>

                <button
                  type="submit"
                  disabled={isSubmitting || !depotName.trim() || !selectedLocation}
                  className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition"
                >
                  {isSubmitting ? 'Creating Depot...' : 'Create Depot'}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
