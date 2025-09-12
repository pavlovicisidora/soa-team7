import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RecommendedUsersComponent } from './recommended-users/recommended-users.component';
import { FormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    RecommendedUsersComponent
  ],
  imports: [
    CommonModule, 
    FormsModule
  ],
  exports: [
    RecommendedUsersComponent
  ]
})
export class FollowerModule { }