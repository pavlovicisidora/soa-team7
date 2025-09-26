import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { TourExecution } from '../tour-execution.model';
import { TourService } from '../tour.service';
// Ažuriramo model da odgovara tvojoj definiciji
import { UserProfile } from '../../stakeholder/user-profile.model'; 
import { StakeholderService } from '../../stakeholder/stakeholder.service';
import 'leaflet-routing-machine';
import { Subject, interval, takeUntil } from 'rxjs';
import * as L from 'leaflet';
import { Keypoint } from '../keypoint.model';

@Component({
  selector: 'app-tour-execution',
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit, OnDestroy {
  execution: TourExecution | null = null;
 
  currentUser: UserProfile = { id: '', username: '', mail: '', role: '' }; //Kupimo usera
  executionId: number;

  private map: L.Map | undefined;
  private userMarker: L.Marker | undefined;
  private keypointMarkers: L.Marker[] = [];

  allKeypoints: Keypoint[] = [];
  completedKeypoints: Keypoint[] = [];
  uncompletedKeypoints: Keypoint[] = [];
  
  private isUpdatingExecution = false;
  private unsubscribe$: Subject<void> = new Subject<void>();

  constructor(
    private router: Router,
    private tourService: TourService,
    private stakeholderService: StakeholderService
  ) {
    const idFromStorage = localStorage.getItem('executionId');
    if (!idFromStorage) {
      console.error('Execution ID nije pronađen u local storage-u!');
      alert('Došlo je do greške, niste na aktivnoj turi.');
      this.router.navigate(['/']); 
      this.executionId = -1;
    } else {
      this.executionId = parseInt(idFromStorage, 10);
    }
  }

  ngOnInit(): void {
    if (this.executionId === -1) return;

    this.loadCompletedFromStorage();
    this.initializeMap();
    this.loadExecutionDetails();
    this.loadUserProfile(); 
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
    if (this.map) {
      this.map.remove();
    }
  }

  // --- FAZA INICIJALIZACIJE ---

  loadCompletedFromStorage(): void {
    const storedKeypoints = localStorage.getItem(`completedKeypoints_${this.executionId}`);
    if (storedKeypoints) {
      this.completedKeypoints = JSON.parse(storedKeypoints);
    }
  }

  initializeMap(): void {
    this.map = L.map('map').setView([45.255, 19.845], 13);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '© OpenStreetMap'
    }).addTo(this.map);

    // Klik na mapu sada samo ažurira poziciju na frontendu.
    // Interval će se pobrinuti da je pošalje na backend.
    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const lat = e.latlng.lat;
      const lng = e.latlng.lng;
      
      this.currentUser.latitude = lat;
      this.currentUser.longitude = lng;
      this.updateUserMarker(lat, lng);
    });
  }

  loadExecutionDetails(): void {
    this.tourService.getTourExecution(this.executionId).subscribe({
      next: (execData) => {
        if (execData.start_time && typeof execData.start_time === 'object' && execData.start_time.seconds) {
          execData.start_time = new Date(execData.start_time.seconds * 1000);
        }
        this.execution = execData;
        this.loadKeypointsForTour(this.execution.tour_id);
      },
      error: (err) => console.error('Greška pri učitavanju egzekucije:', err)
    });
  }
  
  // Učitavanje profila korisnika i postavljanje početne pozicije
  loadUserProfile(): void {
    this.stakeholderService.getUser().subscribe({
      next: (profile) => {
        this.currentUser = profile;
        console.log('Korisnik uspešno učitan:', profile);
        // ako ima lokacija sa simulatora nju prikayujem
        if (profile.latitude && profile.longitude) {
          this.updateUserMarker(profile.latitude, profile.longitude);
          this.map?.panTo([profile.latitude, profile.longitude]); 
        }
        //pokretanje intervala
        this.startExecutionUpdateInterval();
      },
      error: (err) => {
        console.error('Greška pri učitavanju korisnika:', err);
        this.startExecutionUpdateInterval();
      }
    });
  }

  loadKeypointsForTour(tourId: number): void {
    this.tourService.getKeypointsByTourId(tourId).pipe(takeUntil(this.unsubscribe$)).subscribe({
      next: (keypoints: Keypoint[]) => {
        this.allKeypoints = keypoints;
        this.updateKeypointListsAndMarkers();
      },
      error: (err: any) => console.error('Greška pri učitavanju ključnih tačaka:', err)
    });
  }

  // ---PERIODICNO AZURIRANJE ---
  startExecutionUpdateInterval(): void {
    interval(10000) // 10 sekundi
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        console.log("--- Tick 10s ---");
        // ako dodje do toga da nema pozicije
        if (!this.currentUser.latitude || !this.currentUser.longitude) {
          console.log("Pozicija korisnika nije postavljena, preskačem tick.");
          return;
        }
        
        // 1. Uvek ažuriraj poziciju korisnika u bazi
        this.updateUserPositionInDb(this.currentUser.latitude, this.currentUser.longitude);

        // 2. Proveri da li treba kompletirati tačku i ažuriraj execution
        this.checkProximityAndupdateExecution(this.currentUser.latitude, this.currentUser.longitude);
      });
  }

  updateUserPositionInDb(lat: number, lng: number): void {
    // Koristimo 'long' kako tvoj servis očekuje
    this.stakeholderService.updatePosition({ lat: lat, long: lng }).subscribe({
      next: () => console.log('Pozicija uspešno sačuvana u bazi.'),
      error: (err) => console.error('Greška pri čuvanju pozicije:', err)
    });
  }

  checkProximityAndupdateExecution(userLat: number, userLng: number): void {
    if (this.isUpdatingExecution) return;

    let keypointIdToComplete = 0; // Podrazumevana vrednost je 0

    for (const keypoint of this.uncompletedKeypoints) {
      const distance = this.calculateDistance(userLat, userLng, keypoint.latitude, keypoint.longitude);
      if (distance <= 100) {
        keypointIdToComplete = keypoint.id;
        break; 
      }
    }
    // Uvek pozivamo backend, bilo sa ID-jem tačke ili sa 0
    this.updateExecutionOnBackend(keypointIdToComplete);
  }

  // radimo update koji radi sve bukv
  updateExecutionOnBackend(keypointId: number): void {
    this.isUpdatingExecution = true;
    if(keypointId !== 0) {
      console.log(`Šaljem zahtev za kompletiranje tačke sa ID: ${keypointId}`);
    } else {
      console.log("Nema bliskih tačaka, šaljem zahtev za update aktivnosti.");
    }

    this.tourService.completeKeypoint(this.executionId, keypointId).subscribe({
      next: (updatedExecution) => {
        // Ako je tačka zaista kompletirana (ID nije 0)
        if (keypointId !== 0) {
            console.log('Ključna tačka uspešno kompletirana na backendu.');
            const completedKeypoint = this.uncompletedKeypoints.find(kp => kp.id === keypointId);
            if (completedKeypoint) {
                this.completedKeypoints.push(completedKeypoint);
                localStorage.setItem(`completedKeypoints_${this.executionId}`, JSON.stringify(this.completedKeypoints));
                this.updateKeypointListsAndMarkers(); // Odmah osveži UI
            }
        } else {
            console.log('Aktivnost uspešno ažurirana na backendu.');
        }

        if (this.execution) {
          this.execution.status = updatedExecution.status;
        }

        // Provera da li je cela tura gotova
        if (this.execution && this.execution.status === 'COMPLETED') {
          this.unsubscribe$.next(); // Zaustavljamo interval
          this.unsubscribe$.complete();
          alert('Čestitamo! Uspešno ste završili turu!');
          localStorage.removeItem(`completedKeypoints_${this.executionId}`);
          localStorage.removeItem('executionId');
          this.router.navigate(['/']);
        }
        
        this.isUpdatingExecution = false;
      },
      error: (err) => {
        console.error('Greška pri ažuriranju egzekucije:', err);
        this.isUpdatingExecution = false;
      }
    });
  }

  
  //Samo keypointi koje nismo presli uzimamo i iscrtavamo
  updateKeypointListsAndMarkers(): void {
    const completedIds = new Set(this.completedKeypoints.map(kp => kp.id));
    this.uncompletedKeypoints = this.allKeypoints.filter(kp => !completedIds.has(kp.id));

    this.keypointMarkers.forEach(marker => marker.removeFrom(this.map!));
    this.keypointMarkers = [];

    this.uncompletedKeypoints.forEach(kp => {
      const marker = L.marker([kp.latitude, kp.longitude])
        .addTo(this.map!)
        .bindPopup(`<b>${kp.name}</b><br>${kp.description}`);
      this.keypointMarkers.push(marker);
    });
  }

  updateUserMarker(lat: number, lng: number): void {
    const userIcon = L.icon({
        iconUrl: 'https://cdn-icons-png.flaticon.com/512/3017/3017122.png',
        iconSize: [35, 35],
    });

    if (this.userMarker) {
      this.userMarker.setLatLng([lat, lng]);
    } else {
      this.userMarker = L.marker([lat, lng], { icon: userIcon })
        .addTo(this.map!)
        .bindPopup('Tvoja trenutna lokacija')
        .openPopup();
    }
  }

  abandon(): void {
    if (!confirm('Da li ste sigurni da želite da napustite turu?')) return;
    this.tourService.abandonTour(this.executionId).subscribe({
      next: () => {
        alert('Tura je napuštena.');
        localStorage.removeItem(`completedKeypoints_${this.executionId}`);
        localStorage.removeItem('executionId');
        this.router.navigate(['/']);
      },
      error: (err) => console.error('Greška:', err)
    });
  }

  private calculateDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const R = 6371e3;
    const φ1 = lat1 * Math.PI / 180;
    const φ2 = lat2 * Math.PI / 180;
    const Δφ = (lat2 - lat1) * Math.PI / 180;
    const Δλ = (lon2 - lon1) * Math.PI / 180;
    const a = Math.sin(Δφ / 2) * Math.sin(Δφ / 2) + Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
    const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
  }
}