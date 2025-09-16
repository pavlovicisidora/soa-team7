export interface OrderItem {
  tour_id: number;
  tour_name: string;
  price: number;
}

export interface ShoppingCart {
  id: string;
  user_id: string;
  items: OrderItem[];
  total_price: number;
}
