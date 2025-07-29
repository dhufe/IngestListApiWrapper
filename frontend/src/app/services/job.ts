import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {ApplicationConstants} from '../application-constants';
import {environment} from '../../environments/environment';
import {catchError, map, tap} from 'rxjs/operators';
import {Observable, of} from 'rxjs';
import {Job} from '../model/job';
import {User} from '../model/user';

@Injectable({
  providedIn: 'root'
})
export class JobService {
  token: string | null = ""

  constructor(private http: HttpClient,) {

  }

  public getJobById(id: string, user: User): Observable<Job> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        'Authorization': 'Bearer ' + user.token
      })
    }
    return this.http.get<Job>(environment.apiBaseUrl + "job/" + id, httpOptions)
  }

  public getJobs(user: User): Observable<Job[]> {
    const httpOptions = {
      headers: new HttpHeaders({
        'Authorization': 'Bearer ' + user.token,
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      })
    }
    return this.http.get<Job[]>(environment.apiBaseUrl + "jobs", httpOptions)
      /*
      .pipe(
        catchError(this.handleError<Job[]>('getJobs', []))
      );

       */
  }

  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {

      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}

