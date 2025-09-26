import { Component, OnInit, ChangeDetectorRef  } from '@angular/core';
import { ShoppingService } from '../shopping.service';
import { ShoppingCart, OrderItem } from '../shopping.model';

@Component({
  selector: 'app-shopping-cart',
  templateUrl: './shopping-cart.component.html',
  styleUrls: ['./shopping-cart.component.css']
})
export class ShoppingCartComponent implements OnInit {
  cart: ShoppingCart = {
    id: '',
    user_id: '',
    items: [],
    total_price: 0
  };
  isLoading = true;

  constructor(
    private shoppingService: ShoppingService,
    private cdr: ChangeDetectorRef
  ) { }

  ngOnInit(): void {
    this.loadCart();
  }

  loadCart(): void {
    this.isLoading = true;
    this.shoppingService.getCart().subscribe({
      next: (data) => {
        if (data) {
          this.cart = data;
        }
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Error loading cart:', err);
        this.isLoading = false;
      }
    });
  }

  removeItem(item: OrderItem): void {
    this.shoppingService.removeFromCart(item.tour_id).subscribe({
    next: (updatedCart) => {
      console.log('Item removed, updated cart:', updatedCart);
      this.cart = updatedCart;
      this.cdr.markForCheck();
    },
    error: (err) => {
      console.error('Error removing item:', err);
    }
  });
  }

  checkout(): void {
    this.shoppingService.checkout().subscribe({
      next: (tokens) => {
        alert(`Successful purchasing! You recieved ${tokens.length} token.`);
        this.loadCart();
      },
      error: (err) => alert('Error purchasing tour: ' + err.error)
    });
  }
}
