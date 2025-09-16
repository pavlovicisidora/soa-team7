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
import { AuthModule } from "src/app/auth/auth.module";
import { TourModule } from './feature-modules/tour/tour.module';
import { StakeholderModule } from './feature-modules/stakeholder/stakeholder.module';
import { FollowerModule } from './feature-modules/follower/follower.module';
import { ShoppingModule } from './feature-modules/shopping/shopping.module';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    NavbarComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    RouterModule,
    HttpClientModule,
    BlogModule,
    AuthModule,
    TourModule,
    StakeholderModule,
    BlogModule,
    FollowerModule,
    ShoppingModule
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
