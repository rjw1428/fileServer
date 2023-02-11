import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map } from 'rxjs';
import { environment } from 'src/environments/environment';
import { FileInfo } from 'src/model/file-info';
import { FileSizePipe } from 'src/pipes/file-size.pipe';

@Injectable({
  providedIn: 'root'
})
export class FileService {
  uri = environment.apiUrl
  filesizePipe = new FileSizePipe()
  constructor(private http: HttpClient) { }

  getFilelist(path: string = '') {
    return this.http.get<FileInfo[]>(`${this.uri}/api/v1/files?path=${path}`).pipe(
      map(files => (files || [])
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

  getFileSource(filename: string, path: string) {
    return `${this.uri}/api/v1/download?path=${encodeURIComponent(path)}&file=${encodeURIComponent(filename)}`
  }

  uploadFile(form: FormData) {
    return this.http.post(`${this.uri}/api/v1/upload`, form)
  }

  createFolder(folderName: string, path: string) {
    return this.http.get<{success: Boolean}>(`${this.uri}/api/v1/createFolder?path=${encodeURIComponent(path)}&file=${encodeURIComponent(folderName)}`)
  }

  mediaType(filename: string) {
    const name = filename.toLowerCase()
    if (name.endsWith('.mp4') || name.endsWith('.avi')) // avi doesnt play
      return 'video'
    if (name.endsWith('.mp3'))
      return 'audio'
    if (name.endsWith('.jpg') || name.endsWith('.jpeg') || name.endsWith('.gif') || name.endsWith('.png'))
      return 'image'
    else
      return 'unknown'
  }
}
