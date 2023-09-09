using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace byteSample {

  class Encoder {

    public SAMPLE_Encoder([][] Personint obj) byte[] {

      ByteBuffer bb = new ByteBuffer();

      bb.PutShort(obj.length);

      for (int i0 = 0; i0 < obj.length; i0++) {
        bb.PutShort(obj[i0].length);
        for (int i1 = 0; i1 < obj[i0].length; i1++) {

          bb.PutLong(obj[i0][i1].Epoch);

          bb.PutShort(obj[i0][i1].Watch.length);

          for (int index00 = 0; index00 < obj[i0][i1].Watch.length; index00++) {

            bb.PutShort(obj[i0][i1].Watch[index00].length);

            for (int index01 = 0; index01 < obj[i0][i1].Watch[index00]; index01++) {

              bb.PutInt(obj[i0][i1].Watch[index00][index01]);

            }

          }

          bb.PutInt(obj[i0][i1].Xyz);

          bb.PutDouble(obj[i0][i1].Salary);

          bb.PutShort(obj[i0][i1].Employee.length);

          for (int index00 = 0; index00 < obj[i0][i1].Employee.length; index00++) {

            bb.PutShort(obj[i0][i1].Employee[index00].length);

            for (int index11 = 0; index11 < obj[i0][i1].Employee[index00].length; index11++) {

              bb.PutShort(obj[i0][i1].Employee[index00][index11].length);

              for (int index22 = 0; index22 < obj[i0][i1].Employee[index00][index11].length; index22++) {

                bb.PutShort(obj[i0][i1].Employee[index00][index11][index22].length);

                for (int index33 = 0; index33 < obj[i0][i1].Employee[index00][index11][index22].length; index33++) {

                  bb.PutString(obj[i0][i1].Employee[index00][index11][index22][index33].Name);

                  bb.PutDouble(obj[i0][i1].Employee[index00][index11][index22][index33].Salary);

                  bb.PutShort(obj[i0][i1].Employee[index00][index11][index22][index33].Student.length);

                  for (int index40 = 0; index40 < obj[i0][i1].Employee[index00][index11][index22][index33].Student.length; index40++) {

                    bb.PutString(obj[i0][i1].Employee[index00][index11][index22][index33].Student[index40].Name);

                  }

                }

              }

            }

          }

        }

      }

      bb.Dispose();

      return bb.ToArray();
    }
  }
}