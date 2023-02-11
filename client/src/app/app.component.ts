import { Component } from '@angular/core';
import { BehaviorSubject, startWith, Subject, switchMap } from 'rxjs';
import { FileInfo } from 'src/model/file-info';
import { FileService } from './file.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'client';
  path$ = new BehaviorSubject<string>('')
  fileList$ = this.path$.pipe(
    switchMap(path => this.fileService.getFilelist(path).pipe(startWith([])))
  )
  selectedFile$ = new Subject<{name: string, path: any, mediaType: string} | void>()
  constructor(private fileService: FileService) {
  }

  loadNext(file: {name: string, is_dir: boolean, path: any, mediaType: string}, event: Event) {
    event.preventDefault()
    this.selectedFile$.next()
    if (file.is_dir) {
      this.path$.next(file.path+'/'+file.name)
    }
    else {
      this.selectedFile$.next(file)
      console.log(file)
    }
  }

  getFileSource(file: {name: string, path: any, mediaType: string}) {
    return this.fileService.getFileSource(file.name, file.path)
  }

  uploadFile(event: any): void {
    const files: FileList = event.target.files
    console.log(files)
    const len = files.length
    for (let i = 0; i<len; i++) {
      const file = files.item(i)!
      const formData = new FormData();
      formData.append('file', file)
      this.fileService.uploadFile(formData).subscribe(console.log)
    }
  }
}
