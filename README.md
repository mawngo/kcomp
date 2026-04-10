# K-Compressor

Reduce the number of colors used in an image using k-mean clustering.

## Installation

Require go 1.26+

```shell
go install github.com/mawngo/kcomp@latest
```

Alternately, check the [Releases](https://github.com/mawngo/kcomp/releases) page for pre-built binaries.

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
  -n, --colors int        number of colors to use (default 15)
  -o, --out string        output directory name (default ".")
  -s, --series int        number of image to generate, series of output with increasing number of colors up util reached --colors parameter
  -q, --quick             increase speed in exchange of accuracy
  -w, --overwrite         overwrite output if exists
  -i, --round int         maximum number of round before stop adjusting (number of kmeans iterations) (default 100)
  -d, --delta float       delta threshold of convergence (delta between kmeans old and new centroid’s values) (default 0.005)
  -t, --concurrency int   maximum number image process at a time [0=auto] (default 4)
      --kcpu int          maximum cpu used processing each image [0=auto] (default 6)
      --dalgo string      distance algo for kmeans [EuclideanDistance,EuclideanDistanceSquared] (default "EuclideanDistance")
      --jpeg int          specify quality of output jpeg compression [0-100] (default 0 - output png)
      --palette           generate an additional palette image
      --debug             enable debug mode
  -h, --help              help for kcomp
```

## Examples

```shell
> kcomp .\chika.jpeg --colors=5
```

```shell
4:05PM INF Processing cp=5 round=100 img=chika.jpeg dimension=200x200 format=jpeg
4:05PM INF Compress completed out=chika.kcp100n5.png took=33.1119ms iter=3
4:05PM INF Processing completed took=34.7113ms
```

| Original                  | 5 Colors                                  | 4 Colors                                  |
|---------------------------|-------------------------------------------|-------------------------------------------|
| ![chika.jpeg](chika.jpeg) | ![chika.kcp100n5.png](chika.kcp100n5.png) | ![chika.kcp100n4.png](chika.kcp100n4.png) |

