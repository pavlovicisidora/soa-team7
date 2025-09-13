import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TourExecution } from '../tour-execution.model';
import { TourService } from '../tour.service';


@Component({
  selector: 'app-tour-execution',
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit {
  execution: TourExecution | null = null;
  executionId: number;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: TourService
  ) {
    this.executionId = Number(this.route.snapshot.paramMap.get('id'));
  }

  ngOnInit(): void {
    console.log('Active tour session ID:', this.executionId);
  }

  abandon(): void {
    if (!confirm('Are you sure you want to abandon tour?')) return;
    
    this.tourService.abandonTour(this.executionId).subscribe({
      next: (exec) => {
        console.log('Tour abandoned.', exec);
        alert('Tour successfully abandoned.');
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
