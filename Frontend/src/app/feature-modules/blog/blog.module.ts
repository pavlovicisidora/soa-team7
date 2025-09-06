import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogListComponent } from './blog-list/blog-list.component';
import { BlogCreateComponent } from './blog-create/blog-create.component';
import { FormsModule } from '@angular/forms';
import { NgxEditorModule } from 'ngx-editor';


@NgModule({
  declarations: [
    BlogListComponent,
    BlogCreateComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    NgxEditorModule
  ],
  exports: [
    BlogListComponent,
    BlogCreateComponent
  ]
})
export class BlogModule { }
