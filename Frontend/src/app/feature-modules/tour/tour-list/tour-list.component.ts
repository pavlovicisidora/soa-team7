import { Component, OnInit } from '@angular/core';
import { Tour } from '../tour.model';
import { TourService } from '../tour.service';

@Component({
  selector: 'app-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours: Tour[]= [];
  constructor(private tourService: TourService){}

  ngOnInit(): void {
    this.loadTours();
  }
  loadTours(): void {
      this.tourService.getToursByAuthor().subscribe({
        next: (data) => {
          this.tours = data;
        },
        error: (err) => {
          console.error("An error occurred while fetching the tours:", err);
        }
      })
  }
}
