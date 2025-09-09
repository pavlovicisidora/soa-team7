import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserListComponent } from './user-list/user-list.component';
import { ProfileComponent } from './profile/profile.component';



@NgModule({
  declarations: [
    UserListComponent,
    ProfileComponent
  ],
  imports: [
    CommonModule
  ],
  exports: [
    UserListComponent,
    ProfileComponent
  ]
})
export class StakeholderModule { }
