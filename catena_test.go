package catena_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codemodus/catena"
)

var (
	bTxt0   = "0"
	bTxt1   = "1"
	bTxtA   = "A"
	bTxtEnd = "_END_"
)

func Example() {
	// Each wrapper writes either "0", "1", or "A" to the response body before
	// and after ServeHTTP() is called.
	// endHandler writes "_END_" to the response body and returns.
	catena00 := catena.New(handlerWrapper0, handlerWrapper0)
	catena00A1 := catena00.Append(handlerWrapperA, handlerWrapper1)

	catena100A1 := catena.New(handlerWrapper1)
	catena100A1 = catena100A1.Merge(catena00A1)

	mux := http.NewServeMux()
	mux.Handle("/path_implies_body/00_End", catena00.EndFn(endHandler))
	mux.Handle("/path_implies_body/00A1_End", catena00A1.EndFn(endHandler))
	mux.Handle("/path_implies_body/100A1_End", catena100A1.EndFn(endHandler))

	server := httptest.NewServer(mux)

	rBody0, err := getReqBody(server.URL + "/path_implies_body/00_End")
	if err != nil {
		fmt.Println(err)
	}

	rBody1, err := getReqBody(server.URL + "/path_implies_body/00A1_End")
	if err != nil {
		fmt.Println(err)
	}

	rBody2, err := getReqBody(server.URL + "/path_implies_body/100A1_End")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Catena 0 Body:", rBody0)
	fmt.Println("Catena 1 Body:", rBody1)
	fmt.Println("Catena 2 Body:", rBody2)

	// Output:
	// Catena 0 Body: 00_END_00
	// Catena 1 Body: 00A1_END_1A00
	// Catena 2 Body: 100A1_END_1A001
}

func TestCatena(t *testing.T) {
	c0 := catena.New(handlerWrapper0)
	c1 := c0.Append(handlerWrapper1, handlerWrapperA)
	cBefore0 := catena.New(handlerWrapper1)
	c0 = cBefore0.Merge(c0)
	m := http.NewServeMux()
	r0 := "/0"
	r1 := "/1"
	m.Handle(r0, c0.EndFn(endHandler))
	m.Handle(r1, c1.EndFn(endHandler))
	s := httptest.NewServer(m)

	tMap := map[string]string{
		"/0": bTxt1 + bTxt0 + bTxtEnd + bTxt0 + bTxt1,
		"/1": bTxt0 + bTxt1 + bTxtA + bTxtEnd + bTxtA + bTxt1 + bTxt0,
	}

	for k, v := range tMap {
		rb, err := getReqBody(s.URL + k)
		if err != nil {
			t.Error(err)
		}
		want := v
		got := rb
		if got != want {
			t.Errorf("Body = %v, want %v", got, want)
		}
	}
}

func TestNilEnd(t *testing.T) {
	c0 := catena.New(emptyHandlerWrapper)
	m := http.NewServeMux()
	r0 := "/0"
	r1 := "/1"
	m.Handle(r0, c0.End(nil))
	m.Handle(r1, c0.EndFn(nil))
	s := httptest.NewServer(m)

	tMap := map[string]int{
		"/0": http.StatusOK,
		"/1": http.StatusOK,
	}

	for k, v := range tMap {
		rs, err := getReqStatus(s.URL + k)
		if err != nil {
			t.Error(err)
		}
		want := v
		got := rs
		if got != want {
			t.Errorf("Status Code = %v, want %v", got, want)
		}
	}
}

func getReqBody(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()
	return string(body), nil
}

func getReqStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp.StatusCode, nil
}

func handlerWrapper0(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(bTxt0))
		n.ServeHTTP(w, r)
		_, _ = w.Write([]byte(bTxt0))
	})
}

func handlerWrapper1(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(bTxt1))
		n.ServeHTTP(w, r)
		_, _ = w.Write([]byte(bTxt1))
	})
}

func handlerWrapperA(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(bTxtA))
		n.ServeHTTP(w, r)
		_, _ = w.Write([]byte(bTxtA))
	})
}

func emptyHandlerWrapper(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n.ServeHTTP(w, r)
	})
}

func endHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(bTxtEnd))
	return
}

func nilHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func BenchmarkCatena10(b *testing.B) {
	c0 := catena.New(emptyHandlerWrapper,
		emptyHandlerWrapper, emptyHandlerWrapper, emptyHandlerWrapper,
		emptyHandlerWrapper, emptyHandlerWrapper, emptyHandlerWrapper,
		emptyHandlerWrapper, emptyHandlerWrapper, emptyHandlerWrapper)
	m := http.NewServeMux()
	m.Handle("/", c0.EndFn(nilHandler))
	s := httptest.NewServer(m)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		re0, err := http.Get(s.URL + "/")
		if err != nil {
			b.Error(err)
		}
		_ = re0.Body.Close()
	}
}

func BenchmarkNest10(b *testing.B) {
	h := emptyHandlerWrapper(emptyHandlerWrapper(
		emptyHandlerWrapper(emptyHandlerWrapper(
			emptyHandlerWrapper(emptyHandlerWrapper(
				emptyHandlerWrapper(emptyHandlerWrapper(
					emptyHandlerWrapper(emptyHandlerWrapper(
						http.HandlerFunc(nilHandler)))))))))))
	m := http.NewServeMux()
	m.Handle("/", h)
	s := httptest.NewServer(m)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		re0, err := http.Get(s.URL + "/")
		if err != nil {
			b.Error(err)
		}
		_ = re0.Body.Close()
	}
}
