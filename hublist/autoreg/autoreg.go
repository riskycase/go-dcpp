package autoreg

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/direct-connect/go-dcpp/nmdc"
	"github.com/direct-connect/go-dcpp/version"
)

const DefaultPort = 2501

type Info struct {
	Name  string
	Host  string
	Desc  string
	Users int
	Share uint64
}

const timeout = time.Second * 5

var (
	escaper = strings.NewReplacer(
		"|", "&#124;",
	)
)

// Register the hub on the hublist.
func Register(ctx context.Context, addr string, info Info) error {
	u, err := url.Parse(addr)
	if err != nil || u.Host == "" {
		u = &url.URL{Host: addr}
	}
	if _, port, _ := net.SplitHostPort(u.Host); port == "" {
		u.Host += ":" + strconv.Itoa(DefaultPort)
	}

	c, err := net.Dial("tcp", u.Host)
	if err != nil {
		return err
	}
	defer c.Close()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(timeout)
	}
	if err = c.SetDeadline(deadline); err != nil {
		return err
	}
	return registerOn(c, info, u.Host)
}

func registerOn(c io.ReadWriteCloser, info Info, host string) error {
	b := make([]byte, 512)
	i, _, err := scanOneTo(c, b)
	if err == io.ErrShortBuffer {
		return errors.New("$Lock command is too long")
	} else if err != nil {
		return err
	} else if i < 0 {
		return errors.New("expected NMDC $Lock")
	}
	b = b[:i]
	if !bytes.HasPrefix(b, []byte("$Lock ")) {
		return fmt.Errorf("expected $Lock, got: %q", string(b))
	}
	var lock nmdc.Lock
	if err := lock.UnmarshalNMDC(b[6:]); err != nil {
		return err
	}
	key := lock.Key()

	buf := bytes.NewBuffer(b[:0])
	buf.Reset()
	buf.WriteString("$Key ")
	data, err := key.MarshalNMDC()
	if err != nil {
		return err
	}
	buf.Write(data)
	buf.WriteByte('|')

	buf.WriteString(escaper.Replace(info.Name))
	buf.WriteByte('|')
	buf.WriteString(info.Host)
	buf.WriteByte('|')
	// TODO: min share
	buf.WriteString(escaper.Replace(info.Desc))
	buf.WriteByte('|')
	buf.WriteString(strconv.Itoa(info.Users))
	buf.WriteByte('|')
	buf.WriteString(strconv.FormatUint(info.Share, 10))
	buf.WriteByte('|')

	if host != "" {
		buf.WriteString(host)
		buf.WriteByte('|')
	}
	if _, err := c.Write(buf.Bytes()); err != nil {
		return err
	}
	return c.Close()
}

func scanOneTo(r io.Reader, buf []byte) (int, int, error) {
	for i := 0; i < len(buf); {
		n, err := r.Read(buf[i:])
		if err == io.EOF {
			return -1, 0, io.ErrUnexpectedEOF
		} else if err != nil {
			return -1, 0, err
		}
		if j := bytes.Index(buf[i:i+n], []byte("|")); j > 0 {
			return i + j, i + n, nil
		}
		i += n
	}
	return -1, 0, io.ErrShortBuffer
}

// Registry is an interface for a hublist registry.
type Registry interface {
	RegisterHub(info Info) error
}

// NewServer creates a new hub auto-registration server.
func NewServer(r Registry) *Server {
	return &Server{r: r}
}

type Server struct {
	r Registry
}

func (s *Server) Serve(lis net.Listener) error {
	for {
		c, err := lis.Accept()
		if err != nil {
			return err
		}
		go s.ServeConn(c)
	}
}

func (s *Server) ServeConn(c net.Conn) error {
	if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
		c.Close()
		return err
	}
	return s.serve(c)
}

func (s *Server) serve(c io.ReadWriteCloser) error {
	defer c.Close()
	lock := &nmdc.Lock{
		Lock: "hubAutoReg", // TODO: randomize
		PK:   version.Name + " " + version.Vers,
	}
	data, err := lock.MarshalNMDC()
	if err != nil {
		return err
	}
	b := make([]byte, 2048)
	i := copy(b, "$Lock ")
	i += copy(b[i:], data)
	b[i] = '|'
	i++
	_, err = c.Write(b[:i])
	if err != nil {
		return err
	}
	r := bufio.NewScanner(c)
	r.Buffer(b, cap(b))
	r.Split(func(data []byte, atEOF bool) (advance int, token []byte, _ error) {
		i := bytes.Index(data, []byte("|"))
		if i >= 0 {
			return i + 1, data[:i], nil
		} else if atEOF {
			return 0, nil, io.ErrUnexpectedEOF
		}
		// advance
		return 0, nil, nil
	})
	if !r.Scan() {
		return r.Err()
	}
	data = r.Bytes()
	if !bytes.HasPrefix(data, []byte("$Key ")) {
		return errors.New("expected $Key")
	}
	var key nmdc.Key
	if err := key.UnmarshalNMDC(data[5:]); err != nil {
		return err
	} else if lock.Key().Key != key.Key {
		return errors.New("wrong key")
	}
	var info Info

	readStringRaw := func() (string, error) {
		if !r.Scan() {
			return "", r.Err()
		}
		return r.Text(), nil
	}
	readString := func() (string, error) {
		s, err := readStringRaw()
		if err != nil {
			return "", err
		}
		return nmdc.Unescape(s), nil
	}

	info.Name, err = readString()
	if err != nil {
		return err
	}
	info.Host, err = readStringRaw()
	if err != nil {
		return err
	}
	info.Desc, err = readString()
	if err != nil {
		return err
	}
	str, err := readStringRaw()
	if err != nil {
		return err
	}
	info.Users, err = strconv.Atoi(str)
	if err != nil {
		return err
	}
	str, err = readStringRaw()
	if err != nil {
		return err
	}
	info.Share, err = strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	return s.r.RegisterHub(info)
}

var _ Registry = (RegistryFunc)(nil)

type RegistryFunc func(info Info) error

func (f RegistryFunc) RegisterHub(info Info) error {
	return f(info)
}