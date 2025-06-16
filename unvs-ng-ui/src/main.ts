import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './app/app.config';
import { AppComponent } from './app/app.component'; // Import AppComponent thay vì App

bootstrapApplication(AppComponent, appConfig) // Bootstrap AppComponent
  .catch((err) => console.error(err));
