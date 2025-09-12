import { Component, OnInit } from '@angular/core';
import { Tour } from '../tour.model';
import { TourService } from '../tour.service';

@Component({
  selector: 'app-public-tour-list',
  templateUrl: './public-tour-list.component.html',
  styleUrls: ['./public-tour-list.component.css']
})
export class PublicTourListComponent implements OnInit{
  tours: Tour[] = [];
  constructor(private tourService: TourService){}
  ngOnInit(): void {
    this.loadTours();
  }
  loadTours(): void{
    this.tourService.getAllTours().subscribe({
      next: (data) => {
        this.tours = data;

      },
      error: (err) => {
        console.error("An error occurred while fetching the tours:", err);
      }
    });
  }
}
