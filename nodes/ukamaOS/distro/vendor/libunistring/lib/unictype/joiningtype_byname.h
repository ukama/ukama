/* ANSI-C code produced by gperf version 3.2 */
/* Command-line: gperf -m 10 ./unictype/joiningtype_byname.gperf  */
/* Computed positions: -k'1' */

#if !((' ' == 32) && ('!' == 33) && ('"' == 34) && ('#' == 35) \
      && ('%' == 37) && ('&' == 38) && ('\'' == 39) && ('(' == 40) \
      && (')' == 41) && ('*' == 42) && ('+' == 43) && (',' == 44) \
      && ('-' == 45) && ('.' == 46) && ('/' == 47) && ('0' == 48) \
      && ('1' == 49) && ('2' == 50) && ('3' == 51) && ('4' == 52) \
      && ('5' == 53) && ('6' == 54) && ('7' == 55) && ('8' == 56) \
      && ('9' == 57) && (':' == 58) && (';' == 59) && ('<' == 60) \
      && ('=' == 61) && ('>' == 62) && ('?' == 63) && ('A' == 65) \
      && ('B' == 66) && ('C' == 67) && ('D' == 68) && ('E' == 69) \
      && ('F' == 70) && ('G' == 71) && ('H' == 72) && ('I' == 73) \
      && ('J' == 74) && ('K' == 75) && ('L' == 76) && ('M' == 77) \
      && ('N' == 78) && ('O' == 79) && ('P' == 80) && ('Q' == 81) \
      && ('R' == 82) && ('S' == 83) && ('T' == 84) && ('U' == 85) \
      && ('V' == 86) && ('W' == 87) && ('X' == 88) && ('Y' == 89) \
      && ('Z' == 90) && ('[' == 91) && ('\\' == 92) && (']' == 93) \
      && ('^' == 94) && ('_' == 95) && ('a' == 97) && ('b' == 98) \
      && ('c' == 99) && ('d' == 100) && ('e' == 101) && ('f' == 102) \
      && ('g' == 103) && ('h' == 104) && ('i' == 105) && ('j' == 106) \
      && ('k' == 107) && ('l' == 108) && ('m' == 109) && ('n' == 110) \
      && ('o' == 111) && ('p' == 112) && ('q' == 113) && ('r' == 114) \
      && ('s' == 115) && ('t' == 116) && ('u' == 117) && ('v' == 118) \
      && ('w' == 119) && ('x' == 120) && ('y' == 121) && ('z' == 122) \
      && ('{' == 123) && ('|' == 124) && ('}' == 125) && ('~' == 126))
/* The character set is not based on ISO-646.  */
#error "gperf generated tables don't work with this execution character set. Please report a bug to <bug-gperf@gnu.org>."
#endif

#line 25 "./unictype/joiningtype_byname.gperf"
struct named_joining_type { int name; int joining_type; };

#define TOTAL_KEYWORDS 17
#define MIN_WORD_LENGTH 1
#define MAX_WORD_LENGTH 13
#define MIN_HASH_VALUE 1
#define MAX_HASH_VALUE 21
/* maximum key range = 21, duplicates = 0 */

#ifndef GPERF_DOWNCASE
#define GPERF_DOWNCASE 1
static const unsigned char gperf_downcase[256] =
  {
      0,   1,   2,   3,   4,   5,   6,   7,   8,   9,  10,  11,  12,  13,  14,
     15,  16,  17,  18,  19,  20,  21,  22,  23,  24,  25,  26,  27,  28,  29,
     30,  31,  32,  33,  34,  35,  36,  37,  38,  39,  40,  41,  42,  43,  44,
     45,  46,  47,  48,  49,  50,  51,  52,  53,  54,  55,  56,  57,  58,  59,
     60,  61,  62,  63,  64,  97,  98,  99, 100, 101, 102, 103, 104, 105, 106,
    107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121,
    122,  91,  92,  93,  94,  95,  96,  97,  98,  99, 100, 101, 102, 103, 104,
    105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119,
    120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134,
    135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147, 148, 149,
    150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164,
    165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177, 178, 179,
    180, 181, 182, 183, 184, 185, 186, 187, 188, 189, 190, 191, 192, 193, 194,
    195, 196, 197, 198, 199, 200, 201, 202, 203, 204, 205, 206, 207, 208, 209,
    210, 211, 212, 213, 214, 215, 216, 217, 218, 219, 220, 221, 222, 223, 224,
    225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 235, 236, 237, 238, 239,
    240, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 251, 252, 253, 254,
    255
  };
#endif

#ifndef GPERF_CASE_STRCMP
#define GPERF_CASE_STRCMP 1
static int
gperf_case_strcmp (register const char *s1, register const char *s2)
{
  for (;;)
    {
      unsigned char c1 = gperf_downcase[(unsigned char)*s1++];
      unsigned char c2 = gperf_downcase[(unsigned char)*s2++];
      if (c1 != 0 && c1 == c2)
        continue;
      return (int)c1 - (int)c2;
    }
}
#endif

#ifdef __GNUC__
__inline
#else
#ifdef __cplusplus
inline
#endif
#endif
static unsigned int
joining_type_hash (register const char *str, register size_t len)
{
  static const unsigned char asso_values[] =
    {
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22,  3,  4, 22,
      22, 22, 22, 22,  9, 22,  0, 22,  8, 22,
      22, 22,  1, 22,  6,  2, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22,  3,
       4, 22, 22, 22, 22, 22,  9, 22,  0, 22,
       8, 22, 22, 22,  1, 22,  6,  2, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
      22, 22, 22, 22, 22, 22
    };
  return len + asso_values[(unsigned char)str[0]];
}

struct joining_type_stringpool_t
  {
    char joining_type_stringpool_str1[sizeof("L")];
    char joining_type_stringpool_str2[sizeof("R")];
    char joining_type_stringpool_str3[sizeof("U")];
    char joining_type_stringpool_str4[sizeof("C")];
    char joining_type_stringpool_str5[sizeof("D")];
    char joining_type_stringpool_str7[sizeof("T")];
    char joining_type_stringpool_str11[sizeof("LeftJoining")];
    char joining_type_stringpool_str12[sizeof("Left Joining")];
    char joining_type_stringpool_str13[sizeof("RightJoining")];
    char joining_type_stringpool_str14[sizeof("Right Joining")];
    char joining_type_stringpool_str15[sizeof("DualJoining")];
    char joining_type_stringpool_str16[sizeof("Dual Joining")];
    char joining_type_stringpool_str17[sizeof("Transparent")];
    char joining_type_stringpool_str18[sizeof("NonJoining")];
    char joining_type_stringpool_str19[sizeof("Non Joining")];
    char joining_type_stringpool_str20[sizeof("JoinCausing")];
    char joining_type_stringpool_str21[sizeof("Join Causing")];
  };
static const struct joining_type_stringpool_t joining_type_stringpool_contents =
  {
    "L",
    "R",
    "U",
    "C",
    "D",
    "T",
    "LeftJoining",
    "Left Joining",
    "RightJoining",
    "Right Joining",
    "DualJoining",
    "Dual Joining",
    "Transparent",
    "NonJoining",
    "Non Joining",
    "JoinCausing",
    "Join Causing"
  };
#define joining_type_stringpool ((const char *) &joining_type_stringpool_contents)

static const struct named_joining_type joining_type_names[] =
  {
    {-1},
#line 39 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str1, UC_JOINING_TYPE_L},
#line 40 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str2, UC_JOINING_TYPE_R},
#line 42 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str3, UC_JOINING_TYPE_U},
#line 37 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str4, UC_JOINING_TYPE_C},
#line 38 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str5, UC_JOINING_TYPE_D},
    {-1},
#line 41 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str7, UC_JOINING_TYPE_T},
    {-1}, {-1}, {-1},
#line 48 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str11, UC_JOINING_TYPE_L},
#line 47 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str12, UC_JOINING_TYPE_L},
#line 50 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str13, UC_JOINING_TYPE_R},
#line 49 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str14, UC_JOINING_TYPE_R},
#line 46 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str15, UC_JOINING_TYPE_D},
#line 45 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str16, UC_JOINING_TYPE_D},
#line 51 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str17, UC_JOINING_TYPE_T},
#line 53 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str18, UC_JOINING_TYPE_U},
#line 52 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str19, UC_JOINING_TYPE_U},
#line 44 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str20, UC_JOINING_TYPE_C},
#line 43 "./unictype/joiningtype_byname.gperf"
    {(int)(size_t)&((struct joining_type_stringpool_t *)0)->joining_type_stringpool_str21, UC_JOINING_TYPE_C}
  };

static const struct named_joining_type *
uc_joining_type_lookup (register const char *str, register size_t len)
{
  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = joining_type_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        {
          register int o = joining_type_names[key].name;
          if (o >= 0)
            {
              register const char *s = o + joining_type_stringpool;

              if ((((unsigned char)*str ^ (unsigned char)*s) & ~32) == 0 && !gperf_case_strcmp (str, s))
                return &joining_type_names[key];
            }
        }
    }
  return 0;
}
