import {Component, OnInit} from '@angular/core';
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {DomSanitizer} from "@angular/platform-browser";

@Component({
  selector: 'app-file-download',
  templateUrl: './file-download.component.html',
  styleUrls: ['./file-download.component.css']
})
export class FileDownloadComponent implements OnInit {

  message!: string | null;
  safeUrl!: any;

  constructor(private http: HttpClient, private sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
    this.message = null;
    const url = `${environment.apiUrl}`;
    this.http.get<Node>(url + "graph").subscribe({
      next: data => {
        const blob = new Blob([JSON.stringify(data, null, 2)], {type : 'application/json'});
        const url = URL.createObjectURL(blob);
        this.safeUrl = this.sanitizer.bypassSecurityTrustUrl(url);
      },
      error: error => {
        this.message = error.status + ' ' + error.statusText;
      }
    });
  }
}
