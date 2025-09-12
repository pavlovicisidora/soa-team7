import { ReviewUser } from "./review-user.model";

export interface Review {
  id: number;
  rating: number;
  comment: string;
  visitDate: string;
  createdDate: Date;
  tourist: ReviewUser;
  images?: string[];
}