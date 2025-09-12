import { Component } from '@angular/core';
import { Profile } from 'src/app/auth/profile.model';
import { StakeholderService } from '../stakeholder.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent {
profile: Profile | null = null;
  statusMessage: string = 'Loading profile...';

  showEditForm = false;
  editableProfile!: Profile;

  constructor(private stakeholderService: StakeholderService) { }

  ngOnInit(): void {
    this.loadProfile();
  }

  loadProfile(): void {
    this.stakeholderService.fetchProfile().subscribe({
      next: (data) => {
        this.profile = data;
        if (!data.name && !data.surname) {
          this.statusMessage = 'Profile not found or not yet created. Please edit your profile to add information.';
        }
      },
      error: (err) => {
        console.error('Failed to fetch profile:', err);
        if (err.status === 404) {
          this.statusMessage = 'Profile not found. Please edit your profile to create one.';
        } else {
          this.statusMessage = 'An error occurred while fetching your profile. Please try again later.';
        }
        this.profile = null;
      }
    });
  }

   toggleEdit() {
    this.showEditForm = !this.showEditForm;
    if (this.showEditForm && this.profile) {
      this.editableProfile = { ...this.profile };
    }
  }

  saveProfile() {
    this.stakeholderService.updateProfile(this.editableProfile).subscribe({
      next: (updated: Profile) => {
        this.profile = updated;
        this.showEditForm = false;
        this.loadProfile();
      },
      error: (err) => {
        console.error('Error updating profile', err);
      }
    });
  }
}
