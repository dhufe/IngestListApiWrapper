import {Routes} from '@angular/router';
import {LoginComponent} from './components/login/login';
import {JobListComponent} from './components/job/job-list/job-list';
import {JobDetails} from './components/job/job-details/job-details';

export const routes: Routes = [
  {path: 'login', component: LoginComponent},
  {path: 'jobs', component: JobListComponent},
  {path: 'jobs/editJob/:id', component: JobDetails}
];
