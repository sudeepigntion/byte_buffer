using System;
using System.IO;
using System.Text;

namespace csharp_byte_buf
{
    public class ByteBuffer
    {
        private MemoryStream stream;
        private BinaryWriter writer;
        private BinaryReader reader;

        public ByteBuffer()
        {
            stream = new MemoryStream();
            writer = new BinaryWriter(stream, Encoding.UTF8);
            reader = new BinaryReader(stream, Encoding.UTF8);
        }

        // Put methods

        public void PutInt(int value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            Array.Reverse(bytes);
            writer.Write(bytes);
        }

        public void PutShort(short value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            Array.Reverse(bytes);           
            writer.Write(bytes);
        }

        public void PutLong(long value)
        {
            byte[] bytes = BitConverter.GetBytes(value);
            Array.Reverse(bytes);            
            writer.Write(bytes);
        }

        public void PutString(string value)
        {           
            long strLen = value.Length;           
            if (strLen > 0)
            {
                if (strLen < 128)
                {
                    PutByte(1);
                    PutByte((byte)strLen);
                }
                else if (strLen < 32768)
                {
                    PutByte(2);
                    PutShort(Convert.ToInt16(strLen));
                }
                else if (strLen < 2147483648)
                {
                    PutByte(3);
                    PutInt(Convert.ToInt32(strLen));
                }
                else
                {                    
                    PutByte(4);
                    PutLong(strLen);                    
                }
            }
            else
            {
                PutByte(1);
                PutByte(1);
            }



            byte[] bytes = Encoding.UTF8.GetBytes(value);
            //Array.Reverse(bytes);
            writer.Write(bytes);
        }

        public void PutByte(byte value)
        {
            byte[] bytes = new byte[] { value };//BitConverter.GetBytes(value);
            Array.Reverse(bytes);
            writer.Write(bytes);
        }

        public void PutFloat(float value)
        {
            value *= 10000; // Multiply double value by 10,000
            byte[] bytes = BitConverter.GetBytes((long)value);
            Array.Reverse(bytes);
            writer.Write(bytes);
        }

        public void PutDouble(double value)
        {
            value *= 10000; // Multiply double value by 10,000
            byte[] bytes = BitConverter.GetBytes((long)value);
            Array.Reverse(bytes);            
            writer.Write(bytes);
        }

        public void PutBoolean(bool value)
        {
            byte[] bytes = new byte[1];
            //int boolVal = 0;

            //if (value)
            //{
            //    boolVal = 1;
            //}
            // Multiply double value by 10,000
            if(value)
            {
                bytes[0] = 1;
            }
            else
            {
                bytes[0] = 0;
            }
            Array.Reverse(bytes);            
            writer.Write(bytes);
        }

        public bool GetBoolean()
        {
            //byte[] readBytes = reader.ReadBytes(1);
            //Array.Reverse(readBytes); // Reverse the byte order
            //int readValue = BitConverter.ToInt32(readBytes, 0);

            //if (readValue == 1)
            //{
            //    return true;
            //}


            return reader.ReadBoolean();
        }

        // Get methods

        public int GetInt()
        {
            byte[] readBytes = reader.ReadBytes(4);
            Array.Reverse(readBytes); // Reverse the byte order
            int readValue = BitConverter.ToInt32(readBytes, 0);
            return readValue;
        }

        public short GetShort()
        {
            byte[] readBytes = reader.ReadBytes(2);
            Array.Reverse(readBytes); // Reverse the byte order
            short readValue = BitConverter.ToInt16(readBytes, 0);
            return readValue;
        }

        public long GetLong()
        {
            byte[] readBytes = reader.ReadBytes(8);
            Array.Reverse(readBytes); // Reverse the byte order
            long readValue = BitConverter.ToInt16(readBytes, 0);
            return readValue;
        }

        public string GetString()
        {
            int typeString = (int)GetByte();
            string stringData = "";
            byte[] readBytes = null;

            if (typeString == 1)
            {
                int strLen = (int)GetByte();
                readBytes = reader.ReadBytes(strLen);
                Array.Reverse(readBytes);
            }
            else if (typeString == 2)
            {
                short strLen = GetShort();
                readBytes = reader.ReadBytes(strLen);
                Array.Reverse(readBytes);
            }
            else if (typeString == 3)
            {
                int strLen = GetInt();
                readBytes = reader.ReadBytes(strLen);
                Array.Reverse(readBytes);
            }
            else if (typeString == 4)
            {
                long strLen = GetLong();
                readBytes = reader.ReadBytes(Convert.ToInt32(strLen));
                Array.Reverse(readBytes);
            }
            else
            {
                Console.WriteLine("Invalid string length type...");
            }

            Array.Reverse(readBytes); // Reverse the byte order
            stringData = Encoding.UTF8.GetString(readBytes);

            if (stringData == "X")
            {
                return "";
            }

            return stringData;
        }

        public byte GetByte()
        {
            return reader.ReadByte();
        }

        public float GetFloat()
        {
            byte[] readBytes = reader.ReadBytes(8);
            Array.Reverse(readBytes); // Reverse the byte order
            long readValue = BitConverter.ToInt64(readBytes, 0);
            return (float)(readValue / 10000);
        }

        public double GetDouble()
        {
            byte[] readBytes = reader.ReadBytes(8);
            Array.Reverse(readBytes); // Reverse the byte order
            long readValue = BitConverter.ToInt64(readBytes, 0);
            return (double)(readValue / 10000);
        }

        public void Wrap(byte[] data)
        {
            stream = new MemoryStream(data);
            writer = new BinaryWriter(stream);
            reader = new BinaryReader(stream);
        }

        public byte[] ToArray()
        {
            return stream.ToArray();
        }

        public void Dispose()
        {
            writer.Dispose();
            reader.Dispose();
            stream.Dispose();
        }
    }
}
