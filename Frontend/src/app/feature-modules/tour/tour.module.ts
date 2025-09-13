import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourCreateComponent } from './tour-create/tour-create.component';
import { TourListComponent } from './tour-list/tour-list.component';
import { ReactiveFormsModule } from '@angular/forms';
import { PublicTourListComponent } from './public-tour-list/public-tour-list.component';
import { TourDetailComponent } from './tour-detail/tour-detail.component';
import { RouterModule } from '@angular/router';
import { KeypointManageComponent } from './keypoint-manage/keypoint-manage.component';
import { TourExecutionComponent } from './tour-execution/tour-execution.component';



@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent,
    PublicTourListComponent,
    TourDetailComponent,
    KeypointManageComponent,
    TourExecutionComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule
  ],
  exports: [
    TourCreateComponent,
    TourListComponent
  ]
})
export class TourModule { }
