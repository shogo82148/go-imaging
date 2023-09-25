package icc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"math"
	"slices"

	"github.com/shogo82148/go-imaging/fp16"
	"github.com/shogo82148/go-imaging/fp16/fp16color"
	"github.com/shogo82148/go-imaging/internal/parallels"
)

func (p *Profile) DecodeTone(img image.Image) *fp16.NRGBAh {
	cr := p.Get(TagRedTRC).(Curve)
	cg := p.Get(TagGreenTRC).(Curve)
	cb := p.Get(TagBlueTRC).(Curve)

	bounds := img.Bounds()
	ret := fp16.NewNRGBAh(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			fr := float64(r) / 0xffff
			fg := float64(g) / 0xffff
			fb := float64(b) / 0xffff
			fa := float64(a) / 0xffff
			if a != 0 {
				fr /= fa
				fg /= fa
				fb /= fa
			}
			fr = cr.DecodeTone(fr)
			fg = cg.DecodeTone(fg)
			fb = cb.DecodeTone(fb)
			ret.SetNRGBAh(x, y, fp16color.NewNRGBAh(fr, fg, fb, fa))
		}
	})
	return ret
}

func (p *Profile) EncodeTone(img *fp16.NRGBAh) *image.NRGBA {
	cr := p.Get(TagRedTRC).(Curve)
	cg := p.Get(TagGreenTRC).(Curve)
	cb := p.Get(TagBlueTRC).(Curve)

	bounds := img.Bounds()
	ret := image.NewNRGBA(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			fr := cr.EncodeTone(rgba.B.Float64())
			fg := cg.EncodeTone(rgba.G.Float64())
			fb := cb.EncodeTone(rgba.B.Float64())
			ret.SetNRGBA(x, y, color.NRGBA{
				R: uint8(fr * 0xff),
				G: uint8(fg * 0xff),
				B: uint8(fb * 0xff),
				A: uint8(rgba.A.Float64() * 0xff),
			})
		}
	})
	return ret
}

func (p *Profile) EncodeTone16(img *fp16.NRGBAh) *image.NRGBA64 {
	cr := p.Get(TagRedTRC).(Curve)
	cg := p.Get(TagGreenTRC).(Curve)
	cb := p.Get(TagBlueTRC).(Curve)

	bounds := img.Bounds()
	ret := image.NewNRGBA64(bounds)
	parallels.Parallel(bounds.Min.Y, bounds.Max.Y, func(y int) {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba := img.NRGBAhAt(x, y)
			fr := cr.EncodeTone(rgba.B.Float64())
			fg := cg.EncodeTone(rgba.G.Float64())
			fb := cb.EncodeTone(rgba.B.Float64())
			ret.SetNRGBA64(x, y, color.NRGBA64{
				R: uint16(fr * 0xffff),
				G: uint16(fg * 0xffff),
				B: uint16(fb * 0xffff),
				A: uint16(rgba.A.Float64() * 0xffff),
			})
		}
	})
	return ret
}

// Curve is a tone reproduction curve.
type Curve interface {
	EncodeTone(x float64) float64
	DecodeTone(x float64) float64
}

var _ TagContent = (*TagContentCurve)(nil)

type TagContentCurve struct {
	Data []uint16
}

type tagContentCurve struct {
	TagType  TagType
	Reserved uint32
	Count    uint32
}

func (t *TagContentCurve) TagType() TagType { return TagTypeCurve }

func (t *TagContentCurve) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 12+len(t.Data)*2))
	curve := tagContentCurve{
		TagType: t.TagType(),
		Count:   uint32(len(t.Data)),
	}
	if err := binary.Write(buf, binary.BigEndian, curve); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, t.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *TagContentCurve) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	var curve tagContentCurve
	if err := binary.Read(r, binary.BigEndian, &curve); err != nil {
		return err
	}
	t.Data = make([]uint16, curve.Count)
	if err := binary.Read(r, binary.BigEndian, &t.Data); err != nil {
		return err
	}
	return nil
}

func (t *TagContentCurve) DecodeTone(x float64) float64 {
	x = max(0, min(1, x)) // clip to [0.0, 1.0]
	if len(t.Data) == 0 {
		return x
	}
	if len(t.Data) == 1 {
		gamma := U8Fixed8Number(t.Data[0]).Float64()
		return math.Pow(x, gamma)
	}

	i, f := math.Modf(x * float64(len(t.Data)-1))
	i0 := int(i)
	y0 := float64(t.Data[i0]) / 0xffff
	if i0 == len(t.Data)-1 {
		return y0
	}

	// linear interpolation
	i1 := i0 + 1
	y1 := float64(t.Data[i1]) / 0xffff
	return y0 + f*(y1-y0)
}

func (t *TagContentCurve) EncodeTone(y float64) float64 {
	y = max(0, min(1, y)) // clip to [0.0, 1.0]
	if len(t.Data) == 0 {
		return y
	}
	if len(t.Data) == 1 {
		gamma := U8Fixed8Number(t.Data[0]).Float64()
		return math.Pow(y, 1/gamma)
	}

	i := uint16(y * 0xffff)
	idx, ok := slices.BinarySearch(t.Data, i)
	x0 := float64(idx) / float64(len(t.Data)-1)
	if ok || idx == len(t.Data)-1 {
		return x0
	}
	y0 := float64(t.Data[idx]) / 0xffff
	y1 := float64(t.Data[idx+1]) / 0xffff
	if y0 == y1 {
		return x0
	}
	f := (y - y0) / (y1 - y0)
	return x0 + f/float64(len(t.Data)-1)
}

type TagContentParametricCurve struct {
	FunctionType uint16
	Params       [8]S15Fixed16Number // this is not a slice because to avoid extra boundary check.
}

func (t *TagContentParametricCurve) params() ([]S15Fixed16Number, error) {
	switch t.FunctionType {
	case 0x0000:
		return t.Params[:1], nil
	case 0x0001:
		return t.Params[:3], nil
	case 0x0002:
		return t.Params[:4], nil
	case 0x0003:
		return t.Params[:5], nil
	case 0x0004:
		return t.Params[:7], nil
	default:
		return nil, errors.New("icc: unknown parametric curve function type")
	}
}

type tagContentParametricCurve struct {
	TagType      TagType
	_            uint32
	FunctionType uint16
	_            uint16
}

func (t *TagContentParametricCurve) TagType() TagType { return TagTypeParametricCurve }

func (t *TagContentParametricCurve) MarshalBinary() ([]byte, error) {
	// write the header
	w := new(bytes.Buffer)
	curve := tagContentParametricCurve{
		TagType:      TagTypeParametricCurve,
		FunctionType: t.FunctionType,
	}
	if err := binary.Write(w, binary.BigEndian, curve); err != nil {
		return nil, err
	}

	// write the parameters
	params, err := t.params()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.BigEndian, params); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (t *TagContentParametricCurve) UnmarshalBinary(data []byte) error {
	// read the header
	r := bytes.NewReader(data)
	var curve tagContentParametricCurve
	if err := binary.Read(r, binary.BigEndian, &curve); err != nil {
		return err
	}

	// read the parameters
	params, err := t.params()
	if err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, params); err != nil {
		return err
	}
	return nil
}

func (t *TagContentParametricCurve) DecodeTone(x float64) float64 {
	x = max(0, min(1, x)) // clip to [0.0, 1.0]
	switch t.FunctionType {
	// Y = X^g
	case 0x0000:
		g := t.Params[0].Float64()
		return math.Pow(x, g)

	// CIE122-1966
	// Y = (aX + b)^g   if X >= -b/a
	// Y = 0            if X <  -b/a
	case 0x0001:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()

		y := a*x + b
		if y < 0 {
			return 0
		}
		y = math.Pow(y, g)
		return y

	// IEC 61966‐3
	// Y = (aX + b)^g + c  if X >= -b/a
	// Y = c               if X <  -b/a
	case 0x0002:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()

		y := a*x + b
		if y < 0 {
			return c
		}
		y = math.Pow(y, g) + c
		return max(0, min(1, y)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g     if X >= d
	// Y = cX             if X <  d
	case 0x0003:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()

		if x < d {
			return c * x
		}
		y := math.Pow(a*x+b, g)
		return y

	// Y = (aX + b)^g + e  if X >= d
	// Y = cX + f          if X <  d
	case 0x0004:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()
		e := t.Params[5].Float64()
		f := t.Params[6].Float64()

		if x < d {
			return c*x + f
		}
		y := math.Pow(a*x+b, g) + e
		return max(0, min(1, y)) // clip to [0.0, 1.0]
	}
	return x
}

func (t *TagContentParametricCurve) EncodeTone(y float64) float64 {
	y = max(0, min(1, y)) // clip to [0.0, 1.0]
	switch t.FunctionType {
	// Y = X^g
	case 0x0000:
		g := t.Params[0].Float64()
		return math.Pow(y, 1/g)

	// Y = (aX + b)^g   if X >= -b/a
	// Y = 0            if X <  -b/a
	case 0x0001:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()

		x := (math.Pow(y, 1/g) - b) / a
		return x

	// Y = (aX + b)^g + c  if X >= -b/a
	// Y = c               if X <  -b/a
	case 0x0002:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()

		if y < c {
			return -b / a
		}
		y = max(0, y-c)
		x := (math.Pow(y, 1/g) - b) / a
		return max(0, min(1, x)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g     if X >= d
	// Y = cX             if X <  d
	case 0x0003:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()

		x := (math.Pow(y, 1/g) - b) / a
		if x < d {
			return y / c
		}
		return max(0, min(1, x)) // clip to [0.0, 1.0]

	// Y = (aX + b)^g + e  if X >= d
	// Y = cX + f          if X <  d
	case 0x0004:
		g := t.Params[0].Float64()
		a := t.Params[1].Float64()
		b := t.Params[2].Float64()
		c := t.Params[3].Float64()
		d := t.Params[4].Float64()
		e := t.Params[5].Float64()
		f := t.Params[6].Float64()

		x := (y - f) / c
		if x < d {
			return max(0, min(1, x)) // clip to [0.0, 1.0]
		}
		y = max(0, y-e)
		x = (math.Pow(y, 1/g) - b) / a
		return max(0, min(1, x)) // clip to [0.0, 1.0]
	}
	return y
}
