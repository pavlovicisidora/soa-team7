import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Route, Router } from '@angular/router';
import { Profile } from '../profile.model';
import { User } from '../user.model';
import { Editor } from 'ngx-editor';
import { AuthService } from '../auth.service';


@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  editor: Editor;

  profile: Profile = {
    name: '',
    surname: '',
    profilePic: '',
    bio: '',
    motto: ''
  };
  user: User = {
    username: '',
    password: '',
    mail: '',
    role: '',
    blocked: false,
    profile: this.profile,
  }    
  isUploading = false;

  constructor(private authService: AuthService, private router: Router) {
     this.editor = new Editor();
  }

  ngOnDestroy(): void {
    this.editor.destroy();
  }

  registerUser():void{
    if(this.isUploading){
      alert('Please wait for uploading images to be done.');
      return;
    }
    this.authService.registerUser(this.user).subscribe({
      next: (createdUser) =>{
        console.log('User successfully registered!',createdUser);
        this.router.navigate(['/login'])
      },
      error: (err) => {
        console.error('Error registrating user:', err);
      }
    });
  }

 onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      const file = input.files[0];
      this.isUploading = true;
      
      this.authService.uploadImage(file).subscribe({
        next: (response: { filePath: string }) => {
          console.log('Image successfully uploaded, URL path:', response.filePath);
          this.profile.profilePic = response.filePath;  
          this.isUploading = false;
        },
        error: (err) => {
          console.error('Error uploading image:', err);
          alert('Error occurred during image upload.');
          this.isUploading = false;
        }
      });
    }
  }
  
}