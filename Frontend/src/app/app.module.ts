import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './components/home/home.component';
import { BlogModule } from "src/app/feature-modules/blog/blog.module";
import { RouterModule } from '@angular/router';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { JwtInterceptor } from './core/interceptors/jwt.interceptor';
import { NavbarComponent } from './components/navbar/navbar.component';
import { RegisterComponent } from './auth/register/register.component';
import { LoginComponent } from './auth/login/login.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    NavbarComponent,
    RegisterComponent,
    LoginComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    RouterModule,
    HttpClientModule,
    BlogModule
],
  providers: [
    { 
      provide: HTTP_INTERCEPTORS, 
      useClass: JwtInterceptor, 
      multi: true 
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
