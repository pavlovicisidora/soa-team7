import * as L from 'leaflet';

declare module 'leaflet' {
  namespace Routing {
    interface RoutingControlOptions {
      waypoints?: L.LatLngExpression[];
      fitSelectedRoutes?: 'smart' | 'true' | 'false';
      routeWhileDragging?: boolean;
      router?: any; // Možete ovo precizirati ako znate tip routera
      formatter?: any; // Možete ovo precizirati
      geocoder?: any; // Možete ovo precizirati
      language?: string;
      unit?: string;
      use      : boolean;
      alternateRouteColors?: string[];
      showAlternatives?: boolean;
      waypointMode?: 'connect' | 'snap';
      show?: boolean;
      collapsed?: boolean;
      minimized?: boolean;
      createMarker?: (i: number, waypoint: L.LatLng, n: number) => L.Marker;
      routeLine?: (route: any, options: any) => L.Polyline; // Možete precizirati route i options
      containerClassName?: string;
      waypointIcon: boolean;
      lineOptions?: L.PolylineOptions & { styles?: L.PathOptions[] }; // Dodato styles za liniju rute
    }

    interface Control extends L.Control {
      setWaypoints(waypoints: L.LatLngExpression[]): Control;
      spliceWaypoints(index: number, n: number, ...waypoints: L.LatLngExpression[]): Control;
      getWaypoints(): L.LatLng[];
      on(type: string, fn: (e: any) => void, context?: any): this;
      off(type: string, fn: (e: any) => void, context?: any): this;
    }

    function control(options?: RoutingControlOptions): Control;
  }
}