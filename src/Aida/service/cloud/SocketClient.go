package cloud

//import (
//	"net"
//	"common"
//	"strconv"
//	"encoding/json"
//	"utils"
//	"time"
//	"fmt"
//	"os"
//)

//import (
//	"./utils"
//	"common"
//	"encoding/json"
//	"fmt"
//	"net"
//	"os"
//	"strconv"
//	"time"
//)

//type Msg struct {
//	Meta   map[string]interface{} `json:"meta"`
//	Content interface{}            `json:"content"`
//}
//
//func send(conn net.Conn) {
//	for i := 0; i < 1; i++ {
//		session := GetSession()
//		message := &common.Msg{
//			Meta: map[string]interface{}{
//				"meta": "test",
//				"ID":   strconv.Itoa(i),
//			},
//			Content: common.Msg{
//				Meta: map[string]interface{}{
//					"author": "nucky lu",
//				},
//				Content: session,
//			},
//		}
//		result, _ := json.Marshal(message)
//		conn.Write(utils.Enpack((result)))
//		//conn.Write([]byte(message))
//		time.Sleep(1 * time.Second)
//	}
//	fmt.Println("send over")
//	defer conn.Close()
//}
//
//func GetSession() string {
//	gs1 := time.Now().Unix()
//	gs2 := strconv.FormatInt(gs1, 10)
//	return gs2
//}
//
//func main() {
//	server := "localhost:1024"
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
//		os.Exit(1)
//	}
//
//	conn, err := net.DialTCP("tcp", nil, tcpAddr)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
//		os.Exit(1)
//	}
//
//	fmt.Println("connect success")
//	send(conn)
//}
