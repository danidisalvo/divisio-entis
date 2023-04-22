import {Component, OnInit} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-file-upload',
  templateUrl: './file-upload.component.html',
  styleUrls: ['./file-upload.component.css']
})
export class FileUploadComponent implements OnInit {

  fileName = '';
  message!: string | null;

  constructor(private http: HttpClient) {
  }

  ngOnInit(): void {
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];

    if (file) {
      this.fileName = file.name;
      const formData = new FormData();
      formData.append("file", file);

      this.http.post(`${environment.apiUrl}` + 'upload', formData).subscribe({
        next: () => {
          this.message = "File uploaded";
        },
        error: error => {
          console.error(error);
          this.message = error.status + ' ' + error.statusText;
        }
      });
    }
  }
}
