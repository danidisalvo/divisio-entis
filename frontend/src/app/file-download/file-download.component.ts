import {Component, OnInit} from '@angular/core';
import {environment} from "../../environments/environment";
import {HttpClient} from "@angular/common/http";
import {DomSanitizer, SafeUrl} from "@angular/platform-browser";

import {Node} from '../graph/graph.component';

@Component({
  selector: 'app-file-download',
  templateUrl: './file-download.component.html',
  styleUrls: ['./file-download.component.css']
})
export class FileDownloadComponent implements OnInit {

  message!: string | null;
  jblobSafeUrl!: SafeUrl;
  tblobSafeUrl!: SafeUrl;

  constructor(private http: HttpClient, private sanitizer: DomSanitizer) {
  }

  ngOnInit(): void {
    this.message = null;
    const url = `${environment.apiUrl}`;
    this.http.get<Node>(url + "graph").subscribe({
      next: graph => {
        const jblob = new Blob([JSON.stringify(graph, null, 2)], {type : 'application/json'});
        this.jblobSafeUrl = this.sanitizer.bypassSecurityTrustUrl(URL.createObjectURL(jblob));

        this.http.get(url + "graph/print", {responseType: 'text'}).subscribe(
          text => {
            const tblob = new Blob([text], {type : 'text/plain'});
            this.tblobSafeUrl = this.sanitizer.bypassSecurityTrustUrl(URL.createObjectURL(tblob));
          }
        );
      },
      error: error => {
        this.message = error.status + ' ' + error.statusText;
      }
    });
  }
}
