package resize

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"
	"resizeimage/internal/s3"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/sync/errgroup"
)

const timeout = 10 * time.Second

func HandleS3Event(event events.S3EventRecord) error {
	log.Printf("Start to handle S3 event. [bucket=%s, key=%s]\n", event.S3.Bucket.Name, event.S3.Object.Key)
	// get meta info
	headObj, err := s3.HeadObject(event.S3.Bucket.Name, event.S3.Object.Key, timeout)
	if err != nil {
		return err
	}
	log.Printf("The object size is %v byte.\n", *headObj.ContentLength)

	// check object size
	if *headObj.ContentLength > MaxByteSize {
		return errors.New(fmt.Sprintf("The size exceeded the limit(%v byte).", MaxByteSize))
	}

	// get the actual data
	obj, err := s3.GetObject(event.S3.Bucket.Name, event.S3.Object.Key, timeout)
	if err != nil {
		return err
	}
	defer obj.Body.Close()

	// copy the body data because if it is read once, it will end.
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, obj.Body)
	if err != nil {
		return err
	}

	// check format
	config, err := jpeg.DecodeConfig(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}

	// decode
	data, err := jpeg.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		log.Printf("Failed to decode: %v\n", err)
		return err
	}

	// resize
	g := new(errgroup.Group)
	sizes := []Spec{Small, Medium, Large}
	for _, size := range sizes {
		s := size // don't pass the 'size' variable directly to resize()
		g.Go(func() error {
			return resize(data, s, config.ColorModel, event.S3.Bucket.Name, event.S3.Object.Key, s.Name)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func resize(src image.Image, spec Spec, colorModel color.Model, bucketName, key, name string) error {
	log.Printf("Start to resize to %s. [src.w=%d, src.h=%d, dst.w=%d, dst.h=%d]\n",
		spec.Name, src.Bounds().Size().X, src.Bounds().Size().Y, spec.Width, spec.Height)

	// resize image data
	var dst image.Image
	switch colorModel {
	case color.GrayModel:
		dst = ResizeJpegGrayscale(src, spec.Width, spec.Height)
	case color.CMYKModel:
		dst = ResizeJpegCmyk(src, spec.Width, spec.Height)
	default:
		dst = ResizeJpegRGBA(src, spec.Width, spec.Height)
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, dst, nil)
	if err != nil {
		return err
	}

	dstKey, err := convKey(key, name, ".jpg")
	if err != nil {
		return err
	}

	// put the resized image to S3
	err = s3.PutObject(bucketName, dstKey, bytes.NewReader(buf.Bytes()), "image/jpeg", timeout)
	if err != nil {
		return err
	}

	log.Printf("Finished to resize to %s.", spec.Name)
	return nil
}

func convKey(original, name, suffix string) (string, error) {
	if str, found := strings.CutSuffix(original, suffix); found {
		if str, found = strings.CutPrefix(str, SrcKeyPrefix+"/"); found {
			return strings.Join([]string{DstKeyPrefix, str, name + suffix}, "/"), nil
		}
	}
	return "", errors.New(fmt.Sprintf("'original' is illegal: %s", original))
}
