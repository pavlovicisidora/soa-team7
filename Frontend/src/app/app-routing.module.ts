import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogListComponent } from './feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from './feature-modules/blog/blog-create/blog-create.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { TourCreateComponent } from './feature-modules/tour/tour-create/tour-create.component';
import { TourListComponent } from './feature-modules/tour/tour-list/tour-list.component';
import { PublicTourListComponent } from './feature-modules/tour/public-tour-list/public-tour-list.component';
import { TourDetailComponent } from './feature-modules/tour/tour-detail/tour-detail.component';
import { UserListComponent } from './feature-modules/stakeholder/user-list/user-list.component';
import { ProfileComponent } from './feature-modules/stakeholder/profile/profile.component';

const routes: Routes = [
  { path: 'blog-list', component: BlogListComponent },
  { path: 'create-blog', component: BlogCreateComponent }, 
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'navbar', component: NavbarComponent},
  { path: 'create-tour', component: TourCreateComponent},
  { path: 'tour-list', component: TourListComponent},
  { path: 'public-tour-list', component: PublicTourListComponent},
  { path: 'tours/:id', component: TourDetailComponent }, 
  { path: 'user-list', component: UserListComponent},
  { path: 'profile', component: ProfileComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
