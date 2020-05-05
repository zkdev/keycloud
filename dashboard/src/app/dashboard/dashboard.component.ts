import {Component, OnInit} from '@angular/core';
import {MatDialog, MatDialogConfig} from '@angular/material/dialog';
import {DialogComponent} from '../dialog/dialog.component';
import {MatSnackBar} from '@angular/material/snack-bar';
import {UsernamePasswordUrl} from '../models/username-password-url';
import {UserService} from '../services/user.service';
import {UsernameUrl} from '../models/username-url';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  header = ['id', 'username', 'password', 'url', 'Delete'];
  entries: UsernamePasswordUrl[] = [];

  constructor(
    private dialog: MatDialog,
    private popOver: MatSnackBar,
    private userService: UserService,
  ) {
    this.userService.getListOfPasswords().subscribe(
      resp => {
        if (resp.status === 200) {
          const body = JSON.parse(resp.body);
          body.forEach(item => {
            const newEntry = new UsernamePasswordUrl(item.Id, item.Username, item.Password, item.Url);
            this.entries.push(newEntry);
          });
        } else {
          this.popOver.open(`${resp.status}`, '', {duration: 2000});
        }
      }, error => {
        this.popOver.open(`Something went wrong! If this error persists, please contact us with the following error: ${error.error}`,
          '', {duration: 5000});
      }
    );
    console.log(this.entries);
  }

  ngOnInit() { }

  openDialog() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;
    dialogConfig.autoFocus = true;

    const dialogRef = this.dialog.open(DialogComponent, dialogConfig);

    dialogRef.afterClosed().subscribe(
      data => this.saveNewEntry(data)
    );
  }

  removeEntry(item) {
    const index = this.entries.indexOf(item);
    const body = new UsernameUrl(item.username, item.url);
    this.userService.deletePassword(body).subscribe(
      resp => {
        this.entries.splice(index, 1);
        this.popOver.open('Deleted!', '', {duration: 2000});
      }, error => {
        this.popOver.open(`Something went wrong! If this error persists, please contact us with the following error: ${error.error}`,
          '', {duration: 5000});
      }
    );
  }

  copyToClipboard(item) {
    const index = this.entries.indexOf(item);
    const password = this.entries[index].password;
    const selBox = document.createElement('textarea');
    selBox.style.position = 'fixed';
    selBox.style.left = '0';
    selBox.style.top = '0';
    selBox.style.opacity = '0';
    selBox.value = password;
    document.body.appendChild(selBox);
    selBox.focus();
    selBox.select();
    document.execCommand('copy');
    document.body.removeChild(selBox);
    this.popOver.open('Copied!', '', {duration: 2000});
  }

  saveNewEntry(data) {
    const newEntry = new UsernamePasswordUrl(
      'Reload to view id',
      data.username,
      data.password,
      data.url
    );
    this.userService.addPassword(newEntry).subscribe(
      resp => {
        if (resp.status === 200) {
          this.entries.push(newEntry);
          this.popOver.open('Saved', '', {duration: 2000});
        } else {
          this.popOver.open(`${resp.status}`, '', {duration: 2000});
        }
      }, error => {
        this.popOver.open(`Something went wrong! If this error persists, please contact us with the following error: ${error.error}`,
          '', {duration: 5000});
      }
    );
  }

}
