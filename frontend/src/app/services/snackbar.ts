import {inject, Injectable, OnDestroy} from '@angular/core';
import {MatSnackBar, MatSnackBarConfig, MatSnackBarRef, SimpleSnackBar} from '@angular/material/snack-bar';

export class SnackBarMessage  {
  message: string="";
  action: string|undefined;
  config: MatSnackBarConfig|undefined;
}

@Injectable({
  providedIn: 'root'
})

export class SnackbarService implements OnDestroy {
  private snackBarRef: MatSnackBarRef<SimpleSnackBar> | undefined;
  private msgQueue: Array<SnackBarMessage> = [];
  private isInstanceVisible = false;
  private duration: number = 10000 //ms;

  private _snackBar = inject(MatSnackBar);

  constructor(){
  }

  showNext() {
    if (this.msgQueue.length === 0) {
      return;
    }
    let message = this.msgQueue.shift();
    if (message) {
      this.isInstanceVisible = true;
      let config = null;
      if (message.config) {
        config = message.config;
        if (!config.duration) {
          config.duration = this.duration;
        }
      } else {
        config = {duration: this.duration};
      }
      this.snackBarRef = this._snackBar.open(message.message, message.action, config);
      this.snackBarRef.afterDismissed().subscribe(() => {
        this.isInstanceVisible = false;
        this.showNext();
      });
    }
  }

  ngOnDestroy() {
    // this.subscription.unsubscribe();
  }

  /**
   * Add a message
   * @param message The message to show in the snackbar.
   * @param action The label for the snackbar action.
   * @param config Additional configuration options for the snackbar.
   */
  add(message: string, action?: string, config?: MatSnackBarConfig): void{
    let sbMessage = new SnackBarMessage();
    sbMessage.message = message;
    sbMessage.action = action;
    sbMessage.config = config;

    this.msgQueue.push(sbMessage);
    if (!this.isInstanceVisible) {
      this.showNext();
    }
  }
}
