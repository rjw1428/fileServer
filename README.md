# fileServer

## endpoints:

### GET files
Lists the files in a given directory
- Array of subfolders can be added in the post body to navigate to different folders

### GET upload
Upload a file to the root directory
- body should contain the key `file` that containes the file data

### GET download
Request a file from a given directory
- File should be passed as the last item of an array
- Array of subfolders can be added in the post body to navigate to different folders

### GET createFolder
Creates a folder at the given location


# Client
## Build
When building client, use `npm build --prod`