import {Component, OnInit} from '@angular/core';
import {JobComponent} from '../job';

@Component({
  selector: 'app-job-create',
  imports: [],
  templateUrl: './job-create.html',
  styleUrl: './job-create.scss'
})
export class JobCreate extends JobComponent implements OnInit {

  ngOnInit() {

  }

}
