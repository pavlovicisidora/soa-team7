import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogListComponent } from './feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from './feature-modules/blog/blog-create/blog-create.component';

const routes: Routes = [
  { path: 'blog-list', component: BlogListComponent },
  { path: 'create-blog', component: BlogCreateComponent }, 
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
