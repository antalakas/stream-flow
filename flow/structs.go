package flow

//{
//	"IN_BYTES":7443,
//	"IN_PKTS":7,
//	"PROTOCOL":0,
//	"TCP_FLAGS":0,
//	"L4_SRC_PORT":80,
//	"IPV4_SRC_ADDR":"98.139.134.187",
//	"L4_DST_PORT":38670,
//	"IPV4_DST_ADDR":"172.16.133.132",
//	"LAST_SWITCHED":1530695222,
//	"FIRST_SWITCHED":1530695221,
//	"ICMP_TYPE":0,
//	"IP_PROTOCOL_VERSION":4,
//	"FLOW_ID":952,
//	"42":1646
//}
type NetflowEvent struct {
IPV4SrcAddr        		string	`json:"IPV4_SRC_ADDR"`
	L4SrcPort         	int 	`json:"L4_SRC_PORT"`
	IPV4DstAddr         string 	`json:"IPV4_DST_ADDR"`
	L4DstPort       	int 	`json:"L4_DST_PORT"`
	IpProtocolVersion   int 	`json:"IP_PROTOCOL_VERSION"`
	Protocol            int 	`json:"PROTOCOL"`
	TcpFlags     		int 	`json:"TCP_FLAGS"`
	IcmpType            int 	`json:"ICMP_TYPE"`
	FirstSwitched       int64 	`json:"FIRST_SWITCHED"`
	LastSwitched        int64 	`json:"LAST_SWITCHED"`
	InBytes             int 	`json:"IN_BYTES"`
	InPkts        		int 	`json:"IN_PKTS"`
	FlowId             	int 	`json:"FLOW_ID"`
	FortyTwo       	 	int 	`json:"42"`
}