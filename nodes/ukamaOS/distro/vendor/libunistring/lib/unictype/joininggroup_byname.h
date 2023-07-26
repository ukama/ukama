/* ANSI-C code produced by gperf version 3.2 */
/* Command-line: gperf -m 10 ./unictype/joininggroup_byname.gperf  */
/* Computed positions: -k'1-2,10-12,$' */

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

#line 25 "./unictype/joininggroup_byname.gperf"
struct named_joining_group { int name; int joining_group; };

#define TOTAL_KEYWORDS 160
#define MIN_WORD_LENGTH 1
#define MAX_WORD_LENGTH 24
#define MIN_HASH_VALUE 19
#define MAX_HASH_VALUE 363
/* maximum key range = 345, duplicates = 0 */

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
joining_group_hash (register const char *str, register size_t len)
{
  static const unsigned short asso_values[] =
    {
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364,   6,   4, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364,  15,  18,   6, 131,   9,
       54,  50,   4,  35, 118, 162,  16,   5,   4,  91,
       39, 149, 120, 114,  28,  80,  19,  50,   8,  72,
      129,   4, 364, 364, 364, 364, 364,  15,  18,   6,
      131,   9,  54,  50,   4,  35, 118, 162,  16,   5,
        4,  91,  39, 149, 120, 114,  28,  80,  19,  50,
        8,  72, 129,   4, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364, 364, 364, 364,
      364, 364, 364, 364, 364, 364, 364
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[11]+1];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 11:
        hval += asso_values[(unsigned char)str[10]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 10:
        hval += asso_values[(unsigned char)str[9]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 9:
      case 8:
      case 7:
      case 6:
      case 5:
      case 4:
      case 3:
      case 2:
        hval += asso_values[(unsigned char)str[1]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 1:
        hval += asso_values[(unsigned char)str[0]];
        break;
    }
  return hval + asso_values[(unsigned char)str[len - 1]];
}

struct joining_group_stringpool_t
  {
    char joining_group_stringpool_str19[sizeof("E")];
    char joining_group_stringpool_str20[sizeof("Heh")];
    char joining_group_stringpool_str21[sizeof("Heth")];
    char joining_group_stringpool_str23[sizeof("Meem")];
    char joining_group_stringpool_str24[sizeof("He")];
    char joining_group_stringpool_str26[sizeof("Hah")];
    char joining_group_stringpool_str34[sizeof("Beh")];
    char joining_group_stringpool_str35[sizeof("Beth")];
    char joining_group_stringpool_str36[sizeof("HehGoal")];
    char joining_group_stringpool_str37[sizeof("Heh Goal")];
    char joining_group_stringpool_str39[sizeof("Lam")];
    char joining_group_stringpool_str40[sizeof("Alaph")];
    char joining_group_stringpool_str41[sizeof("Lamadh")];
    char joining_group_stringpool_str43[sizeof("ThinYeh")];
    char joining_group_stringpool_str44[sizeof("Thin Yeh")];
    char joining_group_stringpool_str45[sizeof("Teth")];
    char joining_group_stringpool_str48[sizeof("Mim")];
    char joining_group_stringpool_str50[sizeof("Tah")];
    char joining_group_stringpool_str53[sizeof("Manichaean Mem")];
    char joining_group_stringpool_str54[sizeof("Manichaean Zayin")];
    char joining_group_stringpool_str55[sizeof("Manichaean Beth")];
    char joining_group_stringpool_str56[sizeof("Manichaean Lamedh")];
    char joining_group_stringpool_str57[sizeof("Ain")];
    char joining_group_stringpool_str59[sizeof("Pe")];
    char joining_group_stringpool_str60[sizeof("Manichaean Daleth")];
    char joining_group_stringpool_str61[sizeof("Manichaean Dhamedh")];
    char joining_group_stringpool_str62[sizeof("Malayalam Nga")];
    char joining_group_stringpool_str63[sizeof("ManichaeanAleph")];
    char joining_group_stringpool_str64[sizeof("ManichaeanNun")];
    char joining_group_stringpool_str65[sizeof("Manichaean Kaph")];
    char joining_group_stringpool_str66[sizeof("Manichaean Gimel")];
    char joining_group_stringpool_str67[sizeof("Manichaean Ayin")];
    char joining_group_stringpool_str68[sizeof("Manichaean Aleph")];
    char joining_group_stringpool_str70[sizeof("Feh")];
    char joining_group_stringpool_str73[sizeof("MalayalamNna")];
    char joining_group_stringpool_str74[sizeof("Fe")];
    char joining_group_stringpool_str75[sizeof("Malayalam Lla")];
    char joining_group_stringpool_str76[sizeof("Malayalam Llla")];
    char joining_group_stringpool_str77[sizeof("TehMarbuta")];
    char joining_group_stringpool_str78[sizeof("ManichaeanLamedh")];
    char joining_group_stringpool_str79[sizeof("Manichaean Samekh")];
    char joining_group_stringpool_str81[sizeof("FarsiYeh")];
    char joining_group_stringpool_str82[sizeof("Farsi Yeh")];
    char joining_group_stringpool_str83[sizeof("Manichaean Sadhe")];
    char joining_group_stringpool_str84[sizeof("Manichaean Heth")];
    char joining_group_stringpool_str85[sizeof("MalayalamLlla")];
    char joining_group_stringpool_str86[sizeof("Gamal")];
    char joining_group_stringpool_str87[sizeof("MalayalamBha")];
    char joining_group_stringpool_str88[sizeof("Yeh")];
    char joining_group_stringpool_str89[sizeof("Alef")];
    char joining_group_stringpool_str91[sizeof("Nun")];
    char joining_group_stringpool_str92[sizeof("Manichaean One")];
    char joining_group_stringpool_str93[sizeof("HanifiRohingyaPa")];
    char joining_group_stringpool_str94[sizeof("Nya")];
    char joining_group_stringpool_str96[sizeof("Taw")];
    char joining_group_stringpool_str97[sizeof("MalayalamLla")];
    char joining_group_stringpool_str98[sizeof("HanifiRohingyaKinnaYa")];
    char joining_group_stringpool_str99[sizeof("YehBarree")];
    char joining_group_stringpool_str100[sizeof("ManichaeanHeth")];
    char joining_group_stringpool_str101[sizeof("ManichaeanMem")];
    char joining_group_stringpool_str102[sizeof("Manichaean Waw")];
    char joining_group_stringpool_str103[sizeof("Noon")];
    char joining_group_stringpool_str104[sizeof("Manichaean Five")];
    char joining_group_stringpool_str106[sizeof("Teh Marbuta")];
    char joining_group_stringpool_str107[sizeof("Malayalam Bha")];
    char joining_group_stringpool_str108[sizeof("ManichaeanThamedh")];
    char joining_group_stringpool_str109[sizeof("Yeh Barree")];
    char joining_group_stringpool_str111[sizeof("VerticalTail")];
    char joining_group_stringpool_str114[sizeof("ManichaeanBeth")];
    char joining_group_stringpool_str116[sizeof("Teh Marbuta Goal")];
    char joining_group_stringpool_str118[sizeof("Waw")];
    char joining_group_stringpool_str119[sizeof("MalayalamNga")];
    char joining_group_stringpool_str121[sizeof("MalayalamTta")];
    char joining_group_stringpool_str122[sizeof("Gaf")];
    char joining_group_stringpool_str123[sizeof("ManichaeanTen")];
    char joining_group_stringpool_str124[sizeof("ManichaeanTeth")];
    char joining_group_stringpool_str126[sizeof("Shin")];
    char joining_group_stringpool_str128[sizeof("Manichaean Ten")];
    char joining_group_stringpool_str129[sizeof("Manichaean Teth")];
    char joining_group_stringpool_str131[sizeof("Seen")];
    char joining_group_stringpool_str132[sizeof("Manichaean Thamedh")];
    char joining_group_stringpool_str133[sizeof("ManichaeanTaw")];
    char joining_group_stringpool_str134[sizeof("Semkath")];
    char joining_group_stringpool_str136[sizeof("Reh")];
    char joining_group_stringpool_str138[sizeof("ManichaeanPe")];
    char joining_group_stringpool_str139[sizeof("Manichaean Nun")];
    char joining_group_stringpool_str141[sizeof("MalayalamNya")];
    char joining_group_stringpool_str142[sizeof("Zhain")];
    char joining_group_stringpool_str143[sizeof("Sadhe")];
    char joining_group_stringpool_str147[sizeof("MalayalamNnna")];
    char joining_group_stringpool_str148[sizeof("ManichaeanTwenty")];
    char joining_group_stringpool_str149[sizeof("Malayalam Nna")];
    char joining_group_stringpool_str150[sizeof("Malayalam Nnna")];
    char joining_group_stringpool_str152[sizeof("Zain")];
    char joining_group_stringpool_str153[sizeof("ManichaeanYodh")];
    char joining_group_stringpool_str155[sizeof("ManichaeanWaw")];
    char joining_group_stringpool_str157[sizeof("ReversedPe")];
    char joining_group_stringpool_str159[sizeof("YehWithTail")];
    char joining_group_stringpool_str160[sizeof("Yudh")];
    char joining_group_stringpool_str162[sizeof("Malayalam Tta")];
    char joining_group_stringpool_str163[sizeof("Manichaean Resh")];
    char joining_group_stringpool_str164[sizeof("DalathRish")];
    char joining_group_stringpool_str165[sizeof("Dal")];
    char joining_group_stringpool_str167[sizeof("YudhHe")];
    char joining_group_stringpool_str168[sizeof("Yudh He")];
    char joining_group_stringpool_str169[sizeof("Manichaean Qoph")];
    char joining_group_stringpool_str171[sizeof("TehMarbutaGoal")];
    char joining_group_stringpool_str172[sizeof("Qaph")];
    char joining_group_stringpool_str174[sizeof("Manichaean Taw")];
    char joining_group_stringpool_str175[sizeof("Khaph")];
    char joining_group_stringpool_str176[sizeof("ManichaeanSamekh")];
    char joining_group_stringpool_str178[sizeof("Manichaean Yodh")];
    char joining_group_stringpool_str179[sizeof("MalayalamJa")];
    char joining_group_stringpool_str180[sizeof("ManichaeanSadhe")];
    char joining_group_stringpool_str181[sizeof("MalayalamRa")];
    char joining_group_stringpool_str182[sizeof("Hanifi Rohingya Pa")];
    char joining_group_stringpool_str183[sizeof("FinalSemkath")];
    char joining_group_stringpool_str184[sizeof("KnottedHeh")];
    char joining_group_stringpool_str185[sizeof("Kaph")];
    char joining_group_stringpool_str186[sizeof("ManichaeanAyin")];
    char joining_group_stringpool_str187[sizeof("Malayalam Nya")];
    char joining_group_stringpool_str188[sizeof("Hanifi Rohingya Kinna Ya")];
    char joining_group_stringpool_str189[sizeof("Malayalam Ja")];
    char joining_group_stringpool_str190[sizeof("ManichaeanZayin")];
    char joining_group_stringpool_str191[sizeof("Malayalam Ra")];
    char joining_group_stringpool_str193[sizeof("ManichaeanDaleth")];
    char joining_group_stringpool_str194[sizeof("Knotted Heh")];
    char joining_group_stringpool_str195[sizeof("ManichaeanHundred")];
    char joining_group_stringpool_str196[sizeof("Malayalam Ssa")];
    char joining_group_stringpool_str197[sizeof("Reversed Pe")];
    char joining_group_stringpool_str199[sizeof("Manichaean Twenty")];
    char joining_group_stringpool_str201[sizeof("Manichaean Pe")];
    char joining_group_stringpool_str210[sizeof("No Joining Group")];
    char joining_group_stringpool_str211[sizeof("ManichaeanDhamedh")];
    char joining_group_stringpool_str214[sizeof("Manichaean Hundred")];
    char joining_group_stringpool_str216[sizeof("ManichaeanResh")];
    char joining_group_stringpool_str218[sizeof("Vertical Tail")];
    char joining_group_stringpool_str219[sizeof("ManichaeanFive")];
    char joining_group_stringpool_str221[sizeof("Qaf")];
    char joining_group_stringpool_str222[sizeof("ManichaeanKaph")];
    char joining_group_stringpool_str223[sizeof("ManichaeanGimel")];
    char joining_group_stringpool_str226[sizeof("SwashKaf")];
    char joining_group_stringpool_str227[sizeof("Swash Kaf")];
    char joining_group_stringpool_str228[sizeof("ManichaeanOne")];
    char joining_group_stringpool_str230[sizeof("ManichaeanQoph")];
    char joining_group_stringpool_str234[sizeof("Kaf")];
    char joining_group_stringpool_str239[sizeof("RohingyaYeh")];
    char joining_group_stringpool_str245[sizeof("SyriacWaw")];
    char joining_group_stringpool_str263[sizeof("Sad")];
    char joining_group_stringpool_str268[sizeof("StraightWaw")];
    char joining_group_stringpool_str271[sizeof("Yeh with tail")];
    char joining_group_stringpool_str277[sizeof("Straight Waw")];
    char joining_group_stringpool_str279[sizeof("Dalath Rish")];
    char joining_group_stringpool_str287[sizeof("BurushaskiYehBarree")];
    char joining_group_stringpool_str293[sizeof("MalayalamSsa")];
    char joining_group_stringpool_str296[sizeof("Syriac Waw")];
    char joining_group_stringpool_str298[sizeof("Burushaski Yeh Barree")];
    char joining_group_stringpool_str343[sizeof("Rohingya Yeh")];
    char joining_group_stringpool_str357[sizeof("NoJoiningGroup")];
    char joining_group_stringpool_str363[sizeof("Final Semkath")];
  };
static const struct joining_group_stringpool_t joining_group_stringpool_contents =
  {
    "E",
    "Heh",
    "Heth",
    "Meem",
    "He",
    "Hah",
    "Beh",
    "Beth",
    "HehGoal",
    "Heh Goal",
    "Lam",
    "Alaph",
    "Lamadh",
    "ThinYeh",
    "Thin Yeh",
    "Teth",
    "Mim",
    "Tah",
    "Manichaean Mem",
    "Manichaean Zayin",
    "Manichaean Beth",
    "Manichaean Lamedh",
    "Ain",
    "Pe",
    "Manichaean Daleth",
    "Manichaean Dhamedh",
    "Malayalam Nga",
    "ManichaeanAleph",
    "ManichaeanNun",
    "Manichaean Kaph",
    "Manichaean Gimel",
    "Manichaean Ayin",
    "Manichaean Aleph",
    "Feh",
    "MalayalamNna",
    "Fe",
    "Malayalam Lla",
    "Malayalam Llla",
    "TehMarbuta",
    "ManichaeanLamedh",
    "Manichaean Samekh",
    "FarsiYeh",
    "Farsi Yeh",
    "Manichaean Sadhe",
    "Manichaean Heth",
    "MalayalamLlla",
    "Gamal",
    "MalayalamBha",
    "Yeh",
    "Alef",
    "Nun",
    "Manichaean One",
    "HanifiRohingyaPa",
    "Nya",
    "Taw",
    "MalayalamLla",
    "HanifiRohingyaKinnaYa",
    "YehBarree",
    "ManichaeanHeth",
    "ManichaeanMem",
    "Manichaean Waw",
    "Noon",
    "Manichaean Five",
    "Teh Marbuta",
    "Malayalam Bha",
    "ManichaeanThamedh",
    "Yeh Barree",
    "VerticalTail",
    "ManichaeanBeth",
    "Teh Marbuta Goal",
    "Waw",
    "MalayalamNga",
    "MalayalamTta",
    "Gaf",
    "ManichaeanTen",
    "ManichaeanTeth",
    "Shin",
    "Manichaean Ten",
    "Manichaean Teth",
    "Seen",
    "Manichaean Thamedh",
    "ManichaeanTaw",
    "Semkath",
    "Reh",
    "ManichaeanPe",
    "Manichaean Nun",
    "MalayalamNya",
    "Zhain",
    "Sadhe",
    "MalayalamNnna",
    "ManichaeanTwenty",
    "Malayalam Nna",
    "Malayalam Nnna",
    "Zain",
    "ManichaeanYodh",
    "ManichaeanWaw",
    "ReversedPe",
    "YehWithTail",
    "Yudh",
    "Malayalam Tta",
    "Manichaean Resh",
    "DalathRish",
    "Dal",
    "YudhHe",
    "Yudh He",
    "Manichaean Qoph",
    "TehMarbutaGoal",
    "Qaph",
    "Manichaean Taw",
    "Khaph",
    "ManichaeanSamekh",
    "Manichaean Yodh",
    "MalayalamJa",
    "ManichaeanSadhe",
    "MalayalamRa",
    "Hanifi Rohingya Pa",
    "FinalSemkath",
    "KnottedHeh",
    "Kaph",
    "ManichaeanAyin",
    "Malayalam Nya",
    "Hanifi Rohingya Kinna Ya",
    "Malayalam Ja",
    "ManichaeanZayin",
    "Malayalam Ra",
    "ManichaeanDaleth",
    "Knotted Heh",
    "ManichaeanHundred",
    "Malayalam Ssa",
    "Reversed Pe",
    "Manichaean Twenty",
    "Manichaean Pe",
    "No Joining Group",
    "ManichaeanDhamedh",
    "Manichaean Hundred",
    "ManichaeanResh",
    "Vertical Tail",
    "ManichaeanFive",
    "Qaf",
    "ManichaeanKaph",
    "ManichaeanGimel",
    "SwashKaf",
    "Swash Kaf",
    "ManichaeanOne",
    "ManichaeanQoph",
    "Kaf",
    "RohingyaYeh",
    "SyriacWaw",
    "Sad",
    "StraightWaw",
    "Yeh with tail",
    "Straight Waw",
    "Dalath Rish",
    "BurushaskiYehBarree",
    "MalayalamSsa",
    "Syriac Waw",
    "Burushaski Yeh Barree",
    "Rohingya Yeh",
    "NoJoiningGroup",
    "Final Semkath"
  };
#define joining_group_stringpool ((const char *) &joining_group_stringpool_contents)

static const struct named_joining_group joining_group_names[] =
  {
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1},
#line 49 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str19, UC_JOINING_GROUP_E},
#line 60 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str20, UC_JOINING_GROUP_HEH},
#line 63 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str21, UC_JOINING_GROUP_HETH},
    {-1},
#line 71 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str23, UC_JOINING_GROUP_MEEM},
#line 59 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str24, UC_JOINING_GROUP_HE},
    {-1},
#line 58 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str26, UC_JOINING_GROUP_HAH},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 42 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str34, UC_JOINING_GROUP_BEH},
#line 43 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str35, UC_JOINING_GROUP_BETH},
#line 62 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str36, UC_JOINING_GROUP_HEH_GOAL},
#line 61 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str37, UC_JOINING_GROUP_HEH_GOAL},
    {-1},
#line 69 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str39, UC_JOINING_GROUP_LAM},
#line 40 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str40, UC_JOINING_GROUP_ALAPH},
#line 70 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str41, UC_JOINING_GROUP_LAMADH},
    {-1},
#line 194 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str43, UC_JOINING_GROUP_THIN_YEH},
#line 193 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str44, UC_JOINING_GROUP_THIN_YEH},
#line 97 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str45, UC_JOINING_GROUP_TETH},
    {-1}, {-1},
#line 72 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str48, UC_JOINING_GROUP_MIM},
    {-1},
#line 91 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str50, UC_JOINING_GROUP_TAH},
    {-1}, {-1},
#line 139 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str53, UC_JOINING_GROUP_MANICHAEAN_MEM},
#line 123 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str54, UC_JOINING_GROUP_MANICHAEAN_ZAYIN},
#line 115 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str55, UC_JOINING_GROUP_MANICHAEAN_BETH},
#line 133 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str56, UC_JOINING_GROUP_MANICHAEAN_LAMEDH},
#line 39 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str57, UC_JOINING_GROUP_AIN},
    {-1},
#line 76 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str59, UC_JOINING_GROUP_PE},
#line 119 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str60, UC_JOINING_GROUP_MANICHAEAN_DALETH},
#line 135 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str61, UC_JOINING_GROUP_MANICHAEAN_DHAMEDH},
#line 167 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str62, UC_JOINING_GROUP_MALAYALAM_NGA},
#line 114 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str63, UC_JOINING_GROUP_MANICHAEAN_ALEPH},
#line 142 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str64, UC_JOINING_GROUP_MANICHAEAN_NUN},
#line 131 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str65, UC_JOINING_GROUP_MANICHAEAN_KAPH},
#line 117 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str66, UC_JOINING_GROUP_MANICHAEAN_GIMEL},
#line 145 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str67, UC_JOINING_GROUP_MANICHAEAN_AYIN},
#line 113 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str68, UC_JOINING_GROUP_MANICHAEAN_ALEPH},
    {-1},
#line 53 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str70, UC_JOINING_GROUP_FEH},
    {-1}, {-1},
#line 176 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str73, UC_JOINING_GROUP_MALAYALAM_NNA},
#line 52 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str74, UC_JOINING_GROUP_FE},
#line 183 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str75, UC_JOINING_GROUP_MALAYALAM_LLA},
#line 185 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str76, UC_JOINING_GROUP_MALAYALAM_LLLA},
#line 94 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str77, UC_JOINING_GROUP_TEH_MARBUTA},
#line 134 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str78, UC_JOINING_GROUP_MANICHAEAN_LAMEDH},
#line 143 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str79, UC_JOINING_GROUP_MANICHAEAN_SAMEKH},
    {-1},
#line 51 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str81, UC_JOINING_GROUP_FARSI_YEH},
#line 50 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str82, UC_JOINING_GROUP_FARSI_YEH},
#line 149 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str83, UC_JOINING_GROUP_MANICHAEAN_SADHE},
#line 125 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str84, UC_JOINING_GROUP_MANICHAEAN_HETH},
#line 186 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str85, UC_JOINING_GROUP_MALAYALAM_LLLA},
#line 57 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str86, UC_JOINING_GROUP_GAMAL},
#line 180 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str87, UC_JOINING_GROUP_MALAYALAM_BHA},
#line 99 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str88, UC_JOINING_GROUP_YEH},
#line 41 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str89, UC_JOINING_GROUP_ALEF},
    {-1},
#line 74 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str91, UC_JOINING_GROUP_NUN},
#line 157 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str92, UC_JOINING_GROUP_MANICHAEAN_ONE},
#line 190 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str93, UC_JOINING_GROUP_HANIFI_ROHINGYA_PA},
#line 75 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str94, UC_JOINING_GROUP_NYA},
    {-1},
#line 92 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str96, UC_JOINING_GROUP_TAW},
#line 184 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str97, UC_JOINING_GROUP_MALAYALAM_LLA},
#line 192 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str98, UC_JOINING_GROUP_HANIFI_ROHINGYA_KINNA_YA},
#line 101 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str99, UC_JOINING_GROUP_YEH_BARREE},
#line 126 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str100, UC_JOINING_GROUP_MANICHAEAN_HETH},
#line 140 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str101, UC_JOINING_GROUP_MANICHAEAN_MEM},
#line 121 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str102, UC_JOINING_GROUP_MANICHAEAN_WAW},
#line 73 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str103, UC_JOINING_GROUP_NOON},
#line 159 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str104, UC_JOINING_GROUP_MANICHAEAN_FIVE},
    {-1},
#line 93 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str106, UC_JOINING_GROUP_TEH_MARBUTA},
#line 179 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str107, UC_JOINING_GROUP_MALAYALAM_BHA},
#line 138 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str108, UC_JOINING_GROUP_MANICHAEAN_THAMEDH},
#line 100 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str109, UC_JOINING_GROUP_YEH_BARREE},
    {-1},
#line 196 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str111, UC_JOINING_GROUP_VERTICAL_TAIL},
    {-1}, {-1},
#line 116 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str114, UC_JOINING_GROUP_MANICHAEAN_BETH},
    {-1},
#line 95 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str116, UC_JOINING_GROUP_TEH_MARBUTA_GOAL},
    {-1},
#line 98 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str118, UC_JOINING_GROUP_WAW},
#line 168 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str119, UC_JOINING_GROUP_MALAYALAM_NGA},
    {-1},
#line 174 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str121, UC_JOINING_GROUP_MALAYALAM_TTA},
#line 56 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str122, UC_JOINING_GROUP_GAF},
#line 162 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str123, UC_JOINING_GROUP_MANICHAEAN_TEN},
#line 128 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str124, UC_JOINING_GROUP_MANICHAEAN_TETH},
    {-1},
#line 86 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str126, UC_JOINING_GROUP_SHIN},
    {-1},
#line 161 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str128, UC_JOINING_GROUP_MANICHAEAN_TEN},
#line 127 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str129, UC_JOINING_GROUP_MANICHAEAN_TETH},
    {-1},
#line 84 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str131, UC_JOINING_GROUP_SEEN},
#line 137 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str132, UC_JOINING_GROUP_MANICHAEAN_THAMEDH},
#line 156 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str133, UC_JOINING_GROUP_MANICHAEAN_TAW},
#line 85 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str134, UC_JOINING_GROUP_SEMKATH},
    {-1},
#line 79 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str136, UC_JOINING_GROUP_REH},
    {-1},
#line 148 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str138, UC_JOINING_GROUP_MANICHAEAN_PE},
#line 141 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str139, UC_JOINING_GROUP_MANICHAEAN_NUN},
    {-1},
#line 172 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str141, UC_JOINING_GROUP_MALAYALAM_NYA},
#line 108 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str142, UC_JOINING_GROUP_ZHAIN},
#line 83 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str143, UC_JOINING_GROUP_SADHE},
    {-1}, {-1}, {-1},
#line 178 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str147, UC_JOINING_GROUP_MALAYALAM_NNNA},
#line 164 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str148, UC_JOINING_GROUP_MANICHAEAN_TWENTY},
#line 175 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str149, UC_JOINING_GROUP_MALAYALAM_NNA},
#line 177 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str150, UC_JOINING_GROUP_MALAYALAM_NNNA},
    {-1},
#line 107 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str152, UC_JOINING_GROUP_ZAIN},
#line 130 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str153, UC_JOINING_GROUP_MANICHAEAN_YODH},
    {-1},
#line 122 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str155, UC_JOINING_GROUP_MANICHAEAN_WAW},
    {-1},
#line 81 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str157, UC_JOINING_GROUP_REVERSED_PE},
    {-1},
#line 103 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str159, UC_JOINING_GROUP_YEH_WITH_TAIL},
#line 104 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str160, UC_JOINING_GROUP_YUDH},
    {-1},
#line 173 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str162, UC_JOINING_GROUP_MALAYALAM_TTA},
#line 153 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str163, UC_JOINING_GROUP_MANICHAEAN_RESH},
#line 48 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str164, UC_JOINING_GROUP_DALATH_RISH},
#line 46 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str165, UC_JOINING_GROUP_DAL},
    {-1},
#line 106 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str167, UC_JOINING_GROUP_YUDH_HE},
#line 105 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str168, UC_JOINING_GROUP_YUDH_HE},
#line 151 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str169, UC_JOINING_GROUP_MANICHAEAN_QOPH},
    {-1},
#line 96 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str171, UC_JOINING_GROUP_TEH_MARBUTA_GOAL},
#line 78 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str172, UC_JOINING_GROUP_QAPH},
    {-1},
#line 155 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str174, UC_JOINING_GROUP_MANICHAEAN_TAW},
#line 66 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str175, UC_JOINING_GROUP_KHAPH},
#line 144 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str176, UC_JOINING_GROUP_MANICHAEAN_SAMEKH},
    {-1},
#line 129 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str178, UC_JOINING_GROUP_MANICHAEAN_YODH},
#line 170 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str179, UC_JOINING_GROUP_MALAYALAM_JA},
#line 150 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str180, UC_JOINING_GROUP_MANICHAEAN_SADHE},
#line 182 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str181, UC_JOINING_GROUP_MALAYALAM_RA},
#line 189 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str182, UC_JOINING_GROUP_HANIFI_ROHINGYA_PA},
#line 55 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str183, UC_JOINING_GROUP_FINAL_SEMKATH},
#line 68 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str184, UC_JOINING_GROUP_KNOTTED_HEH},
#line 65 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str185, UC_JOINING_GROUP_KAPH},
#line 146 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str186, UC_JOINING_GROUP_MANICHAEAN_AYIN},
#line 171 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str187, UC_JOINING_GROUP_MALAYALAM_NYA},
#line 191 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str188, UC_JOINING_GROUP_HANIFI_ROHINGYA_KINNA_YA},
#line 169 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str189, UC_JOINING_GROUP_MALAYALAM_JA},
#line 124 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str190, UC_JOINING_GROUP_MANICHAEAN_ZAYIN},
#line 181 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str191, UC_JOINING_GROUP_MALAYALAM_RA},
    {-1},
#line 120 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str193, UC_JOINING_GROUP_MANICHAEAN_DALETH},
#line 67 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str194, UC_JOINING_GROUP_KNOTTED_HEH},
#line 166 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str195, UC_JOINING_GROUP_MANICHAEAN_HUNDRED},
#line 187 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str196, UC_JOINING_GROUP_MALAYALAM_SSA},
#line 80 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str197, UC_JOINING_GROUP_REVERSED_PE},
    {-1},
#line 163 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str199, UC_JOINING_GROUP_MANICHAEAN_TWENTY},
    {-1},
#line 147 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str201, UC_JOINING_GROUP_MANICHAEAN_PE},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 37 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str210, UC_JOINING_GROUP_NONE},
#line 136 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str211, UC_JOINING_GROUP_MANICHAEAN_DHAMEDH},
    {-1}, {-1},
#line 165 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str214, UC_JOINING_GROUP_MANICHAEAN_HUNDRED},
    {-1},
#line 154 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str216, UC_JOINING_GROUP_MANICHAEAN_RESH},
    {-1},
#line 195 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str218, UC_JOINING_GROUP_VERTICAL_TAIL},
#line 160 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str219, UC_JOINING_GROUP_MANICHAEAN_FIVE},
    {-1},
#line 77 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str221, UC_JOINING_GROUP_QAF},
#line 132 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str222, UC_JOINING_GROUP_MANICHAEAN_KAPH},
#line 118 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str223, UC_JOINING_GROUP_MANICHAEAN_GIMEL},
    {-1}, {-1},
#line 88 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str226, UC_JOINING_GROUP_SWASH_KAF},
#line 87 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str227, UC_JOINING_GROUP_SWASH_KAF},
#line 158 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str228, UC_JOINING_GROUP_MANICHAEAN_ONE},
    {-1},
#line 152 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str230, UC_JOINING_GROUP_MANICHAEAN_QOPH},
    {-1}, {-1}, {-1},
#line 64 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str234, UC_JOINING_GROUP_KAF},
    {-1}, {-1}, {-1}, {-1},
#line 110 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str239, UC_JOINING_GROUP_ROHINGYA_YEH},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 90 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str245, UC_JOINING_GROUP_SYRIAC_WAW},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 82 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str263, UC_JOINING_GROUP_SAD},
    {-1}, {-1}, {-1}, {-1},
#line 112 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str268, UC_JOINING_GROUP_STRAIGHT_WAW},
    {-1}, {-1},
#line 102 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str271, UC_JOINING_GROUP_YEH_WITH_TAIL},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 111 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str277, UC_JOINING_GROUP_STRAIGHT_WAW},
    {-1},
#line 47 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str279, UC_JOINING_GROUP_DALATH_RISH},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 45 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str287, UC_JOINING_GROUP_BURUSHASKI_YEH_BARREE},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 188 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str293, UC_JOINING_GROUP_MALAYALAM_SSA},
    {-1}, {-1},
#line 89 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str296, UC_JOINING_GROUP_SYRIAC_WAW},
    {-1},
#line 44 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str298, UC_JOINING_GROUP_BURUSHASKI_YEH_BARREE},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 109 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str343, UC_JOINING_GROUP_ROHINGYA_YEH},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1},
#line 38 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str357, UC_JOINING_GROUP_NONE},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 54 "./unictype/joininggroup_byname.gperf"
    {(int)(size_t)&((struct joining_group_stringpool_t *)0)->joining_group_stringpool_str363, UC_JOINING_GROUP_FINAL_SEMKATH}
  };

static const struct named_joining_group *
uc_joining_group_lookup (register const char *str, register size_t len)
{
  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = joining_group_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        {
          register int o = joining_group_names[key].name;
          if (o >= 0)
            {
              register const char *s = o + joining_group_stringpool;

              if ((((unsigned char)*str ^ (unsigned char)*s) & ~32) == 0 && !gperf_case_strcmp (str, s))
                return &joining_group_names[key];
            }
        }
    }
  return 0;
}
