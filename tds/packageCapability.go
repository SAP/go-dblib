// SPDX-FileCopyrightText: 2020 - 2025 SAP SE
//
// SPDX-License-Identifier: Apache-2.0

package tds

import (
	"fmt"
	"math"
)

//go:generate stringer -type=CapabilityType

// CapabilityType is the type of capabilities.
type CapabilityType byte

// Types of capabilities.
const (
	CapabilityRequest CapabilityType = iota + 1
	CapabilityResponse
	CapabilitySecurity
)

//go:generate stringer -type=RequestCapability

// RequestCapability is the type for request capabilities.
type RequestCapability int

// Types of request capabilities.
const (
	TDS_REQ_LANG RequestCapability = iota + 1
	TDS_REQ_RPC
	TDS_REQ_EVT
	TDS_REQ_MSTMT
	TDS_REQ_BCP
	TDS_REQ_CURSOR
	TDS_REQ_DYNF
	TDS_REQ_MSG
	TDS_REQ_PARAM
	TDS_DATA_INT1
	TDS_DATA_INT2
	TDS_DATA_INT4
	TDS_DATA_BIT
	TDS_DATA_CHAR
	TDS_DATA_VCHAR
	TDS_DATA_BIN
	TDS_DATA_VBIN
	TDS_DATA_MNY8
	TDS_DATA_MNY4
	TDS_DATA_DATE8
	TDS_DATA_DATE4
	TDS_DATA_FLT4
	TDS_DATA_FLT8
	TDS_DATA_NUM
	TDS_DATA_TEXT
	TDS_DATA_IMAGE
	TDS_DATA_DEC
	TDS_DATA_LCHAR
	TDS_DATA_LBIN
	TDS_DATA_INTN
	TDS_DATA_DATETIMEN
	TDS_DATA_MONEYN
	TDS_CSR_PREV
	TDS_CSR_FIRST
	TDS_CSR_LAST
	TDS_CSR_ABS
	TDS_CSR_REL
	TDS_CSR_MULTI
	TDS_CON_OOB
	TDS_CON_INBAND
	TDS_CON_LOGICAL
	TDS_PROTO_TEXT
	TDS_PROTO_BULK
	TDS_REQ_URGEVT
	TDS_DATA_SENSITIVITY
	TDS_DATA_BOUNDARY
	TDS_PROTO_DYNAMIC
	TDS_PROTO_DYNPROC
	TDS_DATA_FLTN
	TDS_DATA_BITN
	TDS_DATA_INT8
	TDS_DATA_VOID
	TDS_DOL_BULK
	TDS_OBJECT_JAVA1
	TDS_OBJECT_CHAR
	TDS_REQ_RESERVED1
	TDS_OBJECT_BINARY
	TDS_DATA_COLUMNSTATUS
	TDS_WIDETABLES
	TDS_REQ_RESERVED2
	TDS_DATA_UINT2
	TDS_DATA_UINT4
	TDS_DATA_UINT8
	TDS_DATA_UINTN
	TDS_CUR_IMPLICIT
	TDS_DATA_NLBIN
	TDS_IMAGE_NCHAR
	TDS_BLOB_NCHAR_16
	TDS_BLOB_NCHAR_8
	TDS_BLOB_NCHAR_SCSU
	TDS_DATA_DATE
	TDS_DATA_TIME
	TDS_DATA_INTERVAL
	TDS_CSR_SCROLL
	TDS_CSR_SENSITIVE
	TDS_CSR_INSENSITIVE
	TDS_CSR_SEMISENSITIVE
	TDS_CSR_KEYSETDRIVEN
	TDS_REQ_SRVPKTSIZE
	TDS_DATA_UNITEXT
	TDS_CAP_CLUSTERFAILOVER
	TDS_DATA_SINT1
	TDS_REQ_LARGEIDENT
	TDS_REQ_BLOB_NCHAR_16
	TDS_DATA_XML
	TDS_REQ_CURINFO3
	TDS_REQ_DBRPC2
	TDS_UNUSED_REQ
	TDS_REQ_MIGRATE
	TDS_MULTI_REQUESTS
	TDS_REQ_OPTIONCMD2
	TDS_REQ_LOGINFO
	TDS_DATA_BIGDATETIME
	TDS_DATA_USECS
	TDS_RPCPARAM_LOB
	TDS_REQ_INSTID
	TDS_REQ_GRID
	TDS_REQ_DYN_BATCH
	TDS_REQ_LANG_BATCH
	TDS_REQ_RPC_BATCH
	TDS_DATA_LOBLOCATOR
	TDS_REQ_ROWCOUNT_FOR_SELECT
	TDS_REQ_LOGPARAMS
	TDS_REQ_DYNAMIC_SUPPRESS_PARAMFMT
	TDS_REQ_READONLY
	TDS_REQ_COMMAND_ENCRYPTION
)

//go:generate stringer -type=ResponseCapability

// ResponseCapability is the type for response capabilities.
type ResponseCapability int

// Types of response capabilities.
const (
	TDS_RES_NOMSG ResponseCapability = iota + 1
	TDS_RES_NOEED
	TDS_RES_NOPARAM
	TDS_DATA_NOINT1
	TDS_DATA_NOINT2
	TDS_DATA_NOINT4
	TDS_DATA_NOBIT
	TDS_DATA_NOCHAR
	TDS_DATA_NOVCHAR
	TDS_DATA_NOBIN
	TDS_DATA_NOVBIN
	TDS_DATA_NOMNY8
	TDS_DATA_NOMNY4
	TDS_DATA_NODATE8
	TDS_DATA_NODATE4
	TDS_DATA_NOFLT4
	TDS_DATA_NOFLT8
	TDS_DATA_NONUM
	TDS_DATA_NOTEXT
	TDS_DATA_NOIMAGE
	TDS_DATA_NODEC
	TDS_DATA_NOLCHAR
	TDS_DATA_NOLBIN
	TDS_DATA_NOINTN
	TDS_DATA_NODATETIMEN
	TDS_DATA_NOMONEYN
	TDS_CON_NOOOB
	TDS_CON_NOINBAND
	TDS_PROTO_NOTEXT
	TDS_PROTO_NOBULK
	TDS_DATA_NOSENSITIVITY
	TDS_DATA_NOBOUNDARY
	TDS_RES_NOTDSDEBUG
	TDS_RES_NOSTRIPBLANKS
	TDS_DATA_NOINT8
	TDS_OBJECT_NOJAVA1
	TDS_OBJECT_NOCHAR
	TDS_DATA_NOCOLUMNSTATUS
	TDS_OBJECT_NOBINARY
	TDS_RES_RESERVED
	TDS_DATA_NOUINT2
	TDS_DATA_NOUINT4
	TDS_DATA_NOUINT8
	TDS_DATA_NOUINTN
	TDS_NOWIDETABLES
	TDS_DATA_NONLBIN
	TDS_IMAGE_NONCHAR
	TDS_BLOB_NONCHAR_16
	TDS_BLOB_NONCHAR_8
	TDS_BLOB_NONCHAR_SCSU
	TDS_DATA_NODATE
	TDS_DATA_NOTIME
	TDS_DATA_NOINTERVAL
	TDS_DATA_NOUNITEXT
	TDS_DATA_NOSINT1
	TDS_NO_LARGEIDENT
	TDS_NO_BLOB_NCHAR_16
	TDS_NO_SRVPKTSIZE
	TDS_DATA_NOXML
	TDS_NONINT_RETURN_VALUE
	TDS_RES_NOXNLMETADATA
	TDS_RES_SUPPRESS_FMT
	TDS_RES_SUPPRESS_DONEINPROC
	TDS_UNUSED_RES
	TDS_DATA_NOBIGDATETIME
	TDS_DATA_NOUSECS
	TDS_RES_NO_TDSCONTROL
	TDS_RPCPARAM_NOLOB
	TDS_DATA_NOLOBLOCATOR
	TDS_RES_NOROWCOUNT_FOR_SELECT
	TDS_RES_CUMULATIVE_DONE
	TDS_RES_LIST_DR_MAP
	TDS_RES_DR_NOKILL
)

// SecurityCapability is the type for security capabilities.
type SecurityCapability int

// CapabilityPackage is used to communicate capabilities between clients
// and servers.
type CapabilityPackage struct {
	Capabilities map[CapabilityType]*valueMask
}

// NewCapabilityPackage returns a prepared CapabilityPackage.
func NewCapabilityPackage(request []RequestCapability, response []ResponseCapability,
	security []SecurityCapability) (*CapabilityPackage, error) {
	pkg := &CapabilityPackage{
		Capabilities: make(map[CapabilityType]*valueMask, 3),
	}
	pkg.Capabilities[CapabilityRequest] = newValueMask(int(TDS_REQ_COMMAND_ENCRYPTION))
	pkg.Capabilities[CapabilityResponse] = newValueMask(int(TDS_RES_DR_NOKILL))
	pkg.Capabilities[CapabilitySecurity] = newValueMask(0)

	for _, capa := range request {
		if err := pkg.SetRequestCapability(capa, true); err != nil {
			return nil, err
		}
	}

	for _, capa := range response {
		if err := pkg.SetResponseCapability(capa, true); err != nil {
			return nil, err
		}
	}

	for _, capa := range security {
		if err := pkg.SetSecurityCapability(capa, true); err != nil {
			return nil, err
		}
	}

	return pkg, nil
}

// SetRequestCapability sets the requested capability.
func (pkg *CapabilityPackage) SetRequestCapability(capability RequestCapability, enable bool) error {
	return pkg.Capabilities[CapabilityRequest].setCapability(int(capability), enable)
}

// SetResponseCapability sets the response capability.
func (pkg *CapabilityPackage) SetResponseCapability(capability ResponseCapability, enable bool) error {
	return pkg.Capabilities[CapabilityResponse].setCapability(int(capability), enable)
}

// SetSecurityCapability sets the security capability.
func (pkg *CapabilityPackage) SetSecurityCapability(capability SecurityCapability, enable bool) error {
	return pkg.Capabilities[CapabilitySecurity].setCapability(int(capability), enable)
}

// HasCapability returns whether the package has the requested capability.
func (pkg *CapabilityPackage) HasCapability(capabilityType CapabilityType, capability int) bool {
	return pkg.Capabilities[capabilityType].getCapability(int(capability))
}

// HasRequestCapability returns whether the package has the requested
// request capability.
func (pkg *CapabilityPackage) HasRequestCapability(capability RequestCapability) bool {
	return pkg.HasCapability(CapabilityRequest, int(capability))
}

// HasResponseCapabilities returns whether the package has the requested
// response capability.
func (pkg *CapabilityPackage) HasResponseCapability(capability ResponseCapability) bool {
	return pkg.HasCapability(CapabilityResponse, int(capability))
}

// HasSecurityCapabilities returns whether the package has the requested
// security capability.
func (pkg *CapabilityPackage) HasSecurityCapability(capability SecurityCapability) bool {
	return pkg.HasCapability(CapabilitySecurity, int(capability))
}

// ReadFrom implements the tds.Package interface.
func (pkg *CapabilityPackage) ReadFrom(ch BytesChannel) error {
	totalLength, err := ch.Uint16()
	if err != nil {
		return ErrNotEnoughBytes
	}

	// Read out each capability and its value mask
	length := 0
	for length < int(uint(totalLength)) {
		b, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		length++
		capType := CapabilityType(b)

		capLength, err := ch.Uint8()
		if err != nil {
			return ErrNotEnoughBytes
		}
		length++

		bs, err := ch.Bytes(int(capLength))
		if err != nil {
			return ErrNotEnoughBytes
		}
		length += int(capLength)

		pkg.Capabilities[capType] = parseValueMask(bs)
	}

	if length > int(uint(totalLength)) {
		return fmt.Errorf("read %d bytes instead of %d", length, totalLength)
	}

	return nil
}

// WriteTo implements the tds.Package interface.
func (pkg CapabilityPackage) WriteTo(ch BytesChannel) error {
	if err := ch.WriteByte(byte(TDS_CAPABILITY)); err != nil {
		return fmt.Errorf("failed to write token: %w", err)
	}

	// write length
	bytesToWrite := 0
	for _, vm := range pkg.Capabilities {
		// If no capabilities are set the capability type will be
		// skipped
		if vm.isEmpty() {
			continue
		}

		// Capability type byte and the type's value mask
		bytesToWrite += 2 + len(vm.Bytes())
	}

	if err := ch.WriteUint16(uint16(bytesToWrite)); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	writtenBytes := 0
	for typ, vm := range pkg.Capabilities {
		if vm.isEmpty() {
			continue
		}

		// Write type
		if err := ch.WriteByte(byte(typ)); err != nil {
			return fmt.Errorf("error writing capability type %s: %w", typ, err)
		}
		writtenBytes++

		bs := vm.Bytes()

		// Write length
		if err := ch.WriteUint8(uint8(len(bs))); err != nil {
			return fmt.Errorf("error writing value mask length: %w", err)
		}
		writtenBytes++

		// Write value mask
		if err := ch.WriteBytes(bs); err != nil {
			return fmt.Errorf("error writing value mask: %w", err)
		}
		writtenBytes += len(bs)
	}

	if writtenBytes != bytesToWrite {
		return fmt.Errorf("expected to write %d bytes, wrote %d bytes instead",
			bytesToWrite, writtenBytes)
	}

	return nil
}

func (pkg CapabilityPackage) String() string {
	return fmt.Sprintf("Capabilities: %#v", pkg.Capabilities)
}

var (
	// Used to parse out value masks sent by Open Server applications
	valueMaskBitMasks = []byte{
		0b00000001,
		0b00000010,
		0b00000100,
		0b00001000,
		0b00010000,
		0b00100000,
		0b01000000,
		0b10000000,
	}
)

// valueMasks are ASE versions of bitmasks and are used to communicate
// capabilities.
// A valueMask may extend over multiple bytes.
type valueMask struct {
	// map capabilities to their state
	capabilities []bool
}

func newValueMask(max int) *valueMask {
	vm := &valueMask{}

	// The max passed into newValueMask is not the capacity of the
	// valueMask but the maximum value of the capability type.  Since
	// ASE expects the index of a capability to be their numeric value
	// a capability with the value 12 will also be placed at the index
	// 12 - which would cause issues if the capacity  would be 12 as
	// well.
	vm.capabilities = make([]bool, max+1)

	return vm
}

func (vm *valueMask) isEmpty() bool {
	if len(vm.capabilities) == 1 {
		return true
	}

	for _, capa := range vm.capabilities {
		if capa {
			return false
		}
	}

	return true
}

func (vm *valueMask) setCapability(capability int, state bool) error {
	if capability > len(vm.capabilities) {
		return fmt.Errorf("invalid capability: %d", capability)
	}

	vm.capabilities[capability] = state

	return nil
}

func (vm *valueMask) getCapability(capability int) bool {
	if capability >= len(vm.capabilities) {
		return false
	}

	return vm.capabilities[capability]
}

func parseValueMask(bs []byte) *valueMask {
	max := len(bs) * 8

	vm := newValueMask(max)

	cur := 0
	// walk through bs from last to first byte
	for i := len(bs) - 1; i >= 0; i-- {
		// walk through a single byte (from least to most)
		for j := 0; j < 8; j++ {
			// check if bitmask signals the capability as true
			if bs[i]&valueMaskBitMasks[j] == valueMaskBitMasks[j] {
				// set as true
				vm.capabilities[cur] = true
			} else {
				vm.capabilities[cur] = false
			}
			cur++
		}
	}

	return vm
}

// Bytes returns a byte-slice of valueMasks.
func (vm valueMask) Bytes() []byte {
	// Calculate how many bytes the value mask requires
	max := int(math.Ceil(float64(len(vm.capabilities)) / 8))
	bs := make([]byte, max)

	cur := -1
outer:
	for i := max - 1; i >= 0; i-- {
		for j := 0; j < 8; j++ {
			cur++
			if cur >= len(vm.capabilities) {
				break outer
			}
			if !vm.capabilities[cur] {
				continue
			}
			bs[i] |= valueMaskBitMasks[j]
		}
	}

	return bs
}
