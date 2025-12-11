# K-Compressor

Reduce the number of colors used in an image using k-mean clustering.

## Installation

Require go 1.25+

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
  -o, --out string        Output directory name (default ".")
  -s, --series int        Number of image to generate, series of output with increasing number of colors up util reached --colors parameter
  -q, --quick             Increase speed in exchange of accuracy
  -w, --overwrite         Overwrite output if exists
  -i, --round int         Maximum number of round before stop adjusting (number of kmeans iterations) (default 100)
  -d, --delta float       Delta threshold of convergence (delta between kmeans old and new centroid’s values) (default 0.005)
  -t, --concurrency int   Maximum number image process at a time [0=auto] (default 3)
      --kcpu int          Maximum cpu used processing each image [0=auto] (default 8)
      --dalgo string      Distance algo for kmeans [EuclideanDistance,EuclideanDistanceSquared] (default "EuclideanDistance")
      --jpeg int          Specify quality of output jpeg compression [0-100] (default 0 - output png)
      --palette           Generate an additional palette image
      --debug             Enable debug mode
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

