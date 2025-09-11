import { Component, OnInit, OnDestroy, ViewChild, ElementRef, Input } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { TourService } from '../tour.service'; // Pretpostavka da imate KeypointService
import { Keypoint } from '../keypoint.model';
import { ActivatedRoute } from '@angular/router';


@Component({
  selector: 'app-keypoint-manage',
  templateUrl: './keypoint-manage.component.html',
  styleUrls: ['./keypoint-manage.component.css']
})
export class KeypointManageComponent implements OnInit, OnDestroy {
  private tourId: number | undefined; // Ulazni parametar za ID ture
  @ViewChild('mapContainer', { static: true }) mapContainer!: ElementRef;

  private map: L.Map | undefined;
  private markers: L.Marker[] = [];
  private routingControl: L.Routing.Control | undefined;
  private unsubscribe$: Subject<void> = new Subject<void>();

  keypoints: Keypoint[] = [];
  keypointForm: FormGroup;
  editingKeypointId: number | null = null; // Koristi se za praćenje da li je u režimu izmene

  constructor(
    private fb: FormBuilder,
    private keypointService: TourService,
    private route: ActivatedRoute // Ubacivanje ActivatedRoute
  ) {
    this.keypointForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
      longitude: ['', [Validators.required, Validators.pattern(/^-?\d+\.?\d*$/)]],
      latitude: ['', [Validators.required, Validators.pattern(/^-?\d+\.?\d*$/)]],
      image: ['']
    });
  }

  ngOnInit(): void {
    // Čitanje tourId iz URL parametra
    this.route.paramMap.pipe(
      takeUntil(this.unsubscribe$)
    ).subscribe(params => {
      const id = params.get('tourId');
      if (id) {
        this.tourId = +id; // Konvertuj string u number
        this.initializeMap(); // Inicijalizacija mape nakon što dobijemo tourId
        this.loadKeypointsForTour(this.tourId);
      } else {
        console.error('Tour ID not provided in route parameters.');
        // Možete ovde preusmeriti korisnika ili prikazati poruku o grešci
      }
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
    if (this.map) {
      this.map.remove();
    }
  }

  initializeMap(): void {
    this.map = L.map(this.mapContainer.nativeElement).setView([44.7866, 20.4489], 10); // Centrirano na Beograd, zum 10

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(this.map);

    this.map.on('click', (e: L.LeafletMouseEvent) => {
      this.keypointForm.patchValue({
        latitude: e.latlng.lat,
        longitude: e.latlng.lng
      });
    });
  }

  loadKeypointsForTour(tourId: number): void {
    // Pozivanje servisa za dobijanje ključnih tačaka za datu turu
    this.keypointService.getKeypointsByTourId(tourId)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe({
        next: (keypoints: Keypoint[]) => {
          this.keypoints = keypoints;
          this.updateMapAndRoute();
        },
        error: (err: any) => console.error('Error loading keypoints:', err)
      });
  }

  updateMapAndRoute(): void {
    this.clearMarkersAndRoute();
    this.keypoints.forEach(kp => this.addKeypointMarker(kp));
    this.drawRoute();
  }

  addKeypointMarker(keypoint: Keypoint): void {
    const marker = L.marker([keypoint.latitude, keypoint.longitude]).addTo(this.map!)
      .bindPopup(`<b>${keypoint.name}</b><br>${keypoint.description}`);
    this.markers.push(marker);
  }

  drawRoute(): void {
    if (this.routingControl) {
      this.map!.removeControl(this.routingControl);
    }

    if (this.keypoints.length > 1) {
      const waypoints = this.keypoints.map(kp => L.latLng(kp.latitude, kp.longitude));
      this.routingControl = L.Routing.control({
        waypoints: waypoints,
        routeWhileDragging: true,
        showAlternatives: false,
        lineOptions: {
          styles: [{ color: 'red', weight: 4 }]
        },
        use: false,
        waypointIcon: false
      }).addTo(this.map!);
    }
  }

  clearMarkersAndRoute(): void {
    this.markers.forEach(marker => this.map!.removeLayer(marker));
    this.markers = [];
    if (this.routingControl) {
      this.map!.removeControl(this.routingControl);
      this.routingControl = undefined;
    }
  }

  onSaveKeypoint(): void {
    if (this.keypointForm.invalid) {
      alert('Molimo popunite sva obavezna polja.');
      return;
    }

    const newKeypointData = {
      ...this.keypointForm.value,
      tour_Id: this.tourId // Dodajemo tour_Id
    };

    if (this.editingKeypointId) {
      // Režim izmene
      this.keypointService.updateKeypoint(this.editingKeypointId, newKeypointData)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe({
          next: (updatedKeypoint: Keypoint) => {
            const index = this.keypoints.findIndex(kp => kp.id === updatedKeypoint.id);
            if (index !== -1) {
              this.keypoints[index] = updatedKeypoint;
            }
            this.clearFormAndEditingState();
            this.updateMapAndRoute();
          },
          error: (err: any) => console.error('Error updating keypoint:', err)
        });
    } else {
      // Režim kreiranja
      this.keypointService.createKeypoint(newKeypointData)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe({
          next: (createdKeypoint: Keypoint) => {
            this.keypoints.push(createdKeypoint);
            this.keypointForm.reset();
            this.updateMapAndRoute();
          },
          error: (err: any) => console.error('Error creating keypoint:', err)
        });
    }
  }

  onEditKeypoint(keypoint: Keypoint): void {
    this.editingKeypointId = keypoint.id;
    this.keypointForm.patchValue({
      name: keypoint.name,
      description: keypoint.description,
      longitude: keypoint.longitude,
      latitude: keypoint.latitude,
      image: keypoint.image
    });
  }

  onDeleteKeypoint(keypointId: number): void {
    if (confirm('Da li ste sigurni da želite da obrišete ovu ključnu tačku?')) {
      this.keypointService.deleteKeypoint(keypointId)
        .pipe(takeUntil(this.unsubscribe$))
        .subscribe({
          next: () => {
            this.keypoints = this.keypoints.filter(kp => kp.id !== keypointId);
            this.updateMapAndRoute();
          },
          error: (err : any) => console.error('Error deleting keypoint:', err)
        });
    }
  }

  clearFormAndEditingState(): void {
    this.keypointForm.reset();
    this.editingKeypointId = null;
  }
}