import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { FileInfo } from 'src/model/file-info';
import { FileSizePipe } from 'src/pipes/file-size.pipe';

@Injectable({
  providedIn: 'root'
})
export class FileService {
  // uri = "http://test.ryanwilk.com"
  uri = 'http://localhost:4200'
  filesizePipe = new FileSizePipe()
  constructor(private http: HttpClient) { }

  getFilelist(path: string = '') {
    return this.http.get<FileInfo[]>(`${this.uri}/api/v1/files?path=${path}`).pipe(
      map(files => files
        .filter(file => !file.name.startsWith('.') && !file.name.startsWith('$'))
        .map(file => ({
          ...file,
          size: this.filesizePipe.transform(file.size),
          modified: new Date(file.modified * 1000),
          path,
          mediaType: this.mediaType(file.name)
        })))
    )
  }


  // getFile(filename: string, path: string) {
  //   return this.http.get(`${this.uri}/api/v1/download?$`)
  // }

  getFileSource(filename: string, path: string) {
    return `${this.uri}/api/v1/download?path=${encodeURIComponent(path)}&file=${encodeURIComponent(filename)}`
  }

  uploadFile(form: FormData) {
    return this.http.post(`${this.uri}/api/v1/upload`, form)
  }

  mediaType(filename: string) {
    if (filename.endsWith('.mp4'))
      return 'video'
    if (filename.endsWith('.jpg') || filename.endsWith('.jpeg') || filename.endsWith('.gif') || filename.endsWith('.png'))
      return 'image'
    else
      return 'unknown'
  }
}
