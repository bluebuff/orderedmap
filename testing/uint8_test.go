package orderedmap_test

import (
	"github.com/shengmingzhu/orderedmap"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"sort"
	"testing"
	"time"
)

const (
	testCountUint8 uint8 = 1 << 4
	rangeLenUint8        = 1 << 7
)

func TestNewUint8(t *testing.T) {
	Convey("NewUint8 and Put and Len", t, func() {
		m := orderedmap.NewUint8()
		hm := make(map[uint8]struct{}, testCountUint8)
		sl := make([]uint8, 0, testCountUint8)
		rand.Seed(time.Now().UnixNano())
		for i := uint8(0); i < testCountUint8; {
			key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
			if _, ok := hm[key]; !ok {
				hm[key] = struct{}{}
				m.Put(key, key<<1)
				sl = append(sl, key)
				i++
			}
		}
		So(m.Len(), ShouldEqual, len(hm))
		sort.Slice(sl, func(i, j int) bool {
			return sl[i] < sl[j]
		})

		Convey("Get", func() {
			for k := range hm {
				key, ok := m.Get(k)
				So(ok, ShouldEqual, true)
				So(key, ShouldEqual, k<<1)
			}
		})

		Convey("Min", func() {
			key, _ := m.Min()
			So(key, ShouldEqual, sl[0])
		})

		Convey("Max", func() {
			key, _ := m.Max()
			So(key, ShouldEqual, sl[testCountUint8-1])
		})

		Convey("Keys", func() {
			keys := m.Keys()
			So(len(keys), ShouldEqual, len(sl))
			for i := range keys {
				So(keys[i], ShouldEqual, sl[i])
			}
		})

		Convey("Values", func() {
			values := m.Values()
			So(len(values), ShouldEqual, len(sl))
			for i := range values {
				So(values[i], ShouldEqual, sl[i]<<1)
			}
		})

		Convey("RangeAll", func() {
			pairs := m.RangeAll()
			So(len(pairs), ShouldEqual, len(sl))
			for i := range pairs {
				So(pairs[i].Key, ShouldEqual, sl[i])
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("RangeAllDesc", func() {
			pairs := m.RangeAllDesc()
			So(len(pairs), ShouldEqual, len(sl))
			for i := range pairs {
				So(pairs[i].Key, ShouldEqual, sl[int(testCountUint8)-i-1])
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("Range", func() {
			rand.Seed(time.Now().UnixNano())
			iKey1 := uint8(rand.Int63n(int64(testCountUint8)>>1)) & uint8((int64(testCountUint8)>>1)-1)
			iKey2 := uint8(rand.Int63n(int64(testCountUint8)>>1)) & uint8((int64(testCountUint8)>>1)-1)
			if iKey2 < 1 {
				iKey2 = 1
			}
			pairs := m.Range(sl[iKey1], sl[iKey1+iKey2])
			So(len(pairs), ShouldEqual, iKey2+1)
			for i := range pairs {
				So(pairs[i].Key, ShouldEqual, sl[int(iKey1)+i])
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("RangeDesc", func() {
			rand.Seed(time.Now().UnixNano())
			iKey1 := uint8(rand.Int63n(int64(testCountUint8)>>1)) & uint8((int64(testCountUint8)>>1)-1)
			iKey2 := uint8(rand.Int63n(int64(testCountUint8)>>1)) & uint8((int64(testCountUint8)>>1)-1)
			if iKey2 < 1 {
				iKey2 = 1
			}
			pairs := m.RangeDesc(sl[iKey1], sl[iKey1+iKey2])
			So(len(pairs), ShouldEqual, iKey2+1)
			for i := range pairs {
				So(pairs[i].Key, ShouldEqual, sl[int(iKey1+iKey2)-i])
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("RangeN", func() {
			rand.Seed(time.Now().UnixNano())
			iKey1 := 4
			iKey2 := 3
			pairs := m.RangeN(int(iKey2), sl[iKey1]-1)
			So(len(pairs), ShouldEqual, iKey2)
			for i := range pairs {
				if sl[iKey1]-1 == sl[iKey1-1] {
					So(pairs[i].Key, ShouldEqual, sl[int(iKey1)+i-1])
				} else {
					So(pairs[i].Key, ShouldEqual, sl[int(iKey1)+i])
				}
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("RangeDescN", func() {
			rand.Seed(time.Now().UnixNano())
			iKey1 := 4
			iKey2 := 3
			pairs := m.RangeDescN(int(iKey2), sl[iKey1]+1)
			So(len(pairs), ShouldEqual, iKey2)
			for i := range pairs {
				if sl[iKey1]+1 == sl[iKey1+1] {
					So(pairs[i].Key, ShouldEqual, sl[int(iKey1)-i+1])
				} else {
					So(pairs[i].Key, ShouldEqual, sl[int(iKey1)-i])
				}

				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("String", func() {
			str := m.String()
			//fmt.Println()
			//fmt.Println(str)
			So(len(str), ShouldBeGreaterThan, 0)
		})

		Convey("PopMin", func() {
			for i := 0; i < 4; i++ {
				k, v := m.PopMin()
				So(k, ShouldEqual, sl[i])
				So(v, ShouldEqual, k<<1)
				delete(hm, k)
			}
		})

		Convey("PopMax", func() {
			for i := 0; i < 4; i++ {
				k, v := m.PopMax()
				So(k, ShouldEqual, sl[len(sl)-1-i])
				So(v, ShouldEqual, k<<1)
				delete(hm, k)
			}
		})

		Convey("Delete", func() {
			count := len(hm)
			for key := range hm {
				//fmt.Println()
				//fmt.Println(m)
				m.Delete(key)
				count--
				So(m.Len(), ShouldEqual, count)
			}
		})

		Convey("EmptyMap", func() {
			m := orderedmap.NewUint8()
			_, ok := m.Get(1)
			So(ok, ShouldEqual, false)
			m.Delete(2)
			k, v := m.Min()
			So(k, ShouldEqual, 0)
			So(v, ShouldEqual, nil)
			k, v = m.Max()
			So(k, ShouldEqual, 0)
			So(v, ShouldEqual, nil)
			k, v = m.PopMin()
			So(k, ShouldEqual, 0)
			So(v, ShouldEqual, nil)
			k, v = m.PopMax()
			So(k, ShouldEqual, 0)
			So(v, ShouldEqual, nil)
			res := m.RangeAll()
			So(len(res), ShouldEqual, 0)
			res = m.RangeAllDesc()
			So(len(res), ShouldEqual, 0)
			res = m.Range(1, 10)
			So(len(res), ShouldEqual, 0)
			res = m.RangeDesc(1, 10)
			So(len(res), ShouldEqual, 0)
		})
	})
}

func BenchmarkUint8_Put(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewUint8()
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Put(sl[i], struct{}{})
	}
}

func BenchmarkHashMapUint8_Put(b *testing.B) {
	b.StopTimer()
	m := make(map[uint8]struct{})
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[sl[i]] = struct{}{}
	}
}

func BenchmarkUint8_Get(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewUint8()
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Get(sl[i])
	}
}

func BenchmarkHashMapUint8_Get(b *testing.B) {
	b.StopTimer()
	m := make(map[uint8]struct{})
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m[sl[i]]
	}
}

func BenchmarkUint8_Delete(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewUint8()
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(sl[i])
	}
}

func BenchmarkHashMapUint8_Delete(b *testing.B) {
	b.StopTimer()
	m := make(map[uint8]struct{})
	sl := make([]uint8, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		delete(m, sl[i])
	}
}

func BenchmarkUint8_RangeAll(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewUint8()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenUint8; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m.RangeAll()
	}
}

func BenchmarkHashMapUint8_RangeAllNoSort(b *testing.B) {
	b.StopTimer()
	m := make(map[uint8]struct{})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenUint8; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		arr := make([]*orderedmap.Uint8KeyValue, 0)
		for key, value := range m {
			pair := &orderedmap.Uint8KeyValue{Key: key, Value: value}
			arr = append(arr, pair)
		}
	}
}

func BenchmarkHashMapUint8_RangeAllAndSort(b *testing.B) {
	b.StopTimer()
	m := make(map[uint8]struct{})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenUint8; i++ {
		key := uint8(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		arr := make([]*orderedmap.Uint8KeyValue, 0)
		for key, value := range m {
			pair := &orderedmap.Uint8KeyValue{Key: key, Value: value}
			arr = append(arr, pair)
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].Key < arr[j].Key
		})
	}
}
