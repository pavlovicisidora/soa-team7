import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface Profile {
  name: string;
  surname: string;
  picture: string;
  bio: string;
  motto: string;
}

interface User {
  username: string;
  password: string;
  mail: string;
  role: string;
  blocked: boolean;
  profile: Profile;
}

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  user: User = {
    username: '',
    password: '',
    mail: '',
    role: 'VODIC',
    blocked: false,
    profile: {
      name: '',
      surname: '',
      picture: '',
      bio: '',
      motto: ''
    }
  };

  constructor(private http: HttpClient) {}

  register() {
    console.log('Registering user:', this.user);
    this.http.post('http://localhost:8080/api/register', this.user)
      .subscribe({
        next: (res) => console.log('User registered successfully', res),
        error: (err) => console.error('Error registering user', err)
      });
  }
}