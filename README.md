# K-Compressor

Reduce number of color used in image using k-mean clustering.

# Usage

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
  kcomp [file] [flags]
  
Flags:
  -n, --colors int        Number of colors to use (default 12)
  -t, --concurrency int   Maximum number image process at a time (default 8)
      --dalgo string      Distance algo for kmeans [EuclideanDistance,EuclideanDistanceSquared,Squared] (default "EuclideanDistance")
      --debug             Enable debug mode
  -h, --help              help for kcomp
  -o, --out string        Output directory name (default "kcompressed")
  -O, --out-current-dir   Output on current directory, same as --out=.
  -w, --overwrite         Overwrite output if exists
  -i, --round int         Maximum number of round before stop adjusting (number of kmeans iterations) (default 100)
  -s, --series int        Number of image to generate, series of output with increasing number of colors up util reached --colors parameter (default 1)
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

