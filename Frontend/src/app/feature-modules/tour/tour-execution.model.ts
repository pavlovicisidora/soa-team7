import { Tour } from "./tour.model";

export interface TourExecution {
  id: number;
  tour_id: number;
  tourist_id: string;
  status: 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';
  tour?: Tour;
  start_time?: any; 
}