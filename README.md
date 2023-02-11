# fileServer

## endpoints:

### POST files
Lists the files in a given directory
- Array of subfolders can be added in the post body to navigate to different folders

### POST upload
Upload a file to the root directory
- body should contain the key `file` that containes the file data

### POST download
Request a file from a given directory
- File should be passed as the last item of an array
- Array of subfolders can be added in the post body to navigate to different folders