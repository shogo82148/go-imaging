package exif

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"

	"github.com/shogo82148/pointer"
)

type decodeState struct {
	data      []byte
	byteOrder binary.ByteOrder
}

func Decode(r io.Reader) (*TIFF, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	d := &decodeState{data: data}
	return d.decode()
}

func (d *decodeState) decode() (*TIFF, error) {
	// skip Exif marker
	if len(d.data) < 6 {
		return nil, errors.New("exif: invalid data header")
	}
	if string(d.data[:6]) == "Exif\x00\x00" {
		d.data = d.data[6:]
	}

	// parse TIFF header
	if len(d.data) < 8 {
		return nil, errors.New("exif: invalid TIFF header")
	}
	if d.data[0] == 'M' && d.data[1] == 'M' {
		d.byteOrder = binary.BigEndian
	} else if d.data[0] == 'I' && d.data[1] == 'I' {
		d.byteOrder = binary.LittleEndian
	} else {
		return nil, errors.New("exif: invalid TIFF header")
	}
	if d.byteOrder.Uint16(d.data[2:4]) != 0x002a {
		return nil, errors.New("exif: invalid TIFF header")
	}

	// parse IFD0
	offsetIFD := d.byteOrder.Uint32(d.data[4:8])
	idf0, err := d.decodeIFD(offsetIFD)
	if err != nil {
		return nil, err
	}
	var tiff TIFF
	for _, entry := range idf0.entries {
		switch entry.tag {
		case tagImageWidth:
		case tagImageLength:
		case tagBitsPerSample:
		case tagCompression:
		case tagPhotometricInterpretation:
		case tagImageDescription:
			if entry.dataType == dataTypeAscii {
				tiff.ImageDescription = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.ImageDescription = pointer.String(entry.utf8data)
			}
		case tagMake:
			if entry.dataType == dataTypeAscii {
				tiff.Make = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.Make = pointer.String(entry.utf8data)
			}
		case tagModel:
			if entry.dataType == dataTypeAscii {
				tiff.Model = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.Model = pointer.String(entry.utf8data)
			}
		case tagStripOffsets:
		case tagOrientation:
			if entry.dataType == dataTypeShort && len(entry.shortData) == 1 {
				tiff.Orientation = Orientation(entry.shortData[0])
			}
		case tagSamplesPerPixel:
		case tagRowsPerStrip:
		case tagStripByteCounts:
		case tagXResolution:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 1 {
				tiff.XResolution = pointer.Ptr(entry.rationalData[0])
			}
		case tagYResolution:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 1 {
				tiff.YResolution = pointer.Ptr(entry.rationalData[0])
			}
		case tagPlanarConfiguration:
		case tagResolutionUnit:
			if entry.dataType == dataTypeShort && len(entry.shortData) == 1 {
				tiff.ResolutionUnit = ResolutionUnit(entry.shortData[0])
			}
		case tagTransferFunction:
		case tagSoftware:
			if entry.dataType == dataTypeAscii {
				tiff.Software = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.Software = pointer.String(entry.utf8data)
			}
		case tagDateTime:
			if entry.dataType == dataTypeAscii {
				tiff.DateTime = pointer.String(entry.asciiData)
			}
		case tagArtist:
			if entry.dataType == dataTypeAscii {
				tiff.Artist = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.Artist = pointer.String(entry.utf8data)
			}
		case tagWhitePoint:
		case tagPrimaryChromaticities:
		case tagJPEGInterchangeFormat:
		case tagJPEGInterchangeFormatLength:
		case tagYCbCrCoefficients:
		case tagYCbCrSubSampling:
		case tagYCbCrPositioning:
		case tagReferenceBlackWhite:
		case tagCopyright:
			if entry.dataType == dataTypeAscii {
				tiff.Copyright = pointer.String(entry.asciiData)
			} else if entry.dataType == dataTypeUTF8 {
				tiff.Copyright = pointer.String(entry.utf8data)
			}
		case tagExifIFDPointer:
			if entry.dataType == dataTypeLong && len(entry.longData) == 1 {
				exif, err := d.decodeExif(entry.longData[0])
				if err != nil {
					return nil, err
				}
				tiff.Exif = exif
			}
		case tagGPSInfoIFDPointer:
			if entry.dataType == dataTypeLong && len(entry.longData) == 1 {
				gps, err := d.decodeGPS(entry.longData[0])
				if err != nil {
					return nil, err
				}
				tiff.GPS = gps
			}
		}
	}
	return &tiff, nil
}

func (d *decodeState) decodeExif(offset uint32) (*Exif, error) {
	idfExif, err := d.decodeIFD(offset)
	if err != nil {
		return nil, err
	}

	var exif Exif
	for _, entry := range idfExif.entries {
		switch entry.tag {
		case tagExposureTime:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 1 {
				exif.ExposureTime = pointer.Ptr(entry.rationalData[0])
			}
		case tagFNumber:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 1 {
				exif.FNumber = pointer.Ptr(entry.rationalData[0])
			}
		case tagExposureProgram:
			if entry.dataType == dataTypeShort && len(entry.shortData) == 1 {
				exif.ExposureProgram = ExposureProgram(entry.shortData[0])
			}
		case tagSpectralSensitivity:
			if entry.dataType == dataTypeAscii {
				exif.SpectralSensitivity = pointer.String(entry.asciiData)
			}
		case tagISOSpeedRatings:
			if entry.dataType == dataTypeShort {
				exif.ISOSpeedRatings = entry.shortData
			}
		case tagOECF:
		case tagExifVersion:
		case tagDateTimeOriginal:
			if entry.dataType == dataTypeAscii {
				exif.DateTimeOriginal = pointer.String(entry.asciiData)
			}
		case tagDateTimeDigitized:
			if entry.dataType == dataTypeAscii {
				exif.DateTimeDigitized = pointer.String(entry.asciiData)
			}
		case tagComponentsConfiguration:
		case tagCompressedBitsPerPixel:
		case tagShutterSpeedValue:
			if entry.dataType == dataTypeSRational && len(entry.sRationalData) == 1 {
				exif.ShutterSpeedValue = pointer.Ptr(entry.sRationalData[0])
			}
		case tagApertureValue:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 1 {
				exif.ApertureValue = pointer.Ptr(entry.rationalData[0])
			}
		case tagBrightnessValue:
			if entry.dataType == dataTypeSRational && len(entry.sRationalData) == 1 {
				exif.BrightnessValue = pointer.Ptr(entry.sRationalData[0])
			}
		case tagExposureBiasValue:
			if entry.dataType == dataTypeSRational && len(entry.sRationalData) == 1 {
				exif.ExposureBiasValue = pointer.Ptr(entry.sRationalData[0])
			}

		// TODO: implement
		case tagMaxApertureValue:
		case tagSubjectDistance:
		case tagMeteringMode:
		case tagLightSource:
		case tagFlash:
		case tagFocalLength:
		case tagSubjectArea:
		case tagMakerNote:
		case tagUserComment:
		case tagSubsecTime:
		case tagSubsecTimeOriginal:
		case tagSubsecTimeDigitized:
		case tagFlashpixVersion:
		case tagColorSpace:
		case tagPixelXDimension:
		case tagPixelYDimension:
		case tagRelatedSoundFile:
		case tagFlashEnergy:
		case tagSpatialFrequencyResponse:
		case tagFocalPlaneXResolution:
		case tagFocalPlaneYResolution:
		case tagFocalPlaneResolutionUnit:
		case tagSubjectLocation:
		case tagExposureIndex:
		case tagSensingMethod:
		case tagFileSource:
		case tagSceneType:
		case tagCFAPattern:
		case tagCustomRendered:
		case tagExposureMode:
		case tagWhiteBalance:
		case tagDigitalZoomRatio:
		case tagFocalLengthIn35mmFilm:
		case tagSceneCaptureType:
		case tagGainControl:
		case tagContrast:
		case tagSaturation:
		case tagSharpness:
		case tagDeviceSettingDescription:
		case tagSubjectDistanceRange:
		case tagImageUniqueID:
		}
	}
	return &exif, nil
}

func (d *decodeState) decodeGPS(offset uint32) (*GPS, error) {
	idfGPS, err := d.decodeIFD(offset)
	if err != nil {
		return nil, err
	}

	var gps GPS
	for _, entry := range idfGPS.entries {
		switch entry.tag {
		case tagGPSVersionID:
		case tagGPSLatitudeRef:
			if entry.dataType == dataTypeAscii {
				gps.LatitudeRef = pointer.String(entry.asciiData)
			}
		case tagGPSLatitude:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 3 {
				copy(gps.Latitude[:], entry.rationalData)
			}
		case tagGPSLongitudeRef:
			if entry.dataType == dataTypeAscii {
				gps.LongitudeRef = pointer.String(entry.asciiData)
			}
		case tagGPSLongitude:
			if entry.dataType == dataTypeRational && len(entry.rationalData) == 3 {
				copy(gps.Longitude[:], entry.rationalData)
			}

		// TODO: implement
		case tagGPSAltitudeRef:
		case tagGPSAltitude:
		case tagGPSTimeStamp:
		case tagGPSSatellites:
		case tagGPSStatus:
		case tagGPSMeasureMode:
		case tagGPSDOP:
		case tagGPSSpeedRef:
		case tagGPSSpeed:
		case tagGPSTrackRef:
		case tagGPSTrack:
		case tagGPSImgDirectionRef:
		case tagGPSImgDirection:
		case tagGPSMapDatum:
		case tagGPSDestLatitudeRef:
		case tagGPSDestLatitude:
		case tagGPSDestLongitudeRef:
		case tagGPSDestLongitude:
		case tagGPSDestBearingRef:
		case tagGPSDestBearing:
		case tagGPSDestDistanceRef:
		case tagGPSDestDistance:
		case tagGPSProcessingMethod:
		case tagGPSAreaInformation:
		case tagGPSDateStamp:
		case tagGPSDifferential:
		case tagGPSHPositioningError:
		}
	}
	return &gps, nil
}

func (d *decodeState) decodeIFD(offset uint32) (*idf, error) {
	count := d.byteOrder.Uint16(d.data[offset : offset+2])
	entries := make([]*idfEntry, count)
	for i := 0; i < int(count); i++ {
		entryOffset := offset + 2 + 12*uint32(i)
		entry, err := d.decodeIFDEntry(entryOffset)
		if err != nil {
			return nil, err
		}
		entries[i] = entry
	}
	next := d.byteOrder.Uint32(d.data[offset+2+12*uint32(count):])
	return &idf{
		entries:    entries,
		nextOffset: next,
	}, nil
}

func (d *decodeState) decodeIFDEntry(offset uint32) (*idfEntry, error) {
	tag := tag(d.byteOrder.Uint16(d.data[offset : offset+2]))
	typ := dataType(d.byteOrder.Uint16(d.data[offset+2 : offset+4]))
	count := d.byteOrder.Uint32(d.data[offset+4 : offset+8])
	valueOffset := d.byteOrder.Uint32(d.data[offset+8 : offset+12])
	entry := &idfEntry{
		tag:      tag,
		dataType: typ,
	}
	switch typ {
	case dataTypeByte:
		if count <= 4 {
			entry.byteData = d.data[offset+8 : offset+8+count]
		} else {
			entry.byteData = d.data[valueOffset : valueOffset+count]
		}

	case dataTypeAscii:
		if count <= 4 {
			entry.asciiData = bytes2ascii(d.data[offset+8 : offset+8+count])
		} else {
			entry.asciiData = bytes2ascii(d.data[valueOffset : valueOffset+count])
		}

	case dataTypeShort:
		if count <= 2 {
			entry.shortData = d.decodeShort(offset+8, count)
		} else {
			entry.shortData = d.decodeShort(valueOffset, count)
		}

	case dataTypeLong:
		if count <= 1 {
			entry.longData = d.decodeLong(offset+8, count)
		} else {
			entry.longData = d.decodeLong(valueOffset, count)
		}

	case dataTypeRational:
		if count > 0 {
			entry.rationalData = d.decodeRational(valueOffset, count)
		}

	case dataTypeSByte:
		if count <= 4 {
			entry.sByteData = d.decodeSByte(offset+8, count)
		} else {
			entry.sByteData = d.decodeSByte(valueOffset, count)
		}

	case dataTypeUndefined:
		if count <= 4 {
			entry.undefinedData = d.data[offset+8 : offset+8+count]
		} else {
			entry.undefinedData = d.data[valueOffset : valueOffset+count]
		}

	case dataTypeSShort:
		if count <= 2 {
			entry.sShortData = d.decodeSShort(offset+8, count)
		} else {
			entry.sShortData = d.decodeSShort(valueOffset, count)
		}

	case dataTypeSLong:
		if count <= 1 {
			entry.sLongData = d.decodeSLong(offset+8, count)
		} else {
			entry.sLongData = d.decodeSLong(valueOffset, count)
		}

	case dataTypeSRational:
		if count > 0 {
			entry.sRationalData = d.decodeSRational(valueOffset, count)
		}

	case dataTypeFloat:
		if count <= 1 {
			entry.floatData = d.decodeFloat(offset+8, count)
		} else {
			entry.floatData = d.decodeFloat(valueOffset, count)
		}

	case dataTypeDouble:
		if count > 0 {
			entry.doubleData = d.decodeDouble(valueOffset, count)
		}

	case dataTypeUTF8:
		if count <= 4 {
			entry.utf8data = bytes2ascii(d.data[offset+8 : offset+8+count])
		} else {
			entry.utf8data = bytes2ascii(d.data[valueOffset : valueOffset+count])
		}
	}
	return entry, nil
}

func (d *decodeState) decodeShort(offset, count uint32) []uint16 {
	var ret []uint16
	for i := 0; i < int(count); i++ {
		ret = append(ret, d.byteOrder.Uint16(d.data[offset+uint32(i)*2:]))
	}
	return ret
}

func (d *decodeState) decodeLong(offset, count uint32) []uint32 {
	var ret []uint32
	for i := 0; i < int(count); i++ {
		ret = append(ret, d.byteOrder.Uint32(d.data[offset+uint32(i)*4:]))
	}
	return ret
}

func (d *decodeState) decodeRational(offset, count uint32) []Rational {
	var ret []Rational
	for i := 0; i < int(count); i++ {
		ret = append(ret, Rational{
			Numerator:   d.byteOrder.Uint32(d.data[offset+uint32(i)*8:]),
			Denominator: d.byteOrder.Uint32(d.data[offset+uint32(i)*8+4:]),
		})
	}
	return ret
}

func (d *decodeState) decodeSByte(offset, count uint32) []int8 {
	var ret []int8
	for i := 0; i < int(count); i++ {
		ret = append(ret, int8(d.data[offset+uint32(i)]))
	}
	return ret
}

func (d *decodeState) decodeSShort(offset, count uint32) []int16 {
	var ret []int16
	for i := 0; i < int(count); i++ {
		ret = append(ret, int16(d.byteOrder.Uint16(d.data[offset+uint32(i)*2:])))
	}
	return ret
}

func (d *decodeState) decodeSLong(offset, count uint32) []int32 {
	var ret []int32
	for i := 0; i < int(count); i++ {
		ret = append(ret, int32(d.byteOrder.Uint32(d.data[offset+uint32(i)*4:])))
	}
	return ret
}

func (d *decodeState) decodeSRational(offset, count uint32) []SRational {
	var ret []SRational
	for i := 0; i < int(count); i++ {
		ret = append(ret, SRational{
			Numerator:   int32(d.byteOrder.Uint32(d.data[offset+uint32(i)*8:])),
			Denominator: int32(d.byteOrder.Uint32(d.data[offset+uint32(i)*8+4:])),
		})
	}
	return ret
}

func (d *decodeState) decodeFloat(offset, count uint32) []float32 {
	var ret []float32
	for i := 0; i < int(count); i++ {
		ret = append(ret, math.Float32frombits(d.byteOrder.Uint32(d.data[offset+uint32(i)*4:])))
	}
	return ret
}

func (d *decodeState) decodeDouble(offset, count uint32) []float64 {
	var ret []float64
	for i := 0; i < int(count); i++ {
		ret = append(ret, math.Float64frombits(d.byteOrder.Uint64(d.data[offset+uint32(i)*8:])))
	}
	return ret
}

// bytes2ascii converts a null terminated byte slice to a string.
func bytes2ascii(b []byte) string {
	idx := bytes.IndexByte(b, 0x00)
	if idx < 0 {
		return string(b)
	}
	return string(b[:idx])
}
