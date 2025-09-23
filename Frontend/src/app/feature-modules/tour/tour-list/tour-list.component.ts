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
  editingTourId: number | null = null;
  tourToEdit: Tour | null = null;
  tourStatusOptions: ('DRAFT' | 'PUBLISHED' | 'ARCHIVED')[] = ['DRAFT', 'PUBLISHED', 'ARCHIVED'];
  // Ubacivanje Router-a u konstruktor
  constructor(private tourService: TourService, private router: Router){}

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
      this.tourService.getToursByAuthor().subscribe({
        next: (data) => {
      // Mapiranje u lokalni model gde su datumi pravi JS Date
      this.tours = data.map(tour => ({
        ...tour,
        archived_date_time: tour.archived_date_time
          ? new Date(tour.archived_date_time.seconds * 1000)
          : null,
        published_date_time: tour.published_date_time
          ? new Date(tour.published_date_time.seconds * 1000)
          : null,
      }));
    },
    error: (err) => {
      console.error("An error occurred while fetching the tours:", err);
    }
  })
  }


  // Poziva se klikom na "Update"
  startEdit(tour: Tour): void {
    this.editingTourId = tour.id;
    // Kreiramo plitku kopiju objekta. Ovo je ključno da izmene
    // u formi ne utiču odmah na prikaz u listi.
    this.tourToEdit = { ...tour }; 
  }

  // Poziva se klikom na "Cancel"
  cancelEdit(): void {
    this.editingTourId = null;
    this.tourToEdit = null;
  }
  

  // Poziva se klikom na "Save"
  saveUpdate(): void {
  if (!this.tourToEdit) {
    return;
  }

  const updated = {
    ...this.tourToEdit,
    archived_date_time: this.tourToEdit.archived_date_time
      ? this.dateToTimestamp(this.tourToEdit.archived_date_time)
      : null,
    published_date_time: this.tourToEdit.published_date_time
      ? this.dateToTimestamp(this.tourToEdit.published_date_time)
      : null,
  };

  this.tourService.updateTour(updated).subscribe({
    next: (updatedTour) => {
      const index = this.tours.findIndex(t => t.id === updatedTour.id);
      if (index !== -1) {
        this.tours[index] = {
          ...updatedTour,
          archived_date_time: updatedTour.archived_date_time
            ? new Date(updatedTour.archived_date_time.seconds * 1000)
            : null,
          published_date_time: updatedTour.published_date_time
            ? new Date(updatedTour.published_date_time.seconds * 1000)
            : null,
        };
      }
      this.cancelEdit();
    },
    error: (err) => {
      console.error("Failed to update tour:", err);
    }
  });
}
  // Nova metoda za redirekciju na Keypoints stranicu
  goToKeypoints(tourId: number): void {
    // Navigacija do 'keypoint-manage' rute, prosleđujući tourId kao parametar
    this.router.navigate(['/keypoint-manage', tourId]);
    // Alternativno, ako želite kao query parametar:
    // this.router.navigate(['/keypoint-manage'], { queryParams: { tourId: tourId } });
  }

  dateToTimestamp(date: Date | null): any {
    if (!date) return null;
    const seconds = Math.floor(date.getTime() / 1000);
    const nanos = (date.getTime() % 1000) * 1e6;
    return { seconds, nanos };
  }

}