package parser

type BannedIPs struct {
	IPs []IP
}

type IP struct {
	IP          string
	Type_banned string
}
