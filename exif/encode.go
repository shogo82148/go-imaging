package exif

import (
	"cmp"
	"encoding/binary"
	"fmt"
	"io"
	"slices"
)

type encodeState struct {
	data      []byte
	byteOrder binary.ByteOrder
}

func Encode(w io.Writer, t *TIFF) error {
	var err error
	var idfTIFF, idfExif, idfGPS *idf
	idfTIFF, err = convertTIFFToIDF(t)
	if err != nil {
		return err
	}
	if t.Exif != nil {
		idfExif, err = convertExifToIDF(t.Exif)
		if err != nil {
			return err
		}
	}
	if t.GPS != nil {
		idfGPS, err = convertGPSInfoToIDF(t.GPS)
		if err != nil {
			return err
		}
	}

	e := &encodeState{
		data:      []byte{},
		byteOrder: binary.BigEndian,
	}
	if err := e.encode(idfTIFF, idfExif, idfGPS); err != nil {
		return err
	}

	_, err = io.WriteString(w, "Exif\x00\x00")
	if err != nil {
		return err
	}
	_, err = w.Write(e.data)
	return err
}

func (e *encodeState) encode(idfTIFF, idfExif, idfGPS *idf) error {
	size := 8
	if idfTIFF != nil {
		size += 2 + 12*len(idfTIFF.entries) + 4
	}
	if idfExif != nil {
		size += 2 + 12*len(idfExif.entries) + 4
	}
	if idfGPS != nil {
		size += 2 + 12*len(idfGPS.entries) + 4
	}
	e.extend(size)

	if e.byteOrder == binary.BigEndian {
		e.data[0] = 'M'
		e.data[1] = 'M'
	} else {
		e.data[0] = 'I'
		e.data[1] = 'I'
	}
	e.byteOrder.PutUint16(e.data[2:4], 0x002a)
	e.byteOrder.PutUint32(e.data[4:8], 0x0008)

	var offset uint32 = 8
	var err error
	if idfTIFF != nil {
		offset, err = e.encodeIDF(idfTIFF, offset)
		if err != nil {
			return err
		}
	}
	if idfExif != nil {
		offset, err = e.encodeIDF(idfExif, offset)
		if err != nil {
			return err
		}
	}
	if idfGPS != nil {
		offset, err = e.encodeIDF(idfGPS, offset)
		if err != nil {
			return err
		}
	}
	if offset != uint32(size) {
		panic(fmt.Sprintf("internal error: offset != size: %d != %d", offset, size))
	}
	return nil
}

func (e *encodeState) grow(n int) {
	e.data = slices.Grow(e.data, n)
}

func (e *encodeState) extend(n int) {
	l := len(e.data)
	e.grow(n)
	e.data = e.data[:l+n]
	clear(e.data[l:])
}

func convertTIFFToIDF(t *TIFF) (*idf, error) {
	entries := []*idfEntry{}
	if t.ImageDescription != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagImageDescription, *t.ImageDescription),
		)
	}
	if t.Make != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagMake, *t.Make),
		)
	}
	if t.Model != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagModel, *t.Model),
		)
	}
	if t.Orientation != 0 {
		entries = append(entries, &idfEntry{
			tag:      tagOrientation,
			dataType: dataTypeShort,
			shortData: []uint16{
				uint16(t.Orientation),
			},
		})
	}
	if t.XResolution != nil {
		entries = append(entries, &idfEntry{
			tag:      tagXResolution,
			dataType: dataTypeRational,
			rationalData: []Rational{
				*t.XResolution,
			},
		})
	}
	if t.YResolution != nil {
		entries = append(entries, &idfEntry{
			tag:      tagYResolution,
			dataType: dataTypeRational,
			rationalData: []Rational{
				*t.YResolution,
			},
		})
	}
	if t.ResolutionUnit != 0 {
		entries = append(entries, &idfEntry{
			tag:      tagResolutionUnit,
			dataType: dataTypeShort,
			shortData: []uint16{
				uint16(t.ResolutionUnit),
			},
		})
	}

	slices.SortFunc(entries, func(a, b *idfEntry) int {
		return cmp.Compare(a.tag, b.tag)
	})

	return &idf{
		entries: entries,
	}, nil
}

func convertExifToIDF(t *Exif) (*idf, error) {
	return &idf{
		entries: []*idfEntry{},
	}, nil
}

func convertGPSInfoToIDF(t *GPS) (*idf, error) {
	return &idf{
		entries: []*idfEntry{},
	}, nil
}

func convertAsciiOrUTF8(t tag, s string) *idfEntry {
	if isAscii(s) {
		return &idfEntry{
			tag:       t,
			dataType:  dataTypeAscii,
			asciiData: s,
		}
	} else {
		return &idfEntry{
			tag:      t,
			dataType: dataTypeUTF8,
			utf8data: s,
		}
	}
}

func isAscii(s string) bool {
	for _, r := range s {
		if r > 0x7f {
			return false
		}
	}
	return true
}

func (e *encodeState) encodeIDF(idf *idf, offset uint32) (uint32, error) {
	count := uint16(len(idf.entries))
	e.byteOrder.PutUint16(e.data[offset:offset+2], count)
	offset += 2
	for _, entry := range idf.entries {
		var err error
		offset, err = e.encodeIDFEntry(entry, offset)
		if err != nil {
			return 0, err
		}
	}
	e.byteOrder.PutUint32(e.data[offset:offset+4], idf.nextOffset)
	offset += 4
	return offset, nil
}

func (e *encodeState) encodeIDFEntry(entry *idfEntry, offset uint32) (uint32, error) {
	e.byteOrder.PutUint16(e.data[offset:offset+2], uint16(entry.tag))
	offset += 2
	e.byteOrder.PutUint16(e.data[offset:offset+2], uint16(entry.dataType))
	offset += 2
	switch entry.dataType {
	case dataTypeByte:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.byteData)))
		offset += 4
	case dataTypeAscii:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.asciiData)))
		offset += 4
		if len(entry.asciiData) <= 3 {
			n := copy(e.data[offset:offset+4], entry.asciiData)
			e.data[offset+uint32(n)] = '\x00'
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.asciiData) + 1)
			copy(e.data[l:], entry.asciiData)
			e.data[l+len(entry.asciiData)] = '\x00'
		}
		offset += 4
	case dataTypeShort:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.shortData)))
		offset += 4
		if len(entry.shortData) <= 2 {
			for i, v := range entry.shortData {
				e.byteOrder.PutUint16(e.data[offset+2*uint32(i):offset+2*uint32(i)+2], v)
			}
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(2 * len(entry.shortData))
			for i, v := range entry.shortData {
				e.byteOrder.PutUint16(e.data[l+2*i:l+2*i+2], v)
			}
		}
		offset += 4
	case dataTypeLong:
	case dataTypeRational:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.rationalData)))
		offset += 4
		if len(entry.rationalData) > 0 {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(8 * len(entry.rationalData))
			for i, v := range entry.rationalData {
				e.byteOrder.PutUint32(e.data[l+8*i:l+8*i+4], v.Numerator)
				e.byteOrder.PutUint32(e.data[l+8*i+4:l+8*i+8], v.Denominator)
			}
		}
		offset += 4
	case dataTypeSByte:
	case dataTypeUndefined:
	case dataTypeSShort:
	case dataTypeSLong:
	case dataTypeSRational:
	case dataTypeFloat:
	case dataTypeDouble:
	case dataTypeUTF8:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.utf8data)))
		offset += 4
		if len(entry.utf8data) <= 3 {
			n := copy(e.data[offset:offset+4], entry.utf8data)
			e.data[offset+uint32(n)] = '\x00'
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.utf8data) + 1)
			copy(e.data[l:], entry.utf8data)
			e.data[l+len(entry.utf8data)] = '\x00'
		}
		offset += 4
	}
	return offset, nil
}
