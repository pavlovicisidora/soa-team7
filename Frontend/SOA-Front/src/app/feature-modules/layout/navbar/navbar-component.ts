import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {
  
  currentUser: any; // Varijabla koja će držati korisničke podatke
  
  //constructor(private authService: AuthenticationService) {}

  ngOnInit(): void{
    
    
   // this.authService.currentUser$.subscribe( (user) => {this.currentUser = user; console.log(this.currentUser);});
    
    
  }

  logout() {
    //this.authService.logout();
  }
}
