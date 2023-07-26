/* ANSI-C code produced by gperf version 3.2 */
/* Command-line: gperf -m 10 ./unictype/combiningclass_byname.gperf  */
/* Computed positions: -k'1,6,$' */

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

#line 25 "./unictype/combiningclass_byname.gperf"
struct named_combining_class { int name; int combining_class; };

#define TOTAL_KEYWORDS 53
#define MIN_WORD_LENGTH 1
#define MAX_WORD_LENGTH 20
#define MIN_HASH_VALUE 3
#define MAX_HASH_VALUE 66
/* maximum key range = 64, duplicates = 0 */

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
combining_class_hash (register const char *str, register size_t len)
{
  static const unsigned char asso_values[] =
    {
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 28, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67,  4,  1, 67, 31,  1,
      67,  2,  8,  6, 67, 47, 15, 67, 11,  1,
      67, 67,  9,  6, 10,  3,  2, 22, 67, 13,
      67, 67, 67, 67, 67, 67, 67,  4,  1, 67,
      31,  1, 67,  2,  8,  6, 67, 47, 15, 67,
      11,  1, 67, 67,  9,  6, 10,  3,  2, 22,
      67, 13, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67, 67, 67, 67, 67,
      67, 67, 67, 67, 67, 67
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[5]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 5:
      case 4:
      case 3:
      case 2:
      case 1:
        hval += asso_values[(unsigned char)str[0]];
        break;
    }
  return hval + asso_values[(unsigned char)str[len - 1]];
}

struct combining_class_stringpool_t
  {
    char combining_class_stringpool_str3[sizeof("B")];
    char combining_class_stringpool_str5[sizeof("OV")];
    char combining_class_stringpool_str8[sizeof("ATB")];
    char combining_class_stringpool_str9[sizeof("A")];
    char combining_class_stringpool_str10[sizeof("Above")];
    char combining_class_stringpool_str11[sizeof("ATA")];
    char combining_class_stringpool_str12[sizeof("BR")];
    char combining_class_stringpool_str13[sizeof("VR")];
    char combining_class_stringpool_str14[sizeof("IS")];
    char combining_class_stringpool_str15[sizeof("AR")];
    char combining_class_stringpool_str16[sizeof("Virama")];
    char combining_class_stringpool_str17[sizeof("ATAR")];
    char combining_class_stringpool_str18[sizeof("BL")];
    char combining_class_stringpool_str19[sizeof("R")];
    char combining_class_stringpool_str20[sizeof("Nukta")];
    char combining_class_stringpool_str21[sizeof("AL")];
    char combining_class_stringpool_str22[sizeof("NR")];
    char combining_class_stringpool_str23[sizeof("ATBL")];
    char combining_class_stringpool_str24[sizeof("Right")];
    char combining_class_stringpool_str25[sizeof("Overlay")];
    char combining_class_stringpool_str26[sizeof("AttachedAbove")];
    char combining_class_stringpool_str27[sizeof("Attached Above")];
    char combining_class_stringpool_str28[sizeof("Below")];
    char combining_class_stringpool_str29[sizeof("Left")];
    char combining_class_stringpool_str30[sizeof("BelowRight")];
    char combining_class_stringpool_str31[sizeof("L")];
    char combining_class_stringpool_str32[sizeof("IotaSubscript")];
    char combining_class_stringpool_str33[sizeof("AboveRight")];
    char combining_class_stringpool_str34[sizeof("DB")];
    char combining_class_stringpool_str35[sizeof("BelowLeft")];
    char combining_class_stringpool_str36[sizeof("Iota Subscript")];
    char combining_class_stringpool_str37[sizeof("DA")];
    char combining_class_stringpool_str38[sizeof("AboveLeft")];
    char combining_class_stringpool_str39[sizeof("AttachedBelowLeft")];
    char combining_class_stringpool_str40[sizeof("AttachedAboveRight")];
    char combining_class_stringpool_str41[sizeof("Attached Below Left")];
    char combining_class_stringpool_str42[sizeof("Attached Above Right")];
    char combining_class_stringpool_str44[sizeof("DoubleAbove")];
    char combining_class_stringpool_str45[sizeof("Double Above")];
    char combining_class_stringpool_str47[sizeof("AttachedBelow")];
    char combining_class_stringpool_str48[sizeof("Attached Below")];
    char combining_class_stringpool_str49[sizeof("Below Left")];
    char combining_class_stringpool_str50[sizeof("Below Right")];
    char combining_class_stringpool_str51[sizeof("KV")];
    char combining_class_stringpool_str52[sizeof("Above Left")];
    char combining_class_stringpool_str53[sizeof("Above Right")];
    char combining_class_stringpool_str55[sizeof("NotReordered")];
    char combining_class_stringpool_str56[sizeof("Not Reordered")];
    char combining_class_stringpool_str60[sizeof("NK")];
    char combining_class_stringpool_str61[sizeof("KanaVoicing")];
    char combining_class_stringpool_str63[sizeof("Kana Voicing")];
    char combining_class_stringpool_str65[sizeof("DoubleBelow")];
    char combining_class_stringpool_str66[sizeof("Double Below")];
  };
static const struct combining_class_stringpool_t combining_class_stringpool_contents =
  {
    "B",
    "OV",
    "ATB",
    "A",
    "Above",
    "ATA",
    "BR",
    "VR",
    "IS",
    "AR",
    "Virama",
    "ATAR",
    "BL",
    "R",
    "Nukta",
    "AL",
    "NR",
    "ATBL",
    "Right",
    "Overlay",
    "AttachedAbove",
    "Attached Above",
    "Below",
    "Left",
    "BelowRight",
    "L",
    "IotaSubscript",
    "AboveRight",
    "DB",
    "BelowLeft",
    "Iota Subscript",
    "DA",
    "AboveLeft",
    "AttachedBelowLeft",
    "AttachedAboveRight",
    "Attached Below Left",
    "Attached Above Right",
    "DoubleAbove",
    "Double Above",
    "AttachedBelow",
    "Attached Below",
    "Below Left",
    "Below Right",
    "KV",
    "Above Left",
    "Above Right",
    "NotReordered",
    "Not Reordered",
    "NK",
    "KanaVoicing",
    "Kana Voicing",
    "DoubleBelow",
    "Double Below"
  };
#define combining_class_stringpool ((const char *) &combining_class_stringpool_contents)

static const struct named_combining_class combining_class_names[] =
  {
    {-1}, {-1}, {-1},
#line 47 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str3, UC_CCC_B},
    {-1},
#line 38 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str5, UC_CCC_OV},
    {-1}, {-1},
#line 43 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str8, UC_CCC_ATB},
#line 52 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str9, UC_CCC_A},
#line 81 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str10, UC_CCC_A},
#line 44 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str11, UC_CCC_ATA},
#line 48 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str12, UC_CCC_BR},
#line 41 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str13, UC_CCC_VR},
#line 56 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str14, UC_CCC_IS},
#line 53 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str15, UC_CCC_AR},
#line 63 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str16, UC_CCC_VR},
#line 45 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str17, UC_CCC_ATAR},
#line 46 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str18, UC_CCC_BL},
#line 50 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str19, UC_CCC_R},
#line 60 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str20, UC_CCC_NK},
#line 51 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str21, UC_CCC_AL},
#line 37 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str22, UC_CCC_NR},
#line 42 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str23, UC_CCC_ATBL},
#line 78 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str24, UC_CCC_R},
#line 59 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str25, UC_CCC_OV},
#line 69 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str26, UC_CCC_ATA},
#line 68 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str27, UC_CCC_ATA},
#line 74 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str28, UC_CCC_B},
#line 77 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str29, UC_CCC_L},
#line 76 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str30, UC_CCC_BR},
#line 49 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str31, UC_CCC_L},
#line 89 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str32, UC_CCC_IS},
#line 83 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str33, UC_CCC_AR},
#line 54 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str34, UC_CCC_DB},
#line 73 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str35, UC_CCC_BL},
#line 88 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str36, UC_CCC_IS},
#line 55 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str37, UC_CCC_DA},
#line 80 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str38, UC_CCC_AL},
#line 65 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str39, UC_CCC_ATBL},
#line 71 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str40, UC_CCC_ATAR},
#line 64 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str41, UC_CCC_ATBL},
#line 70 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str42, UC_CCC_ATAR},
    {-1},
#line 87 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str44, UC_CCC_DA},
#line 86 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str45, UC_CCC_DA},
    {-1},
#line 67 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str47, UC_CCC_ATB},
#line 66 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str48, UC_CCC_ATB},
#line 72 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str49, UC_CCC_BL},
#line 75 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str50, UC_CCC_BR},
#line 40 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str51, UC_CCC_KV},
#line 79 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str52, UC_CCC_AL},
#line 82 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str53, UC_CCC_AR},
    {-1},
#line 58 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str55, UC_CCC_NR},
#line 57 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str56, UC_CCC_NR},
    {-1}, {-1}, {-1},
#line 39 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str60, UC_CCC_NK},
#line 62 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str61, UC_CCC_KV},
    {-1},
#line 61 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str63, UC_CCC_KV},
    {-1},
#line 85 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str65, UC_CCC_DB},
#line 84 "./unictype/combiningclass_byname.gperf"
    {(int)(size_t)&((struct combining_class_stringpool_t *)0)->combining_class_stringpool_str66, UC_CCC_DB}
  };

static const struct named_combining_class *
uc_combining_class_lookup (register const char *str, register size_t len)
{
  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = combining_class_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        {
          register int o = combining_class_names[key].name;
          if (o >= 0)
            {
              register const char *s = o + combining_class_stringpool;

              if ((((unsigned char)*str ^ (unsigned char)*s) & ~32) == 0 && !gperf_case_strcmp (str, s))
                return &combining_class_names[key];
            }
        }
    }
  return 0;
}
