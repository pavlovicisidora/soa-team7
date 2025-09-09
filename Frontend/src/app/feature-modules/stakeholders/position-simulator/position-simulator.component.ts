import { Component, OnInit, AfterViewInit } from '@angular/core';
import * as L from 'leaflet';
import { StakeholdersService } from '../stakeholders.service';
import { UserProfile } from '../user-profile.model';

const iconRetinaUrl = 'assets/marker-icon-2x-blue.png';
const iconUrl = 'assets/marker-icon-blue.png';
const shadowUrl = 'assets/marker-shadow.png';
const iconDefault = L.icon({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  tooltipAnchor: [16, -28],
  shadowSize: [41, 41]
});
L.Marker.prototype.options.icon = iconDefault;

@Component({
  selector: 'app-position-simulator',
  templateUrl: './position-simulator.component.html',
  styleUrls: ['./position-simulator.component.css']
})
export class PositionSimulatorComponent implements AfterViewInit {
  private map: any;
  private marker: any;
  public currentUser: UserProfile | null = null;

  constructor(private stakeholderService: StakeholdersService) { }

  ngAfterViewInit(): void {
    this.initMap();
    this.loadUserProfile();
  }

  private initMap(): void {
    this.map = L.map('map', {
      center: [44.7866, 20.4489],
      zoom: 13
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 18,
      attribution: '© OpenStreetMap'
    }).addTo(this.map);

    this.map.on('click', (e: L.LeafletMouseEvent) => {
      const coords = e.latlng;
      this.updatePosition(coords.lat, coords.lng);
    });
  }
  
  private loadUserProfile(): void {
    this.stakeholderService.getUser().subscribe({
      next: (profile) => {
        this.currentUser = profile;
        console.log('User successfully loaded:', profile);
        if (profile.latitude && profile.longitude) {
          this.addMarker(profile.latitude, profile.longitude);
          this.map.panTo([profile.latitude, profile.longitude]); 
        }
      },
      error: (err) => {
        console.error('Error loading user:', err);
      }
    });
  }

  private updatePosition(lat: number, lng: number): void {
    this.addMarker(lat, lng);
    
    this.stakeholderService.updatePosition({ lat: lat, long: lng }).subscribe({
      next: () => {
        console.log('Position successfully saved in database.');
        if (this.currentUser) {
          this.currentUser.latitude = lat;
          this.currentUser.longitude = lng;
        }
      },
      error: (err) => console.error('Error saving position:', err)
    });
  }
  
  private addMarker(lat: number, lng: number): void {
    if (this.marker) {
      this.map.removeLayer(this.marker);
    }
    this.marker = L.marker([lat, lng])
      .addTo(this.map)
      .bindPopup(`Your selected location.`)
      .openPopup();
  }
}