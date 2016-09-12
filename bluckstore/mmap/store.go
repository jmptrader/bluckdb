package memap

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"fmt"
	"bytes"
	"encoding/gob"
	"io/ioutil"
)

type MmapKVStore struct {
	Dir *Directory
}

const FILE_NAME = "bluck.data"
const META_FILE_NAME = "bluck.meta"
const DB_DIRECTORY = "/tmp/"

func (store *MmapKVStore) Open() {

	f, err := os.OpenFile(DB_DIRECTORY + FILE_NAME, os.O_RDWR, 0644)
	defer f.Close()

	if err != nil {
		f, _ = os.Create(DB_DIRECTORY + FILE_NAME)
		fmt.Println("OpenFile DB error : " + err.Error())
		store.Dir = &Directory{
			Gd:    0,
			Table: make([]int, 1),
		}
		f.Write(make([]byte, 4096))
	} else {
		meta, err := ioutil.ReadFile(DB_DIRECTORY + META_FILE_NAME)

		if err != nil {
			fmt.Println("OpenFile Metadata error : " + err.Error())
		} else {
			store.Dir = DecodeMeta(bytes.NewBuffer(meta))
		}
	}

	store.Dir.dataFile = f
	store.Dir.data, _ = mmap.Map(store.Dir.dataFile, mmap.RDWR, 0644)
}

func (s *MmapKVStore) Get(k string) string {
	return s.Dir.get(k)
}

func (s *MmapKVStore) Put(k, v string) {
	s.Dir.put(k, v)
}

func (s *MmapKVStore) Close() {
	s.Dir.data.Unmap()
	s.Dir.dataFile.Close()
	mFile, err := os.OpenFile(DB_DIRECTORY + META_FILE_NAME, os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0644)
	defer mFile.Close()

	if err != nil {
		fmt.Println(err)
	}
	mFile.WriteAt(EncodeMeta(s.Dir).Bytes(), 0)

}

func (s *MmapKVStore) Rm() {
	os.Remove(DB_DIRECTORY + FILE_NAME)
	os.Remove(DB_DIRECTORY + META_FILE_NAME)
}

func (s *MmapKVStore) RestoreMETA() {
	/*

		stat, _ := f.Stat()
		numBuckets := FindBucketNumber(stat.Size())
		tableSize := NextPowerOfTwo(uint(numBuckets))

		store.Dir = &Directory{
			Gd: FindTwoToPowerOfN(uint(numBuckets)),
			Table: make([]int, tableSize),
			LastPageId: int(numBuckets) - 1,
		}

		for i := 0; i < int(tableSize); i ++ {
			if i >= int(numBuckets) {
				store.Dir.Table[i] = 0
			} else {
				store.Dir.Table[i] = i
			}
		}
	 */
}

func FindBucketNumber(fileSize int64) int64 {
	return fileSize / int64(4096)
}

func FindTwoToPowerOfN(v uint) uint {
	for i :=  uint(1); ;i++ {
		if (v >> i) <= 0 {
			return i
			break
		}
	}
	return 0 // never
}

func NextPowerOfTwo(v uint) uint {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func DecodeMeta(buff *bytes.Buffer) *Directory {
	dec := gob.NewDecoder(buff)
	var dir Directory
	dec.Decode(&dir)
	return &dir
}

func EncodeMeta(dir *Directory) *bytes.Buffer{
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	enc.Encode(&dir)
	return &buff
}