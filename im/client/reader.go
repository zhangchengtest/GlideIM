package client

import (
	"go_im/im/conn"
	"go_im/im/message"
	"sync"
)

var messageReader MessageReader

// recyclePool 回收池, 减少临时对象, 回收复用 ReaderRes
var recyclePool sync.Pool

func init() {
	recyclePool = sync.Pool{
		New: func() interface{} {
			return &ReaderRes{}
		},
	}
	SetMessageReader(&defaultReader{})
}

func SetMessageReader(s MessageReader) {
	messageReader = s
}

type ReaderRes struct {
	err error
	m   *message.Message
}

// Recycle 回收当前对象, 一定要在用完后调用这个方法, 否则无法回收
func (r *ReaderRes) Recycle() {
	r.m = nil
	r.err = nil
	recyclePool.Put(r)
}

// MessageReader 表示一个从连接中(Connection)读取消息的读取者, 可以用于定义如何从连接中读取并解析消息.
type MessageReader interface {

	// Read 阻塞读取, 会阻塞当前协程
	Read(conn conn.Connection) (*message.Message, error)

	// ReadCh 返回两个管道, 第一个用于读取内容, 第二个用于发送停止读取, 停止读取时切记要发送停止信号
	ReadCh(conn conn.Connection) (<-chan *ReaderRes, chan<- struct{})
}

type defaultReader struct{}

func (d *defaultReader) ReadCh(conn conn.Connection) (<-chan *ReaderRes, chan<- struct{}) {
	c := make(chan *ReaderRes)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				goto DONE
			default:
				m, err := d.Read(conn)
				res := recyclePool.Get().(*ReaderRes)
				res.err = err
				res.m = m
				c <- res
			}
		}
	DONE:
		close(done)
		close(c)
	}()
	return c, done
}

func (d *defaultReader) Read(conn conn.Connection) (*message.Message, error) {
	m := message.Message{}
	bytes, err := conn.Read()
	if err != nil {
		return nil, err
	}
	err = m.Deserialize(bytes)
	return &m, err
}
