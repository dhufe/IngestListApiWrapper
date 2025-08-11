import {Component, OnInit} from '@angular/core';
import {JobComponent} from '../job';
import {Job} from '../../../model/job';
import {JobService} from '../../../services/job';
import {ActivatedRoute} from '@angular/router';
import {FormsModule} from '@angular/forms';
import {MatFormField, MatInputModule, MatLabel} from '@angular/material/input';

@Component({
  selector: 'app-job-details',
  imports: [
    FormsModule,
    MatInputModule,
    MatFormField,
    MatLabel
  ],
  templateUrl: './job-details.html',
  styleUrl: './job-details.scss',
})
export class JobDetails extends JobComponent implements OnInit {

  job: Job = new Job();


  constructor(private route: ActivatedRoute, private jobService: JobService) {
    super();
    this.route.params.subscribe(params => {
      this.job = params['id'];
      if (params && params['id']) {
        const user = this.currentUser
        if (user) {
          this.jobService.getJobById(params['id'], user)
            .subscribe(job => this.job = job);
        }
      }
    });
  }

  ngOnInit() {

  }

  updateHost() {

  }
}



