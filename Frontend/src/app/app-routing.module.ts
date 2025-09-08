import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogListComponent } from './feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from './feature-modules/blog/blog-create/blog-create.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { TourCreateComponent } from './feature-modules/tour/tour-create/tour-create.component';
import { TourListComponent } from './feature-modules/tour/tour-list/tour-list.component';

const routes: Routes = [
  { path: 'blog-list', component: BlogListComponent },
  { path: 'create-blog', component: BlogCreateComponent }, 
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'navbar', component: NavbarComponent},
  { path: 'create-tour', component: TourCreateComponent},
  { path: 'tour-list', component: TourListComponent}, 

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
