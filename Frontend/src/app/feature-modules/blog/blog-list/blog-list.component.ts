import { Component, OnInit } from '@angular/core';
import { Blog } from '../blog.model';
import { BlogService } from '../blog.service';

@Component({
  selector: 'app-blog-list',
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.css']
})
export class BlogListComponent implements OnInit {
  blogs: Blog[] = []; 

  constructor(private blogService: BlogService) { }

  ngOnInit(): void {
    this.loadBlogs();
  }

  loadBlogs(): void {
    this.blogService.getAllBlogs().subscribe(
      (data: Blog[]) => {
        this.blogs = data.map(blog => {
          if (blog.created_at && blog.created_at.seconds) {
            blog.created_at = new Date(blog.created_at.seconds * 1000);
          }
          return blog;
        });
        
        console.log('Successfully loaded blogs!', this.blogs);
      },
      (error) => {
        console.error('Error loading blogs:', error);
      }
    );
  }

  like(blogId: string): void {
    this.blogService.likeBlog(blogId).subscribe(
      (updatedBlog) => {
        console.log('Blog is liked', updatedBlog);
        this.updateBlogInList(updatedBlog);
      }
    );
  }
  
  unlike(blogId: string): void {
    this.blogService.unlikeBlog(blogId).subscribe(
      (updatedBlog) => {
        console.log('Blog is disliked', updatedBlog);
        this.updateBlogInList(updatedBlog);
      }
    );
  }

  private updateBlogInList(updatedBlog: Blog): void {
    const index = this.blogs.findIndex(b => b.id === updatedBlog.id);
    if (index !== -1) {
      this.blogs[index] = updatedBlog;
    }
  }
}
