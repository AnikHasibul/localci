package localci

// Write implements io.Writer
// appends bytes to buffer
func (ci ciObj) Write(p []byte) (n int, err error) {
	if ci.writeToStdout {
		ci.append(p)
	}
	return len(p), nil
}

// append appends the given byte to buffer
func (ci *ciObj) append(p []byte) {
	ci.wBuff = append(ci.wBuff, p...)
}

// flush flushes the buffer
func (ci *ciObj) flush() []byte {
	defer func() {
		ci.wBuff = nil
	}()
	return ci.wBuff
}
