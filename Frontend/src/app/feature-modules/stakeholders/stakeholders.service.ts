import { Injectable } from '@angular/core';
import { UserProfile } from './user-profile.model';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersService {
  private userUrl = 'http://localhost:8080/api/stakeholders/user'; 
  private positionUrl = 'http://localhost:8080/api/stakeholders/position';
  constructor(private http: HttpClient) { }
  
  getUser(): Observable<UserProfile> {
    return this.http.get<UserProfile>(this.userUrl);
  }
  
  updatePosition(coords: { lat: number, long: number }): Observable<any> {
    return this.http.put(this.positionUrl, coords);
  }
}
