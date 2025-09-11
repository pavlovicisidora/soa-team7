import { Component, OnInit } from '@angular/core';
import { Tour } from '../tour.model';
import { TourService } from '../tour.service';
import { Router } from '@angular/router'; // Importovanje Router-a

@Component({
  selector: 'app-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours: Tour[]= [];

  // Ubacivanje Router-a u konstruktor
  constructor(private tourService: TourService, private router: Router){}

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

  // Nova metoda za redirekciju na Keypoints stranicu
  goToKeypoints(tourId: number): void {
    // Navigacija do 'keypoint-manage' rute, prosleđujući tourId kao parametar
    this.router.navigate(['/keypoint-manage', tourId]);
    // Alternativno, ako želite kao query parametar:
    // this.router.navigate(['/keypoint-manage'], { queryParams: { tourId: tourId } });
  }
}