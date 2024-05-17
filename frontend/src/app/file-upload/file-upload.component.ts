import {Component, OnInit, ViewChild} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import { ToastContainerDirective, ToastrService } from 'ngx-toastr';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-file-upload',
  templateUrl: './file-upload.component.html',
  styleUrls: ['./file-upload.component.css']
})
export class FileUploadComponent implements OnInit {
  @ViewChild(ToastContainerDirective, { static: true })
  toastContainer: ToastContainerDirective | undefined;

  fileName = '';

  constructor(private http: HttpClient, private toastrService: ToastrService) {
  }

  ngOnInit(): void {
    this.toastrService.overlayContainer = this.toastContainer;
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];

    if (file) {
      this.fileName = file.name;
      const formData = new FormData();
      formData.append("file", file);

      this.http.post(`${environment.apiUrl}` + 'upload', formData).subscribe({
        next: () => {
          this.toastrService.success("File uploaded")
        },
        error: error => {
          this.toastrService.error(error.status + ' ' + error.statusTex)
        }
      });
    }
  }
}
