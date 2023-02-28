package parser

type BannedIPs struct {
	IPs []IP
}

type IP struct {
	IP          string
	type_banned string
}
