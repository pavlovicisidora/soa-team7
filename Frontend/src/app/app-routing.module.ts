import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogListComponent } from './feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from './feature-modules/blog/blog-create/blog-create.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { NavbarComponent } from './components/navbar/navbar.component';

const routes: Routes = [
  { path: 'blog-list', component: BlogListComponent },
  { path: 'create-blog', component: BlogCreateComponent }, 
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'navbar', component: NavbarComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
