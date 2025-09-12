import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { User } from 'src/app/auth/user.model';
import { UserAccount } from './user-account.model';
import { Profile } from 'src/app/auth/profile.model';
import { UserProfile } from './user-profile.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholderService {
  private apiUsersUrl = 'api/users';
  private apiProfileUrl = 'api/profile'; 
  private userUrl = 'api/users/user'; 
  private positionUrl = 'api/users/position';
  constructor(private http: HttpClient) { }
  getAllUsers() : Observable<UserAccount[]>{
    return this.http.get<UserAccount[]>(`${this.apiUsersUrl}`) 
  }
  fetchProfile(): Observable<Profile>{
    return this.http.get<Profile>(`${this.apiProfileUrl}`)
  }

  getUser(): Observable<UserProfile> {
    return this.http.get<UserProfile>(this.userUrl);
  }
  
  updatePosition(coords: { lat: number, long: number }): Observable<any> {
    return this.http.put(this.positionUrl, coords);
  }

   updateProfile(profile: Profile): Observable<Profile> {
    return this.http.patch<Profile>(`${this.apiProfileUrl}`, profile);
  }

  blockUser(username: string): Observable<any> {
  return this.http.post(`${this.apiUsersUrl}/block?username=${username}`, {});
}



}
