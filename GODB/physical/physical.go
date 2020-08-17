package physical

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type size_t int64

const (
	IntegerLen = 8 // 8 bytes for int64
	//ByteDir    = binary.LittleEndian
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Storage struct {
	_f     *os.File
	locked bool
}

func (s *Storage) assert_not_closed() bool {
	if s._f == nil {
		return false
	}
	return true
}

func (s *Storage) _seek_end() int64 {
	// Seek(offset, relative_ori_pos: 0 top 1 any 2 rear)
	pos, ok := s._f.Seek(0, 2)
	if ok != nil {
		fmt.Println("seek err!")
		return 0
	}
	return pos
}

// cur pos read
func (s *Storage) read_integer() int64 {
	// test
	//s._f.Seek(0, 0)

	cot := make([]byte, IntegerLen)
	_, err := s._f.Read(cot)
	check(err)
	buf := bytes.NewReader(cot)
	var ret int64
	err = binary.Read(buf, binary.LittleEndian, &ret)
	check(err)
	return ret
}

// cur pos write
func (s *Storage) write_integer(theInt int64) error {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, theInt)
	if err != nil {
		return err
	}

	//buf.WriteTo(s._f)
	_, err = s._f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) read(address int64) []byte {
	s._f.Seek(address, 0)
	datalen := s.read_integer() // read_integer 本身也会移动

	fmt.Println("reading: ", datalen)
	cot := make([]byte, datalen)
	_, err := s._f.Read(cot)
	check(err)
	// return bytes
	return cot
}

// data formar:
// len of data bytes | data bytes
//
func (s *Storage) write(data []byte) int64 {
	// lock()
	pos, err := s._f.Seek(0, 2)
	check(err)
	datalen := int64(len(data))
	s.write_integer(datalen)

	s._f.Write(data)
	return pos
}

func (s *Storage) commit_root_address(root_address int64) {
	// lock()
	s._f.Sync()
	// seek_superblock
	s._f.Seek(0, 0)
	err := s.write_integer(root_address)
	check(err)
	s._f.Sync()
	// unlock()
}

func (s *Storage) get_root_address() int64 {
	// seek_superblock
	s._f.Seek(0, 0)
	pos := s.read_integer()

	return pos
}

func (s *Storage) close() {
	// unlock()
	s._f.Close()
}

func NewStorage(f *os.File) *Storage {
	s := new(Storage)
	s._f = f
	s.locked = false
	return s
}
