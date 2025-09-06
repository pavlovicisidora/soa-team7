import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'Frontend';
  
  ngOnInit(): void {
    // Nabavi validan token iz Postman-a i nalepi ga ovde dok ne budemo imali autentifikaciju
    const FAKE_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjhhOTliMzI5N2I5OWUyMTliMDk5Yjk2Iiwicm9sZSI6IlRVUklTVEEiLCJleHAiOjE3NTcxOTU3OTIsImlhdCI6MTc1NzEwOTM5Mn0.b8FwAlrh2uP5bNHrFvOTVWnJY65ufukq2p9fek8XqjI";
    
    // Sačuvaj ga u localStorage, odakle će ga JwtInterceptor pročitati
    localStorage.setItem('jwt_token', FAKE_TOKEN);

    console.log('!!! Privremeni JWT token je postavljen za razvoj! Ukloniti!!!');
  }
}

