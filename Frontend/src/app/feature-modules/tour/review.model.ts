export interface Review {
  id: number;
  rating: number;
  comment: string;
  visitDate: string;
  createdDate: Date;
  touristName: string;
  images?: string[];
}