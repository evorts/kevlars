/**
 * @Author: steven
 * @Description:
 * @File: cache_test
 * @Date: 29/09/23 10.46
 */

package inmemory

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"time"

	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	rm  Manager
	rme Manager
)

func assertPanic(t *testing.T, f func(ctx context.Context) Manager) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("code not panic")
		}
	}()
	f(context.Background())
}

func assertNotPanic(t *testing.T, f func(ctx context.Context) Manager) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("code panic")
		}
	}()
	f(context.Background())
}

func TestMain(m *testing.M) {
	var (
		mr  *miniredis.Miniredis
		err error
	)
	if mr, err = miniredis.Run(); err != nil {
		log.Fatalf("error '%s' wasn't expected when initiated minimal redis for test", err.Error())
	}
	rm = NewRedis(mr.Addr())
	rme = NewRedis("invalid-address")
	code := m.Run()
	os.Exit(code)
}

func TestManagerConnect(t *testing.T) {
	// emulate when using invalid configuration
	err := rme.Connect(context.Background())
	assert.Error(t, err)
	// emulate when using proper configuration
	err = rm.Connect(context.Background())
	assert.Nil(t, err)
}

func TestManagerMustConnect(t *testing.T) {
	assertPanic(t, rme.MustConnect)
	assertNotPanic(t, rm.MustConnect)
}

func TestManagerSetString(t *testing.T) {
	assert.Nil(t, rm.SetString(context.Background(), "some_key", "some_value", 1*time.Second))
	assert.Eventually(t, func() bool {
		var test = rm.GetString(context.Background(), "some_key")
		fmt.Print(test)
		return rm.GetString(context.Background(), "some_key") == "some_value"
	}, 2*time.Second, 1*time.Second, "testing")
}

func TestManagerGetString(t *testing.T) {
	assert.Empty(t, rm.GetString(context.Background(), "key_not_exist"))
}
