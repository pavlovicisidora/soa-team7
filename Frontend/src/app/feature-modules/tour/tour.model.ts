export interface Tour{
    id: number;
    name: string;
    description?: string;
    difficulty: string;
    tags: string[];
    status: 'DRAFT' | 'PUBLISHED' | 'ARCHIVED';
    price: number;
    authorId: string;
    distance_in_km: number;
    published_date_time?: any;
    archived_date_time?: any;
}