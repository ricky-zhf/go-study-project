package bloom_filter

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

}

func getData(id int64) string {
	fmt.Println("query...")
	time.Sleep(5 * time.Second) // 模拟一个比较耗时的操作
	return "liwenzhou.com"
}

func TestSingleFlight(t *testing.T) {
	g := new(singleflight.Group)

	// 第1次调用
	go func() {
		v1, _, shared := g.Do("getData", func() (interface{}, error) {
			ret := getData(1)
			return ret, nil
		})
		fmt.Printf("1st call: v1:%v, shared:%v\n", v1, shared)
	}()

	time.Sleep(2 * time.Second)

	// 第2次调用（第1次调用已开始但未结束）
	v2, _, shared := g.Do("getData", func() (interface{}, error) {
		ret := getData(1)
		return ret, nil
	})
	fmt.Printf("2nd call: v2:%v, shared:%v\n", v2, shared)
}

func doChanGetData(ctx context.Context, g *singleflight.Group, id int64) (string, error) {
	ch := g.DoChan("getData", func() (interface{}, error) {
		ret := getData(id)
		return ret, nil
	})
	select {
	case <-ctx.Done():
		fmt.Println("end...")
		return "", ctx.Err()
	case ret := <-ch:
		return ret.Val.(string), ret.Err
	}
}

func TestSingleFlightChannel(t *testing.T) {
	g := new(singleflight.Group)

	// 第1次调用
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		v1, err := doChanGetData(ctx, g, 1)
		fmt.Printf("v1:%v err:%v\n", v1, err)
	}()

	time.Sleep(2 * time.Second)

	// 第2次调用（第1次调用已开始但未结束）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	v2, err := doChanGetData(ctx, g, 1)
	fmt.Printf("v2:%v err:%v\n", v2, err)
}
