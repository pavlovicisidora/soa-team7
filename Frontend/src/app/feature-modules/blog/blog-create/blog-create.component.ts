import { Component, OnDestroy } from '@angular/core';
import { BlogService, CreateBlogDTO } from '../blog.service';
import { Router } from '@angular/router';
import { Editor } from 'ngx-editor'; 

@Component({
  selector: 'app-blog-create',
  templateUrl: './blog-create.component.html',
  styleUrls: ['./blog-create.component.css']
})
export class BlogCreateComponent implements OnDestroy {
  editor: Editor;

  blog: CreateBlogDTO = {
    title: '',
    content: '', 
    images: []
  };
  isUploading = false;

  constructor(private blogService: BlogService, private router: Router) {
    this.editor = new Editor();
  }

  ngOnDestroy(): void {
    this.editor.destroy();
  }

  onContentChange(newContent: string): void {
    this.blog.content = newContent;
  }

  createBlog(): void {
    if (this.isUploading) {
      alert('Please wait for uploading images to be done.');
      return;
    }
    this.blogService.createBlog(this.blog).subscribe({
      next: (createdBlog) => {
        console.log('Blog successfully created!', createdBlog);
        this.router.navigate(['/blog-list']);
      },
      error: (err) => {
        console.error('Error creating blog:', err);
      }
    });
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      const file = input.files[0];
      this.isUploading = true;
      
      this.blogService.uploadImage(file).subscribe({
        next: (response: { filePath: string }) => {
          console.log('Image successfully uploaded, URL path:', response.filePath);
          this.blog.images.push({ url: response.filePath });
          this.isUploading = false;
        },
        error: (err) => {
          console.error('Error uploading image:', err);
          alert('Error occured during uploading image.');
          this.isUploading = false;
        }
      });
    }
  }
}
