import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';     
import { LoginResponse } from './loginResponse.model';
import { User } from './user.model';



@Injectable({
    providedIn: 'root'
})
export class AuthService{
    private uploadUrl = '/api/images/upload'; 
    
    constructor(private http: HttpClient) { }

    loginUser(username: string, password: string): Observable<LoginResponse> {
        return this.http.post<LoginResponse>('http://localhost:8081/login', { username, password });
    }

    registerUser(userData: User) :Observable<User>{
        return this.http.post<User>('http://localhost:8081/register', userData);
    }


    uploadImage(file: File): Observable<{ filePath: string }> {
        const formData = new FormData();
        formData.append('image', file, file.name);
        return this.http.post<{ filePath: string }>(this.uploadUrl, formData);
    }

    setToken(token: string): void {
        localStorage.setItem('jwt_token', token);
    }

    getToken(): string | null {
        return localStorage.getItem('jwt_token');
    }

    getUsername(): string | null {
         return localStorage.getItem('username');
    }

    getRole(): string | null {
        return localStorage.getItem('role');
    }

    logout(): void {
        localStorage.removeItem('jwt_token');
        localStorage.removeItem('username');
        localStorage.removeItem('role');
    }


}