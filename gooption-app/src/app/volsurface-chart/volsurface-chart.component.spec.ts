import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VolsurfaceChartComponent } from './volsurface-chart.component';

describe('VolsurfaceChartComponent', () => {
  let component: VolsurfaceChartComponent;
  let fixture: ComponentFixture<VolsurfaceChartComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VolsurfaceChartComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VolsurfaceChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should be created', () => {
    expect(component).toBeTruthy();
  });
});
