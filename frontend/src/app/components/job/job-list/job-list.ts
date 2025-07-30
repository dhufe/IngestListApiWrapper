import {Component, OnInit, signal, ViewEncapsulation, ViewChild} from '@angular/core';
import {MatButtonModule} from '@angular/material/button';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';
import {FormsModule} from '@angular/forms';
import {MatCardModule} from '@angular/material/card';
import {MatIconModule} from "@angular/material/icon";
import {JobService} from '../../../services/job';
import {MatSort, MatSortModule, Sort} from '@angular/material/sort';
import {Job} from '../../../model/job';
import {MatTableDataSource, MatTableModule} from '@angular/material/table';
import {JobComponent} from '../job';
import {RouterModule} from '@angular/router';


@Component({
  selector: 'app-job',
  templateUrl: './job-list.html',
  styleUrls: ['./job-list.scss'],
  encapsulation: ViewEncapsulation.None,
  imports: [MatCardModule, FormsModule, MatInputModule,
    MatFormFieldModule, MatButtonModule, MatIconModule, MatTableModule,
    MatPaginatorModule, MatPaginator, MatSort, MatSortModule, RouterModule
  ],
  providers: [
    JobService
  ],
})
export class JobListComponent extends JobComponent implements OnInit {

  dataSource: MatTableDataSource<Job>;
  jobs: Job[] = [];
  displayedColumns: string[] = ['Id', 'FileName', 'Type', 'Status', 'Aktion'];

  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;

  constructor(private jobService: JobService) {
    super();


    // Assign the data to the data source for the table to render
    this.dataSource = new MatTableDataSource(this.jobs);
    this.paginator = new MatPaginator();

    this.sort = new MatSort();
  }

  ngOnInit() {
    this.getJobs();
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.paginator.pageSize = 10;
    this.dataSource.sort = this.sort;
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  announceSortChange(sortState: Sort) {
    if (sortState.active) {
      /*
      this.hostService.search(this.buildQueryParams(this.pageIndex, sortState.active, sortState.direction)).subscribe({
        next: data => {
          this.hostDataSource.data = data.resultList;
          this.totalItems = data.totalCount;
        },
        error: err => console.log(err)
      });
       */
    }
  }

  openEditJob(job: Job) {

  }

  openDetails(job: Job): void {
    // Erst Dialog erstellen
/*
    const dialogData = new DeleteDialogData('Job', 'den Job ' + deleteHost.hostname);
    const dialogRef = this.dialog.open(DeleteDialogComponent, {
      data: dialogData
    });
*/
    // Dann via Service Aktion ausführen
    /*
    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.hostService.delete(deleteHost).subscribe(
          (deleted) =>  {
            if (deleted) {
              this.search();
              this.showSuccessMessage('Host gelöscht.');
            } else {
              this.showErrorMessage('Host konnte nicht gelöscht werden.');
            }
          });
      }
    });
     */
  }

  getJobs() {
    const user = this.currentUser
    if (user) {
      this.jobService.getJobs(user)
        .subscribe(jobs => this.dataSource.data = jobs)
    }
  }
}
