package sample;

import java.io.*;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.charset.StandardCharsets;

public class JavaBuffer {
    private ByteArrayOutputStream stream;
    private DataOutputStream writer;
    private DataInputStream reader;
    private ByteBuffer buffer;

    public JavaBuffer() {
        stream = new ByteArrayOutputStream();
        writer = new DataOutputStream(stream);
        reader = new DataInputStream(new ByteArrayInputStream(stream.toByteArray()));
    }

    // Put methods
    public void putInt(int value) throws IOException {
        buffer = ByteBuffer.allocate(Integer.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putInt(value);
        writer.write(buffer.array());
    }

    public void putShort(short value) throws IOException {
        buffer = ByteBuffer.allocate(Short.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putShort(value);
        writer.write(buffer.array());
    }

    public void putLong(long value) throws IOException {
        buffer = ByteBuffer.allocate(Long.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putLong(value);
        writer.write(buffer.array());
    }

    public void putString(String value) throws IOException {
        long strLen = value.length();

        if(strLen > 0)
        {
            if (strLen < 128) {
                putByte((byte)1);
                putByte((byte)strLen);
            } else if (strLen < 32768) {
                putByte((byte)2);
                putShort((short) strLen);
            } else if (strLen < 2147483647) {
                putByte((byte)3);
                putInt((int)strLen);
            } else {
                putByte((byte)4);
                putLong(strLen);
            }
            byte[] bytes = value.getBytes(StandardCharsets.UTF_8);
            writer.write(bytes);
        }
        else {
            putByte((byte)1);
            putByte((byte)1);
            writer.write("X".getBytes(StandardCharsets.UTF_8));
        }
    }

    public void putByte(byte value) throws IOException {
        buffer = ByteBuffer.allocate(Byte.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.put(value);
        writer.write(buffer.array());
    }

    public void putFloat(float value) throws IOException {
        buffer = ByteBuffer.allocate(Long.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putLong((long) (value * 10000.0));
        writer.write(buffer.array());
    }

    public void putFloatUsingIntEncoding(float value) throws IOException {
        buffer = ByteBuffer.allocate(Long.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putLong((long) (value * 10000.0));
        writer.write(buffer.array());
    }

    public void putDouble(double value) throws IOException {
        buffer = ByteBuffer.allocate(Double.BYTES);
        buffer.order(ByteOrder.BIG_ENDIAN);
        buffer.putDouble(value);
        writer.write(buffer.array());
    }

    public void putBoolean(boolean value) throws IOException {
        if(value)
        {
            buffer = ByteBuffer.allocate(Byte.BYTES);
            buffer.order(ByteOrder.BIG_ENDIAN);
            buffer.put((byte)1);
        }
        else
        {
            buffer = ByteBuffer.allocate(Byte.BYTES);
            buffer.order(ByteOrder.BIG_ENDIAN);
            buffer.put((byte)0);
        }
        writer.write(buffer.array());
    }

    // Get methods
    public int getInt() throws IOException {
        byte[] bytes = new byte[Integer.BYTES];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getInt();
    }

    public short getShort() throws IOException {
        byte[] bytes = new byte[Short.BYTES];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getShort();
    }

    public long getLong() throws IOException {
        byte[] bytes = new byte[Long.BYTES];
        reader.readFully(bytes);
        ByteBuffer buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getLong();
    }

    public long getLongForStr() throws IOException {
        byte[] bytes = new byte[4];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getLong();
    }

    public String getString() throws IOException {
            int typeString = getByte();
            String stringData;
            byte[] readBytes;
            if(typeString == 1)
            {
                int strLen = getByte();
                readBytes = new byte[strLen];
                reader.readFully(readBytes);
            }
            else if(typeString == 2){
                short strLen = getShort();
                readBytes = new byte[strLen];
                reader.readFully(readBytes);
            }
            else if(typeString == 3){
                int strLen = getInt();
                readBytes = new byte[strLen];
                reader.readFully(readBytes);
            }
            else if(typeString == 4){
                long strLen = getLongForStr();
                readBytes = new byte[(int) strLen];
                reader.readFully(readBytes);
            }
            else{
                throw new IOException("Invalid string length type...");
            }

            stringData = new String(readBytes, StandardCharsets.UTF_8);

            if(stringData.equals("X")){
                return "";
            }

            return stringData;
    }

    public byte getByte() throws IOException {
        return reader.readByte();
    }

    public float getFloat() throws IOException {
        byte[] bytes = new byte[Float.BYTES];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getFloat();
    }

    public float getFloatUsingIntEncoding() throws IOException {
        byte[] bytes = new byte[Long.BYTES];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);

        return (float) (buffer.getLong()/10000);
    }

    public double getDouble() throws IOException {
        byte[] bytes = new byte[Double.BYTES];
        reader.readFully(bytes);
        buffer = ByteBuffer.wrap(bytes);
        buffer.order(ByteOrder.BIG_ENDIAN);
        return buffer.getDouble();
    }

    public boolean getBool() throws IOException {
        return reader.readBoolean();
    }

    public void wrap(byte[] data) {
        stream = new ByteArrayOutputStream();
        writer = new DataOutputStream(stream);
        reader = new DataInputStream(new ByteArrayInputStream(data));
    }

    public byte[] toArray() {
        return stream.toByteArray();
    }

    public void close() throws IOException {
        writer.close();
        reader.close();
        stream.close();
    }
}

