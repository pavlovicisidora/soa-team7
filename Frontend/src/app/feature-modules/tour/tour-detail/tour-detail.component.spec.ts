import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TourDetailComponent } from './tour-detail.component';

describe('TourDetailComponent', () => {
  let component: TourDetailComponent;
  let fixture: ComponentFixture<TourDetailComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [TourDetailComponent]
    });
    fixture = TestBed.createComponent(TourDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
