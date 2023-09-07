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


    
    It is 2-3 times faster than protobuff and flatbuffer. Depending on the structure of the class
    Currently supporting only Golang, upcoming language support: C#, Java, Javascript and Rust
    For sample example check the file generated using sample.bb (sample.go, sample_encoder.go and sample_decoder.go)

    To run the generator
    1. Install Golang
    2. Create a bytebuffer class with .bb extension for sample check (sample.bb)
    3. Then run go run bytebuffer.go -fileName=yourbytebufferfile.bb -language=golang -package=yourpackagename
