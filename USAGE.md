# Usage
## sf6vid
```
Usage:
  sf6vid [command]

Available Commands:
  censor      Censor the player information in a video
  help        Help about any command
  shrink      Reduces the size of the video including frame by the indicated percentage, and compresses the quality in other ways
  trim        Trim the video for the provided start and/or end times

Flags:
      --debug     More verbose logging
  -h, --help      help for sf6vid
  -v, --version   version for sf6vid

Use "sf6vid [command] --help" for more information about a command.
```
### censor
```
Censor either the player 1 or player 2 identifying information in the video.
If the output path already exists, it will be replaced.

Usage:
  sf6vid censor [flags]

Flags:
      --p1               Censor player 1 side
      --p2               Censor player 2 side
  -i, --input string     Path to input file
  -o, --output string    Path to output file
      --open             Open the file after running this command
      --blur int         Custom blur value for when the box blur is used (requires --box-blur flag otherwise this value will be ignored) (default 6)
      --box-blur         Use the box blur filter instead of the new pixelize filter (pixelize requires ffmpeg 6+)
      --start duration   Optional start time for trimming the video
      --end duration     Optional end time for trimming the video
  -h, --help             help for censor

Global Flags:
      --debug   More verbose logging
```
### trim
```
You can provide one or both flags --start and --end.
If you omit --start, the original start time of the video will be used.
If you omit --end, the original end time of the video will be used.
At least one is required.
--start and --end use duration syntax, e.g. 5m30s for 5 minutes and 30 seconds

Usage:
  sf6vid trim [flags]

Flags:
  -i, --input string     Path to input file
  -o, --output string    Path to output file
      --open             Open the file after running this command
      --start duration   Start time for trimming the video
      --end duration     End time for trimming the video
  -h, --help             help for trim

Global Flags:
      --debug   More verbose logging
```
