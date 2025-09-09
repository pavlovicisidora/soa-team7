import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../auth.service';
import { LoginResponse } from '../loginResponse.model';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  username: string = '';
  password: string = '';

  constructor(private authService: AuthService, private router: Router) {}

  login(): void {
  this.authService.loginUser(this.username, this.password).subscribe({
    next: (res: LoginResponse) => {
      console.log('Login successful:', res);

      // Sačuvaj token
      this.authService.setToken(res.token);

      // Sačuvaj username/role
      localStorage.setItem('role', res.role);
      localStorage.setItem('username', res.username);

      // Redirect nakon logina
      // this.router.navigate(['/navbar']);
    },
    error: (err) => {
      console.error('Login failed:', err);
      alert('Invalid username or password.');
    }
  });
}
}