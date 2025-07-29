import {Component, OnInit, signal, ViewEncapsulation} from '@angular/core';
import {User} from '../../model/user';
import {App} from '../../app';
import {Router, RouterLink} from '@angular/router';
import {ComponentUtils} from '../component.utils';
import {MatButtonModule} from '@angular/material/button';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {FormsModule} from '@angular/forms';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from "@angular/material/icon";
import {AuthService} from '../../services/auth';


@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  styleUrls: ['./login.scss'],
  encapsulation: ViewEncapsulation.None,
  imports: [MatCardModule, FormsModule, MatInputModule,
    MatFormFieldModule, MatButtonModule, MatIconModule
  ],
  providers: [
    AuthService
  ],
})
export class LoginComponent extends ComponentUtils implements OnInit {
  user: User = new User();
  constructor(private authService: AuthService, private router: Router,
              private appComponent: App) {
    super();
  }

  ngOnInit() {
  }

  hide = signal(true);
  clickEvent(event: MouseEvent) {
    let e = event as PointerEvent;
    if (e.pointerType === 'mouse') {
      this.hide.set(!this.hide());
      event.stopPropagation();
    }

  }

  login() {
    this.authService.login(this.user).subscribe({
      next: data => {
        if (data) {
          this.router.navigate(['/']).then();
          //this.appComponent.updateUser();
        } else {
          this.showErrorMessage('Ungültige E-Mail oder Passwort');
        }
      },
      error: err => {
        this.showErrorMessage('Fehler:  Anmeldung konnte nicht erfolgen');
      }
    });
  }

}
