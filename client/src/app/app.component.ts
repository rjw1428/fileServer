import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, filter, map, skip, startWith, Subject, switchMap, take, tap } from 'rxjs';
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
  selectedFile$ = new Subject<{ name: string, path: any, mediaType: string } | void>()
  constructor(
    private fileService: FileService,
    private route: ActivatedRoute,
    private router: Router,
  ) {
    this.route.queryParamMap.pipe(
      filter(params => !!params.keys.length && params.has('path') && !!params.get('path')),
      map(params => params.get('path')!)
    ).subscribe(path => {
      console.log(path)
      this.path$.next(path)
    })

    this.path$.pipe(skip(1)).subscribe(path => this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { path },
    }))
  }

  loadNext(file: { name: string, is_dir: boolean, path: any, mediaType: string } | null, event: Event) {
    event.preventDefault()
    this.selectedFile$.next()
    console.log(file)

    if (!file) {
      return this.path$.next('')
    }
    if (file.is_dir) {
      this.path$.next(file.path + '/' + file.name)
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

  getFileSource(file: { name: string, path: any, mediaType: string }) {
    return this.fileService.getFileSource(file.name, file.path)
  }

  uploadFile(event: any): void {
    const files: FileList = event.target.files
    console.log(files)
    const len = files.length
    for (let i = 0; i < len; i++) {
      const file = files.item(i)!
      const formData = new FormData();
      formData.append('file', file)
      this.fileService.uploadFile(formData).subscribe(console.log)
    }
  }

  createFolder(folderName: string) {
    this.path$.pipe(
      take(1),
      switchMap(path => this.fileService.createFolder(folderName, path).pipe(
        tap(console.log),
        filter(resp => !!resp.success),
        map(() => path)
      ))
    ).subscribe(path => this.path$.next(path))
    
  }
}
