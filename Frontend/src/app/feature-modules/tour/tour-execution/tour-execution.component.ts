import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TourExecution } from '../tour-execution.model';
import { TourService } from '../tour.service';
import { UserProfile } from '../../stakeholder/user-profile.model';
import { StakeholderService } from '../../stakeholder/stakeholder.service';


@Component({
  selector: 'app-tour-execution',
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit {
  execution: TourExecution | null = null;
  currentUser: UserProfile | null = null;
  executionId: number;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: TourService,
    private stakeholderService: StakeholderService
  ) {
    this.executionId = Number(this.route.snapshot.paramMap.get('id'));
  }

  ngOnInit(): void {
    console.log('Active tour session ID:', this.executionId);
    this.loadExecutionDetails();
    this.loadUserPosition();
  }

  loadExecutionDetails(): void {
    this.tourService.getTourExecution(this.executionId).subscribe({
      next: (execData) => {
        if (execData.start_time && execData.start_time.seconds) {
          execData.start_time = new Date(execData.start_time.seconds * 1000);
        }
        this.execution = execData;
        console.log('Loaded execution details:', this.execution);
      },
      error: (err) => console.error('Error loading execution details:', err)
    });
  }

  loadUserPosition(): void {
    this.stakeholderService.getUser().subscribe({
      next: (profile) => {
        this.currentUser = profile;
      },
      error: (err) => console.error('Error loading user:', err)
    });
  }

  abandon(): void {
    if (!confirm('Are you sure you want to abandon tour?')) return;


    
    this.tourService.abandonTour(this.executionId).subscribe({
      next: (exec) => {
        console.log('Tour abandoned.', exec);
        alert('Tour successfully abandoned.');
        localStorage.removeItem('executionId');
        this.router.navigate(['/']);
      },
      error: (err) => console.error('Error:', err)
    });
  }

  complete(): void {
    this.tourService.completeTour(this.executionId).subscribe({
      next: (exec) => {
        console.log('Tour completed.', exec);
        alert('Successfully completed tour.');
        this.router.navigate(['/']);
      },
      error: (err) => console.error('Error:', err)
    });
  }
}
