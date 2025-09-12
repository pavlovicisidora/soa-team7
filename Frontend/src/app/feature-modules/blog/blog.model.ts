export interface Image {
  url: string;
}

export interface Blog {
  id: string;
  title: string;
  content: string;
  created_at: any;
  user_id: string;
  images: Image[];
  like_count: number;
}