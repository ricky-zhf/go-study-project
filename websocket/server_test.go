package websocket

import (
	"testing"
)

func Test_runServer(t *testing.T) {

	runServer()
	select {}

	/*请求测试curl
	curl -i -N \
	    -H "Connection: Upgrade" \
	    -H "Upgrade: websocket" \
	    -H "Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==" \
	    -H "Sec-WebSocket-Version: 13" \
	    http://localhost:8080/ws
	*/
}
