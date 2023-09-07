# byte_buffer
ByteBuffer Generator

go run bytebuffer.go -fileName=ByteBufferClass -package=packageName -language=golang

It will generate 3 files (1 class file, 1 encoder file, and 1 decoder file)

Encoding........................
bytebuffer encoding:  3.01875ms
bytebuffer length:  2200076
json encoding:  16.064125ms
json length:  7100252
xml encoding:  41.675542ms
xml length:  6600232
Decoding........................
bytebuffer decoding:  884.25µs
json decoding:  69.766083ms
xml decoding:  188.862208ms

It is 2-3 times faster than protobuff and flatbuffer. Depending on the structure of the class
