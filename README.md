# K-Compressor

Reduce number of color used in image using k-mean clustering.

# Usage

```shell
> kcomp .\my-image.jpeg --colors=4
```

### Options

```
> kcomp -h  
Reduce number of colors used in image

Usage:
  kcomp [file] [flags]

Flags:
      --colors int        Number of colors to use (default 20)
      --concurrency int   Maximum number image process at a time (default 4)
      --debug             Enable debug mode
  -h, --help              help for kcomp
      --out string        Output directory name (default ".")
      --overwrite         Overwrite output if exists
      --round int         Maximum number of round before stop adjusting (default 50)

```