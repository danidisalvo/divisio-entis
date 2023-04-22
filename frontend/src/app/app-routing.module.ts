import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {GraphComponent} from './graph/graph.component';
import {HelpPageComponent} from './help-page/help-page.component';
import {HomeComponent} from './home/home.component';
import {FileUploadComponent} from "./file-upload/file-upload.component";
import {FileDownloadComponent} from "./file-download/file-download.component";
import {PageNotFoundComponent} from "./page-not-found/page-not-found.component";

const routes: Routes = [
  {
    path: '',
    component: HomeComponent
  },
  {
    path: 'graph',
    component: GraphComponent
  },
  {
    path: 'download',
    component: FileDownloadComponent
  },
  {
    path: 'upload',
    component: FileUploadComponent
  },
  {
    path: 'help',
    component: HelpPageComponent
  },
  {
    path: '**',
    component: PageNotFoundComponent
  }
];

const isIframe = window !== window.parent && !window.opener;

@NgModule({
  imports: [RouterModule.forRoot(routes, {
    initialNavigation: !isIframe ? 'enabledBlocking' : 'disabled' // Don't perform initial navigation in iframes
  })],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
