import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from './blog.model';
import { BlogComment } from './blog-comment.model';

export interface CreateBlogDTO {
  title: string;
  content: string; 
  images: { url: string }[];
}

@Injectable({
  providedIn: 'root'
})
export class BlogService {
  private apiUrl = '/api/blogs'; 
  private uploadUrl = '/api/images/upload'; 

  constructor(private http: HttpClient) { }

  getAllBlogs(): Observable<Blog[]> {
    return this.http.get<Blog[]>(this.apiUrl);
  }

  getBlogById(id: string): Observable<Blog> {
    return this.http.get<Blog>(`${this.apiUrl}/${id}`);
  }

  likeBlog(id: string): Observable<Blog> {
    return this.http.post<Blog>(`${this.apiUrl}/${id}/like`, {}); 
  }

  unlikeBlog(id: string): Observable<Blog> {
    return this.http.delete<Blog>(`${this.apiUrl}/${id}/like`);
  }

  createBlog(blogData: CreateBlogDTO): Observable<Blog> {
    return this.http.post<Blog>(this.apiUrl, blogData);
  }

  uploadImage(file: File): Observable<{ filePath: string }> {
    const formData = new FormData();
    formData.append('image', file, file.name);
    return this.http.post<{ filePath: string }>(this.uploadUrl, formData);
  }

  addCommentOnBlog(blogId: string, text: string): Observable<BlogComment> {
  return this.http.post<BlogComment>('/api/comments', { blog_id: blogId, text });
}

getCommentsForBlog(blogId: string): Observable<BlogComment[]> {
  return this.http.get<BlogComment[]>(`/api/comments/${blogId}`);
}


}