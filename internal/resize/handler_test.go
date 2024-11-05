package resize

const jpegImageFilePath = "../../testdata/sakura.jpg"

// const jpegImageFilePath = "../testdata/grayscale.jpg"

//const jpegImageFilePath = "../testdata/png.png"

//func TestRun(t *testing.T) {
//	err := Run(jpegImageFilePath)
//	if err != nil {
//		t.Errorf("Failed: %v", err)
//	}
//}

//func TestResize(t *testing.T) {
//	img, err := LoadImage(jpegImageFilePath)
//	if err != nil {
//		t.Errorf("Failed: %v", err)
//	}
//
//	testcases := []struct {
//		w, h int
//	}{
//		{200, 150},
//		{120, 400},
//		{1200, 800},
//	}
//
//	for _, testcase := range testcases {
//
//		// resize image
//		newImage := ResizeRGBA(img, testcase.w, testcase.h)
//
//		// check width and height of the resized image
//		if newImage.Bounds().Size().X != testcase.w || newImage.Bounds().Size().Y != testcase.h {
//			t.Errorf("The size does not match the expected size. expected: w=%v, h=%v, actual: w=%v, h=%v",
//				testcase.w, testcase.h, newImage.Bounds().Size().X, newImage.Bounds().Size().Y)
//		}
//	}
//}
