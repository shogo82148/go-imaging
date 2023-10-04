package exif

import (
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
	e.setLen(size)

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

func (e *encodeState) setLen(n int) {
	if n > cap(e.data) {
		e.grow(n - len(e.data))
	}
	e.data = e.data[:n]
}

func convertTIFFToIDF(t *TIFF) (*idf, error) {
	return &idf{
		entries: []*idfEntry{},
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

	}
	return offset + 12, nil
}
