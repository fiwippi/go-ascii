## fluid-sim

Takes an input image and simulates all coloured pixels as fluids falling down.
**ffmpeg must be in your path**.

### Options
```
-fontsize float
    Font size in points (NOT pixels) (default 10)
-input string
    Path to the image you want to make fluid (all non-black pixels are treated as fluid)
-output string
    Name of the output video you want to make e.g. 'test.mp4'
-overwrite
    Whether to automatically overwrite the output file if one already exists without prompting
-random-colours
    Whether each fluid should be a random colour, if not a predetermined pattern is used instead (default true)
-speed int
    How quickly the fluid should fall, in the range 1-9,  1 is slowest, 9 is fastest (default 4)
```

### Example
#### Input
![input](assets/spiral.png)

#### Output (Random Colours)
![output](assets/spiral.gif)

### Build
1. Run `make`
