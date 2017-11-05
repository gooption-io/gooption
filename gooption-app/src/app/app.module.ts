import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { HorizonService } from './horizon.service';
import { GobsService } from './gobs.service';

import { AppComponent } from './app.component';
import { GoogleChartComponent } from './google-chart/google-chart.component';
import { VolsurfaceChartComponent } from './volsurface-chart/volsurface-chart.component';

@NgModule({
  declarations: [
    AppComponent,
    GoogleChartComponent,
    VolsurfaceChartComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule
  ],
  providers: [HorizonService, GobsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
