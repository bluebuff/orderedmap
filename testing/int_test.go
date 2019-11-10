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
	testCountInt int = 1 << 4
	rangeLenInt      = 1 << 10
)

func TestNewInt(t *testing.T) {
	Convey("NewInt and Put and Len", t, func() {
		m := orderedmap.NewInt()
		hm := make(map[int]struct{}, testCountInt)
		sl := make([]int, 0, testCountInt)
		rand.Seed(time.Now().UnixNano())
		for i := int(0); i < testCountInt; {
			key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
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
				So(m.Get(k), ShouldEqual, k<<1)
			}
		})

		Convey("Min", func() {
			key, _ := m.Min()
			So(key, ShouldEqual, sl[0])
		})

		Convey("Max", func() {
			key, _ := m.Max()
			So(key, ShouldEqual, sl[testCountInt-1])
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
				So(pairs[i].Key, ShouldEqual, sl[int(testCountInt)-i-1])
				So(pairs[i].Value, ShouldEqual, pairs[i].Key<<1)
			}
		})

		Convey("Range", func() {
			rand.Seed(time.Now().UnixNano())
			iKey1 := int(rand.Int63n(int64(testCountInt)>>1)) & int((int64(testCountInt)>>1)-1)
			iKey2 := int(rand.Int63n(int64(testCountInt)>>1)) & int((int64(testCountInt)>>1)-1)
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
			iKey1 := int(rand.Int63n(int64(testCountInt)>>1)) & int((int64(testCountInt)>>1)-1)
			iKey2 := int(rand.Int63n(int64(testCountInt)>>1)) & int((int64(testCountInt)>>1)-1)
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
			m := orderedmap.NewInt()
			v := m.Get(1)
			So(v, ShouldEqual, nil)
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

func BenchmarkInt_Put(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewInt()
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Put(sl[i], struct{}{})
	}
}

func BenchmarkHashMapInt_Put(b *testing.B) {
	b.StopTimer()
	m := make(map[int]struct{})
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[sl[i]] = struct{}{}
	}
}

func BenchmarkInt_Get(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewInt()
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Get(sl[i])
	}
}

func BenchmarkHashMapInt_Get(b *testing.B) {
	b.StopTimer()
	m := make(map[int]struct{})
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m[sl[i]]
	}
}

func BenchmarkInt_Delete(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewInt()
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Delete(sl[i])
	}
}

func BenchmarkHashMapInt_Delete(b *testing.B) {
	b.StopTimer()
	m := make(map[int]struct{})
	sl := make([]int, 0, b.N)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		sl = append(sl, key)
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		delete(m, sl[i])
	}
}

func BenchmarkInt_RangeAll(b *testing.B) {
	b.StopTimer()
	m := orderedmap.NewInt()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenInt; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m.Put(key, struct{}{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = m.RangeAll()
	}
}

func BenchmarkHashMapInt_RangeAllNoSort(b *testing.B) {
	b.StopTimer()
	m := make(map[int]struct{})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenInt; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		arr := make([]*orderedmap.IntKeyValue, 0)
		for key, value := range m {
			pair := &orderedmap.IntKeyValue{Key: key, Value: value}
			arr = append(arr, pair)
		}
	}
}

func BenchmarkHashMapInt_RangeAllAndSort(b *testing.B) {
	b.StopTimer()
	m := make(map[int]struct{})
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rangeLenInt; i++ {
		key := int(uint32(rand.Int31n(0x40000000)) & (0x40000000 - 1))
		m[key] = struct{}{}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		arr := make([]*orderedmap.IntKeyValue, 0)
		for key, value := range m {
			pair := &orderedmap.IntKeyValue{Key: key, Value: value}
			arr = append(arr, pair)
		}
		sort.Slice(arr, func(i, j int) bool {
			return arr[i].Key < arr[j].Key
		})
	}
}