module github.com/3ter/iMagine

go 1.15

require (
	github.com/faiface/beep v1.0.2
	github.com/faiface/pixel v0.10.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5
)

replace "github.com/3ter/iMagine/internal/utils" => ./internal/utils