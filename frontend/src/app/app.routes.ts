import { Routes } from '@angular/router';
import { LoginComponent} from './components/login/login';
import { JobListComponent} from './components/job/job-list/job-list';

export const routes: Routes = [
  {path: 'login', component: LoginComponent},
  {path: 'jobs', component: JobListComponent }
];
