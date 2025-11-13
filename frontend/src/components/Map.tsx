import React, { useEffect, useRef } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Fix for default markers in Leaflet with webpack
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon-2x.png',
  iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-icon.png',
  shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/images/marker-shadow.png',
});

interface MapProps {
  center: [number, number];
  zoom: number;
  depots?: Array<{
    id: number;
    name: string;
    latitude: number;
    longitude: number;
  }>;
  buses?: Array<{
    id: number;
    name: string;
    latitude: number;
    longitude: number;
    status: string;
  }>;
  routes?: Array<{
    id: number;
    origin_lat: number;
    origin_lng: number;
    dest_lat: number;
    dest_lng: number;
    origin: string;
    destination: string;
  }>;
  onDepotClick?: (depot: any) => void;
  onBusClick?: (bus: any) => void;
  onMapClick?: (lat: number, lng: number) => void;
}

const Map: React.FC<MapProps> = ({
  center,
  zoom,
  depots = [],
  buses = [],
  routes = [],
  onDepotClick,
  onBusClick,
  onMapClick,
}) => {
  const mapRef = useRef<HTMLDivElement>(null);
  const mapInstanceRef = useRef<L.Map | null>(null);
  const markersRef = useRef<L.Layer[]>([]);

  useEffect(() => {
    if (!mapRef.current || mapInstanceRef.current) return;

    // Initialize map
    const map = L.map(mapRef.current).setView(center, zoom);

    // Add tile layer (OpenStreetMap)
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
    }).addTo(map);

    mapInstanceRef.current = map;

    // Handle map click
    if (onMapClick) {
      map.on('click', (e: L.LeafletMouseEvent) => {
        onMapClick(e.latlng.lat, e.latlng.lng);
      });
    }

    return () => {
      if (mapInstanceRef.current) {
        mapInstanceRef.current.remove();
        mapInstanceRef.current = null;
      }
    };
  }, []);

  useEffect(() => {
    if (!mapInstanceRef.current) return;

    // Clear existing markers
    markersRef.current.forEach(marker => {
      mapInstanceRef.current?.removeLayer(marker);
    });
    markersRef.current = [];

    // Add depot markers
    depots.forEach(depot => {
      const marker = L.marker([depot.latitude, depot.longitude])
        .addTo(mapInstanceRef.current!)
        .bindPopup(`
          <div>
            <strong>${depot.name}</strong><br>
            Depot ID: ${depot.id}<br>
            <button onclick="window.selectDepot(${depot.id})">Select Depot</button>
          </div>
        `);

      if (onDepotClick) {
        marker.on('click', () => onDepotClick(depot));
      }

      markersRef.current.push(marker);
    });

    // Add bus markers
    buses.forEach(bus => {
      const busIcon = L.divIcon({
        html: `
          <div style="
            background-color: ${bus.status === 'available' ? '#4CAF50' : bus.status === 'on_trip' ? '#FF9800' : '#F44336'};
            color: white;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            font-weight: bold;
            border: 2px solid white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.3);
          ">
            🚌
          </div>
        `,
        className: 'bus-marker',
        iconSize: [20, 20],
        iconAnchor: [10, 10],
      });

      const marker = L.marker([bus.latitude, bus.longitude], { icon: busIcon })
        .addTo(mapInstanceRef.current!)
        .bindPopup(`
          <div>
            <strong>${bus.name}</strong><br>
            Status: ${bus.status}<br>
            Bus ID: ${bus.id}<br>
            <button onclick="window.selectBus(${bus.id})">Select Bus</button>
          </div>
        `);

      if (onBusClick) {
        marker.on('click', () => onBusClick(bus));
      }

      markersRef.current.push(marker);
    });

    // Add route lines
    routes.forEach(route => {
      const routeLine = L.polyline(
        [
          [route.origin_lat, route.origin_lng],
          [route.dest_lat, route.dest_lng]
        ],
        {
          color: '#2196F3',
          weight: 3,
          opacity: 0.7,
          dashArray: '10, 10',
        }
      ).addTo(mapInstanceRef.current!)
        .bindPopup(`
          <div>
            <strong>Route: ${route.origin} → ${route.destination}</strong><br>
            Route ID: ${route.id}
          </div>
        `);

      markersRef.current.push(routeLine);
    });

  }, [depots, buses, routes, onDepotClick, onBusClick]);

  // Make functions available globally for popup buttons
  useEffect(() => {
    (window as any).selectDepot = (depotId: number) => {
      const depot = depots.find(d => d.id === depotId);
      if (depot && onDepotClick) {
        onDepotClick(depot);
      }
    };

    (window as any).selectBus = (busId: number) => {
      const bus = buses.find(b => b.id === busId);
      if (bus && onBusClick) {
        onBusClick(bus);
      }
    };
  }, [depots, buses, onDepotClick, onBusClick]);

  return (
    <div
      ref={mapRef}
      style={{
        height: '100%',
        width: '100%',
        minHeight: '400px',
        borderRadius: '8px',
        overflow: 'hidden',
      }}
    />
  );
};

export default Map;
