import { Component, OnInit } from '@angular/core';
import { Tour } from '../tour.model';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TourService } from '../tour.service';
import { ReviewCreationDto, ReviewService } from '../review.service';
import { forkJoin, of, switchMap } from 'rxjs';
import { Review } from '../review.model';
import { Validators } from 'ngx-editor';
import { StakeholderService } from '../../stakeholder/stakeholder.service';
import { ShoppingService } from '../../shopping/shopping.service';
import { AuthService } from 'src/app/auth/auth.service';
import { Keypoint } from '../keypoint.model';

@Component({
  selector: 'app-tour-detail',
  templateUrl: './tour-detail.component.html',
  styleUrls: ['./tour-detail.component.css']
})
export class TourDetailComponent implements OnInit{
  tour: Tour | undefined;
  reviews: Review[] = [];
  reviewForm!: FormGroup;
   selectedFiles: File[] = [];
  isUploading: boolean = false;
  isStartingTour = false;
  isTourist = false;
  hasPurchasedTour = false;

  firstKeypoint: Keypoint | null = null; 

  constructor(private route: ActivatedRoute, private tourService: TourService, private reviewService: ReviewService, private fb: FormBuilder, private stakeholderService: StakeholderService, private router: Router, private shoppingService: ShoppingService, private authService: AuthService){}
  ngOnInit(): void {
    this.isTourist = this.authService.isTourist();
    const tourId = Number(this.route.snapshot.paramMap.get('id'))
    if(tourId){
      this.loadTourAndReviews(tourId);
    }
    this.reviewForm = this.fb.group({
      rating: [5, [Validators.required]],
      comment: ['', Validators.required],
      visitingdate: [null, Validators.required] 
    });
    
    if (this.authService.getUsername != null) { 
      this.checkPurchaseStatus(tourId);
    }
  }

  checkPurchaseStatus(tourId: number): void {
    this.shoppingService.checkToken(tourId).subscribe(response => {
      this.hasPurchasedTour = response.hasToken;
    });
  }

  loadTourAndReviews(tourId: number): void {
    forkJoin({
      tour: this.tourService.getTourById(tourId),
      reviews: this.reviewService.getReviewsForTour(tourId),
      keypoints: this.tourService.getKeypointsByTourId(tourId)
    }).subscribe({
      next: ({ tour, reviews, keypoints }) => {
         this.tour = {
        ...tour,
        published_date_time: tour.published_date_time
          ? new Date((tour.published_date_time as any).seconds * 1000)
          : null,
        archived_date_time: tour.archived_date_time
          ? new Date((tour.archived_date_time as any).seconds * 1000)
          : null
      };
        this.reviews = reviews.map(review => {
        const timestamp = review.createdDate as any; 
        return { ...review,  createdDate: new Date(timestamp.seconds * 1000)  };
      }).sort((a, b) => b.createdDate.getTime() - a.createdDate.getTime());

      this.firstKeypoint = keypoints && keypoints.length > 0 ? keypoints[0] : null;
    },
      error: (err) => {
        console.error("Error loading the tour and reviews:", err);
      }
    });
  }
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files) {
      this.selectedFiles.push(...Array.from(input.files));
    }
  }
  removeFile(index: number): void {
    this.selectedFiles.splice(index, 1);
  }
  submitReview(): void {
    if (this.reviewForm.invalid || !this.tour) {
      this.reviewForm.markAllAsTouched();
      return;
    }
    this.isUploading = true;
    const reviewData: ReviewCreationDto = {
      ...this.reviewForm.value,
      rating: Number(this.reviewForm.value.rating)
    };
    const uploadObservables = this.selectedFiles.map(file => this.tourService.uploadImage(file));
    const imageUploads$ = uploadObservables.length > 0 ? forkJoin(uploadObservables) : of([]);
    imageUploads$.pipe(
      switchMap(uploadResults => {
        const imageUrls = uploadResults.map(result => result.filePath);
        const formValue = this.reviewForm.value;

        const reviewData: ReviewCreationDto = {
          rating: Number(formValue.rating),
          comment: formValue.comment,
          visitingdate: formValue.visitingdate,
          images: imageUrls
        };

        return this.reviewService.createReview(this.tour!.id, reviewData);
      })
    ).subscribe({
      next: (newReview) => {
        console.log('Review successfully added', newReview);
      

        if (this.tour) {
          this.loadTourAndReviews(this.tour.id); 
        }

        this.isUploading = false;
        this.selectedFiles = [];
        this.reviewForm.reset({
          rating: 5,
          comment: '',
          visitingdate: null
        });
      },
    error: (err) => {
      console.error('Error while adding the review', err)
      this.isUploading = false;
    }
  });
  }

  startTour(): void {
    if (!this.tour) return;

    if(localStorage.getItem('executionId') != null){
      alert('You already have one tour started!')
      return
    }
     this.isStartingTour = true;

    this.stakeholderService.getUser().subscribe({
      next: (profile) => {
        if (!profile.latitude || !profile.longitude) {
          alert('Please set your position in Position Simulator.');
          this.isStartingTour = false;
          this.router.navigate(['/position-simulator']);
          return;
        }

        const startLocation = {
          latitude: profile.latitude,
          longitude: profile.longitude
        };

        this.tourService.startTour(this.tour!.id, startLocation).subscribe({
          next: (execution) => {
            console.log('Successfully started the tour!', execution);
            this.isStartingTour = false;
            localStorage.setItem('executionId', execution.id.toString());
            this.router.navigate(['/tour-execution', execution.id]);
          },
          error: (err) => {
            console.error('Error starting tour:', err);
            this.isStartingTour = false;
          }
        });
      },
      error: (err) => {
        console.error('Error fetching user`s position:', err);
        this.isStartingTour = false;
      }
    });
  }

  addToCart(tourId: number): void {
    this.shoppingService.addToCart(tourId).subscribe({
      next: (cart) => alert(`${cart.items[cart.items.length - 1].tour_name} is added to cart!`),
      error: (err) => alert('Error adding to cart: ' + err.error)
    });
  }

  goToKeypoints(tourId: number): void {
    this.router.navigate(['/keypoint-manage', tourId]);
  }
}
