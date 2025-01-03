package rabbitmq

import (
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

var queue_name = os.Getenv("QUEUE_NAME")
var exchange_name = os.Getenv("EXCHANGE_NAME")
var rabbit_mq = os.Getenv("RABBIT_MQ")
var rabbit_username = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_vhost = os.Getenv("RABBIT_VHOST")
var rabbit_port = os.Getenv("RABBIT_PORT")

func Init() {
	fmt.Println("rabbitmq init")
	exchange()
	queue()
	bind()
}

func bind() {
	url_bind := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/bindings/" + rabbit_vhost + "/e/" + exchange_name +
		"/q/" + queue_name

	resp, err := FastGet(url_bind)

	if err == nil {
		if string(resp.Body()) == "[]" {
			CreateBind()
			return
		}
	}
}

func CreateBind() {
	url_bind := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/bindings/" + rabbit_vhost + "/e/" + exchange_name +
		"/q/" + queue_name

	jsonObj := gabs.New()

	resp, err := FastPost(url_bind, []byte(jsonObj.String()))

	if err == nil {
		fmt.Println(string(resp.Body()))
	}
}

func queue() {
	url_queue := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/queues/" + rabbit_vhost + "/" + queue_name

	resp, err := FastGet(url_queue)

	if err == nil {
		ok := false
		jsonParsed, _ := gabs.ParseJSON(resp.Body())

		error, ok := jsonParsed.Path("error").Data().(string)

		if ok {
			fmt.Println(error)
			CreateQueue()
			return
		}

		auto_delete, _ := jsonParsed.Path("auto_delete").Data().(bool)
		durable, _ := jsonParsed.Path("durable").Data().(bool)

		if auto_delete != false || durable != true {
			DeleteQueue()
			CreateQueue()
			return
		}
	}
}

func exchange() {
	url_get_exchange := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/exchanges/" + rabbit_vhost + "/" + exchange_name

	fmt.Println("url get exchange")
	fmt.Println(url_get_exchange)

	resp, err := FastGet(url_get_exchange)

	if err == nil {
		ok := false
		jsonParsed, _ := gabs.ParseJSON(resp.Body())

		error, ok := jsonParsed.Path("error").Data().(string)

		if ok {
			fmt.Println(error)
			CreateExchange()
			return
		}

		exchange_type, _ := jsonParsed.Path("type").Data().(string)
		auto_delete, _ := jsonParsed.Path("auto_delete").Data().(bool)
		durable, _ := jsonParsed.Path("durable").Data().(bool)

		if exchange_type != "fanout" || auto_delete != false || durable != true {
			DeleteExchange()
			CreateExchange()
			return
		}
	}
}

func CreateQueue() {
	url_queue := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/queues/" + rabbit_vhost + "/" + queue_name

	jsonObj := gabs.New()
	jsonObj.Set(false, "auto_delete")
	jsonObj.Set(true, "durable")

	resp, err := FastPut(url_queue, []byte(jsonObj.String()))

	if err == nil {
		fmt.Println(string(resp.Body()))
	}
}

func CreateExchange() {
	url_exchange := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/exchanges/" + rabbit_vhost + "/" + exchange_name

	jsonObj := gabs.New()
	jsonObj.Set("fanout", "type")
	jsonObj.Set(false, "auto_delete")
	jsonObj.Set(true, "durable")

	resp, err := FastPut(url_exchange, []byte(jsonObj.String()))

	if err == nil {
		fmt.Println(string(resp.Body()))
	}
}

func DeleteQueue() {
	url_queue := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/queues/" + rabbit_vhost + "/" + queue_name + "queue"

	resp, err := FastDelete(url_queue)

	if err == nil {
		fmt.Println(resp)
	}
}

func DeleteExchange() {
	url_exchange := "https://" +
		rabbit_username + ":" + rabbit_password + "@" +
		rabbit_host + "/api/exchanges/" + rabbit_vhost + "/" + exchange_name

	resp, err := FastDelete(url_exchange)

	if err == nil {
		fmt.Println(resp)
	}
}

func FastGet(url string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func FastPut(url string, data []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("PUT")
	req.SetBody(data)

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func FastPost(url string, data []byte) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetContentType("application/json")
	req.Header.SetMethod("POST")
	req.SetBody(data)

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

func FastDelete(url string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.SetRequestURI(url)
	req.Header.SetMethod("DELETE")

	timeOut := 3 * time.Second
	var err = fasthttp.DoTimeout(req, resp, timeOut)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	out := fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}