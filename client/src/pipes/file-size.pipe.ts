import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'fileSize'
})
export class FileSizePipe implements PipeTransform {

  transform(size: number) {
    const kb = size / 1024
    if (kb < 500) {
      return kb.toFixed(2) + 'KB'
    }
    const mb = kb / 1024
    if (mb < 500) {
      return mb.toFixed(2) + 'MB'
    }
    return (mb / 1024).toFixed(2) + 'GB'
  }

}
