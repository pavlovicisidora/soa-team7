import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourCreateComponent } from './tour-create/tour-create.component';
import { TourListComponent } from './tour-list/tour-list.component';
import { ReactiveFormsModule } from '@angular/forms';
import { PublicTourListComponent } from './public-tour-list/public-tour-list.component';
import { TourDetailComponent } from './tour-detail/tour-detail.component';
import { RouterModule } from '@angular/router';



@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent,
    PublicTourListComponent,
    TourDetailComponent
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
