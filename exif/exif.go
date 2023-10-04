package exif

type Exif struct {
	// Orientation is the orientation of the image.
	Orientation Orientation

	// XResolution is the image resolution in width direction.
	XResolution *Rational

	// YResolution is the image resolution in height direction.
	YResolution *Rational

	// ResolutionUnit is the unit of XResolution and YResolution.
	ResolutionUnit ResolutionUnit

	// ImageDescription is the image title, that is title of the image given by the photographer.
	ImageDescription *string

	// Make is the manufacturer of the recording equipment.
	Make *string

	// Model is the model name or model number of the equipment.
	Model *string

	// Software is the software used for image processing.
	Software *string

	// DateTime is the date and time of image creation.
	DateTime *string

	// Artist is the person who created the image.
	Artist *string

	// Copyright is the copyright.
	Copyright *string
}

type Orientation int

const (
	OrientationUnknown     Orientation = 0
	OrientationTopLeft     Orientation = 1
	OrientationTopRight    Orientation = 2
	OrientationBottomRight Orientation = 3
	OrientationBottomLeft  Orientation = 4
	OrientationLeftTop     Orientation = 5
	OrientationRightTop    Orientation = 6
	OrientationRightBottom Orientation = 7
	OrientationLeftBottom  Orientation = 8
)

type ResolutionUnit int

const (
	ResolutionUnitUnknown    ResolutionUnit = 0
	ResolutionUnitInch       ResolutionUnit = 2
	ResolutionUnitCentimeter ResolutionUnit = 3
)

type idf struct {
	entries    []*idfEntry
	nextOffset uint32
}

type tag uint16

// TIFF metadata tags.
const (
	tagImageWidth                  tag = 0x0100
	tagImageLength                 tag = 0x0101
	tagBitsPerSample               tag = 0x0102
	tagCompression                 tag = 0x0103
	tagPhotometricInterpretation   tag = 0x0106
	tagImageDescription            tag = 0x010e
	tagMake                        tag = 0x010f
	tagModel                       tag = 0x0110
	tagStripOffsets                tag = 0x0111
	tagOrientation                 tag = 0x0112
	tagSamplesPerPixel             tag = 0x0115
	tagRowsPerStrip                tag = 0x0116
	tagStripByteCounts             tag = 0x0117
	tagXResolution                 tag = 0x011a
	tagYResolution                 tag = 0x011b
	tagPlanarConfiguration         tag = 0x011c
	tagResolutionUnit              tag = 0x0128
	tagTransferFunction            tag = 0x012d
	tagSoftware                    tag = 0x0131
	tagDateTime                    tag = 0x0132
	tagArtist                      tag = 0x013b
	tagWhitePoint                  tag = 0x013e
	tagPrimaryChromaticities       tag = 0x013f
	tagJPEGInterchangeFormat       tag = 0x0201
	tagJPEGInterchangeFormatLength tag = 0x0202
	tagYCbCrCoefficients           tag = 0x0211
	tagYCbCrSubSampling            tag = 0x0212
	tagYCbCrPositioning            tag = 0x0213
	tagReferenceBlackWhite         tag = 0x0214
	tagCopyright                   tag = 0x8298

	// tagExposureTime                tag = 0x829a
	// tagFNumber                     tag = 0x829d
	// tagExifIFDPointer              tag = 0x8769
	// tagGPSInfoIFDPointer           tag = 0x8825
	// tagInteroperabilityIFDPointer  tag = 0xa005
	// tagExposureProgram             tag = 0x8822
	// tagSpectralSensitivity         tag = 0x8824
)

// Exif metadata tags.
const (
	tagExifIFDPointer    tag = 0x8769
	tagGPSInfoIFDPointer tag = 0x8825
)

type idfEntry struct {
	tag           tag
	dataType      dataType
	byteData      []byte
	asciiData     string
	shortData     []uint16
	longData      []uint32
	rationalData  []Rational
	sByteData     []int8
	undefinedData []byte
	sShortData    []int16
	sLongData     []int32
	sRationalData []SRational
	floatData     []float32
	doubleData    []float64
	utf8data      string
}

type dataType uint16

const (
	dataTypeByte      dataType = 0x0001
	dataTypeAscii     dataType = 0x0002
	dataTypeShort     dataType = 0x0003
	dataTypeLong      dataType = 0x0004
	dataTypeRational  dataType = 0x0005
	dataTypeSByte     dataType = 0x0006
	dataTypeUndefined dataType = 0x0007
	dataTypeSShort    dataType = 0x0008
	dataTypeSLong     dataType = 0x0009
	dataTypeSRational dataType = 0x000a
	dataTypeFloat     dataType = 0x000b
	dataTypeDouble    dataType = 0x000c
	dataTypeUTF8      dataType = 0x0081
)

type Rational struct {
	Numerator   uint32
	Denominator uint32
}

type SRational struct {
	Numerator   int32
	Denominator int32
}
