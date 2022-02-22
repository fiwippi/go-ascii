## video

### pre-requisites
- `ffmpeg`

### build
```console
$ go build
```

### usage
⚠️ - args shouldn't contain spaces since they're interpreted as a delimiter, e.g. `-args="-metadata title='HAS SPACE'"`
```console
$ ./video --help
Usage: ./video -i in.mp4 out.mp4
  -args string
        specify extra args for ffmpeg
  -fontsize float
        fontsize of the ascii characters (default 14)
  -i string
        path to video to convert to ascii (default "explosion.mkv")
  -y    automatically overwrites the output file if it exists

$ ./video -i assets/explosion.mkv -fontsize 22 -args="-c:v libx264 -crf 24" -y out.mp4
```

### example
```console
$ go build && ./video -i assets/explosion.mkv -fontsize 22 out.mp4
```

![example](assets/explosion.gif)