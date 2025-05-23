// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"sync"
)

// Interface satisfaction check.
var _ BytesChannel = (*PacketQueue)(nil)

// PacketQueue is loosely modeled after bytes.Buffer.
// It supports automatically writing data into Packets, generating new
// Packets as required and reading over Packet boundaries.
type PacketQueue struct {
	sync.Mutex
	queue                  []*Packet
	indexPacket, indexData int
	recvEOM                bool

	// packetSize should be a function returning the currently used
	// packetSize.
	packetSize func() int
}

// NewPacketQueue returns an initialized PacketQueue.
func NewPacketQueue(packetSize func() int) *PacketQueue {
	queue := &PacketQueue{
		packetSize: packetSize,
	}
	queue.Reset()
	return queue
}

// Reset resets a PacketQueue as if it were newly initialized.
// Note: All queued packets will be discarded.
func (queue *PacketQueue) Reset() {
	queue.Lock()
	defer queue.Unlock()

	queue.queue = []*Packet{}
	queue.indexPacket = 0
	queue.indexData = 0
	queue.recvEOM = false
}

// AddPacket adds a packet to the queue.
func (queue *PacketQueue) AddPacket(packet *Packet) {
	queue.Lock()
	defer queue.Unlock()

	queue.queue = append(queue.queue, packet)
	if packet.Header.Status&TDS_BUFSTAT_EOM == TDS_BUFSTAT_EOM {
		queue.recvEOM = true
	}
}

// Position returns the two indizes used by PacketQueue to store its
// position in the queue and their respective data.
//
// The first returned integer is the packet index. Note that the packet
// index can change in both directions - it grows when bytes are read or
// written and it shrinks when DiscardUntilCurrentPosition is called.
//
// The second returned integer is the data index. The data index points
// to the last unread or unwritten byte of the packet the packet index
// points to.
// The data index only grows when bytes are read or written to the
// queue. It may shrink when DiscardUntilCurrentPosition is called.
func (queue *PacketQueue) Position() (int, int) {
	return queue.indexPacket, queue.indexData
}

// SetPosition sets the two indizes used by PacketQueue.
// See Position for more details.
func (queue *PacketQueue) SetPosition(indexPacket, indexData int) {
	queue.Lock()
	defer queue.Unlock()

	queue.indexPacket = indexPacket
	queue.indexData = indexData
}

// DiscardUntilCurrentPosition discards all consumed packets, indicated
// by the position indizes.
// See Position for more details regarding positions.
func (queue *PacketQueue) DiscardUntilCurrentPosition() {
	queue.Lock()
	defer queue.Unlock()

	// shift queue
	queue.queue = queue.queue[queue.indexPacket:]
	queue.indexPacket = 0

	// indexPacket points to no packet in the queue, reset indexData and
	// return.
	if queue.indexPacket >= len(queue.queue) {
		queue.indexData = 0
		return
	}

	// If indexData is the end of the indexPacket the packet itself can
	// be discarded as well.
	if queue.indexData >= len(queue.queue[queue.indexPacket].Data) {
		queue.queue = queue.queue[1:]
		queue.indexData = 0
	}
}

// AllPacketsConsumend returns true if all packets have been consumed.
func (queue *PacketQueue) AllPacketsConsumed() bool {
	if len(queue.queue) == 0 && queue.indexPacket == 0 && queue.indexData == 0 {
		// No packets in queue also means that all packets have been
		// consumed
		return true
	}

	if queue.indexPacket >= len(queue.queue) {
		return true
	}

	return queue.indexPacket == len(queue.queue)-1 && queue.indexData == len(queue.queue[queue.indexPacket].Data)
}

// IsEOM returns true if all packets have been consumed and it is the
// end of the message.
func (queue *PacketQueue) IsEOM() bool {
	return queue.AllPacketsConsumed() && queue.recvEOM
}

// Read satisfies the io.Reader interface.
func (queue *PacketQueue) Read(p []byte) (int, error) {
	var err error
	p, err = queue.Bytes(len(p))
	return len(p), err
}

// Write satisfies the io.Writer interface.
func (queue *PacketQueue) Write(p []byte) (int, error) {
	return len(p), queue.WriteBytes(p)
}

// Read methods

// Bytes returns a slice of bytes from the queue.
//
// The returned byte slice will always be of length n.
//
// If there aren't enough bytes to read n bytes Bytes will return
// a wrapped io.EOF. The returned byte slice will still be of length n.
func (queue *PacketQueue) Bytes(n int) ([]byte, error) {
	queue.Lock()
	defer queue.Unlock()

	if n == 0 {
		return []byte{}, nil
	}

	bs := make([]byte, n)
	// bsOffset is the index in the return slice where data still needs
	// to be written.
	bsOffset := 0

	for {
		if queue.AllPacketsConsumed() {
			// All available packets have been consumed
			return bs, ErrNotEnoughBytes
		}
		data := queue.queue[queue.indexPacket].Data

		startIndex := queue.indexData
		// (n - bsOffset) is the amount of bytes that still need to be
		// read.
		endIndex := queue.indexData + (n - bsOffset)
		if endIndex > len(data) {
			endIndex = len(data)
		}

		copy(bs[bsOffset:], data[startIndex:endIndex])
		bsOffset += endIndex - startIndex

		queue.indexData = endIndex
		// Move indizes forward if the current packet is consumed
		// entirely.
		if queue.indexData == len(data) {
			queue.indexPacket += 1
			queue.indexData = 0
		}

		if bsOffset == n {
			break
		}
	}

	return bs, nil
}

// Byte implements the tds.BytesChannel interface.
func (queue *PacketQueue) Byte() (byte, error) {
	bs, err := queue.Bytes(1)
	return bs[0], err
}

// Uint8 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Uint8() (uint8, error) {
	b, err := queue.Byte()
	return uint8(b), err
}

// Int8 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Int8() (int8, error) {
	b, err := queue.Byte()
	return int8(b), err
}

// Uint16 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Uint16() (uint16, error) {
	bs, err := queue.Bytes(2)
	return endian.Uint16(bs), err
}

// Int16 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Int16() (int16, error) {
	i, err := queue.Uint16()
	return int16(i), err
}

// Uint32 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Uint32() (uint32, error) {
	bs, err := queue.Bytes(4)
	return endian.Uint32(bs), err
}

// Int32 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Int32() (int32, error) {
	i, err := queue.Uint32()
	return int32(i), err
}

// Uint64 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Uint64() (uint64, error) {
	bs, err := queue.Bytes(8)
	return endian.Uint64(bs), err
}

// Int64 implements the tds.BytesChannel interface.
func (queue *PacketQueue) Int64() (int64, error) {
	i, err := queue.Uint64()
	return int64(i), err
}

// String implements the tds.BytesChannel interface.
func (queue *PacketQueue) String(size int) (string, error) {
	bs, err := queue.Bytes(size)
	return string(bs), err
}

// Write methods

// WriteBytes writes a slice of bytes.
//
// The returned integer is the size of bs, the returned error is always nil.
func (queue *PacketQueue) WriteBytes(bs []byte) error {
	queue.Lock()
	defer queue.Unlock()

	if len(bs) == 0 {
		return nil
	}

	bsOffset := 0

	for bsOffset < len(bs) {
		// Add new packet if the index points to no packet
		if queue.indexPacket == len(queue.queue) {
			queue.queue = append(queue.queue, NewPacket(queue.packetSize()))
		}

		// Retrieve current package and calculate how many bytes can
		// still be written to it.
		curPacket := queue.queue[queue.indexPacket]
		freeBytes := int(curPacket.Header.Length) - PacketHeaderSize - queue.indexData

		// No free bytes, add a new packet.
		if freeBytes == 0 {
			curPacket = NewPacket(queue.packetSize())
			queue.queue = append(queue.queue, curPacket)
			queue.indexPacket++
			queue.indexData = 0
			freeBytes = int(curPacket.Header.Length) - PacketHeaderSize
		}

		// Calculate how many bytes are left in bs if more free bytes
		// are available in the packet than are left in bs.
		if freeBytes > len(bs)-bsOffset {
			freeBytes = len(bs) - bsOffset
		}

		copy(curPacket.Data[queue.indexData:], bs[bsOffset:bsOffset+freeBytes])
		bsOffset += freeBytes
		queue.indexData += freeBytes
	}

	return nil
}

// WriteByte implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteByte(b byte) error {
	return queue.WriteBytes([]byte{b})
}

// WriteUint8 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteUint8(i uint8) error {
	return queue.WriteByte(byte(i))
}

// WriteInt8 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteInt8(i int8) error {
	return queue.WriteUint8(uint8(i))
}

// WriteUint16 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteUint16(i uint16) error {
	bs := make([]byte, 2)
	endian.PutUint16(bs, i)
	return queue.WriteBytes(bs)
}

// WriteInt16 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteInt16(i int16) error {
	return queue.WriteUint16(uint16(i))
}

// WriteUint32 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteUint32(i uint32) error {
	bs := make([]byte, 4)
	endian.PutUint32(bs, i)
	return queue.WriteBytes(bs)
}

// WriteInt32 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteInt32(i int32) error {
	return queue.WriteUint32(uint32(i))
}

// WriteUint64 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteUint64(i uint64) error {
	bs := make([]byte, 8)
	endian.PutUint64(bs, i)
	return queue.WriteBytes(bs)
}

// WriteInt64 implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteInt64(i int64) error {
	return queue.WriteUint64(uint64(i))
}

// WriteString implements the tds.BytesChannel interface.
func (queue *PacketQueue) WriteString(s string) error {
	return queue.WriteBytes([]byte(s))
}
