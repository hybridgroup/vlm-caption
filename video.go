package main

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/hybridgroup/mjpeg"
	"github.com/hybridgroup/yzma/pkg/mtmd"
	"gocv.io/x/gocv"
)

var (
	webcam *gocv.VideoCapture
	img    gocv.Mat
	mutex  sync.Mutex
)

func startCapture(deviceID string, stream *mjpeg.Stream) {
	var err error
	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	img = gocv.NewMat()
	defer img.Close()

	for {
		captureFrame(deviceID, stream)
	}
}

func captureFrame(deviceID string, stream *mjpeg.Stream) {
	mutex.Lock()
	defer mutex.Unlock()

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Device closed: %v\n", deviceID)
		return
	}
	if img.Empty() {
		return
	}

	buf, _ := gocv.IMEncode(".jpg", img)
	stream.UpdateJPEG(buf.GetBytes())
	buf.Close()
}

func matToBitmap(img gocv.Mat) (mtmd.Bitmap, error) {
	mutex.Lock()
	defer mutex.Unlock()

	rgb := gocv.NewMatWithSize(img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	defer rgb.Close()

	gocv.CvtColor(img, &rgb, gocv.ColorBGRToRGB)
	ptr, err := rgb.DataPtrUint8()
	if err != nil {
		return mtmd.Bitmap(0), err
	}

	bitmap := mtmd.BitmapInit(uint32(img.Cols()), uint32(img.Rows()), uintptr(unsafe.Pointer(&ptr[0])))
	return bitmap, nil
}
