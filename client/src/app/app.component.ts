import { Component } from '@angular/core';
import { BehaviorSubject, map, startWith, Subject, switchMap, take } from 'rxjs';
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

  loadNext(file: {name: string, is_dir: boolean, path: any, mediaType: string} | null, event: Event) {
    event.preventDefault()
    this.selectedFile$.next()
    if (!file) {
      return this.path$.next('')
    }
    if (file.is_dir) {
      this.path$.next(file.path+'/'+file.name)
    }
    else {
      this.selectedFile$.next(file)
    }
  }

  goBack() {
    this.path$.pipe(
      take(1),
      map(path => {
        const seg = path.split('/')
        return seg.slice(0, seg.length - 1).join('/')
      })
    ).subscribe(backPath => this.path$.next(backPath))
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
