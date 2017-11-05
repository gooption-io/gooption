import { TestBed, inject } from '@angular/core/testing';

import { GobsService } from './gobs.service';

describe('GobsService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [GobsService]
    });
  });

  it('should be created', inject([GobsService], (service: GobsService) => {
    expect(service).toBeTruthy();
  }));
});
