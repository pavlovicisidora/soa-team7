export interface TourExecution {
  id: number;
  tour_id: number;
  tourist_id: string;
  status: 'IN_PROGRESS' | 'COMPLETED' | 'ABANDONED';
}