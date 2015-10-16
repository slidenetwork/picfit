package engines

type Operation struct {
	Name string
}

var Resize = &Operation{
	"resize",
}

var Thumbnail = &Operation{
	"thumbnail",
}

var Rotate = &Operation{
	"rotate",
}

var Flip = &Operation{
	"flip",
}

var Fit = &Operation{
	"fit",
}

var Original = &Operation{
	"original",
}

var Operations = map[string]*Operation{
	Resize.Name:    Resize,
	Thumbnail.Name: Thumbnail,
	Flip.Name:      Flip,
	Rotate.Name:    Rotate,
	Fit.Name:       Fit,
	Original.Name:  Original,
}
