# Usage
## sf6vid
```
Usage:
  sf6vid [command]

Available Commands:
  censor      Censor the player information in a video
  help        Help about any command

Flags:
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
      --blur int        Custom blur value for when the box blur is used (requires --box-blur flag otherwise this value will be ignored) (default 6)
      --box-blur        Use the box blur filter instead of the new pixelize filter (pixelize requires ffmpeg 6+)
  -h, --help            help for censor
  -i, --input string    Path to input file
      --open            Open the file after running this command
  -o, --output string   Path to output file
      --p1              Censor player 1 side
      --p2              Censor player 2 side
```
