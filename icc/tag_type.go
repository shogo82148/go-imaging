package icc

import "fmt"

type TagType uint32

const (
	TagTypeColorantOrder             TagType = 0x636c726f // 'clro'
	TagTypeCurve                     TagType = 0x63757276 // 'curv'
	TagTypeDataType                  TagType = 0x64617461 // 'data'
	TagTypeDateTime                  TagType = 0x6474696d // 'dtim'
	TagTypeDict                      TagType = 0x64637420 // 'dict'
	TagTypeEmbeddedHeightImage       TagType = 0x6568696d // 'ehim'
	TagTypeEmbeddedNormalImage       TagType = 0x656e696d // 'enim'
	TagTypeFloat16Array              TagType = 0x666c3136 // 'fl16'
	TagTypeFloat32Array              TagType = 0x666c3234 // 'fl32'
	TagTypeFloat64Array              TagType = 0x666c3634 // 'fl64'
	TagTypeLutAtoB                   TagType = 0x6d414220 // 'mAB '
	TagTypeLutBtoA                   TagType = 0x6d424120 // 'mBA '
	TagTypeMeasurement               TagType = 0x6d656173 // 'meas'
	TagTypeMultiLocalizedUnicode     TagType = 0x6d6c7563 // 'mluc'
	TagTypeMultiProcessElements      TagType = 0x6d706574 // 'mpet'
	TagTypeParametricCurve           TagType = 0x70617261 // 'para'
	TagTypeS15Fixed16Array           TagType = 0x73663332 // 'sf32'
	TagTypeSignature                 TagType = 0x73696720 // 'sig '
	TagTypeSparseMatrixArray         TagType = 0x736d6174 // 'smat'
	TagTypeSpectralViewingConditions TagType = 0x7376636e // 'svcn'
	TagTypeTagArrayType              TagType = 0x74617279 // 'tary'
	TagTypeTagStruct                 TagType = 0x74737472 // 'tstr'
	TagTypeU16Fixed16Array           TagType = 0x75663332 // 'uf32'
	TagTypeUint16Array               TagType = 0x75693136 // 'ui16'
	TagTypeUint32Array               TagType = 0x75693332 // 'ui32'
	TagTypeUint64Array               TagType = 0x75693634 // 'ui64'
	TagTypeUint8Array                TagType = 0x75693038 // 'ui08'
	TagTypeUTF16                     TagType = 0x75743136 // 'ut16'
	TagTypeUTF8                      TagType = 0x75746638 // 'utf8'
	TagTypeUTF8Zip                   TagType = 0x7a757438 // 'zut8'
	TagTypeXYZ                       TagType = 0x58595a20 // 'XYZ '
	TagTypeZipXML                    TagType = 0x7a786d6c // 'zxml'
)

func (t TagType) String() string {
	switch t {
	case TagTypeColorantOrder:
		return "ColorantOrder"
	case TagTypeCurve:
		return "Curve"
	case TagTypeDataType:
		return "DataType"
	case TagTypeDateTime:
		return "DateTime"
	case TagTypeDict:
		return "Dict"
	case TagTypeEmbeddedHeightImage:
		return "EmbeddedHeightImage"
	case TagTypeEmbeddedNormalImage:
		return "EmbeddedNormalImage"
	case TagTypeFloat16Array:
		return "Float16Array"
	case TagTypeFloat32Array:
		return "Float32Array"
	case TagTypeFloat64Array:
		return "Float64Array"
	case TagTypeLutAtoB:
		return "LutAtoB"
	case TagTypeLutBtoA:
		return "LutBtoA"
	case TagTypeMeasurement:
		return "Measurement"
	case TagTypeMultiLocalizedUnicode:
		return "MultiLocalizedUnicode"
	case TagTypeMultiProcessElements:
		return "MultiProcessElements"
	case TagTypeParametricCurve:
		return "ParametricCurve"
	case TagTypeS15Fixed16Array:
		return "S15Fixed16Array"
	case TagTypeSignature:
		return "Signature"
	case TagTypeSparseMatrixArray:
		return "SparseMatrixArray"
	case TagTypeSpectralViewingConditions:
		return "SpectralViewingConditions"
	case TagTypeTagArrayType:
		return "TagArrayType"
	case TagTypeTagStruct:
		return "TagStruct"
	case TagTypeU16Fixed16Array:
		return "U16Fixed16Array"
	case TagTypeUint16Array:
		return "Uint16Array"
	case TagTypeUint32Array:
		return "Uint32Array"
	case TagTypeUint64Array:
		return "Uint64Array"
	case TagTypeUint8Array:
		return "Uint8Array"
	case TagTypeUTF16:
		return "UTF16"
	case TagTypeUTF8:
		return "UTF8"
	case TagTypeUTF8Zip:
		return "UTF8Zip"
	case TagTypeXYZ:
		return "XYZ"
	case TagTypeZipXML:
		return "ZipXML"
	default:
		return fmt.Sprintf(
			"Unknown Tag Type(%08xh '%c%c%c%c')",
			uint32(t),
			printable(byte(t>>24)),
			printable(byte(t>>16)),
			printable(byte(t>>8)),
			printable(byte(t)),
		)
	}
}
