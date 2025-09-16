import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ShoppingCart } from './shopping.model';

@Injectable({
  providedIn: 'root'
})
export class ShoppingService {
  private apiUrl = '/api/shopping/cart'; 

  constructor(private http: HttpClient) { }

  getCart(): Observable<ShoppingCart> {
    return this.http.get<ShoppingCart>(this.apiUrl);
  }

  addToCart(tourId: number): Observable<ShoppingCart> {
    return this.http.post<ShoppingCart>(`${this.apiUrl}/${tourId}`, {});
  }

  removeFromCart(tourId: number): Observable<ShoppingCart> {
    return this.http.delete<ShoppingCart>(`${this.apiUrl}/${tourId}`);
  }

  checkout(): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/checkout`, {});
  }

  checkToken(tourId: number): Observable<{ hasToken: boolean }> {
    return this.http.get<{ hasToken: boolean }>(`/api/shopping/token/${tourId}`);
  }
}
