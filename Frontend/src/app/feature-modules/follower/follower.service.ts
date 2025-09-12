import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Follower } from './model/follower.model';


@Injectable({
  providedIn: 'root'
})
export class FollowerService {
  private apiUrl = 'http://localhost:8080/api/follower'; // promeni po potrebi

  constructor(private http: HttpClient) {}

  // ✅ Prati korisnika
  followUser(userId: string): Observable<any> {
    return this.http.post(`${this.apiUrl}/follow/${userId}`, {});
  }

  // ✅ Otprati korisnika
  unfollowUser(userId: string): Observable<any> {
    return this.http.delete(`${this.apiUrl}/unfollow/${userId}`);
  }

  // ✅ Vrati sve pratioce korisnika
  getFollowers(userId: string): Observable<Follower[]> {
    return this.http.get<Follower[]>(`${this.apiUrl}/followers/${userId}`);
  }

  // ✅ Vrati sve koje korisnik prati
  getFollowing(userId: string): Observable<Follower[]> {
    return this.http.get<Follower[]>(`${this.apiUrl}/following/${userId}`);
  }

  // ✅ Vrati preporučene korisnike
  getRecommendedUsers(): Observable<Follower[]> {
    return this.http.get<Follower[]>(`${this.apiUrl}/recommendations`);
  }
}
