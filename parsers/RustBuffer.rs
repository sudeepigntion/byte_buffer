use bytebuffer::ByteBuffer;
use bytebuffer::Endian;

#[derive(Debug, Default)]
pub struct ByteBuff {
    pub multiplier: f64,
    pub endian: String,
    pub buffer: ByteBuffer,
}

impl ByteBuff {
    pub fn init(&mut self, endian: String) {
        if endian == "big" {
            self.buffer.set_endian(Endian::BigEndian);
        } else {
            self.buffer.set_endian(Endian::LittleEndian);
        }

        if self.multiplier == 0.0 {
            self.multiplier = 10000.0;
        }
    }

    pub fn wrap(&mut self, byte_data: Vec<u8>){

        self.buffer = ByteBuffer::from_vec(byte_data);
    }

    pub fn put_short(&mut self, value: i16) {
        self.buffer.write_u16(value as u16);
    }

    pub fn put_int(&mut self, value: i32) {
        self.buffer.write_u32(value as u32);
    }

    pub fn put_long(&mut self, value: i64) {
        self.buffer.write_u64(value as u64);
    }

    pub fn put_float(&mut self, value: f64) {
        self.buffer.write_u64((value * self.multiplier) as u64);
    }

    pub fn put_bool(&mut self, value: bool) {
        if value {
            self.buffer.write_u8(1);
        } else {
            self.buffer.write_u8(0);
        }
    }

    pub fn put_string(&mut self, value: String) {
        let str_len: i64 = value.len() as i64;
        if str_len > 0 {
            if str_len < 128 {
                self.buffer.write_u8(1);
                self.buffer.write_u8(str_len as u8);
            } else if str_len < 32768 {
                self.buffer.write_u8(2);
                self.buffer.write_u16(str_len as u16);
            } else if str_len < 2147483648 {
                self.buffer.write_u8(3);
                self.buffer.write_u32(str_len as u32);
            } else {
                self.buffer.write_u8(4);
                self.buffer.write_u64(str_len as u64);
            }

            self.buffer.write_bytes(value.as_bytes());
        } else {
            self.buffer.write_u8(1);
            self.buffer.write_u8(1);
            self.buffer.write_bytes("X".as_bytes());
        }
    }

    pub fn get_short(&mut self) -> i16 {
        match self.buffer.read_u16() {
            Ok(value) => {
                return value as i16;
            }
            Err(err) => {
                println!("{:?}", err);
                return 0;
            }
        }
    }

    pub fn get_int(&mut self) -> i32 {
        match self.buffer.read_u32() {
            Ok(value) => {
                return value as i32;
            }
            Err(err) => {
                println!("{:?}", err);
                return 0;
            }
        }
    }

    pub fn get_long(&mut self) -> i64 {
        match self.buffer.read_u64() {
            Ok(value) => {
                return value as i64;
            }
            Err(err) => {
                println!("{:?}", err);
                return 0;
            }
        }
    }

    pub fn get_float(&mut self) -> f64 {
        match self.buffer.read_u64() {
            Ok(value) => {
                return value as f64 / self.multiplier;
            }
            Err(err) => {
                println!("{:?}", err);
                return 0.0;
            }
        }
    }

    pub fn get_bool(&mut self) -> bool {
        match self.buffer.read_u8() {
            Ok(value) => {
                if value == 1 {
                    return true;
                } else {
                    return false;
                }
            }
            Err(err) => {
                println!("{:?}", err);
                return false;
            }
        }
    }

    pub fn get_string(&mut self) -> String {
        return match self.buffer.read_u8() {
            Ok(type_string) => {
                let mut string_data = String::from("");

                if type_string == 1 {
                    let str_len = self.buffer.read_u8().unwrap();
                    string_data =
                        String::from_utf8(self.buffer.read_bytes(str_len as usize).unwrap()).unwrap();
                } else if type_string == 2 {
                    let str_len = self.buffer.read_u16().unwrap();
                    string_data =
                        String::from_utf8(self.buffer.read_bytes(str_len as usize).unwrap()).unwrap();
                } else if type_string == 3 {
                    let str_len = self.buffer.read_u32().unwrap();
                    string_data =
                        String::from_utf8(self.buffer.read_bytes(str_len as usize).unwrap()).unwrap();
                } else {
                    let str_len = self.buffer.read_u64().unwrap();
                    string_data =
                        String::from_utf8(self.buffer.read_bytes(str_len as usize).unwrap()).unwrap();
                }

                if string_data == "X" {
                    "".to_string()
                } else {
                    string_data
                }
            }
            Err(err) => {
                println!("{:?}", err);
                "".to_string()
            }
        }
    }

    pub fn to_array(&self) -> Vec<u8> {
        return self.buffer.clone().into_vec();
    }
}
