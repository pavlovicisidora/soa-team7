import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PublicTourListComponent } from './public-tour-list.component';

describe('PublicTourListComponent', () => {
  let component: PublicTourListComponent;
  let fixture: ComponentFixture<PublicTourListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [PublicTourListComponent]
    });
    fixture = TestBed.createComponent(PublicTourListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
