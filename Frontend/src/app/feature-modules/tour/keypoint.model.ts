export interface Keypoint {
  id: number;
  name: string;
  description: string;
  longitude: number;
  latitude: number;
  image: string; // URL slike
  tour_Id: number;
}