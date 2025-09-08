import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourCreateComponent } from './tour-create/tour-create.component';
import { TourListComponent } from './tour-list/tour-list.component';
import { ReactiveFormsModule } from '@angular/forms';



@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent
  ],
  imports: [
    CommonModule,
    ReactiveFormsModule
  ],
  exports: [
    TourCreateComponent,
    TourListComponent
  ]
})
export class TourModule { }
