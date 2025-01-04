# rutubeuploader
upload video to rutube from command line

### Usage:
#### Generate token
```shell
$ rutubeuploader token --user vova@mail.com --password Password
```
You will get `token.json` in current directory which will be used for authentication in other operations. 
Token has to expire in half an hour but checking for expiration doesn't work on rutube service yet at upload action. 
So issued token is never expies until rutube fix error on the server side.

#### Upload video to Rutube
```shell
$ rutubeuploader upload video_list.json
```

or
```shell
$ cat video_list.json | rutubeuploader upload
```

files will be set to queue of processing on Rutube platform.


#### video_list.json input JSON format

```json
[
    {
        "url": "https://user:password@videoserver.su/video/alien.mp4",  // url to video where Rutube will get from. Set necessary credentials to url to download video successfully. Available schema: https/http, ftp
        "quality_report": true,       // if true, notification will be called every time when video will be converted to every step of quality, if false notification will be called once when all convertions will be completed
        "author": 123,                // Identity of author. This author has to have access to upload video to specified channel
        "title": "Alien",             // video title, max 100 runes
        "description": "Alien movie", // description of video, max 5000 runes
        "hidden": false,              // true: private, false: public video 
        "category_id": 13,            // Identity number of category for video. Default 13
        "converter_params": "",       // additional video convertion parameters. Example xml-tagging: converter_params=%7B%22editor_xml%22%3A%22ftp%3A%5C%2F%5C%2Frutube%3pass%4010.122.50.222%5C%2FPR291117-A.xml%22%7D
        "protected": false            // true: video will be queued to DRM checking
    },
    /// etc.
]
```
