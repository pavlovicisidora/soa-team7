import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { TourCreationDto, TourService } from '../tour.service';
import { Router } from '@angular/router';
import { Validators } from 'ngx-editor';

@Component({
  selector: 'app-tour-create',
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent implements OnInit {
  tourForm!: FormGroup;

  constructor(private fb: FormBuilder, private TourService: TourService, private router: Router){}
  ngOnInit(): void {
    this.tourForm = this.fb.group({
      name:['', Validators.required],
      description: ['',Validators.required],
      difficulty: ['easy',Validators.required],
      tags: ['']
    })
  }
  onSubmit():void{
    if(this.tourForm.invalid){
      console.error('Form is not valid');
      return;
    }
    const formValue = this.tourForm.value;
    const tourData: TourCreationDto= {
      name: formValue.name,
      description: formValue.description,
      difficulty: formValue.difficulty,
      tags: formValue.tags.split(',').map((tag: string) => tag.trim())
    };
    this.TourService.createTour(tourData).subscribe({
      next:(createdTour) => {
        console.log('Tour created:',createdTour);
        this.router.navigate(['/tour-list']);
      },
      error: (err) => {
        console.error('An error occurred while creating the tour.',err)
      }
    });
  }
}
