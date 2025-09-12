import { Component, OnInit } from '@angular/core';
import { FollowerService } from '../follower.service';
import { Follower } from '../model/follower.model';
@Component({
  selector: 'app-recommended-users',
  templateUrl: './recommended-users.component.html',
  styleUrls: ['./recommended-users.component.css']
})
export class RecommendedUsersComponent implements OnInit {
  recommendedUsers: Follower[] = [];

  constructor(private followerService: FollowerService) {}

  ngOnInit(): void {
    this.loadRecommendedUsers();
  }

  loadRecommendedUsers(): void {
    this.followerService.getRecommendedUsers().subscribe(
      (users: Follower[]) => {
        console.log(users);
        this.recommendedUsers = users;
      },
      (error) => {
        console.error('Greška pri učitavanju preporučenih korisnika:', error);
      }
    );
  }

  follow(userId: string): void {
    this.followerService.followUser(userId).subscribe(
      () => {
        console.log(`Pratiš korisnika ${userId}`);
        // opcionalno: odmah ukloni iz liste preporučenih
        this.recommendedUsers = this.recommendedUsers.filter(u => u.user_id !== userId);
      },
      (error) => {
        console.error('Greška pri praćenju:', error);
      }
    );
  }
}
