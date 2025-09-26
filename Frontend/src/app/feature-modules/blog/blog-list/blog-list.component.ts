import { Component, OnInit } from '@angular/core';
import { Blog } from '../blog.model';
import { BlogComment } from '../blog-comment.model';
import { BlogService } from '../blog.service';

@Component({
  selector: 'app-blog-list',
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.css']
})
export class BlogListComponent implements OnInit {
  blogs: Blog[] = []; 
  newCommentTexts: { [blog_id: string]: string } = {};
  commentsByBlog: { [blogId: string]: BlogComment[] } = {};
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
    if (updatedBlog.created_at && updatedBlog.created_at.seconds) {
      updatedBlog.created_at = new Date(updatedBlog.created_at.seconds * 1000);
    }
    
    this.blogs[index] = updatedBlog;
  }
}





// ******************************
  // Dodavanje komentara
  // ******************************
 submitComment(blog_id: string): void {
  const text = this.newCommentTexts[blog_id];
  if (!text || text.trim() === '') return;

  this.blogService.addCommentOnBlog(blog_id, text.trim())
    .subscribe((comment: BlogComment) => {
      console.log('Komentar dodat:', comment);
      this.newCommentTexts[blog_id] = '';

      // Dodaj odmah u prikaz komentara
      if (!this.commentsByBlog[blog_id]) {
        this.commentsByBlog[blog_id] = [];
      }
      this.commentsByBlog[blog_id].push(comment);
    }, error => {
      console.error('Greška pri dodavanju komentara:', error);
    });
}


  loadComments(blogId: string): void {
  this.blogService.getCommentsForBlog(blogId).subscribe(
    (comments: BlogComment[]) => {
      // konvertuj created_at u Date
      comments.forEach(c => {
        if (c.created_at && c.created_at.seconds) {
          c.created_at = new Date(c.created_at.seconds * 1000);
        }
      });
      this.commentsByBlog[blogId] = comments;
    },
    (error) => {
      console.error('Greška pri učitavanju komentara:', error);
    }
  );
}

}
