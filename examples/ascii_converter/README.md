## ascii-convertor

Takes an image or video as input and returns it made out of ascii text characters

### Options

```
-charset string
    Type of charset you want to use, 'limited' or 'extended' or 'block' (default "limited")
-ffmpeg string
    Path to the ffmpeg binary (default "ffmpeg")
-ffprobe string
    Path to the ffprobe binary (default "ffprobe")
-font string
    Path to a .ttf font file which the characters will be rendered as. If empty, 'Cascadia Code Mono Bold' is used
-fontsize float
    Font size in points (NOT pixels) (default 14)
-input string
    Path to the image/video you want to make ascii. Images must be jpeg or png
-interpolate
    Whether to use interpolation for video rendering (default true)
-intpf float
    Interpolation factor to use, between 1 (none) to 0 (max) interpolation (default 0.6)
-optimise string
    Whether to optimise 'speed', 'balanced' or 'filesize' for the file encoding (default "balanced")
-output string
    Name of the output image/video you want to make e.g. 'test.jpg'. Images may be jpeg or png
-overwrite
    Whether to automatically overwrite the output file if one already exists without prompting
-transparency
    Whether to use a transparent background for the ascii images
-verbose
    Prints verbose information
```

### Build
1. run `make`
