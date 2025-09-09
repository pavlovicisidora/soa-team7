import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { User } from 'src/app/auth/user.model';
import { UserAccount } from './user-account.model';
import { Profile } from 'src/app/auth/profile.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholderService {
  private apiUserUrl = 'api/users';
  private apiProfileUrl = 'api/profile'; 
  constructor(private http: HttpClient) { }
  getAllUsers() : Observable<UserAccount[]>{
    return this.http.get<UserAccount[]>(`${this.apiUserUrl}`) 
  }
  fetchProfile(): Observable<Profile>{
    return this.http.get<Profile>(`${this.apiProfileUrl}`)
  }

}
