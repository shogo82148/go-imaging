package exif

import (
	"cmp"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"slices"
)

type encodeState struct {
	data      []byte
	byteOrder binary.ByteOrder

	exifOffset uint32 // offset to the pointer for ExifIFD
	gpsOffset  uint32 // offset to the pointer for GPSInfoIFD
}

func Encode(w io.Writer, t *TIFF) error {
	e := &encodeState{
		data:      []byte{},
		byteOrder: binary.BigEndian,
	}
	if err := e.encode(t); err != nil {
		return err
	}

	_, err := io.WriteString(w, "Exif\x00\x00")
	if err != nil {
		return err
	}
	_, err = w.Write(e.data)
	return err
}

func (e *encodeState) encode(t *TIFF) error {
	var err error
	var idfTIFF, idfExif, idfGPS *idf
	idfTIFF, err = e.convertTIFFToIDF(t)
	if err != nil {
		return err
	}
	if t.Exif != nil {
		idfExif, err = e.convertExifToIDF(t.Exif)
		if err != nil {
			return err
		}
	}
	if t.GPS != nil {
		idfGPS, err = e.convertGPSToIDF(t.GPS)
		if err != nil {
			return err
		}
	}

	if err := e.encodeHeader(); err != nil {
		return err
	}
	if err := e.encodeTIFF(idfTIFF); err != nil {
		return err
	}
	if idfExif != nil {
		if err := e.encodeExif(idfExif); err != nil {
			return err
		}
	}
	if idfGPS != nil {
		if err := e.encodeGPS(idfGPS); err != nil {
			return err
		}
	}
	return nil
}

func (e *encodeState) encodeHeader() error {
	e.extend(8)
	if e.byteOrder == binary.BigEndian {
		e.data[0] = 'M'
		e.data[1] = 'M'
	} else {
		e.data[0] = 'I'
		e.data[1] = 'I'
	}
	e.byteOrder.PutUint16(e.data[2:4], 0x002a)
	e.byteOrder.PutUint32(e.data[4:8], 0x0008)
	return nil
}

func (e *encodeState) encodeTIFF(idfTIFF *idf) error {
	offset := uint32(len(e.data))
	e.extend(2 + 12*len(idfTIFF.entries) + 4)
	_, err := e.encodeIDF(idfTIFF, offset)
	return err
}

func (e *encodeState) encodeExif(idfExif *idf) error {
	offset := uint32(len(e.data))
	e.byteOrder.PutUint32(e.data[e.exifOffset:], offset)
	e.extend(2 + 12*len(idfExif.entries) + 4)
	_, err := e.encodeIDF(idfExif, offset)
	return err
}

func (e *encodeState) encodeGPS(idfGPS *idf) error {
	offset := uint32(len(e.data))
	e.byteOrder.PutUint32(e.data[e.gpsOffset:], offset)
	e.extend(2 + 12*len(idfGPS.entries) + 4)
	_, err := e.encodeIDF(idfGPS, offset)
	return err
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

// align aligns the data to the next 2-byte boundary.
func (e *encodeState) align() {
	if len(e.data)%2 != 0 {
		e.extend(1)
	}
}

func (e *encodeState) convertTIFFToIDF(t *TIFF) (*idf, error) {
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
	if t.Software != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagSoftware, *t.Software),
		)
	}
	if t.DateTime != nil {
		entries = append(entries, &idfEntry{
			tag:       tagDateTime,
			dataType:  dataTypeAscii,
			asciiData: *t.DateTime,
		})
	}
	if t.Artist != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagArtist, *t.Artist),
		)
	}
	if t.Copyright != nil {
		entries = append(
			entries,
			convertAsciiOrUTF8(tagCopyright, *t.Copyright),
		)
	}
	slices.SortFunc(entries, func(a, b *idfEntry) int {
		return cmp.Compare(a.tag, b.tag)
	})

	// add dummy entry for Exif and GPS
	offset := uint32(8)
	offset += 2 + 12*uint32(len(entries))
	if t.Exif != nil {
		e.exifOffset = offset + 8
		entries = append(entries, &idfEntry{
			tag:      tagExifIFDPointer,
			dataType: dataTypeLong,
			longData: []uint32{0},
		})
		offset += 12
	}
	if t.GPS != nil {
		e.gpsOffset = offset + 8
		entries = append(entries, &idfEntry{
			tag:      tagGPSInfoIFDPointer,
			dataType: dataTypeLong,
			longData: []uint32{
				0,
			},
		})
		offset += 12
	}

	return &idf{
		entries: entries,
	}, nil
}

func (e *encodeState) convertExifToIDF(exif *Exif) (*idf, error) {
	entries := []*idfEntry{}

	// exif version
	entries = append(entries, &idfEntry{
		tag:           tagExifVersion,
		dataType:      dataTypeUndefined,
		undefinedData: []byte("0300"),
	})

	if exif.ExposureTime != nil {
		entries = append(entries, &idfEntry{
			tag:      tagExposureTime,
			dataType: dataTypeRational,
			rationalData: []Rational{
				*exif.ExposureTime,
			},
		})
	}
	if exif.FNumber != nil {
		entries = append(entries, &idfEntry{
			tag:      tagFNumber,
			dataType: dataTypeRational,
			rationalData: []Rational{
				*exif.FNumber,
			},
		})
	}
	if exif.ExposureProgram != 0 {
		entries = append(entries, &idfEntry{
			tag:      tagExposureProgram,
			dataType: dataTypeShort,
			shortData: []uint16{
				uint16(exif.ExposureProgram),
			},
		})
	}
	if exif.ISOSpeedRatings != nil {
		entries = append(entries, &idfEntry{
			tag:       tagISOSpeedRatings,
			dataType:  dataTypeShort,
			shortData: exif.ISOSpeedRatings,
		})
	}
	if exif.DateTimeOriginal != nil {
		entries = append(entries, &idfEntry{
			tag:       tagDateTimeOriginal,
			dataType:  dataTypeAscii,
			asciiData: *exif.DateTimeOriginal,
		})
	}
	if exif.DateTimeDigitized != nil {
		entries = append(entries, &idfEntry{
			tag:       tagDateTimeDigitized,
			dataType:  dataTypeAscii,
			asciiData: *exif.DateTimeDigitized,
		})
	}
	if exif.ShutterSpeedValue != nil {
		entries = append(entries, &idfEntry{
			tag:      tagShutterSpeedValue,
			dataType: dataTypeSRational,
			sRationalData: []SRational{
				*exif.ShutterSpeedValue,
			},
		})
	}
	if exif.ApertureValue != nil {
		entries = append(entries, &idfEntry{
			tag:      tagApertureValue,
			dataType: dataTypeRational,
			rationalData: []Rational{
				*exif.ApertureValue,
			},
		})
	}
	if exif.BrightnessValue != nil {
		entries = append(entries, &idfEntry{
			tag:      tagBrightnessValue,
			dataType: dataTypeSRational,
			sRationalData: []SRational{
				*exif.BrightnessValue,
			},
		})
	}
	if exif.ExposureBiasValue != nil {
		entries = append(entries, &idfEntry{
			tag:      tagExposureBiasValue,
			dataType: dataTypeSRational,
			sRationalData: []SRational{
				*exif.ExposureBiasValue,
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

func (e *encodeState) convertGPSToIDF(gps *GPS) (*idf, error) {
	entries := []*idfEntry{}
	if gps.LatitudeRef != nil {
		entries = append(entries, &idfEntry{
			tag:       tagGPSLatitudeRef,
			dataType:  dataTypeAscii,
			asciiData: *gps.LatitudeRef,
		})
	}
	if gps.Latitude != [3]Rational{} {
		entries = append(entries, &idfEntry{
			tag:          tagGPSLatitude,
			dataType:     dataTypeRational,
			rationalData: gps.Latitude[:],
		})
	}
	if gps.LongitudeRef != nil {
		entries = append(entries, &idfEntry{
			tag:       tagGPSLongitudeRef,
			dataType:  dataTypeAscii,
			asciiData: *gps.LongitudeRef,
		})
	}
	if gps.Longitude != [3]Rational{} {
		entries = append(entries, &idfEntry{
			tag:          tagGPSLongitude,
			dataType:     dataTypeRational,
			rationalData: gps.Longitude[:],
		})
	}
	slices.SortFunc(entries, func(a, b *idfEntry) int {
		return cmp.Compare(a.tag, b.tag)
	})
	return &idf{
		entries: entries,
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
		if len(entry.byteData) <= 4 {
			copy(e.data[offset:offset+4], entry.byteData)
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.byteData))
			copy(e.data[l:], entry.byteData)
		}
		offset += 4

	case dataTypeAscii:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.asciiData)+1))
		offset += 4
		if len(entry.asciiData) <= 3 {
			n := copy(e.data[offset:offset+4], entry.asciiData)
			e.data[offset+uint32(n)] = '\x00' // null terminator
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.asciiData) + 1)
			copy(e.data[l:], entry.asciiData)
			e.data[l+len(entry.asciiData)] = '\x00' // null terminator
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
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.longData)))
		offset += 4
		if len(entry.longData) <= 1 {
			for i, v := range entry.longData {
				e.byteOrder.PutUint32(e.data[offset+4*uint32(i):offset+4*uint32(i)+4], v)
			}
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(4 * len(entry.longData))
			for i, v := range entry.longData {
				e.byteOrder.PutUint32(e.data[l+4*i:l+4*i+4], v)
			}
		}
		offset += 4

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
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.sByteData)))
		offset += 4
		if len(entry.sByteData) <= 4 {
			copyInt8ToUint8(e.data[offset:offset+4], entry.sByteData)
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.sByteData))
			copyInt8ToUint8(e.data[l:], entry.sByteData)
		}
		offset += 4

	case dataTypeUndefined:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.undefinedData)))
		offset += 4
		if len(entry.undefinedData) <= 4 {
			copy(e.data[offset:offset+4], entry.undefinedData)
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.undefinedData))
			copy(e.data[l:], entry.undefinedData)
		}
		offset += 4

	case dataTypeSShort:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.sShortData)))
		offset += 4
		if len(entry.sShortData) <= 2 {
			for i, v := range entry.sShortData {
				e.byteOrder.PutUint16(e.data[offset+2*uint32(i):offset+2*uint32(i)+2], uint16(v))
			}
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(2 * len(entry.sShortData))
			for i, v := range entry.sShortData {
				e.byteOrder.PutUint16(e.data[l+2*i:l+2*i+2], uint16(v))
			}
		}
		offset += 4

	case dataTypeSLong:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.sLongData)))
		offset += 4
		if len(entry.sLongData) <= 1 {
			for i, v := range entry.sLongData {
				e.byteOrder.PutUint32(e.data[offset+4*uint32(i):offset+4*uint32(i)+4], uint32(v))
			}
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(4 * len(entry.sLongData))
			for i, v := range entry.sLongData {
				e.byteOrder.PutUint32(e.data[l+4*i:l+4*i+4], uint32(v))
			}
		}
		offset += 4

	case dataTypeSRational:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.sRationalData)))
		offset += 4
		if len(entry.sRationalData) > 0 {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(8 * len(entry.sRationalData))
			for i, v := range entry.sRationalData {
				e.byteOrder.PutUint32(e.data[l+8*i:l+8*i+4], uint32(v.Numerator))
				e.byteOrder.PutUint32(e.data[l+8*i+4:l+8*i+8], uint32(v.Denominator))
			}
		}
		offset += 4

	case dataTypeFloat:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.floatData)))
		offset += 4
		if len(entry.floatData) <= 1 {
			for i, v := range entry.floatData {
				e.byteOrder.PutUint32(e.data[offset+4*uint32(i):offset+4*uint32(i)+4], math.Float32bits(v))
			}
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(4 * len(entry.floatData))
			for i, v := range entry.floatData {
				e.byteOrder.PutUint32(e.data[l+4*i:l+4*i+4], math.Float32bits(v))
			}
		}
		offset += 4

	case dataTypeDouble:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.doubleData)))
		offset += 4
		l := len(e.data)
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
		e.extend(8 * len(entry.doubleData))
		for i, v := range entry.doubleData {
			e.byteOrder.PutUint64(e.data[l+8*i:l+8*i+8], math.Float64bits(v))
		}
		offset += 4

	case dataTypeUTF8:
		e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(len(entry.utf8data)+1))
		offset += 4
		if len(entry.utf8data) <= 3 {
			n := copy(e.data[offset:offset+4], entry.utf8data)
			e.data[offset+uint32(n)] = '\x00' // null-terminated
		} else {
			l := len(e.data)
			e.byteOrder.PutUint32(e.data[offset:offset+4], uint32(l))
			e.extend(len(entry.utf8data) + 1)
			copy(e.data[l:], entry.utf8data)
			e.data[l+len(entry.utf8data)] = '\x00' // null-terminated
		}
		offset += 4

	default:
		panic(fmt.Sprintf("internal error: unknown data type: %d", entry.dataType))
	}
	e.align()
	return offset, nil
}

func copyInt8ToUint8(dst []byte, src []int8) {
	for i, v := range src {
		dst[i] = uint8(v)
	}
}
