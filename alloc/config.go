package alloc

type Config struct {
	MaxLimit    int64
	MinSize     int64
	MaxSize     int64
	Spread      float64
	ReMinSize   int64
	ReMaxSize   int64
	ReSpread    float64
	Print       bool
	ReleaseHalf bool
}
