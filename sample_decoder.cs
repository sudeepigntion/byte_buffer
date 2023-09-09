using System;
using System.Collections.Generic;
using System.IO;
using System.Text;

namespace byteSample {

  class Decoder {

    public SAMPLE_Decoder(byte[] byteArr)[][] Person {

      ByteBuffer bb = new ByteBuffer();

      bb.Wrap(byteArr);

      int arrLen = (int) bb.GetShort();
      Person[][] obj = new Person[arrLen][];

      for (int i0 = 0; i0 < arrLen; i0++) {

        int arrLen1 = (int) bb.GetShort();
        obj[i0] = new Person[arrLen1];
        for (int i1 = 0; i1 < arrLen1; i1++) {

          obj[i0][i1].Epoch = bb.GetLong();
          int WatchArrLen0 = (int) bb.GetShort();
          int[][] obj[i0][i1].Watch = new int[WatchArrLen0][];
          for (int index00: = 0; index00 < WatchArrLen0; index00++) {

            int WatchArrLen1: = (int) bb.GetShort();
            int[] obj[i0][i1].Watch[index00] = new int[WatchArrLen1];
            for (int index01: = 0; index01 < WatchArrLen1; index01++) {

              obj[i0][i1].Watch[index00][index01] = bb.GetInt();

            }

          }

          obj[i0][i1].Xyz = bb.GetInt();

          obj[i0][i1].Salary = bb.GetDouble();

          int EmployeeArrLen0 = (int) bb.GetShort();
          Employees obj[i0][i1].Employee = new Employees[EmployeeArrLen0][][][];
          for (int index00 = 0; index00 < EmployeeArrLen0; index00++) {

            int EmployeeArrLen1 = (int) bb.GetShort();
            Employees obj[i0][i1].Employee[index00] = new Employees[EmployeeArrLen1][][];
            for (int index11 = 0; index11 < EmployeeArrLen1; index11++) {

              int EmployeeArrLen2 = (int) bb.GetShort();
              Employees obj[i0][i1].Employee[index00][index11] = new Employees[EmployeeArrLen2][];
              for (int index22 = 0; index22 < EmployeeArrLen2; index22++) {

                int EmployeeArrLen3 = (int) bb.GetShort();
                Employees obj[i0][i1].Employee[index00][index11][index22] = new Employees[EmployeeArrLen3];
                for (int index33 = 0; index33 < EmployeeArrLen3; index33++) {

                  obj[i0][i1].Employee[index00][index11][index22][index33].Name = bb.GetString();

                  obj[i0][i1].Employee[index00][index11][index22][index33].Salary = bb.GetDouble();

                  int StudentArrLen0 = (int) bb.GetShort();
                  StudentClass obj[i0][i1].Employee[index00][index11][index22][index33].Student = new StudentClass[StudentArrLen0];
                  for (int index40 = 0; index40 < StudentArrLen0; index40++) {

                    obj[i0][i1].Employee[index00][index11][index22][index33].Student[index40].Name = bb.GetString();

                  }

                }

              }

            }

          }

        }

      }

      bb.Dispose();

      return obj;
    }
  }
}