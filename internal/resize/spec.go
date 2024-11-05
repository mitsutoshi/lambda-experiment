package resize

type Spec struct {
	Width  int
	Height int
	Name   string
}

var (
	Small  = Spec{Width: 300, Height: 200, Name: "small"}
	Medium = Spec{Width: 500, Height: 400, Name: "medium"}
	Large  = Spec{Width: 800, Height: 600, Name: "large"}
)

const (
	MaxByteSize  int64  = 1024 * 1024 * 10
	SrcKeyPrefix string = "original"
	DstKeyPrefix string = "resized"
)
