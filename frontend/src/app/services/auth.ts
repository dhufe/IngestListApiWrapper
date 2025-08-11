import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {map} from 'rxjs/operators';
import {User} from '../model/user';
import {ApplicationConstants} from '../application-constants';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private http: HttpClient) {

  }

  public logOut() {

  }

  public login(user: User) {
    const httpOptions = {
      headers: new HttpHeaders({
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      })
    }

    const httpBody = {
      'email': user.email,
      'password': user.password
    }

    return this.http.post<any>(environment.apiBaseUrl + "login", httpBody, httpOptions).pipe(
      map((res: Response) => {
        const token = res;

        if (token) {
          localStorage.setItem(ApplicationConstants.CURRENT_USER_TOKEN, JSON.stringify(token));
        }

        return res;
      }));
  }
}
