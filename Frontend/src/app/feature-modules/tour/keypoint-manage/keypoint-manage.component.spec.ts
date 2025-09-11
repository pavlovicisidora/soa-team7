import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KeypointManageComponent } from './keypoint-manage.component';

describe('KeypointManageComponent', () => {
  let component: KeypointManageComponent;
  let fixture: ComponentFixture<KeypointManageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [KeypointManageComponent]
    });
    fixture = TestBed.createComponent(KeypointManageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
