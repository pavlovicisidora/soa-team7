import { Component, OnInit } from '@angular/core';
import { User } from 'src/app/auth/user.model';
import { StakeholderService } from '../stakeholder.service';
import { UserAccount } from '../user-account.model';

@Component({
  selector: 'app-user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})
export class UserListComponent implements OnInit {
  users: UserAccount[] = [];
  errorMessage: string | null = null;

  constructor(private stakeholderService: StakeholderService) { }

  ngOnInit(): void {
    this.loadUsers();
  }

  loadUsers(): void {
    this.stakeholderService.getAllUsers().subscribe({
      next: (data) => {
        this.users = data;
        this.errorMessage = null;
      },
      error: (err) => {
        console.error('Failed to fetch users:', err);
        this.errorMessage = 'You do not have permission to view this page. Administrator access is required.';
        this.users = []; 
      }
    });
  }

  blockUser(user: UserAccount): void {
    this.stakeholderService.blockUser(user.username).subscribe(
      () => {
        user.blocked = true; // odmah ažuriraj status u tabeli
      },
      (error) => console.error('Greška pri blokiranju korisnika', error)
    );
  }

}
