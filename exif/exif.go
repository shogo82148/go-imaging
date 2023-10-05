package exif

import "strconv"

type TIFF struct {
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

	// Exif is the Exif information.
	Exif *Exif

	//  is the GPS information.
	*GPS
}

type Exif struct {
	// ExposureTime is the exposure time, given in seconds.
	ExposureTime *Rational

	// FNumber is the F number.
	FNumber *Rational

	// ExposureProgram is the exposure program.
	ExposureProgram ExposureProgram

	// SpectralSensitivity is the spectral sensitivity of each channel of the camera used.
	SpectralSensitivity *string

	// ISOSpeedRatings is the ISO speed and ISO latitude of the camera or input device as specified in ISO 12232.
	ISOSpeedRatings []uint16

	// DateTimeOriginal is the date and time when the original image data was generated.
	DateTimeOriginal *string

	// DateTimeDigitized is the date and time when the image was stored as digital data.
	DateTimeDigitized *string

	// ShutterSpeedValue is the shutter speed.
	ShutterSpeedValue *SRational

	// ApertureValue is the lens aperture.
	ApertureValue *Rational

	// BrightnessValue is the value of brightness.
	BrightnessValue *SRational

	// ExposureBiasValue is the exposure bias.
	ExposureBiasValue *SRational
}

type GPS struct {
	// VersionID is the version of GPSInfoIFD.
	VersionID [4]byte

	// LatitudeRef is the latitude of the GPS receiver.
	LatitudeRef *string

	// Latitude is the latitude of the GPS receiver.
	Latitude [3]Rational

	// LongitudeRef is the longitude of the GPS receiver.
	LongitudeRef *string

	// Longitude is the longitude of the GPS receiver.
	Longitude [3]Rational

	// AltitudeRef is the altitude of the GPS receiver.
	AltitudeRef *byte

	// Altitude is the altitude of the GPS receiver.
	Altitude *Rational

	// TimeStamp is the time as UTC (Coordinated Universal Time).
	TimeStamp [3]Rational

	// Satellites is the GPS satellites used for measurements.
	Satellites *string

	// Status is the status of the GPS receiver when the image is recorded.
	Status *string

	// MeasureMode is the GPS measurement mode.
	MeasureMode *string

	// DOP is the measurement precision of GPS receiver.
	DOP Rational

	// SpeedRef is the speed unit of GPS receiver.
	SpeedRef *string

	// Speed is the speed of GPS receiver.
	Speed *Rational

	// TrackRef is the reference for giving the direction of GPS receiver movement.
	TrackRef *string

	// Track is the direction of GPS receiver movement.
	Track *Rational

	// ImgDirectionRef is the reference for giving the direction of the image when it is captured.
	ImgDirectionRef *string

	// ImgDirection is the direction of the image when it is captured.
	ImgDirection *Rational

	// MapDatum is the geodetic survey data used by the GPS receiver.
	MapDatum *string

	// DestLatitudeRef is the latitude of the destination point.
	DestLatitudeRef *string

	// DestLatitude is the latitude of the destination point.
	DestLatitude [3]Rational

	// DestLongitudeRef is the longitude of the destination point.
	DestLongitudeRef *string

	// DestLongitude is the longitude of the destination point.
	DestLongitude [3]Rational

	// DestBearingRef is the reference used for giving the bearing to the destination point.
	DestBearingRef *string

	// DestBearing is the bearing to the destination point.
	DestBearing *Rational

	// DestDistanceRef is the unit used for expressing the distance to the destination point.
	DestDistanceRef *string

	// DestDistance is the distance to the destination point.
	DestDistance *Rational

	// ProcessingMethod is the name of the method used for location finding.
	ProcessingMethod []byte

	// AreaInformation is the name of the GPS area.
	AreaInformation []byte

	// DateStamp is the date as UTC (Coordinated Universal Time).
	DateStamp *string

	// Differential is the differential correction applied to the GPS receiver.
	Differential *uint16

	// HPositioningError is the horizontal positioning errors.
	HPositioningError *Rational
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

func (o Orientation) String() string {
	switch o {
	case OrientationTopLeft:
		return "TopLeft"
	case OrientationTopRight:
		return "TopRight"
	case OrientationBottomRight:
		return "BottomRight"
	case OrientationBottomLeft:
		return "BottomLeft"
	case OrientationLeftTop:
		return "LeftTop"
	case OrientationRightTop:
		return "RightTop"
	case OrientationRightBottom:
		return "RightBottom"
	case OrientationLeftBottom:
		return "LeftBottom"
	default:
		return "Unknown(" + strconv.Itoa(int(o)) + ")"
	}
}

type ResolutionUnit int

const (
	ResolutionUnitUnknown    ResolutionUnit = 0
	ResolutionUnitInch       ResolutionUnit = 2
	ResolutionUnitCentimeter ResolutionUnit = 3
)

func (u ResolutionUnit) String() string {
	switch u {
	case ResolutionUnitInch:
		return "Inch"
	case ResolutionUnitCentimeter:
		return "Centimeter"
	default:
		return "Unknown(" + strconv.Itoa(int(u)) + ")"
	}
}

type ExposureProgram int

const (
	ExposureProgramUnknown          ExposureProgram = 0
	ExposureProgramManual           ExposureProgram = 1
	ExposureProgramNormal           ExposureProgram = 2
	ExposureProgramAperturePriority ExposureProgram = 3
	ExposureProgramShutterPriority  ExposureProgram = 4
	ExposureProgramCreative         ExposureProgram = 5
	ExposureProgramAction           ExposureProgram = 6
	ExposureProgramPortrait         ExposureProgram = 7
	ExposureProgramLandscape        ExposureProgram = 8
)

func (p ExposureProgram) String() string {
	switch p {
	case ExposureProgramManual:
		return "Manual"
	case ExposureProgramNormal:
		return "Normal"
	case ExposureProgramAperturePriority:
		return "AperturePriority"
	case ExposureProgramShutterPriority:
		return "ShutterPriority"
	case ExposureProgramCreative:
		return "Creative"
	case ExposureProgramAction:
		return "Action"
	case ExposureProgramPortrait:
		return "Portrait"
	case ExposureProgramLandscape:
		return "Landscape"
	default:
		return "Unknown(" + strconv.Itoa(int(p)) + ")"
	}
}

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
)

// Exif metadata tags.
const (
	tagExifIFDPointer    tag = 0x8769
	tagGPSInfoIFDPointer tag = 0x8825
)

// Exif IFD metadata tags.
const (
	tagExposureTime             tag = 0x829a
	tagFNumber                  tag = 0x829d
	tagExposureProgram          tag = 0x8822
	tagSpectralSensitivity      tag = 0x8824
	tagISOSpeedRatings          tag = 0x8827
	tagOECF                     tag = 0x8828
	tagExifVersion              tag = 0x9000
	tagDateTimeOriginal         tag = 0x9003
	tagDateTimeDigitized        tag = 0x9004
	tagComponentsConfiguration  tag = 0x9101
	tagCompressedBitsPerPixel   tag = 0x9102
	tagShutterSpeedValue        tag = 0x9201
	tagApertureValue            tag = 0x9202
	tagBrightnessValue          tag = 0x9203
	tagExposureBiasValue        tag = 0x9204
	tagMaxApertureValue         tag = 0x9205
	tagSubjectDistance          tag = 0x9206
	tagMeteringMode             tag = 0x9207
	tagLightSource              tag = 0x9208
	tagFlash                    tag = 0x9209
	tagFocalLength              tag = 0x920a
	tagSubjectArea              tag = 0x9214
	tagMakerNote                tag = 0x927c
	tagUserComment              tag = 0x9286
	tagSubsecTime               tag = 0x9290
	tagSubsecTimeOriginal       tag = 0x9291
	tagSubsecTimeDigitized      tag = 0x9292
	tagFlashpixVersion          tag = 0xa000
	tagColorSpace               tag = 0xa001
	tagPixelXDimension          tag = 0xa002
	tagPixelYDimension          tag = 0xa003
	tagRelatedSoundFile         tag = 0xa004
	tagFlashEnergy              tag = 0xa20b
	tagSpatialFrequencyResponse tag = 0xa20c
	tagFocalPlaneXResolution    tag = 0xa20e
	tagFocalPlaneYResolution    tag = 0xa20f
	tagFocalPlaneResolutionUnit tag = 0xa210
	tagSubjectLocation          tag = 0xa214
	tagExposureIndex            tag = 0xa215
	tagSensingMethod            tag = 0xa217
	tagFileSource               tag = 0xa300
	tagSceneType                tag = 0xa301
	tagCFAPattern               tag = 0xa302
	tagCustomRendered           tag = 0xa401
	tagExposureMode             tag = 0xa402
	tagWhiteBalance             tag = 0xa403
	tagDigitalZoomRatio         tag = 0xa404
	tagFocalLengthIn35mmFilm    tag = 0xa405
	tagSceneCaptureType         tag = 0xa406
	tagGainControl              tag = 0xa407
	tagContrast                 tag = 0xa408
	tagSaturation               tag = 0xa409
	tagSharpness                tag = 0xa40a
	tagDeviceSettingDescription tag = 0xa40b
	tagSubjectDistanceRange     tag = 0xa40c
	tagImageUniqueID            tag = 0xa420
)

// IFD metadata tags.
const (
	tagGPSVersionID         tag = 0x0000
	tagGPSLatitudeRef       tag = 0x0001
	tagGPSLatitude          tag = 0x0002
	tagGPSLongitudeRef      tag = 0x0003
	tagGPSLongitude         tag = 0x0004
	tagGPSAltitudeRef       tag = 0x0005
	tagGPSAltitude          tag = 0x0006
	tagGPSTimeStamp         tag = 0x0007
	tagGPSSatellites        tag = 0x0008
	tagGPSStatus            tag = 0x0009
	tagGPSMeasureMode       tag = 0x000a
	tagGPSDOP               tag = 0x000b
	tagGPSSpeedRef          tag = 0x000c
	tagGPSSpeed             tag = 0x000d
	tagGPSTrackRef          tag = 0x000e
	tagGPSTrack             tag = 0x000f
	tagGPSImgDirectionRef   tag = 0x0010
	tagGPSImgDirection      tag = 0x0011
	tagGPSMapDatum          tag = 0x0012
	tagGPSDestLatitudeRef   tag = 0x0013
	tagGPSDestLatitude      tag = 0x0014
	tagGPSDestLongitudeRef  tag = 0x0015
	tagGPSDestLongitude     tag = 0x0016
	tagGPSDestBearingRef    tag = 0x0017
	tagGPSDestBearing       tag = 0x0018
	tagGPSDestDistanceRef   tag = 0x0019
	tagGPSDestDistance      tag = 0x001a
	tagGPSProcessingMethod  tag = 0x001b
	tagGPSAreaInformation   tag = 0x001c
	tagGPSDateStamp         tag = 0x001d
	tagGPSDifferential      tag = 0x001e
	tagGPSHPositioningError tag = 0x001f
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
