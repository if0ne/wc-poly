// Generated by `wit-bindgen-wrpc-go` 0.12.0. DO NOT EDIT!
package wc

import (
	bytes "bytes"
	context "context"
	binary "encoding/binary"
	errors "errors"
	fmt "fmt"
	io "io"
	slog "log/slog"
	math "math"
	sync "sync"
	atomic "sync/atomic"
	utf8 "unicode/utf8"
	wrpc "wrpc.io/go"
)

type ChatMessageRequest struct {
	Role    string
	Content string
}

func (v *ChatMessageRequest) String() string { return "ChatMessageRequest" }

func (v *ChatMessageRequest) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 2)
	slog.Debug("writing field", "name", "role")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Role, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `role` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}
	slog.Debug("writing field", "name", "content")
	write1, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Content, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `content` field: %w", err)
	}
	if write1 != nil {
		writes[1] = write1
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

type ChatRequest struct {
	Model    string
	Messages []*ChatMessageRequest
}

func (v *ChatRequest) String() string { return "ChatRequest" }

func (v *ChatRequest) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 2)
	slog.Debug("writing field", "name", "model")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Model, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `model` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}
	slog.Debug("writing field", "name", "messages")
	write1, err := func(v []*ChatMessageRequest, w interface {
		io.ByteWriter
		io.Writer
	}) (write func(wrpc.IndexWriter) error, err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return nil, fmt.Errorf("list length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing list length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return nil, fmt.Errorf("failed to write list length of %d: %w", n, err)
		}
		slog.Debug("writing list elements")
		writes := make(map[uint32]func(wrpc.IndexWriter) error, n)
		for i, e := range v {
			write, err := (e).WriteToIndex(w)
			if err != nil {
				return nil, fmt.Errorf("failed to write list element %d: %w", i, err)
			}
			if write != nil {
				writes[uint32(i)] = write
			}
		}
		if len(writes) > 0 {
			return func(w wrpc.IndexWriter) error {
				var wg sync.WaitGroup
				var wgErr atomic.Value
				for index, write := range writes {
					wg.Add(1)
					w, err := w.Index(index)
					if err != nil {
						return fmt.Errorf("failed to index nested list writer: %w", err)
					}
					write := write
					go func() {
						defer wg.Done()
						if err := write(w); err != nil {
							wgErr.Store(err)
						}
					}()
				}
				wg.Wait()
				err := wgErr.Load()
				if err == nil {
					return nil
				}
				return err.(error)
			}, nil
		}
		return nil, nil
	}(v.Messages, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `messages` field: %w", err)
	}
	if write1 != nil {
		writes[1] = write1
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

type Metrics struct {
	Total              uint64
	Load               uint64
	PromptEvalCount    uint32
	PromptEvalDuration uint64
	EvalCount          uint32
	EvalDuration       uint64
}

func (v *Metrics) String() string { return "Metrics" }

func (v *Metrics) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 6)
	slog.Debug("writing field", "name", "total")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v uint64, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen64)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u64")
		_, err = w.Write(b[:i])
		return err
	}(v.Total, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `total` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}
	slog.Debug("writing field", "name", "load")
	write1, err := (func(wrpc.IndexWriter) error)(nil), func(v uint64, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen64)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u64")
		_, err = w.Write(b[:i])
		return err
	}(v.Load, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `load` field: %w", err)
	}
	if write1 != nil {
		writes[1] = write1
	}
	slog.Debug("writing field", "name", "prompt-eval-count")
	write2, err := (func(wrpc.IndexWriter) error)(nil), func(v uint32, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen32)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u32")
		_, err = w.Write(b[:i])
		return err
	}(v.PromptEvalCount, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `prompt-eval-count` field: %w", err)
	}
	if write2 != nil {
		writes[2] = write2
	}
	slog.Debug("writing field", "name", "prompt-eval-duration")
	write3, err := (func(wrpc.IndexWriter) error)(nil), func(v uint64, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen64)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u64")
		_, err = w.Write(b[:i])
		return err
	}(v.PromptEvalDuration, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `prompt-eval-duration` field: %w", err)
	}
	if write3 != nil {
		writes[3] = write3
	}
	slog.Debug("writing field", "name", "eval-count")
	write4, err := (func(wrpc.IndexWriter) error)(nil), func(v uint32, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen32)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u32")
		_, err = w.Write(b[:i])
		return err
	}(v.EvalCount, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `eval-count` field: %w", err)
	}
	if write4 != nil {
		writes[4] = write4
	}
	slog.Debug("writing field", "name", "eval-duration")
	write5, err := (func(wrpc.IndexWriter) error)(nil), func(v uint64, w io.Writer) (err error) {
		b := make([]byte, binary.MaxVarintLen64)
		i := binary.PutUvarint(b, uint64(v))
		slog.Debug("writing u64")
		_, err = w.Write(b[:i])
		return err
	}(v.EvalDuration, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `eval-duration` field: %w", err)
	}
	if write5 != nil {
		writes[5] = write5
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

type ChatMessageResponse struct {
	Role     string
	Content  string
	Thinking string
	Images   [][]uint8
}

func (v *ChatMessageResponse) String() string { return "ChatMessageResponse" }

func (v *ChatMessageResponse) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 4)
	slog.Debug("writing field", "name", "role")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Role, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `role` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}
	slog.Debug("writing field", "name", "content")
	write1, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Content, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `content` field: %w", err)
	}
	if write1 != nil {
		writes[1] = write1
	}
	slog.Debug("writing field", "name", "thinking")
	write2, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Thinking, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `thinking` field: %w", err)
	}
	if write2 != nil {
		writes[2] = write2
	}
	slog.Debug("writing field", "name", "images")
	write3, err := func(v [][]uint8, w interface {
		io.ByteWriter
		io.Writer
	}) (write func(wrpc.IndexWriter) error, err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return nil, fmt.Errorf("list length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing list length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return nil, fmt.Errorf("failed to write list length of %d: %w", n, err)
		}
		slog.Debug("writing list elements")
		writes := make(map[uint32]func(wrpc.IndexWriter) error, n)
		for i, e := range v {
			write, err := func(v []uint8, w interface {
				io.ByteWriter
				io.Writer
			}) (write func(wrpc.IndexWriter) error, err error) {
				n := len(v)
				if n > math.MaxUint32 {
					return nil, fmt.Errorf("list length of %d overflows a 32-bit integer", n)
				}
				if err = func(v int, w io.Writer) error {
					b := make([]byte, binary.MaxVarintLen32)
					i := binary.PutUvarint(b, uint64(v))
					slog.Debug("writing list length", "len", n)
					_, err = w.Write(b[:i])
					return err
				}(n, w); err != nil {
					return nil, fmt.Errorf("failed to write list length of %d: %w", n, err)
				}
				slog.Debug("writing list elements")
				writes := make(map[uint32]func(wrpc.IndexWriter) error, n)
				for i, e := range v {
					write, err := (func(wrpc.IndexWriter) error)(nil), func(v uint8, w io.ByteWriter) error {
						slog.Debug("writing u8 byte")
						return w.WriteByte(v)
					}(e, w)
					if err != nil {
						return nil, fmt.Errorf("failed to write list element %d: %w", i, err)
					}
					if write != nil {
						writes[uint32(i)] = write
					}
				}
				if len(writes) > 0 {
					return func(w wrpc.IndexWriter) error {
						var wg sync.WaitGroup
						var wgErr atomic.Value
						for index, write := range writes {
							wg.Add(1)
							w, err := w.Index(index)
							if err != nil {
								return fmt.Errorf("failed to index nested list writer: %w", err)
							}
							write := write
							go func() {
								defer wg.Done()
								if err := write(w); err != nil {
									wgErr.Store(err)
								}
							}()
						}
						wg.Wait()
						err := wgErr.Load()
						if err == nil {
							return nil
						}
						return err.(error)
					}, nil
				}
				return nil, nil
			}(e, w)
			if err != nil {
				return nil, fmt.Errorf("failed to write list element %d: %w", i, err)
			}
			if write != nil {
				writes[uint32(i)] = write
			}
		}
		if len(writes) > 0 {
			return func(w wrpc.IndexWriter) error {
				var wg sync.WaitGroup
				var wgErr atomic.Value
				for index, write := range writes {
					wg.Add(1)
					w, err := w.Index(index)
					if err != nil {
						return fmt.Errorf("failed to index nested list writer: %w", err)
					}
					write := write
					go func() {
						defer wg.Done()
						if err := write(w); err != nil {
							wgErr.Store(err)
						}
					}()
				}
				wg.Wait()
				err := wgErr.Load()
				if err == nil {
					return nil
				}
				return err.(error)
			}, nil
		}
		return nil, nil
	}(v.Images, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `images` field: %w", err)
	}
	if write3 != nil {
		writes[3] = write3
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

type ChatResponse struct {
	Model      string
	CreateAt   string
	Message    *ChatMessageResponse
	DoneReason string
	Metrics    *Metrics
}

func (v *ChatResponse) String() string { return "ChatResponse" }

func (v *ChatResponse) WriteToIndex(w wrpc.ByteWriter) (func(wrpc.IndexWriter) error, error) {
	writes := make(map[uint32]func(wrpc.IndexWriter) error, 5)
	slog.Debug("writing field", "name", "model")
	write0, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.Model, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `model` field: %w", err)
	}
	if write0 != nil {
		writes[0] = write0
	}
	slog.Debug("writing field", "name", "create-at")
	write1, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.CreateAt, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `create-at` field: %w", err)
	}
	if write1 != nil {
		writes[1] = write1
	}
	slog.Debug("writing field", "name", "message")
	write2, err := (v.Message).WriteToIndex(w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `message` field: %w", err)
	}
	if write2 != nil {
		writes[2] = write2
	}
	slog.Debug("writing field", "name", "done-reason")
	write3, err := (func(wrpc.IndexWriter) error)(nil), func(v string, w io.Writer) (err error) {
		n := len(v)
		if n > math.MaxUint32 {
			return fmt.Errorf("string byte length of %d overflows a 32-bit integer", n)
		}
		if err = func(v int, w io.Writer) error {
			b := make([]byte, binary.MaxVarintLen32)
			i := binary.PutUvarint(b, uint64(v))
			slog.Debug("writing string byte length", "len", n)
			_, err = w.Write(b[:i])
			return err
		}(n, w); err != nil {
			return fmt.Errorf("failed to write string byte length of %d: %w", n, err)
		}
		slog.Debug("writing string bytes")
		_, err = w.Write([]byte(v))
		if err != nil {
			return fmt.Errorf("failed to write string bytes: %w", err)
		}
		return nil
	}(v.DoneReason, w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `done-reason` field: %w", err)
	}
	if write3 != nil {
		writes[3] = write3
	}
	slog.Debug("writing field", "name", "metrics")
	write4, err := (v.Metrics).WriteToIndex(w)
	if err != nil {
		return nil, fmt.Errorf("failed to write `metrics` field: %w", err)
	}
	if write4 != nil {
		writes[4] = write4
	}

	if len(writes) > 0 {
		return func(w wrpc.IndexWriter) error {
			var wg sync.WaitGroup
			var wgErr atomic.Value
			for index, write := range writes {
				wg.Add(1)
				w, err := w.Index(index)
				if err != nil {
					return fmt.Errorf("failed to index nested record writer: %w", err)
				}
				write := write
				go func() {
					defer wg.Done()
					if err := write(w); err != nil {
						wgErr.Store(err)
					}
				}()
			}
			wg.Wait()
			err := wgErr.Load()
			if err == nil {
				return nil
			}
			return err.(error)
		}, nil
	}
	return nil, nil
}

type Handler interface {
	Chat(ctx__ context.Context, req *ChatRequest) (*ChatResponse, error)
}

func ServeInterface(s wrpc.Server, h Handler) (stop func() error, err error) {
	stops := make([]func() error, 0, 1)
	stop = func() error {
		for _, stop := range stops {
			if err := stop(); err != nil {
				return err
			}
		}
		return nil
	}

	stop0, err := s.Serve("yasp:giga/wc", "chat", func(ctx context.Context, w wrpc.IndexWriteCloser, r wrpc.IndexReadCloser) {
		defer func() {
			if err := w.Close(); err != nil {
				slog.DebugContext(ctx, "failed to close writer", "instance", "yasp:giga/wc", "name", "chat", "err", err)
			}
		}()
		slog.DebugContext(ctx, "reading parameter", "i", 0)
		p0, err := func(r wrpc.IndexReadCloser, path ...uint32) (*ChatRequest, error) {
			v := &ChatRequest{}
			var err error
			slog.Debug("reading field", "name", "model")
			v.Model, err = func(r interface {
				io.ByteReader
				io.Reader
			}) (string, error) {
				var x uint32
				var s uint8
				for i := 0; i < 5; i++ {
					slog.Debug("reading string length byte", "i", i)
					b, err := r.ReadByte()
					if err != nil {
						if i > 0 && err == io.EOF {
							err = io.ErrUnexpectedEOF
						}
						return "", fmt.Errorf("failed to read string length byte: %w", err)
					}
					if s == 28 && b > 0x0f {
						return "", errors.New("string length overflows a 32-bit integer")
					}
					if b < 0x80 {
						x = x | uint32(b)<<s
						if x == 0 {
							return "", nil
						}
						buf := make([]byte, x)
						slog.Debug("reading string bytes", "len", x)
						_, err = r.Read(buf)
						if err != nil {
							return "", fmt.Errorf("failed to read string bytes: %w", err)
						}
						if !utf8.Valid(buf) {
							return string(buf), errors.New("string is not valid UTF-8")
						}
						return string(buf), nil
					}
					x |= uint32(b&0x7f) << s
					s += 7
				}
				return "", errors.New("string length overflows a 32-bit integer")
			}(r)
			if err != nil {
				return nil, fmt.Errorf("failed to read `model` field: %w", err)
			}
			slog.Debug("reading field", "name", "messages")
			v.Messages, err = func(r wrpc.IndexReadCloser, path ...uint32) ([]*ChatMessageRequest, error) {
				var x uint32
				var s uint
				for i := 0; i < 5; i++ {
					slog.Debug("reading list length byte", "i", i)
					b, err := r.ReadByte()
					if err != nil {
						if i > 0 && err == io.EOF {
							err = io.ErrUnexpectedEOF
						}
						return nil, fmt.Errorf("failed to read list length byte: %w", err)
					}
					if s == 28 && b > 0x0f {
						return nil, errors.New("list length overflows a 32-bit integer")
					}
					if b < 0x80 {
						x = x | uint32(b)<<s
						if x == 0 {
							return nil, nil
						}
						vs := make([]*ChatMessageRequest, x)
						for i := range vs {
							slog.Debug("reading list element", "i", i)
							vs[i], err = func(r wrpc.IndexReadCloser, path ...uint32) (*ChatMessageRequest, error) {
								v := &ChatMessageRequest{}
								var err error
								slog.Debug("reading field", "name", "role")
								v.Role, err = func(r interface {
									io.ByteReader
									io.Reader
								}) (string, error) {
									var x uint32
									var s uint8
									for i := 0; i < 5; i++ {
										slog.Debug("reading string length byte", "i", i)
										b, err := r.ReadByte()
										if err != nil {
											if i > 0 && err == io.EOF {
												err = io.ErrUnexpectedEOF
											}
											return "", fmt.Errorf("failed to read string length byte: %w", err)
										}
										if s == 28 && b > 0x0f {
											return "", errors.New("string length overflows a 32-bit integer")
										}
										if b < 0x80 {
											x = x | uint32(b)<<s
											if x == 0 {
												return "", nil
											}
											buf := make([]byte, x)
											slog.Debug("reading string bytes", "len", x)
											_, err = r.Read(buf)
											if err != nil {
												return "", fmt.Errorf("failed to read string bytes: %w", err)
											}
											if !utf8.Valid(buf) {
												return string(buf), errors.New("string is not valid UTF-8")
											}
											return string(buf), nil
										}
										x |= uint32(b&0x7f) << s
										s += 7
									}
									return "", errors.New("string length overflows a 32-bit integer")
								}(r)
								if err != nil {
									return nil, fmt.Errorf("failed to read `role` field: %w", err)
								}
								slog.Debug("reading field", "name", "content")
								v.Content, err = func(r interface {
									io.ByteReader
									io.Reader
								}) (string, error) {
									var x uint32
									var s uint8
									for i := 0; i < 5; i++ {
										slog.Debug("reading string length byte", "i", i)
										b, err := r.ReadByte()
										if err != nil {
											if i > 0 && err == io.EOF {
												err = io.ErrUnexpectedEOF
											}
											return "", fmt.Errorf("failed to read string length byte: %w", err)
										}
										if s == 28 && b > 0x0f {
											return "", errors.New("string length overflows a 32-bit integer")
										}
										if b < 0x80 {
											x = x | uint32(b)<<s
											if x == 0 {
												return "", nil
											}
											buf := make([]byte, x)
											slog.Debug("reading string bytes", "len", x)
											_, err = r.Read(buf)
											if err != nil {
												return "", fmt.Errorf("failed to read string bytes: %w", err)
											}
											if !utf8.Valid(buf) {
												return string(buf), errors.New("string is not valid UTF-8")
											}
											return string(buf), nil
										}
										x |= uint32(b&0x7f) << s
										s += 7
									}
									return "", errors.New("string length overflows a 32-bit integer")
								}(r)
								if err != nil {
									return nil, fmt.Errorf("failed to read `content` field: %w", err)
								}
								return v, nil
							}(r, append(path, uint32(i))...)
							if err != nil {
								return nil, fmt.Errorf("failed to read list element %d: %w", i, err)
							}
						}
						return vs, nil
					}
					x |= uint32(b&0x7f) << s
					s += 7
				}
				return nil, errors.New("list length overflows a 32-bit integer")
			}(r, append(path, 1)...)
			if err != nil {
				return nil, fmt.Errorf("failed to read `messages` field: %w", err)
			}
			return v, nil
		}(r, []uint32{0}...)

		if err != nil {
			slog.WarnContext(ctx, "failed to read parameter", "i", 0, "instance", "yasp:giga/wc", "name", "chat", "err", err)
			if err := r.Close(); err != nil {
				slog.ErrorContext(ctx, "failed to close reader", "instance", "yasp:giga/wc", "name", "chat", "err", err)
			}
			return
		}
		slog.DebugContext(ctx, "calling `yasp:giga/wc.chat` handler")
		r0, err := h.Chat(ctx, p0)
		if cErr := r.Close(); cErr != nil {
			slog.ErrorContext(ctx, "failed to close reader", "instance", "yasp:giga/wc", "name", "chat", "err", err)
		}
		if err != nil {
			slog.WarnContext(ctx, "failed to handle invocation", "instance", "yasp:giga/wc", "name", "chat", "err", err)
			return
		}

		var buf bytes.Buffer
		writes := make(map[uint32]func(wrpc.IndexWriter) error, 1)

		write0, err := (r0).WriteToIndex(&buf)
		if err != nil {
			slog.WarnContext(ctx, "failed to write result value", "i", 0, "instance", "yasp:giga/wc", "name", "chat", "err", err)
			return
		}
		if write0 != nil {
			writes[0] = write0
		}
		slog.DebugContext(ctx, "transmitting `yasp:giga/wc.chat` result")
		_, err = w.Write(buf.Bytes())
		if err != nil {
			slog.WarnContext(ctx, "failed to write result", "instance", "yasp:giga/wc", "name", "chat", "err", err)
			return
		}
		if len(writes) > 0 {
			for index, write := range writes {
				_ = write
				switch index {
				case 0:
					w, err := w.Index(0)
					if err != nil {
						slog.ErrorContext(ctx, "failed to index result writer", "instance", "yasp:giga/wc", "name", "chat", "err", err)
						return
					}
					write := write
					go func() {
						if err := write(w); err != nil {
							slog.WarnContext(ctx, "failed to write nested result value", "instance", "yasp:giga/wc", "name", "chat", "err", err)
						}
					}()
				}
			}
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to serve `yasp:giga/wc.chat`: %w", err)
	}
	stops = append(stops, stop0)
	return stop, nil
}
