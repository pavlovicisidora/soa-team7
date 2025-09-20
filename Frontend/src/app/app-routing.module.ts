import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogListComponent } from './feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from './feature-modules/blog/blog-create/blog-create.component';
import { PositionSimulatorComponent } from './feature-modules/stakeholder/position-simulator/position-simulator.component';
import { LoginComponent } from './auth/login/login.component';
import { RegisterComponent } from './auth/register/register.component';
import { NavbarComponent } from './components/navbar/navbar.component';
import { TourCreateComponent } from './feature-modules/tour/tour-create/tour-create.component';
import { TourListComponent } from './feature-modules/tour/tour-list/tour-list.component';
import { PublicTourListComponent } from './feature-modules/tour/public-tour-list/public-tour-list.component';
import { TourDetailComponent } from './feature-modules/tour/tour-detail/tour-detail.component';
import { UserListComponent } from './feature-modules/stakeholder/user-list/user-list.component';
import { ProfileComponent } from './feature-modules/stakeholder/profile/profile.component';
import { KeypointManageComponent } from './feature-modules/tour/keypoint-manage/keypoint-manage.component';
import { RecommendedUsersComponent } from './feature-modules/follower/recommended-users/recommended-users.component';
import { TourExecutionComponent } from './feature-modules/tour/tour-execution/tour-execution.component';
import { ShoppingCartComponent } from './feature-modules/shopping/shopping-cart/shopping-cart.component';

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
  { path: 'position-simulator', component: PositionSimulatorComponent },
  { path: 'keypoint-manage/:tourId', component: KeypointManageComponent },
  { path: 'recommended-users', component: RecommendedUsersComponent },
  { path: 'tour-execution/:id', component: TourExecutionComponent },
  { path: 'cart', component: ShoppingCartComponent },
  
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
