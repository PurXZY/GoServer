package gamenet

const (
	preSize  = 0
	initSize = 16
)

type ByteBuffer struct {
	_buffer      []byte
	_prependSize int
	_readerIndex int
	_writerIndex int
}

func NewByteBuffer() *ByteBuffer {
	return &ByteBuffer{
		_buffer:      make([]byte, preSize+initSize),
		_prependSize: preSize,
		_readerIndex: preSize,
		_writerIndex: preSize,
	}
}

func (b *ByteBuffer) Append(buff []byte) {
	size := len(buff)
	if size == 0 {
		return
	}
	b.WrGrow(size)
	copy(b._buffer[b._writerIndex:], buff)
	b.WrFlip(size)
}

func (b *ByteBuffer) WrBuf() []byte {
	if b._writerIndex >= len(b._buffer) {
		return nil
	}
	return b._buffer[b._writerIndex:]
}

func (b *ByteBuffer) WrSize() int {
	return len(b._buffer) - b._writerIndex
}

func (b *ByteBuffer) WrFlip(size int) {
	b._writerIndex += size
}

func (b *ByteBuffer) WrGrow(size int) {
	if size > b.WrSize() {
		b.wrreserve(size)
	}
}

func (b *ByteBuffer) RdBuf() []byte {
	if b._readerIndex >= len(b._buffer) {
		return nil
	}
	return b._buffer[b._readerIndex:]
}

func (b *ByteBuffer) RdReady() bool {
	return b._writerIndex > b._readerIndex
}

func (b *ByteBuffer) RdSize() int {
	return b._writerIndex - b._readerIndex
}

func (b *ByteBuffer) RdFlip(size int) {
	if size < b.RdSize() {
		b._readerIndex += size
	} else {
		b.Reset()
	}
}

func (b *ByteBuffer) Reset() {
	b._readerIndex = b._prependSize
	b._writerIndex = b._prependSize
}

func (b *ByteBuffer) MaxSize() int {
	return len(b._buffer)
}

func (b *ByteBuffer) wrreserve(size int) {
	if b.WrSize()+b._readerIndex < size+b._prependSize {
		newSize := b.MaxSize()
		for newSize < b._writerIndex+size {
			newSize <<= 1
		}
		tmpBuff := make([]byte, newSize)
		copy(tmpBuff, b._buffer)
		b._buffer = tmpBuff
	} else {
		readable := b.RdSize()
		copy(b._buffer[b._prependSize:], b._buffer[b._readerIndex:b._writerIndex])
		b._readerIndex = b._prependSize
		b._writerIndex = b._readerIndex + readable
	}
}

func (b *ByteBuffer) Prepend(buff []byte) bool {
	size := len(buff)
	if b._readerIndex < size {
		return false
	}
	b._readerIndex -= size
	copy(b._buffer[b._readerIndex:], buff)
	return true
}