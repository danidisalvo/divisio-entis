import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule} from "@angular/common/http";
import {MatButtonModule} from '@angular/material/button';
import {MatCheckboxModule} from "@angular/material/checkbox";
import {MatDialogModule} from "@angular/material/dialog";
import {MatFormFieldModule} from "@angular/material/form-field";
import {MatIconModule} from '@angular/material/icon';
import {MatInputModule} from "@angular/material/input";
import {MatListModule} from '@angular/material/list';
import {MatMenuModule} from '@angular/material/menu';
import {MatSelectModule} from "@angular/material/select";
import {MatSliderModule} from "@angular/material/slider";
import {MatSlideToggleModule} from "@angular/material/slide-toggle";
import {MatToolbarModule} from '@angular/material/toolbar';
import {MAT_COLOR_FORMATS, NgxMatColorPickerModule, NGX_MAT_COLOR_FORMATS} from '@angular-material-components/color-picker';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {ConfirmDialogComponent} from './confirm-dialog/confirm-dialog.component';
import {EditNodeDialogComponent} from './edit-node-dialog/edit-node-dialog.component';
import {FileDownloadComponent} from './file-download/file-download.component';
import {FileUploadComponent} from './file-upload/file-upload.component';
import {GraphComponent} from './graph/graph.component';
import {HelpPageComponent} from './help-page/help-page.component';
import {HomeComponent} from './home/home.component';
import {PageNotFoundComponent} from './page-not-found/page-not-found.component';
import {NgModule} from "@angular/core";
import {ToastNoAnimationModule} from "ngx-toastr";

@NgModule({
  declarations: [
    AppComponent,
    ConfirmDialogComponent,
    EditNodeDialogComponent,
    FileDownloadComponent,
    FileUploadComponent,
    GraphComponent,
    HomeComponent,
    PageNotFoundComponent,
    HelpPageComponent,
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    MatButtonModule,
    MatCheckboxModule,
    MatDialogModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatListModule,
    MatMenuModule,
    MatSelectModule,
    MatSliderModule,
    MatSlideToggleModule,
    MatToolbarModule,
    NgxMatColorPickerModule,
    ReactiveFormsModule,
    ToastNoAnimationModule.forRoot(),
  ],
  providers: [
    { provide: MAT_COLOR_FORMATS, useValue: NGX_MAT_COLOR_FORMATS }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
