import { Component, OnDestroy, OnInit } from '@angular/core';
import { Http } from '@angular/http'
import { GoogleChartComponent } from '../google-chart/google-chart.component'

import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';
import { GobsService } from './../gobs.service';

@Component({
  selector: 'volsurface-chart',
  templateUrl: './volsurface-chart.component.html',
  styleUrls: ['./volsurface-chart.component.css']
})
export class VolsurfaceChartComponent implements OnInit {
  constructor(private gobs: GobsService) { }

  public data = [];
  private ngUnsubscribe: Subject<void> = new Subject<void>();

  private options1 = {
    title: 'Prices',
    curveType: 'function',
    legend: { position: 'bottom' }
  };
  private options2 = {
    title: 'Prices',
    curveType: 'function',
    legend: { position: 'bottom' }
  };

  private toIVChartData(slice: any) {
    var data_slice = [['Strike', 'Vol']];
    var quotes = slice.quotes.filter(q => !q.hasOwnProperty('error'));
    quotes.forEach(q => {
      data_slice.push([
        q.input.strike,
        q.vol
      ]);
    });
    
    return data_slice;
  }

  private toBidAskChartData(slice: any) {
    var data_slice = [['Strike', 'Bid', 'Ask']];
    var quotes = slice.quotes.filter(q => !q.hasOwnProperty('error'));
    quotes.forEach(q => {
      data_slice.push([
        q.input.strike,
        q.input.bid,
        q.input.ask        
      ]);
    });
    
    return data_slice;
  }

  ngOnInit() {
    this.gobs.impliedvol()
      .takeUntil(this.ngUnsubscribe)
      .subscribe(slices => {
        console.log(slices);
        slices.map((slice, idx) => {
          this.data.push({
            "id": idx,
            "slice": this.toIVChartData(slice),
            "bidask": this.toBidAskChartData(slice)
          })
        });
        console.log(this.data);
      });
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
