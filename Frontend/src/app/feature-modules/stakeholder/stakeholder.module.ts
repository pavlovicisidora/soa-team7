import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { UserListComponent } from './user-list/user-list.component';
import { ProfileComponent } from './profile/profile.component';
import { PositionSimulatorComponent } from './position-simulator/position-simulator.component';
import { FormsModule } from '@angular/forms';


@NgModule({
  declarations: [
    UserListComponent,
    ProfileComponent,
    PositionSimulatorComponent
  ],
  imports: [
    CommonModule,
    FormsModule
  ],
  exports: [
    UserListComponent,
    ProfileComponent,
    PositionSimulatorComponent
  ]
})
export class StakeholderModule { }
