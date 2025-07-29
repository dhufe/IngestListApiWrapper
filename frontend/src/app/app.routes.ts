import { Routes } from '@angular/router';
import { LoginComponent} from './components/login/login';
import {JobComponent} from './components/job/job';

export const routes: Routes = [
  {path: 'login', component: LoginComponent},
  {path: 'jobs', component: JobComponent }
];
