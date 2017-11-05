import { Injectable } from '@angular/core';
import { Http } from '@angular/http'
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/operator/map';

@Injectable()
export class GobsService {
  private query = `{
    marketdata(func: eq(timestamp, 1509883055), first:1) @cascade { 
      spot {
        ...indexInfo
      }
      vol  {
        ...indexInfo
      }
      rate  {
        ...indexInfo
      }
    } 

    quotes(func: eq(timestamp, 1509883055)) @cascade { 
      expiry
      puts {
      ...quote (orderasc: strike) 
      }
      calls {
      ...quote (orderasc: strike) 
      }
    } 
  }
    
  fragment quote {
    strike
    bid
    ask
    openinterest
  }

  fragment indexInfo {
    index @filter(eq(ticker, "AAPL") or eq(ticker, "USD.FEDFUND")) {
      timestamp
      ticker
      value
    }
  }`

  constructor(private http: Http) { }

  private toIVRequest(dgraph_data: any) : any {
    console.log(JSON.stringify(dgraph_data));    
    return {
      "pricingdate": 1508274400,                  
      "marketdata": {
        "spot": {
          "index": dgraph_data.marketdata[0].spot[0].index[0]
        },
        "rate": {
          "index": dgraph_data.marketdata[0].rate[0].index[0]
        }
      },
      "quotes": dgraph_data.quotes
    }
  }

  private sliceSubject = new BehaviorSubject([]);

  impliedvol()  {
    this.http.post("http://localhost:8080/query", this.query)
    .map(response => this.toIVRequest(response.json().data))
    .subscribe(iv_request => 
      this.http.post("http://localhost:8081/v1/gobs/impliedvol", iv_request)
      .subscribe(volsurf => {
        console.log(volsurf);       
        this.sliceSubject.next(volsurf.json().volsurface.slices);
      })
    );

    return this.sliceSubject.asObservable();
  }
}
