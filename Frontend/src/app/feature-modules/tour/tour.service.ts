import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Tour } from './tour.model';
import { Keypoint } from './keypoint.model';

export interface TourCreationDto {
  name: string;
  description: string;
  difficulty: string; 
  tags: string[];
}
@Injectable({
  providedIn: 'root'
})
export class TourService {
  private apiUrl = "api/tours"
  private uploadUrl = '/api/images/upload'; 

  constructor(private http : HttpClient) { }
  getToursByAuthor(): Observable<Tour[]>{
    return this.http.get<Tour[]>(this.apiUrl);
  }
  createTour(tourData:TourCreationDto): Observable<Tour>{
    return this.http.post<Tour>(this.apiUrl,tourData);
  }
  getAllTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/all`);
  }
  getTourById(id: number): Observable<Tour> {
    return this.http.get<Tour>(`${this.apiUrl}/${id}`);
  }
  uploadImage(file: File): Observable<{ filePath: string }> {
  const formData = new FormData();
  formData.append('image', file, file.name);
  return this.http.post<{ filePath: string }>(this.uploadUrl, formData);
  }

   getKeypointsByTourId(tourId: number): Observable<Keypoint[]> {
    return this.http.get<Keypoint[]>(`http://localhost:8080/api/keypoints/${tourId}`);
  }

  createKeypoint(keypoint: Omit<Keypoint, 'id'>): Observable<Keypoint> {
    return this.http.post<Keypoint>(`http://localhost:8080/api/keypoints/${keypoint.tour_Id}`, keypoint);
  }

  updateKeypoint(id: number, keypoint: Keypoint): Observable<Keypoint> {
    return this.http.put<Keypoint>(`http://localhost:8080/api/keypoints/${id}`,keypoint);
  }

  deleteKeypoint(id: number): Observable<void> {
    return this.http.delete<void>(`http://localhost:8080/api/keypoints/${id}`);
  }
}
