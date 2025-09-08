export interface Tour{
    id: number;
    name: string;
    description?: string;
    difficulty: string;
    tags: string[];
    status: 'DRAFT' | 'PUBLISHED' | 'ARCHIVED';
    price: number;
    authorId: string;
}