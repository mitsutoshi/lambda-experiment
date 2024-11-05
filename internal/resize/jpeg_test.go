package resize

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"testing"
)

const testJpegFilePath = "../../testdata/sakura.jpg"

const (
	jpegPathGrayscale = "../../testdata/Lenna_grayscale_150x150.jpg"
	jpegPathYcbcr     = "../../testdata/Lenna_ycbcr_150x150.jpg"
)

var images = []string{jpegPathGrayscale, jpegPathYcbcr}

func TestResizeRGBA(t *testing.T) {

	testcases := []struct {
		w, h int
	}{
		{100, 200},
		{120, 200},
		{602, 405},
	}

	for _, testcase := range testcases {
		for _, path := range images {

			data, err := read(path)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
			}

			newImage := ResizeJpegRGBA(data, testcase.w, testcase.h)
			if newImage.Bounds().Size().X != testcase.w || newImage.Bounds().Size().Y != testcase.h {
				t.Errorf("The size does not match the expected size. expected: w=%v, h=%v, actual: w=%v, h=%v",
					testcase.w, testcase.h, newImage.Bounds().Size().X, newImage.Bounds().Size().Y)
			}
		}
	}
}

func TestResizeGrayscale(t *testing.T) {

	testcases := []struct {
		w, h int
	}{
		{100, 200},
		{120, 200},
		{602, 405},
	}

	for _, testcase := range testcases {
		for _, path := range images {

			data, err := read(path)
			if err != nil {
				t.Errorf("Failed to read file: %v", err)
			}

			newImage := ResizeJpegGrayscale(data, testcase.w, testcase.h)
			if newImage.Bounds().Size().X != testcase.w || newImage.Bounds().Size().Y != testcase.h {
				t.Errorf("The size does not match the expected size. expected: w=%v, h=%v, actual: w=%v, h=%v",
					testcase.w, testcase.h, newImage.Bounds().Size().X, newImage.Bounds().Size().Y)
			}
		}
	}
}

func read(path string) (image.Image, error) {

	// read image file data
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	fmt.Printf("read file: path=%s, x=%v, y=%v\n", path, data.Bounds().Max.X, data.Bounds().Max.Y)
	return data, nil
}

func TestS(t *testing.T) {
	s := "original/100/fashion/abcdef.jpg"
	//e := "resized/100/fashion/abcdef/small.jpg"
	name := "small.jpg"

	a, _ := strings.CutSuffix(s, ".jpg")
	b, _ := strings.CutPrefix(a, "original/")
	ss := []string{"resized", b, name}
	c := strings.Join(ss, "/")
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
