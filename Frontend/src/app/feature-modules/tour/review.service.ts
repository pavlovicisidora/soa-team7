import { Injectable } from '@angular/core';
import { Review } from './review.model';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';

export interface ReviewCreationDto {
  rating: number;
  comment: string;
  visitingdate: string;
  images?: string[];
}

@Injectable({
  providedIn: 'root'
})
export class ReviewService {
  private apiUrl = 'api/reviews'; 

  constructor(private http: HttpClient) { }
  createReview(tourId: number, review: ReviewCreationDto): Observable<Review> {
    return this.http.post<Review>(`${this.apiUrl}/${tourId}`, review);
  }

  getReviewsForTour(tourId: number): Observable<Review[]> {
    return this.http.get<Review[]>(`${this.apiUrl}/${tourId}`);
  }
}
