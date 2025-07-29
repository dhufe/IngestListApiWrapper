import {inject} from "@angular/core";
import {Router} from "@angular/router";
import {BreakpointObserver, Breakpoints} from "@angular/cdk/layout";

import {SnackbarService} from '../services/snackbar';
import {AuthService} from '../services/auth';

import {Observable} from "rxjs";
import {map, shareReplay} from "rxjs/operators";
import {User} from '../model/user';

import {ApplicationConstants} from '../application-constants';
import {App} from '../app';

export abstract class ComponentUtils {

  private breakpointObserver = inject(BreakpointObserver);

  isHandset$: Observable<boolean> = this.breakpointObserver.observe(Breakpoints.Handset)
    .pipe(
      map(result => result.matches),
      shareReplay()
    );

  private _currentUser: User | null;

  constructor() {
    this._currentUser = null;
  }

  protected _router = inject(Router);
  private _snackbarService = inject(SnackbarService);
  private _authService = inject(AuthService);

  public logout() {
    this._currentUser = null;
    this._authService.logOut();
    this._router.navigate(['/']).then();
  }

  get currentUser(): User | null {
    if (!this._currentUser) {
      this._currentUser = this.getCurrentUser();
    }
    return this._currentUser;
  }

  private getCurrentUser(): User|null {
    const userJson: string|null = localStorage.getItem(ApplicationConstants.CURRENT_USER_TOKEN);
    if (userJson) {
      return JSON.parse(userJson);
    }
    return null;
  }

  protected refreshCurrentUser(refreshUser: User|null): User | null {
    if (refreshUser) {
      localStorage.setItem(ApplicationConstants.CURRENT_USER_TOKEN, JSON.stringify(refreshUser));
    }
    return this.getCurrentUser();
  }

  public showSuccessMessage(message: string, duration?: number) {
    this._snackbarService.add(message, 'OK', {panelClass: ['green-snackbar'], duration: duration});
  }

  public showErrorMessage(message: string, duration?: number) {
    this._snackbarService.add(message, 'OK', {panelClass: ['red-snackbar'], duration: duration});
  }

  public getDateWithoutTime() {
    const date = new Date();
    date.setHours(0, 0, 0, 0);
    return date;
  }

}
