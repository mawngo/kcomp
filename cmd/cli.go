package cmd

import (
	"fmt"
	"github.com/mawngo/kcomp/internal/kmeans"
	"github.com/phsym/console-slog"
	"github.com/spf13/cobra"
	"runtime"

	// Add webp support.
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Init() *slog.LevelVar {
	level := &slog.LevelVar{}
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level:      level,
			TimeFormat: time.Kitchen,
		}))
	slog.SetDefault(logger)
	cobra.EnableCommandSorting = false
	return level
}

type CLI struct {
	command *cobra.Command
}

// NewCLI create new CLI instance and set up application config.
func NewCLI() *CLI {
	level := Init()
	defaultConcurrency := max(1, runtime.NumCPU()/5)
	defaultKConcurrency := max(1, runtime.NumCPU()/defaultConcurrency)

	f := flags{
		Colors:       15,
		Output:       ".",
		Round:        100,
		Concurrency:  defaultConcurrency,
		KConcurrency: defaultKConcurrency,
		DistanceAlgo: "EuclideanDistance",
		Delta:        0.005,
	}

	command := cobra.Command{
		Use:   "kcomp [files...]",
		Short: "Reduce number of colors used in image",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			debug, err := cmd.PersistentFlags().GetBool("debug")
			if err != nil {
				return err
			}
			if debug {
				level.Set(slog.LevelDebug)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			if q, err := cmd.Flags().GetBool("quick"); err == nil && q {
				f.Delta = 0.01
				f.Round = 50
			}

			if _, err := os.Stat(f.Output); err != nil {
				err := os.Mkdir(f.Output, os.ModePerm)
				if err != nil {
					slog.Info("Error creating output directory", slog.Any("dir", f.Output))
					return
				}
			}

			if f.Concurrency < 1 {
				f.Concurrency = defaultConcurrency
			}

			if f.KConcurrency < 1 {
				f.KConcurrency = defaultKConcurrency
			}
			if f.Series < 1 {
				f.Series = 1
			}

			con := make(chan struct{}, f.Concurrency)
			for _, arg := range args {
				ch := scan(arg)
				for img := range ch {
					if s := f.Series; s > 1 {
						step := f.Colors / s
						start := 1
						if step <= 1 {
							start = 2
							step = 1
							s = f.Colors
						}

						for i := start; i < s; i++ {
							sf := f
							sf.Colors = step * i
							process(img, sf, con)
						}
					}
					process(img, f, con)
				}
			}

			for range f.Concurrency {
				con <- struct{}{}
			}
			slog.Info("Processing completed", slog.Duration("took", time.Since(now)))
		},
	}

	command.Flags().IntVarP(&f.Colors, "colors", "n", f.Colors, "Number of colors to use")
	command.Flags().StringVarP(&f.Output, "out", "o", f.Output, "Output directory name")
	command.Flags().IntVarP(&f.Series, "series", "s", f.Series, "Number of image to generate, series of output with increasing number of colors up util reached --colors parameter")
	command.Flags().BoolP("quick", "q", false, "Increase speed in exchange of accuracy")
	command.Flags().BoolVarP(&f.Overwrite, "overwrite", "w", f.Overwrite, "Overwrite output if exists")
	command.Flags().IntVarP(&f.Round, "round", "i", f.Round, "Maximum number of round before stop adjusting (number of kmeans iterations)")
	command.Flags().Float64VarP(&f.Delta, "delta", "d", f.Delta, "Delta threshold of convergence (delta between kmeans old and new centroid’s values)")
	command.Flags().IntVarP(&f.Concurrency, "concurrency", "t", f.Concurrency, "Maximum number image process at a time [0=auto]")
	command.Flags().IntVar(&f.KConcurrency, "kcpu", f.KConcurrency, "Maximum cpu used processing each image [0=auto]")
	command.Flags().StringVar(&f.DistanceAlgo, "dalgo", f.DistanceAlgo, "Distance algo for kmeans [EuclideanDistance,EuclideanDistanceSquared]")
	command.Flags().IntVar(&f.JPEG, "jpeg", f.JPEG, "Specify quality of output jpeg compression [0-100] (default 0 - output png)")
	command.Flags().BoolVar(&f.Palette, "palette", f.Palette, "Generate an additional palette image")
	command.PersistentFlags().Bool("debug", false, "Enable debug mode")
	command.Flags().SortFlags = false
	return &CLI{&command}
}

func process(img DecodedImage, f flags, con chan struct{}) {
	con <- struct{}{}
	go func() {
		defer func() {
			<-con
		}()
		handleImg(img, f)
	}()
}

func handleImg(img DecodedImage, f flags) {
	slog.Info("Processing",
		slog.Any("cp", f.Colors),
		slog.Any("round", f.Round),
		slog.String("img", img.Basename),
		slog.String("dimension", fmt.Sprintf("%dx%d", img.Width, img.Height)),
		slog.String("format", img.Type),
	)

	outExt := ".png"
	if f.JPEG > 0 {
		outExt = ".jpeg"
	}
	outfile := filepath.Join(f.Output, strings.TrimSuffix(img.Basename, filepath.Ext(img.Path))+".kcp"+strconv.Itoa(f.Round)+"n"+strconv.Itoa(f.Colors)+outExt)
	if stats, err := os.Stat(outfile); err == nil {
		slog.Info("File existed",
			slog.Any("path", outfile),
			slog.Bool("isDir", stats.IsDir()),
			slog.Bool("overwrite", f.Overwrite),
		)
		if !f.Overwrite {
			return
		}
		if stats.IsDir() {
			return
		}
	}

	now := time.Now()
	d := make([][]float64, 0, img.Width*img.Height)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			d = append(d, []float64{float64(r >> 8), float64(g >> 8), float64(b >> 8), float64(a >> 8)})
		}
	}

	algo := kmeans.EuclideanDistance
	if f.DistanceAlgo == "EuclideanDistanceSquared" {
		algo = kmeans.EuclideanDistanceSquared
	}

	slog.Debug("Start partitioning",
		slog.Int("cp", f.Colors),
		slog.String("img", img.Basename),
		slog.Int("round", f.Round),
		slog.Duration("elapsed", time.Since(now)),
	)
	m := kmeans.NewTrainer(f.Colors,
		kmeans.WithConcurrency(f.KConcurrency),
		kmeans.WithDistanceFunc(algo),
		kmeans.WithMaxIterations(f.Round),
		kmeans.WithDeltaThreshold(f.Delta)).
		Fit(d)
	rbga := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: img.Width, Y: img.Height}})
	for index, number := range m.Guesses() {
		cluster := m.Cluster(number)
		y := index / img.Width
		x := index % img.Width
		rbga.SetRGBA(x, y, color.RGBA{
			R: round(cluster[0]),
			G: round(cluster[1]),
			B: round(cluster[2]),
			A: round(cluster[3]),
		})
	}
	o, err := os.Create(outfile)
	if err == nil {
		defer func() {
			err := o.Close()
			if err != nil {
				slog.Error("Error closing image file",
					slog.String("out", outfile),
					slog.Any("err", err))
			}
		}()
		if f.JPEG == 0 {
			err = png.Encode(o, rbga)
		} else {
			err = jpeg.Encode(o, rbga, &jpeg.Options{Quality: f.JPEG})
		}
	}
	if err != nil {
		slog.Error("Error writing image",
			slog.String("out", outfile),
			slog.Any("err", err))
		return
	}
	if f.Palette {
		genPalette(m.Centroids(), outfile)
	}
	slog.Info("Compress completed",
		slog.String("out", outfile),
		slog.Duration("took", time.Since(now)),
		slog.Int("iter", m.Iter()))
}

func genPalette(centroids kmeans.Dataset, originalOutFile string) {
	filename := strings.TrimSuffix(originalOutFile, filepath.Ext(originalOutFile)) + ".palette.png"

	swatchWidth := 400
	if len(centroids) > 1 {
		swatchWidth = 200 - min(7*len(centroids)-2, 140)
	}

	width := swatchWidth * len(centroids)
	height := int(float64(width) / math.Phi)
	rect := image.Rect(0, 0, width, height)

	img := image.NewRGBA(rect)
	for i, cluster := range centroids {
		c := color.RGBA{
			R: round(cluster[0]),
			G: round(cluster[1]),
			B: round(cluster[2]),
			A: round(cluster[3]),
		}
		startX := i * swatchWidth
		endX := (i + 1) * swatchWidth
		for y := 0; y < height; y++ {
			for x := startX; x < endX; x++ {
				img.Set(x, y, c)
			}
		}
	}

	o, err := os.Create(filename)
	if err == nil {
		defer func() {
			err := o.Close()
			if err != nil {
				slog.Error("Error closing palette file",
					slog.String("out", filename),
					slog.Any("err", err))
			}
		}()
		err = png.Encode(o, img)
	}
	if err != nil {
		slog.Error("Error writing palette image",
			slog.String("out", filename),
			slog.Any("err", err))
		return
	}
}

func round(f float64) uint8 {
	return uint8(math.Round(f))
}

type flags struct {
	Colors       int
	Output       string
	Round        int
	Overwrite    bool
	Concurrency  int
	KConcurrency int
	DistanceAlgo string
	JPEG         int
	Delta        float64
	Series       int
	Palette      bool
}

func scan(dir string) <-chan DecodedImage {
	ch := make(chan DecodedImage, 1)
	info, err := os.Stat(dir)
	if err != nil {
		slog.Error("Err scanning file(s)", slog.String("path", dir), slog.Any("err", err))
		close(ch)
		return ch
	}

	go func() {
		defer close(ch)
		if !info.IsDir() {
			img, err := decode(dir)
			if err != nil {
				slog.Error("Err decoding image", slog.String("path", dir), slog.Any("err", err))
				return
			}
			ch <- img
			return
		}

		files, err := os.ReadDir(".")
		if err != nil {
			slog.Error("Err scanning dir", slog.String("path", dir), slog.Any("err", err))
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			path := filepath.Join(dir, file.Name())
			img, err := decode(path)
			if err != nil {
				slog.Error("Not an image", slog.String("path", path), slog.Any("err", err))
				continue
			}
			ch <- img
		}
	}()

	return ch
}

func decode(path string) (DecodedImage, error) {
	img := DecodedImage{
		Path:     path,
		Basename: filepath.Base(path),
	}
	f, err := os.Open(path)
	if err != nil {
		return img, err
	}
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	if err != nil {
		return img, err
	}
	img.Config = config

	_, err = f.Seek(0, 0)
	if err != nil {
		panic(err)
	}
	slog.Debug("Decoding image",
		slog.String("path", path),
		slog.String("dimension", fmt.Sprintf("%dx%d", img.Width, img.Height)))
	imageData, imageType, err := image.Decode(f)
	if err != nil {
		return img, err
	}
	img.Type = imageType
	img.Image = imageData

	return img, nil
}

type DecodedImage struct {
	image.Image
	image.Config
	Type     string
	Path     string
	Basename string
}

func (cli *CLI) Execute() {
	if err := cli.command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
