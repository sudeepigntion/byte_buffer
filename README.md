# byte_buffer
ByteBuffer Generator

    go run bytebuffer.go -fileName=ByteBufferClass -package=packageName -language=golang
    
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

    After compression in byte buffer, It's time increases 10ms-18ms more than flat buffer but the size is 900KB where as the flat buffer size was 32MB
    
    Encoding-decoding........................
    bytebuffer length:  943963
    bytebuffer encoding-decoding:  70.762875ms
    flatbuffer length:  32000080
    flatbuffer encoding-decoding:  52.507458ms
    json length:  71000112
    json encoding-decoding:  884.017167ms
    xml length:  66000102
    xml encoding-decoding:  2.366141583s

    
    It is 2-3 times faster than protobuff and flatbuffer. Depending on the structure of the class
    Currently supporting only Golang, upcoming language support: C#, Java, Javascript and Rust
    For sample example check the file generated using sample.bb (sample.go, sample_encoder.go and sample_decoder.go)

    To run the generator
    1. Install Golang
    2. Create a bytebuffer class with .bb extension for sample check (sample.bb)
    3. Then run go run bytebuffer.go -fileName=yourbytebufferfile.bb -language=golang -package=yourpackagename
