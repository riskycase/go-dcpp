package hub

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func measure(m prometheus.Histogram) func() {
	start := time.Now()
	return func() {
		dt := time.Since(start)
		m.Observe(dt.Seconds())
	}
}

var (
	cntConnAccepted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_accepted",
		Help: "The total number of accepted connections",
	})
	cntConnError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_error",
		Help: "The total number of connections failed with an error",
	})
	cntConnOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_conn_open",
		Help: "The number of open connections",
	})

	durConnPeek = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_conn_peek",
		Help: "The time to peek protocol magic from the connection",
	})
	durConnPeekTLS = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_conn_tls_peek",
		Help: "The time to peek protocol magic from the TLS connection",
	})

	cntConnAuto = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_auto",
		Help: "The number of open connections with protocol auto-detection",
	})
	cntConnNMDC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_nmdc",
		Help: "The number of open NMDC connections",
	})
	cntConnADC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_adc",
		Help: "The number of open ADC connections",
	})
	cntConnIRC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_irc",
		Help: "The number of open IRC connections",
	})
	cntConnHTTP1 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_http1",
		Help: "The number of open HTTP1 connections",
	})
	cntConnHTTP2 = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_http2",
		Help: "The number of open HTTP2 connections",
	})

	cntConnNMDCOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_conn_nmdc_open",
		Help: "The number of open NMDC connections",
	})
	cntConnADCOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_conn_adc_open",
		Help: "The number of open ADC connections",
	})
	cntConnIRCOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_conn_irc_open",
		Help: "The number of open IRC connections",
	})
	cntConnHTTPOpen = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_conn_http_open",
		Help: "The number of open HTTP connections",
	})

	cntConnTLS = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_tls",
		Help: "The total number of accepted TLS connections",
	})
	cntConnNMDCS = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_tls_nmdc",
		Help: "The total number of secure NMDC connections",
	})
	cntConnADCS = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_tls_adc",
		Help: "The total number of secure ADC connections",
	})
	cntConnIRCS = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_tls_irc",
		Help: "The total number of secure IRC connections",
	})
	cntConnHTTPS = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_tls_http",
		Help: "The total number of secure HTTP connections",
	})

	cntConnALPN = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_alpn",
		Help: "The total number of TLS connections that support ALPN",
	})
	cntConnAlpnNMDC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_alpn_nmdc",
		Help: "The total number of accepted NMDC connections that support ALPN",
	})
	cntConnAlpnADC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_alpn_adc",
		Help: "The total number of accepted ADC connections that support ALPN",
	})
	cntConnAlpnHTTP = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_conn_alpn_http",
		Help: "The total number of accepted HTTP connections that support ALPN",
	})

	cntPings = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_pings",
		Help: "The total number of pings",
	})
	cntPingsNMDC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_pings_nmdc",
		Help: "The total number of NMDC pings",
	})
	cntPingsADC = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_pings_adc", // TODO
		Help: "The total number of ADC pings",
	})
	cntPingsHTTP = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_pings_http",
		Help: "The total number of HTTP pings",
	})

	cntChatRooms = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_chat_rooms",
		Help: "The number of active chat rooms",
	})
	cntChatMsg = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_chat_msg",
		Help: "The total number of chat messages sent",
	})
	cntChatMsgDropped = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_chat_msg_dropped",
		Help: "The total number of chat messages dropped",
	})
	cntChatMsgPM = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_chat_msg_pm",
		Help: "The total number of private messages sent",
	})

	cntSearch = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dc_search",
		Help: "The total number of search requests processed",
	})
	durSearch = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_search_dur",
		Help: "The time to send the search request",
	})

	sizeNMDCLinesR = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_nmdc_lines_read",
		Help: "The number of bytes of NMDC protocol received",
	})
	sizeNMDCLinesW = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_nmdc_lines_write",
		Help: "The number of bytes of NMDC protocol sent",
	})
	durNMDCHandshake = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_nmdc_handshake_dur",
		Help: "The time to perform NMDC handshake",
	})
	cntNMDCCommandsR = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dc_nmdc_commands_read",
		Help: "The total number of NMDC commands received",
	}, []string{"cmd"})
	cntNMDCCommandsW = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dc_nmdc_commands_write",
		Help: "The total number of NMDC commands sent",
	}, []string{"cmd"})

	sizeADCLinesR = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_adc_lines_read",
		Help: "The number of bytes of ADC protocol received",
	})
	sizeADCLinesW = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_adc_lines_write",
		Help: "The number of bytes of ADC protocol sent",
	})
	durADCHandshake = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "dc_adc_handshake_dur",
		Help: "The time to perform ADC handshake",
	})
	cntADCPackets = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dc_adc_packets",
		Help: "The total number of ADC packets",
	}, []string{"kind"})

	cntPeers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "dc_peers",
		Help: "The number of active peers",
	})
)
