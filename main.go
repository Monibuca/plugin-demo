package demo

import (
	"net/http"

	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
	"m7s.live/engine/v4/track"
)

type DemoConfig struct {
	config.HTTP
	config.Publish
	config.Subscribe
	Foo string `default:"bar"`
}

var demoConfig DemoConfig

var DemoPlugin = InstallPlugin(&demoConfig)

func (conf *DemoConfig) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		// 插件启动
		break
	}
}

// http://localhost:8080/demo/api/test/pub
func (conf *DemoConfig) API_test_pub(rw http.ResponseWriter, r *http.Request) {
	var pub DemoPublisher
	err := DemoPlugin.Publish("demo/test", &pub)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	} else {
		vt := track.NewH264(pub.Stream)
		// 根据实际情况写入视频帧，需要注意pts和dts需要写入正确的值 即毫秒数*90
		vt.WriteAnnexB(0, 0, []byte{0, 0, 0, 1})
	}
	rw.Write([]byte("test_pub"))
}

// http://localhost:8080/demo/api/test/sub
func (conf *DemoConfig) API_test_sub(rw http.ResponseWriter, r *http.Request) {
	var sub DemoSubscriber
	err := DemoPlugin.Subscribe("demo/test", &sub)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	} else {
		sub.PlayRaw()
	}
	rw.Write([]byte("test_sub"))
}

type DemoPublisher struct {
	Publisher
}

func (pub *DemoPublisher) OnEvent(event any) {
	switch event.(type) {
	case IPublisher:
		// 发布成功
	default:
		pub.Publisher.OnEvent(event)
	}
}

type DemoSubscriber struct {
	Subscriber
}

func (sub *DemoSubscriber) OnEvent(event any) {
	switch event.(type) {
	case ISubscriber:
		// 订阅成功
	case AudioFrame:
		// 音频帧处理
	case VideoFrame:
		// 视频帧处理
	default:
		sub.Subscriber.OnEvent(event)
	}
}
