# K-Compressor

Reduce number of color used in image using k-mean clustering.

## Installation

Require go 1.22+

```shell
go install github.com/mawngo/kcomp@latest
```

## Usage

compress image

```shell
> kcomp .\my-image.jpeg
```

or compress directory of images

```shell
> kcomp .\my-dir
```

### Options

```
> kcomp -h  
Reduce number of colors used in image

Usage:
  kcomp [files...] [flags]
  
Flags:
  -n, --colors int        Number of colors to use (default 15)
  -t, --concurrency int   Maximum number image process at a time [min:1] (default 8)
      --dalgo string      Distance algo for kmeans [EuclideanDistance,EuclideanDistanceSquared] (default "EuclideanDistance")
      --debug             Enable debug mode
  -d, --delta float       Delta threshold of convergence (delta between kmeans old and new centroid’s values) (default 0.005)
  -h, --help              help for kcomp
      --jpeg int          Specify quality of output jpeg compression [0-100] (set to 0 to output png)
  -o, --out string        Output directory name (default ".")
  -O, --out-current-dir   Output on current directory (same as --out=.)
  -w, --overwrite         Overwrite output if exists
  -q, --quick             Increase speed in exchange of accuracy
  -i, --round int         Maximum number of round before stop adjusting (number of kmeans iterations) (default 100)
  -s, --series int        Number of image to generate, series of output with increasing number of colors up util reached --colors parameter [min:1] (default 1)
```

## Examples

```shell
> kcomp .\chika.jpeg --colors=5
```

```shell
10:46PM INF Processing cp=5 round=100 img=chika.jpeg dimension=200x200 format=jpeg
10:46PM INF Compress completed out=chika.100cp5.png took=60.239ms
10:46PM INF Processing completed.
```

| Original                  | 5 Colors                              | 4 Colors                              |
|---------------------------|---------------------------------------|---------------------------------------|
| ![chika.jpeg](chika.jpeg) | ![chika.100cp5.png](chika.100cp5.png) | ![chika.100cp4.png](chika.100cp4.png) |

