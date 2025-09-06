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
    const FAKE_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjhhNGRjOWIwZTg5MzU5NTU1ZTQ0MDBkIiwicm9sZSI6IlRVUklTVEEiLCJleHAiOjE3NTcyNTI0NzUsImlhdCI6MTc1NzE2NjA3NX0.MHICkGqRI-lLn_P2T0_77uSBdIgCg9b3pr3yyzq7lL4";
    
    // Sačuvaj ga u localStorage, odakle će ga JwtInterceptor pročitati
    localStorage.setItem('jwt_token', FAKE_TOKEN);

    console.log('!!! Privremeni JWT token je postavljen za razvoj! Ukloniti!!!');
  }
}

