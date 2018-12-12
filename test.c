#include <stdio.h>
#include <stdint.h>

uint16_t CheckSum(uint16_t *buffer, int size)
{
  uint32_t cksum=0;
  printf("# chksum = %x\n", cksum);
  
  while(size >1) {
    cksum += *buffer++;
    size -= sizeof(uint16_t);
    printf("# chksum = %x\n", cksum);
  }
  
  if(size) {
    cksum += *(uint8_t*)buffer;
    printf("# chksum = %x\n", cksum);
  }
  
  cksum = (cksum >> 16) + (cksum & 0xffff);
  printf("# chksum = %x\n", cksum);
  cksum += (cksum >>16);
  printf("# chksum = %x\n", cksum);
  
  return (uint16_t)(~cksum);
}

typedef struct {
  const char* x;
  int len;
} TestCase;

int main(int argc, char* argv[]) {
  const TestCase cases[] = {
			    {"", 0},
			    {"\x55", 1},
			    {"\x55\x88", 2},
			    {"\x55\x88\x99", 3},
			    {"\x7f\x00\x00\x01"
			     "\x7f\x00\x00\x01"
			     "\x00\x06" "\x00\x14"
			     "\x00\x50\x1f\x90\x00"
			     "\x00\x00\x01\x00\x00"
			     "\x00\x00\x50\x02\x00"
			     "\x00\x00\x00\x00\x00",
			     12+20},
			    {NULL, 0},
  };
  for (int i=0; cases[i].x != NULL; ++i) {
    printf("%x\n", CheckSum((uint16_t*)cases[i].x, cases[i].len));
  }
}

/*

0000   02 00 00 00 45 00 00 28 26 b9 00 00 40 06 00 00   ....E..(&¹..@...
0010   7f 00 00 01 7f 00 00 01 00 50 1f 90 00 00 00 01   .........P......
0020   00 00 00 00 50 02 00 00 dd f4 00 00               ....P...Ýô..

*/
