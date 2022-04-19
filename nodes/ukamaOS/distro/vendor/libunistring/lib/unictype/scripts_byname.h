/* ANSI-C code produced by gperf version 3.1 */
/* Command-line: gperf -m 10 ./unictype/scripts_byname.gperf  */
/* Computed positions: -k'1,3,5,8' */

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

#line 4 "./unictype/scripts_byname.gperf"
struct named_script { int name; unsigned int index; };

#define TOTAL_KEYWORDS 137
#define MIN_WORD_LENGTH 2
#define MAX_WORD_LENGTH 22
#define MIN_HASH_VALUE 2
#define MAX_HASH_VALUE 210
/* maximum key range = 209, duplicates = 0 */

#ifdef __GNUC__
__inline
#else
#ifdef __cplusplus
inline
#endif
#endif
static unsigned int
scripts_hash (register const char *str, register size_t len)
{
  static const unsigned char asso_values[] =
    {
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211,   9,   0,  16,  40,  28,
      211,  64,  66,  24,  68,  47,  29,   8,  71,  44,
       16, 211,  61,  26,  19,  46, 102,  23, 211,   0,
      211, 211, 211, 211, 211,  51, 211,   3,  54,  55,
       36,   3,  55,  51,  36,  11,  48,  26,   6,   2,
        0,   1,  70, 211,  16,  64,  11,  18,  46,   3,
      211,  83, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211, 211, 211, 211, 211,
      211, 211, 211, 211, 211, 211
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[7]];
      /*FALLTHROUGH*/
      case 7:
      case 6:
      case 5:
        hval += asso_values[(unsigned char)str[4]];
      /*FALLTHROUGH*/
      case 4:
      case 3:
        hval += asso_values[(unsigned char)str[2]];
      /*FALLTHROUGH*/
      case 2:
      case 1:
        hval += asso_values[(unsigned char)str[0]];
        break;
    }
  return hval;
}

struct script_stringpool_t
  {
    char script_stringpool_str2[sizeof("Yi")];
    char script_stringpool_str9[sizeof("Bamum")];
    char script_stringpool_str10[sizeof("Bengali")];
    char script_stringpool_str11[sizeof("Brahmi")];
    char script_stringpool_str12[sizeof("Mro")];
    char script_stringpool_str14[sizeof("Ahom")];
    char script_stringpool_str15[sizeof("Miao")];
    char script_stringpool_str16[sizeof("Braille")];
    char script_stringpool_str17[sizeof("Balinese")];
    char script_stringpool_str18[sizeof("Mandaic")];
    char script_stringpool_str19[sizeof("Armenian")];
    char script_stringpool_str20[sizeof("Myanmar")];
    char script_stringpool_str21[sizeof("Mongolian")];
    char script_stringpool_str22[sizeof("Adlam")];
    char script_stringpool_str23[sizeof("Cham")];
    char script_stringpool_str24[sizeof("Multani")];
    char script_stringpool_str25[sizeof("Common")];
    char script_stringpool_str26[sizeof("Thai")];
    char script_stringpool_str27[sizeof("Chakma")];
    char script_stringpool_str28[sizeof("Thaana")];
    char script_stringpool_str29[sizeof("Arabic")];
    char script_stringpool_str30[sizeof("Avestan")];
    char script_stringpool_str31[sizeof("Cherokee")];
    char script_stringpool_str32[sizeof("Tamil")];
    char script_stringpool_str33[sizeof("Lao")];
    char script_stringpool_str34[sizeof("Meetei_Mayek")];
    char script_stringpool_str35[sizeof("Mende_Kikakui")];
    char script_stringpool_str36[sizeof("Sinhala")];
    char script_stringpool_str37[sizeof("Anatolian_Hieroglyphs")];
    char script_stringpool_str38[sizeof("Phoenician")];
    char script_stringpool_str39[sizeof("Sharada")];
    char script_stringpool_str40[sizeof("Linear_B")];
    char script_stringpool_str41[sizeof("Carian")];
    char script_stringpool_str42[sizeof("Batak")];
    char script_stringpool_str43[sizeof("Tangut")];
    char script_stringpool_str45[sizeof("Latin")];
    char script_stringpool_str47[sizeof("Shavian")];
    char script_stringpool_str48[sizeof("Modi")];
    char script_stringpool_str49[sizeof("Linear_A")];
    char script_stringpool_str51[sizeof("Syriac")];
    char script_stringpool_str52[sizeof("Cuneiform")];
    char script_stringpool_str53[sizeof("Osmanya")];
    char script_stringpool_str54[sizeof("Limbu")];
    char script_stringpool_str55[sizeof("Osage")];
    char script_stringpool_str56[sizeof("Samaritan")];
    char script_stringpool_str57[sizeof("Kannada")];
    char script_stringpool_str58[sizeof("Caucasian_Albanian")];
    char script_stringpool_str59[sizeof("Tai_Tham")];
    char script_stringpool_str60[sizeof("Tirhuta")];
    char script_stringpool_str61[sizeof("Takri")];
    char script_stringpool_str62[sizeof("Buginese")];
    char script_stringpool_str63[sizeof("Oriya")];
    char script_stringpool_str64[sizeof("Bhaiksuki")];
    char script_stringpool_str65[sizeof("Tai_Le")];
    char script_stringpool_str66[sizeof("Warang_Citi")];
    char script_stringpool_str67[sizeof("Marchen")];
    char script_stringpool_str68[sizeof("Saurashtra")];
    char script_stringpool_str69[sizeof("Han")];
    char script_stringpool_str70[sizeof("Khmer")];
    char script_stringpool_str71[sizeof("Canadian_Aboriginal")];
    char script_stringpool_str72[sizeof("Kharoshthi")];
    char script_stringpool_str73[sizeof("Hanunoo")];
    char script_stringpool_str74[sizeof("Lydian")];
    char script_stringpool_str75[sizeof("Nko")];
    char script_stringpool_str76[sizeof("Manichaean")];
    char script_stringpool_str77[sizeof("Buhid")];
    char script_stringpool_str78[sizeof("Newa")];
    char script_stringpool_str79[sizeof("Bassa_Vah")];
    char script_stringpool_str80[sizeof("Khojki")];
    char script_stringpool_str81[sizeof("Bopomofo")];
    char script_stringpool_str82[sizeof("Telugu")];
    char script_stringpool_str83[sizeof("Tagalog")];
    char script_stringpool_str84[sizeof("Tagbanwa")];
    char script_stringpool_str85[sizeof("Grantha")];
    char script_stringpool_str86[sizeof("Hatran")];
    char script_stringpool_str87[sizeof("Ogham")];
    char script_stringpool_str88[sizeof("Inherited")];
    char script_stringpool_str89[sizeof("Glagolitic")];
    char script_stringpool_str90[sizeof("Hangul")];
    char script_stringpool_str91[sizeof("Tibetan")];
    char script_stringpool_str92[sizeof("Gothic")];
    char script_stringpool_str93[sizeof("Lycian")];
    char script_stringpool_str94[sizeof("Phags_Pa")];
    char script_stringpool_str95[sizeof("Katakana")];
    char script_stringpool_str96[sizeof("Psalter_Pahlavi")];
    char script_stringpool_str97[sizeof("Lisu")];
    char script_stringpool_str98[sizeof("Greek")];
    char script_stringpool_str99[sizeof("Devanagari")];
    char script_stringpool_str100[sizeof("Kaithi")];
    char script_stringpool_str101[sizeof("Cyrillic")];
    char script_stringpool_str102[sizeof("Sundanese")];
    char script_stringpool_str103[sizeof("Coptic")];
    char script_stringpool_str104[sizeof("Cypriot")];
    char script_stringpool_str105[sizeof("Siddham")];
    char script_stringpool_str106[sizeof("Meroitic_Cursive")];
    char script_stringpool_str107[sizeof("Sora_Sompeng")];
    char script_stringpool_str108[sizeof("Old_Permic")];
    char script_stringpool_str109[sizeof("Malayalam")];
    char script_stringpool_str110[sizeof("Meroitic_Hieroglyphs")];
    char script_stringpool_str111[sizeof("Mahajani")];
    char script_stringpool_str112[sizeof("Pau_Cin_Hau")];
    char script_stringpool_str113[sizeof("Khudawadi")];
    char script_stringpool_str114[sizeof("Palmyrene")];
    char script_stringpool_str115[sizeof("Rejang")];
    char script_stringpool_str116[sizeof("Vai")];
    char script_stringpool_str117[sizeof("Gurmukhi")];
    char script_stringpool_str118[sizeof("Tifinagh")];
    char script_stringpool_str119[sizeof("Duployan")];
    char script_stringpool_str120[sizeof("Old_Italic")];
    char script_stringpool_str121[sizeof("Runic")];
    char script_stringpool_str122[sizeof("SignWriting")];
    char script_stringpool_str123[sizeof("Ugaritic")];
    char script_stringpool_str124[sizeof("Georgian")];
    char script_stringpool_str125[sizeof("Javanese")];
    char script_stringpool_str126[sizeof("Syloti_Nagri")];
    char script_stringpool_str127[sizeof("Deseret")];
    char script_stringpool_str128[sizeof("Ethiopic")];
    char script_stringpool_str129[sizeof("Hebrew")];
    char script_stringpool_str132[sizeof("Imperial_Aramaic")];
    char script_stringpool_str133[sizeof("Pahawh_Hmong")];
    char script_stringpool_str134[sizeof("Old_South_Arabian")];
    char script_stringpool_str135[sizeof("Old_Turkic")];
    char script_stringpool_str136[sizeof("Inscriptional_Pahlavi")];
    char script_stringpool_str137[sizeof("Inscriptional_Parthian")];
    char script_stringpool_str141[sizeof("Lepcha")];
    char script_stringpool_str142[sizeof("Egyptian_Hieroglyphs")];
    char script_stringpool_str144[sizeof("Hiragana")];
    char script_stringpool_str147[sizeof("Gujarati")];
    char script_stringpool_str148[sizeof("Nabataean")];
    char script_stringpool_str150[sizeof("Ol_Chiki")];
    char script_stringpool_str151[sizeof("Tai_Viet")];
    char script_stringpool_str153[sizeof("Elbasan")];
    char script_stringpool_str155[sizeof("New_Tai_Lue")];
    char script_stringpool_str171[sizeof("Old_Persian")];
    char script_stringpool_str179[sizeof("Old_North_Arabian")];
    char script_stringpool_str185[sizeof("Kayah_Li")];
    char script_stringpool_str210[sizeof("Old_Hungarian")];
  };
static const struct script_stringpool_t script_stringpool_contents =
  {
    "Yi",
    "Bamum",
    "Bengali",
    "Brahmi",
    "Mro",
    "Ahom",
    "Miao",
    "Braille",
    "Balinese",
    "Mandaic",
    "Armenian",
    "Myanmar",
    "Mongolian",
    "Adlam",
    "Cham",
    "Multani",
    "Common",
    "Thai",
    "Chakma",
    "Thaana",
    "Arabic",
    "Avestan",
    "Cherokee",
    "Tamil",
    "Lao",
    "Meetei_Mayek",
    "Mende_Kikakui",
    "Sinhala",
    "Anatolian_Hieroglyphs",
    "Phoenician",
    "Sharada",
    "Linear_B",
    "Carian",
    "Batak",
    "Tangut",
    "Latin",
    "Shavian",
    "Modi",
    "Linear_A",
    "Syriac",
    "Cuneiform",
    "Osmanya",
    "Limbu",
    "Osage",
    "Samaritan",
    "Kannada",
    "Caucasian_Albanian",
    "Tai_Tham",
    "Tirhuta",
    "Takri",
    "Buginese",
    "Oriya",
    "Bhaiksuki",
    "Tai_Le",
    "Warang_Citi",
    "Marchen",
    "Saurashtra",
    "Han",
    "Khmer",
    "Canadian_Aboriginal",
    "Kharoshthi",
    "Hanunoo",
    "Lydian",
    "Nko",
    "Manichaean",
    "Buhid",
    "Newa",
    "Bassa_Vah",
    "Khojki",
    "Bopomofo",
    "Telugu",
    "Tagalog",
    "Tagbanwa",
    "Grantha",
    "Hatran",
    "Ogham",
    "Inherited",
    "Glagolitic",
    "Hangul",
    "Tibetan",
    "Gothic",
    "Lycian",
    "Phags_Pa",
    "Katakana",
    "Psalter_Pahlavi",
    "Lisu",
    "Greek",
    "Devanagari",
    "Kaithi",
    "Cyrillic",
    "Sundanese",
    "Coptic",
    "Cypriot",
    "Siddham",
    "Meroitic_Cursive",
    "Sora_Sompeng",
    "Old_Permic",
    "Malayalam",
    "Meroitic_Hieroglyphs",
    "Mahajani",
    "Pau_Cin_Hau",
    "Khudawadi",
    "Palmyrene",
    "Rejang",
    "Vai",
    "Gurmukhi",
    "Tifinagh",
    "Duployan",
    "Old_Italic",
    "Runic",
    "SignWriting",
    "Ugaritic",
    "Georgian",
    "Javanese",
    "Syloti_Nagri",
    "Deseret",
    "Ethiopic",
    "Hebrew",
    "Imperial_Aramaic",
    "Pahawh_Hmong",
    "Old_South_Arabian",
    "Old_Turkic",
    "Inscriptional_Pahlavi",
    "Inscriptional_Parthian",
    "Lepcha",
    "Egyptian_Hieroglyphs",
    "Hiragana",
    "Gujarati",
    "Nabataean",
    "Ol_Chiki",
    "Tai_Viet",
    "Elbasan",
    "New_Tai_Lue",
    "Old_Persian",
    "Old_North_Arabian",
    "Kayah_Li",
    "Old_Hungarian"
  };
#define script_stringpool ((const char *) &script_stringpool_contents)

static const struct named_script script_names[] =
  {
    {-1}, {-1},
#line 51 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str2, 36},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 98 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str9, 83},
#line 25 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str10, 10},
#line 108 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str11, 93},
#line 129 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str12, 114},
    {-1},
#line 140 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str14, 125},
#line 113 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str15, 98},
#line 67 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str16, 52},
#line 76 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str17, 61},
#line 109 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str18, 94},
#line 19 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str19, 4},
#line 37 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str20, 22},
#line 46 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str21, 31},
#line 146 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str22, 131},
#line 91 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str23, 76},
#line 143 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str24, 128},
#line 15 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str25, 0},
#line 34 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str26, 19},
#line 110 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str27, 95},
#line 23 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str28, 8},
#line 21 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str29, 6},
#line 94 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str30, 79},
#line 41 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str31, 26},
#line 29 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str32, 14},
#line 35 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str33, 20},
#line 100 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str34, 85},
#line 127 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str35, 112},
#line 33 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str36, 18},
#line 141 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str37, 126},
#line 78 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str38, 63},
#line 114 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str39, 99},
#line 62 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str40, 47},
#line 89 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str41, 74},
#line 107 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str42, 92},
#line 151 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str43, 136},
    {-1},
#line 16 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str45, 1},
    {-1},
#line 64 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str47, 49},
#line 128 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str48, 113},
#line 124 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str49, 109},
    {-1},
#line 22 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str51, 7},
#line 77 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str52, 62},
#line 65 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str53, 50},
#line 60 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str54, 45},
#line 150 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str55, 135},
#line 96 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str56, 81},
#line 31 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str57, 16},
#line 117 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str58, 102},
#line 92 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str59, 77},
#line 138 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str60, 123},
#line 116 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str61, 101},
#line 68 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str62, 53},
#line 28 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str63, 13},
#line 147 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str64, 132},
#line 61 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str65, 46},
#line 139 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str66, 124},
#line 148 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str67, 133},
#line 85 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str68, 70},
#line 50 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str69, 35},
#line 45 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str70, 30},
#line 42 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str71, 27},
#line 75 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str72, 60},
#line 57 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str73, 42},
#line 90 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str74, 75},
#line 80 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str75, 65},
#line 126 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str76, 111},
#line 58 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str77, 43},
#line 149 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str78, 134},
#line 118 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str79, 103},
#line 123 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str80, 108},
#line 49 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str81, 34},
#line 30 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str82, 15},
#line 56 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str83, 41},
#line 59 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str84, 44},
#line 121 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str85, 106},
#line 142 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str86, 127},
#line 43 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str87, 28},
#line 55 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str88, 40},
#line 71 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str89, 56},
#line 39 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str90, 24},
#line 36 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str91, 21},
#line 53 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str92, 38},
#line 88 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str93, 73},
#line 79 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str94, 64},
#line 48 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str95, 33},
#line 135 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str96, 120},
#line 97 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str97, 82},
#line 17 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str98, 2},
#line 24 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str99, 9},
#line 106 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str100, 91},
#line 18 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str101, 3},
#line 81 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str102, 66},
#line 69 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str103, 54},
#line 66 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str104, 51},
#line 136 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str105, 121},
#line 111 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str106, 96},
#line 115 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str107, 100},
#line 134 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str108, 119},
#line 32 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str109, 17},
#line 112 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str110, 97},
#line 125 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str111, 110},
#line 133 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str112, 118},
#line 137 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str113, 122},
#line 132 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str114, 117},
#line 87 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str115, 72},
#line 84 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str116, 69},
#line 26 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str117, 11},
#line 72 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str118, 57},
#line 119 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str119, 104},
#line 52 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str120, 37},
#line 44 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str121, 29},
#line 145 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str122, 130},
#line 63 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str123, 48},
#line 38 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str124, 23},
#line 99 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str125, 84},
#line 73 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str126, 58},
#line 54 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str127, 39},
#line 40 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str128, 25},
#line 20 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str129, 5},
    {-1}, {-1},
#line 101 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str132, 86},
#line 122 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str133, 107},
#line 102 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str134, 87},
#line 105 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str135, 90},
#line 104 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str136, 89},
#line 103 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str137, 88},
    {-1}, {-1}, {-1},
#line 82 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str141, 67},
#line 95 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str142, 80},
    {-1},
#line 47 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str144, 32},
    {-1}, {-1},
#line 27 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str147, 12},
#line 131 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str148, 116},
    {-1},
#line 83 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str150, 68},
#line 93 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str151, 78},
    {-1},
#line 120 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str153, 105},
    {-1},
#line 70 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str155, 55},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 74 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str171, 59},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 130 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str179, 115},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 86 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str185, 71},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 144 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str210, 129}
  };

const struct named_script *
uc_script_lookup (register const char *str, register size_t len)
{
  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = scripts_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        {
          register int o = script_names[key].name;
          if (o >= 0)
            {
              register const char *s = o + script_stringpool;

              if (*str == *s && !strcmp (str + 1, s + 1))
                return &script_names[key];
            }
        }
    }
  return 0;
}
