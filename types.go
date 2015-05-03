package irmagician

type Dump struct {
	Scale  int    `json:"postscale"`
	Format string `json:"format"`
	Freq   int    `json:"freq"`
	Data   []byte `json:"data"`
}
