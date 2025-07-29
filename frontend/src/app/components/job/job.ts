import {Component, OnInit, signal, ViewEncapsulation} from '@angular/core';
import {ComponentUtils} from '../component.utils';
import {MatButtonModule} from '@angular/material/button';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {FormsModule} from '@angular/forms';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from "@angular/material/icon";
import {JobService} from '../../services/job';

import {Job} from '../../model/job';
import {MatTableModule} from '@angular/material/table';


@Component({
  selector: 'app-job',
  templateUrl: './job.html',
  styleUrls: ['./job.scss'],
  encapsulation: ViewEncapsulation.None,
  imports: [MatCardModule, FormsModule, MatInputModule,
    MatFormFieldModule, MatButtonModule, MatIconModule, MatTableModule
  ],
  providers: [
    JobService
  ],
})
export class JobComponent extends ComponentUtils implements OnInit {

  jobs: Job[] = [];
  displayedColumns: string[] = ['Id', 'FileName', 'Type', 'Status'];

  constructor(private jobService: JobService) {
    super();
  }

  ngOnInit() {
    this.getJobs();
  }

  getJobs() {
    const user = this.currentUser
    if (user) {
      this.jobService.getJobs(user)
        .subscribe(jobs => this.jobs = jobs.slice(1, 20))
      console.log(this.jobs)
    }
  }
}
