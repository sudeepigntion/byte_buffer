# Byte Buffer Generator
ByteBuffer Generator

    go run bytebuffer.go -fileName=ByteBufferClass -package=packageName -language=golang
    go run bytebuffer.go -fileName=ByteBufferClass -package=packageName -language=rust
    go run bytebuffer.go -fileName=ByteBufferClass -package=packageName -language=csharp
    
    It will generate 3 files (1 class file, 1 encoder file, and 1 decoder file)
    
    Encoding........................
    bytebuffer encoding:  2.457833ms
    bytebuffer length:  2200034
    flatbuffer encoding:  5.342459ms
    flatbuffer length:  3200080
    json encoding:  16.184833ms
    json length:  7100112
    xml encoding:  43.585417ms
    xml length:  6600102
    Decoding........................
    bytebuffer decoding:  578.708µs
    flatbuffer decoding:  992.292µs
    json decoding:  69.930208ms
    xml decoding:  189.59925ms

    After snappy compression in byte buffer, Its time increases 10ms-18ms more than flat buffer but the size is 900KB whereas the flat buffer size was 32MB
    
    Encoding-decoding........................
    bytebuffer length:  943963
    bytebuffer encoding-decoding:  70.762875ms
    flatbuffer length:  32000080
    flatbuffer encoding-decoding:  52.507458ms
    json length:  71000112
    json encoding-decoding:  884.017167ms
    xml length:  66000102
    xml encoding-decoding:  2.366141583s

    After gzip compression in byte buffer, Its time is almost equal to flat buffer but the size decreased to 480KB the flat buffer size without compression is 320MB
    
    Encoding-decoding........................
    bytebuffer length:  485225
    bytebuffer encoding-decoding:  1.505644792s
    flatbuffer length:  320000080
    flatbuffer encoding-decoding:  1.317980541s
    json length:  710000112
    json encoding-decoding:  15.4076795s
    xml length:  660000102
    xml encoding-decoding:  37.202624375s

    
    It is 2-3 times faster than protobuff and flatbuffer. Depending on the structure of the class
    Currently supporting only Golang, Rust, and C# upcoming language support: Java, Javascript, and NodeJS
    For a sample example check the file generated using sample.bb (sample.go, sample_encoder.go and sample_decoder.go)

    To run the generator
    1. Install Golang
    2. Create a byte buffer class with .bb extension for sample check (sample.bb)
    3. Then run go run bytebuffer.go -fileName=yourbytebufferfile.bb -language=golang -package=yourpackagename

    For Rust, you must have this crate
    https://docs.rs/bytebuffer/latest/bytebuffer/struct.ByteBuffer.html
