import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  username: string = '';
  password: string = '';

  constructor(private http: HttpClient) {}

  login() {
    const payload = {
      username: this.username,
      password: this.password
    };

    console.log('Logging in:', payload);

    this.http.post('http://localhost:8080/api/login', payload)
      .subscribe({
        next: (res) => {
          console.log('Login successful', res);
          // ovde možeš sačuvati token u localStorage/sessionStorage
          // npr. localStorage.setItem('token', res['token']);
        },
        error: (err) => console.error('Login error', err)
      });
  }
}