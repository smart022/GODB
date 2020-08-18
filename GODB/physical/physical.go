package physical

import (
	"bytes"
	"encoding/binary"
	"sync"
	"fmt"
	"os"
)

type size_t int64

const (
	INTEGER_LENGTH = 8 // 8 bytes for int64
	//ByteDir    = binary.LittleEndian
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Storage struct {
	_f     *os.File
	mutex	sync.Mutex
	Locked bool
}

// 原代码里的 lock unlock 不用实现，直接调mutex的

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

	cot := make([]byte, INTEGER_LENGTH)
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
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.LittleEndian, theInt)
	if err != nil {
		
		fmt.Println("binary.Write failed:", err)
		return err
	}

	//buf.WriteTo(s._f)
	_, err = s._f.Write(buff.Bytes())
	if err != nil {
		fmt.Println("file Write failed:", err)
		return err
	}

	return nil
}

func (s *Storage) Read(address int64) []byte {
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
func (s *Storage) Write(data []byte) int64 {
	// lock()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// Seek(0,1) current pos
	// 0,1,2: head,curr,tail
	pos, err := s._f.Seek(0, 1)

	check(err)
	datalen := int64(len(data))
	s.write_integer(datalen)

	s._f.Write(data)
	return pos
}

// 从头改写更新
func (s *Storage) Commit_root_address(root_address int64) {
	// lock()
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s._f.Sync()
	// seek_superblock
	s._f.Seek(0, 0)
	err := s.write_integer(root_address)
	check(err)
	s._f.Sync()
	// unlock()
}

func (s *Storage) Get_root_address() int64 {
	// seek_superblock
	s._f.Seek(0, 0)
	root_pos := s.read_integer()

	return root_pos
}

func (s *Storage) close() {
	s._f.Close()
}

func NewStorage(f *os.File) *Storage {
	s := new(Storage)
	s._f = f
	return s
}
