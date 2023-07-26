/* ANSI-C code produced by gperf version 3.2 */
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

#line 20 "./unictype/scripts_byname.gperf"
struct named_script { int name; unsigned int index; };

#define TOTAL_KEYWORDS 163
#define MIN_WORD_LENGTH 2
#define MAX_WORD_LENGTH 22
#define MIN_HASH_VALUE 3
#define MAX_HASH_VALUE 249
/* maximum key range = 247, duplicates = 0 */

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
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250,  52,  18,   4,  88,  34,
      250,  64,  80,  37,   3,  21,  32,   1,  39,  29,
       29, 250,  37,   6,  16,   1, 134,  79, 250,   1,
       60, 250, 250, 250, 250,  69, 250,   4,  77, 105,
       29,  32,  47,  48,  81,   6,  31,  10,  50,   3,
        1,   1, 114, 250,   6,  63,  12,  29,  29,  47,
      250,  90,   3, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250, 250, 250, 250, 250,
      250, 250, 250, 250, 250, 250
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[7]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 7:
      case 6:
      case 5:
        hval += asso_values[(unsigned char)str[4]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 4:
      case 3:
        hval += asso_values[(unsigned char)str[2]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
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
    char script_stringpool_str3[sizeof("Yi")];
    char script_stringpool_str5[sizeof("Mro")];
    char script_stringpool_str9[sizeof("Miao")];
    char script_stringpool_str12[sizeof("Cham")];
    char script_stringpool_str13[sizeof("Mandaic")];
    char script_stringpool_str14[sizeof("Common")];
    char script_stringpool_str15[sizeof("Myanmar")];
    char script_stringpool_str16[sizeof("Mongolian")];
    char script_stringpool_str17[sizeof("Chakma")];
    char script_stringpool_str18[sizeof("Sinhala")];
    char script_stringpool_str20[sizeof("Carian")];
    char script_stringpool_str21[sizeof("Sharada")];
    char script_stringpool_str22[sizeof("Syriac")];
    char script_stringpool_str23[sizeof("Shavian")];
    char script_stringpool_str24[sizeof("Thai")];
    char script_stringpool_str25[sizeof("Chorasmian")];
    char script_stringpool_str26[sizeof("Cuneiform")];
    char script_stringpool_str27[sizeof("Thaana")];
    char script_stringpool_str28[sizeof("Samaritan")];
    char script_stringpool_str29[sizeof("Bamum")];
    char script_stringpool_str30[sizeof("Bengali")];
    char script_stringpool_str31[sizeof("Brahmi")];
    char script_stringpool_str32[sizeof("Toto")];
    char script_stringpool_str33[sizeof("Kannada")];
    char script_stringpool_str34[sizeof("Modi")];
    char script_stringpool_str35[sizeof("Khmer")];
    char script_stringpool_str36[sizeof("Lao")];
    char script_stringpool_str37[sizeof("Takri")];
    char script_stringpool_str38[sizeof("Khojki")];
    char script_stringpool_str39[sizeof("Yezidi")];
    char script_stringpool_str40[sizeof("Osmanya")];
    char script_stringpool_str43[sizeof("Nko")];
    char script_stringpool_str44[sizeof("Oriya")];
    char script_stringpool_str45[sizeof("Batak")];
    char script_stringpool_str47[sizeof("Phoenician")];
    char script_stringpool_str48[sizeof("Kharoshthi")];
    char script_stringpool_str49[sizeof("Tai_Tham")];
    char script_stringpool_str50[sizeof("Latin")];
    char script_stringpool_str51[sizeof("Bhaiksuki")];
    char script_stringpool_str52[sizeof("Tangut")];
    char script_stringpool_str53[sizeof("Mende_Kikakui")];
    char script_stringpool_str54[sizeof("Canadian_Aboriginal")];
    char script_stringpool_str55[sizeof("Katakana")];
    char script_stringpool_str56[sizeof("Khitan_Small_Script")];
    char script_stringpool_str57[sizeof("Ahom")];
    char script_stringpool_str58[sizeof("Tirhuta")];
    char script_stringpool_str59[sizeof("Caucasian_Albanian")];
    char script_stringpool_str60[sizeof("Tai_Le")];
    char script_stringpool_str61[sizeof("Saurashtra")];
    char script_stringpool_str62[sizeof("Multani")];
    char script_stringpool_str63[sizeof("Linear_B")];
    char script_stringpool_str65[sizeof("Armenian")];
    char script_stringpool_str67[sizeof("Sogdian")];
    char script_stringpool_str68[sizeof("Arabic")];
    char script_stringpool_str69[sizeof("Limbu")];
    char script_stringpool_str70[sizeof("Osage")];
    char script_stringpool_str71[sizeof("Lydian")];
    char script_stringpool_str72[sizeof("Kawi")];
    char script_stringpool_str73[sizeof("Javanese")];
    char script_stringpool_str74[sizeof("Tamil")];
    char script_stringpool_str75[sizeof("Rejang")];
    char script_stringpool_str77[sizeof("Cherokee")];
    char script_stringpool_str78[sizeof("Meetei_Mayek")];
    char script_stringpool_str79[sizeof("Braille")];
    char script_stringpool_str80[sizeof("Tagbanwa")];
    char script_stringpool_str81[sizeof("Makasar")];
    char script_stringpool_str82[sizeof("Anatolian_Hieroglyphs")];
    char script_stringpool_str83[sizeof("Sundanese")];
    char script_stringpool_str84[sizeof("Han")];
    char script_stringpool_str86[sizeof("Tangsa")];
    char script_stringpool_str87[sizeof("Grantha")];
    char script_stringpool_str88[sizeof("Gothic")];
    char script_stringpool_str89[sizeof("Hanunoo")];
    char script_stringpool_str90[sizeof("Newa")];
    char script_stringpool_str91[sizeof("Glagolitic")];
    char script_stringpool_str92[sizeof("Khudawadi")];
    char script_stringpool_str93[sizeof("Old_South_Arabian")];
    char script_stringpool_str94[sizeof("Old_Turkic")];
    char script_stringpool_str95[sizeof("Marchen")];
    char script_stringpool_str96[sizeof("Sora_Sompeng")];
    char script_stringpool_str97[sizeof("Linear_A")];
    char script_stringpool_str98[sizeof("Bassa_Vah")];
    char script_stringpool_str99[sizeof("Lisu")];
    char script_stringpool_str100[sizeof("Old_Permic")];
    char script_stringpool_str101[sizeof("Warang_Citi")];
    char script_stringpool_str102[sizeof("Hatran")];
    char script_stringpool_str103[sizeof("Avestan")];
    char script_stringpool_str104[sizeof("Old_Sogdian")];
    char script_stringpool_str105[sizeof("Nandinagari")];
    char script_stringpool_str106[sizeof("Soyombo")];
    char script_stringpool_str107[sizeof("Buginese")];
    char script_stringpool_str108[sizeof("Phags_Pa")];
    char script_stringpool_str109[sizeof("Balinese")];
    char script_stringpool_str110[sizeof("Adlam")];
    char script_stringpool_str111[sizeof("Greek")];
    char script_stringpool_str112[sizeof("Tibetan")];
    char script_stringpool_str113[sizeof("Gurmukhi")];
    char script_stringpool_str114[sizeof("Kaithi")];
    char script_stringpool_str115[sizeof("Gujarati")];
    char script_stringpool_str116[sizeof("Hangul")];
    char script_stringpool_str117[sizeof("Medefaidrin")];
    char script_stringpool_str118[sizeof("Ogham")];
    char script_stringpool_str119[sizeof("Syloti_Nagri")];
    char script_stringpool_str120[sizeof("Telugu")];
    char script_stringpool_str121[sizeof("Tagalog")];
    char script_stringpool_str122[sizeof("Georgian")];
    char script_stringpool_str123[sizeof("Siddham")];
    char script_stringpool_str124[sizeof("Ugaritic")];
    char script_stringpool_str125[sizeof("Nyiakeng_Puachue_Hmong")];
    char script_stringpool_str126[sizeof("Old_North_Arabian")];
    char script_stringpool_str127[sizeof("Mahajani")];
    char script_stringpool_str128[sizeof("Nag_Mundari")];
    char script_stringpool_str129[sizeof("Psalter_Pahlavi")];
    char script_stringpool_str130[sizeof("Coptic")];
    char script_stringpool_str131[sizeof("Cypriot")];
    char script_stringpool_str132[sizeof("Devanagari")];
    char script_stringpool_str133[sizeof("Buhid")];
    char script_stringpool_str134[sizeof("Meroitic_Cursive")];
    char script_stringpool_str135[sizeof("Elymaic")];
    char script_stringpool_str136[sizeof("Nushu")];
    char script_stringpool_str137[sizeof("Cypro_Minoan")];
    char script_stringpool_str138[sizeof("Meroitic_Hieroglyphs")];
    char script_stringpool_str139[sizeof("Inscriptional_Pahlavi")];
    char script_stringpool_str140[sizeof("Inscriptional_Parthian")];
    char script_stringpool_str141[sizeof("Nabataean")];
    char script_stringpool_str142[sizeof("Pau_Cin_Hau")];
    char script_stringpool_str143[sizeof("Vai")];
    char script_stringpool_str144[sizeof("Bopomofo")];
    char script_stringpool_str145[sizeof("Dogra")];
    char script_stringpool_str146[sizeof("Hiragana")];
    char script_stringpool_str147[sizeof("Lycian")];
    char script_stringpool_str148[sizeof("Runic")];
    char script_stringpool_str149[sizeof("Manichaean")];
    char script_stringpool_str150[sizeof("Old_Uyghur")];
    char script_stringpool_str151[sizeof("Gunjala_Gondi")];
    char script_stringpool_str152[sizeof("Masaram_Gondi")];
    char script_stringpool_str153[sizeof("Tifinagh")];
    char script_stringpool_str154[sizeof("Malayalam")];
    char script_stringpool_str155[sizeof("Old_Italic")];
    char script_stringpool_str156[sizeof("SignWriting")];
    char script_stringpool_str157[sizeof("Egyptian_Hieroglyphs")];
    char script_stringpool_str158[sizeof("Zanabazar_Square")];
    char script_stringpool_str161[sizeof("Old_Persian")];
    char script_stringpool_str164[sizeof("Deseret")];
    char script_stringpool_str165[sizeof("Inherited")];
    char script_stringpool_str167[sizeof("Wancho")];
    char script_stringpool_str170[sizeof("Vithkuqi")];
    char script_stringpool_str173[sizeof("Cyrillic")];
    char script_stringpool_str176[sizeof("Tai_Viet")];
    char script_stringpool_str179[sizeof("Palmyrene")];
    char script_stringpool_str180[sizeof("Hanifi_Rohingya")];
    char script_stringpool_str181[sizeof("Elbasan")];
    char script_stringpool_str182[sizeof("New_Tai_Lue")];
    char script_stringpool_str193[sizeof("Ol_Chiki")];
    char script_stringpool_str195[sizeof("Hebrew")];
    char script_stringpool_str199[sizeof("Old_Hungarian")];
    char script_stringpool_str201[sizeof("Dives_Akuru")];
    char script_stringpool_str206[sizeof("Kayah_Li")];
    char script_stringpool_str212[sizeof("Duployan")];
    char script_stringpool_str223[sizeof("Imperial_Aramaic")];
    char script_stringpool_str229[sizeof("Ethiopic")];
    char script_stringpool_str233[sizeof("Lepcha")];
    char script_stringpool_str249[sizeof("Pahawh_Hmong")];
  };
static const struct script_stringpool_t script_stringpool_contents =
  {
    "Yi",
    "Mro",
    "Miao",
    "Cham",
    "Mandaic",
    "Common",
    "Myanmar",
    "Mongolian",
    "Chakma",
    "Sinhala",
    "Carian",
    "Sharada",
    "Syriac",
    "Shavian",
    "Thai",
    "Chorasmian",
    "Cuneiform",
    "Thaana",
    "Samaritan",
    "Bamum",
    "Bengali",
    "Brahmi",
    "Toto",
    "Kannada",
    "Modi",
    "Khmer",
    "Lao",
    "Takri",
    "Khojki",
    "Yezidi",
    "Osmanya",
    "Nko",
    "Oriya",
    "Batak",
    "Phoenician",
    "Kharoshthi",
    "Tai_Tham",
    "Latin",
    "Bhaiksuki",
    "Tangut",
    "Mende_Kikakui",
    "Canadian_Aboriginal",
    "Katakana",
    "Khitan_Small_Script",
    "Ahom",
    "Tirhuta",
    "Caucasian_Albanian",
    "Tai_Le",
    "Saurashtra",
    "Multani",
    "Linear_B",
    "Armenian",
    "Sogdian",
    "Arabic",
    "Limbu",
    "Osage",
    "Lydian",
    "Kawi",
    "Javanese",
    "Tamil",
    "Rejang",
    "Cherokee",
    "Meetei_Mayek",
    "Braille",
    "Tagbanwa",
    "Makasar",
    "Anatolian_Hieroglyphs",
    "Sundanese",
    "Han",
    "Tangsa",
    "Grantha",
    "Gothic",
    "Hanunoo",
    "Newa",
    "Glagolitic",
    "Khudawadi",
    "Old_South_Arabian",
    "Old_Turkic",
    "Marchen",
    "Sora_Sompeng",
    "Linear_A",
    "Bassa_Vah",
    "Lisu",
    "Old_Permic",
    "Warang_Citi",
    "Hatran",
    "Avestan",
    "Old_Sogdian",
    "Nandinagari",
    "Soyombo",
    "Buginese",
    "Phags_Pa",
    "Balinese",
    "Adlam",
    "Greek",
    "Tibetan",
    "Gurmukhi",
    "Kaithi",
    "Gujarati",
    "Hangul",
    "Medefaidrin",
    "Ogham",
    "Syloti_Nagri",
    "Telugu",
    "Tagalog",
    "Georgian",
    "Siddham",
    "Ugaritic",
    "Nyiakeng_Puachue_Hmong",
    "Old_North_Arabian",
    "Mahajani",
    "Nag_Mundari",
    "Psalter_Pahlavi",
    "Coptic",
    "Cypriot",
    "Devanagari",
    "Buhid",
    "Meroitic_Cursive",
    "Elymaic",
    "Nushu",
    "Cypro_Minoan",
    "Meroitic_Hieroglyphs",
    "Inscriptional_Pahlavi",
    "Inscriptional_Parthian",
    "Nabataean",
    "Pau_Cin_Hau",
    "Vai",
    "Bopomofo",
    "Dogra",
    "Hiragana",
    "Lycian",
    "Runic",
    "Manichaean",
    "Old_Uyghur",
    "Gunjala_Gondi",
    "Masaram_Gondi",
    "Tifinagh",
    "Malayalam",
    "Old_Italic",
    "SignWriting",
    "Egyptian_Hieroglyphs",
    "Zanabazar_Square",
    "Old_Persian",
    "Deseret",
    "Inherited",
    "Wancho",
    "Vithkuqi",
    "Cyrillic",
    "Tai_Viet",
    "Palmyrene",
    "Hanifi_Rohingya",
    "Elbasan",
    "New_Tai_Lue",
    "Ol_Chiki",
    "Hebrew",
    "Old_Hungarian",
    "Dives_Akuru",
    "Kayah_Li",
    "Duployan",
    "Imperial_Aramaic",
    "Ethiopic",
    "Lepcha",
    "Pahawh_Hmong"
  };
#define script_stringpool ((const char *) &script_stringpool_contents)

static const struct named_script script_names[] =
  {
    {-1}, {-1}, {-1},
#line 67 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str3, 36},
    {-1},
#line 145 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str5, 114},
    {-1}, {-1}, {-1},
#line 129 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str9, 98},
    {-1}, {-1},
#line 107 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str12, 76},
#line 125 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str13, 94},
#line 31 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str14, 0},
#line 53 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str15, 22},
#line 62 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str16, 31},
#line 126 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str17, 95},
#line 49 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str18, 18},
    {-1},
#line 105 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str20, 74},
#line 130 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str21, 99},
#line 38 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str22, 7},
#line 80 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str23, 49},
#line 50 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str24, 19},
#line 183 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str25, 152},
#line 93 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str26, 62},
#line 39 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str27, 8},
#line 112 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str28, 81},
#line 114 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str29, 83},
#line 41 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str30, 10},
#line 124 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str31, 93},
#line 190 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str32, 159},
#line 47 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str33, 16},
#line 144 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str34, 113},
#line 61 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str35, 30},
#line 51 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str36, 20},
#line 132 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str37, 101},
#line 139 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str38, 108},
#line 186 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str39, 155},
#line 81 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str40, 50},
    {-1}, {-1},
#line 96 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str43, 65},
#line 44 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str44, 13},
#line 123 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str45, 92},
    {-1},
#line 94 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str47, 63},
#line 91 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str48, 60},
#line 108 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str49, 77},
#line 32 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str50, 1},
#line 163 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str51, 132},
#line 167 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str52, 136},
#line 143 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str53, 112},
#line 58 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str54, 27},
#line 64 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str55, 33},
#line 185 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str56, 154},
#line 156 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str57, 125},
#line 154 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str58, 123},
#line 133 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str59, 102},
#line 77 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str60, 46},
#line 101 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str61, 70},
#line 159 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str62, 128},
#line 78 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str63, 47},
    {-1},
#line 35 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str65, 4},
    {-1},
#line 177 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str67, 146},
#line 37 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str68, 6},
#line 76 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str69, 45},
#line 166 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str70, 135},
#line 106 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str71, 75},
#line 192 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str72, 161},
#line 115 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str73, 84},
#line 45 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str74, 14},
#line 103 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str75, 72},
    {-1},
#line 57 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str77, 26},
#line 116 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str78, 85},
#line 83 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str79, 52},
#line 75 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str80, 44},
#line 174 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str81, 143},
#line 157 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str82, 126},
#line 97 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str83, 66},
#line 66 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str84, 35},
    {-1},
#line 189 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str86, 158},
#line 137 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str87, 106},
#line 69 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str88, 38},
#line 73 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str89, 42},
#line 165 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str90, 134},
#line 87 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str91, 56},
#line 153 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str92, 122},
#line 118 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str93, 87},
#line 121 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str94, 90},
#line 164 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str95, 133},
#line 131 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str96, 100},
#line 140 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str97, 109},
#line 134 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str98, 103},
#line 113 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str99, 82},
#line 150 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str100, 119},
#line 155 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str101, 124},
#line 158 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str102, 127},
#line 110 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str103, 79},
#line 178 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str104, 147},
#line 180 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str105, 149},
#line 170 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str106, 139},
#line 84 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str107, 53},
#line 95 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str108, 64},
#line 92 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str109, 61},
#line 162 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str110, 131},
#line 33 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str111, 2},
#line 52 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str112, 21},
#line 42 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str113, 11},
#line 122 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str114, 91},
#line 43 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str115, 12},
#line 55 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str116, 24},
#line 175 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str117, 144},
#line 59 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str118, 28},
#line 89 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str119, 58},
#line 46 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str120, 15},
#line 72 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str121, 41},
#line 54 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str122, 23},
#line 152 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str123, 121},
#line 79 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str124, 48},
#line 181 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str125, 150},
#line 146 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str126, 115},
#line 141 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str127, 110},
#line 193 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str128, 162},
#line 151 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str129, 120},
#line 85 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str130, 54},
#line 82 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str131, 51},
#line 40 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str132, 9},
#line 74 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str133, 43},
#line 127 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str134, 96},
#line 179 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str135, 148},
#line 169 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str136, 138},
#line 187 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str137, 156},
#line 128 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str138, 97},
#line 120 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str139, 89},
#line 119 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str140, 88},
#line 147 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str141, 116},
#line 149 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str142, 118},
#line 100 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str143, 69},
#line 65 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str144, 34},
#line 172 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str145, 141},
#line 63 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str146, 32},
#line 104 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str147, 73},
#line 60 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str148, 29},
#line 142 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str149, 111},
#line 188 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str150, 157},
#line 173 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str151, 142},
#line 168 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str152, 137},
#line 88 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str153, 57},
#line 48 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str154, 17},
#line 68 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str155, 37},
#line 161 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str156, 130},
#line 111 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str157, 80},
#line 171 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str158, 140},
    {-1}, {-1},
#line 90 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str161, 59},
    {-1}, {-1},
#line 70 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str164, 39},
#line 71 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str165, 40},
    {-1},
#line 182 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str167, 151},
    {-1}, {-1},
#line 191 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str170, 160},
    {-1}, {-1},
#line 34 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str173, 3},
    {-1}, {-1},
#line 109 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str176, 78},
    {-1}, {-1},
#line 148 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str179, 117},
#line 176 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str180, 145},
#line 136 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str181, 105},
#line 86 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str182, 55},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1},
#line 99 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str193, 68},
    {-1},
#line 36 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str195, 5},
    {-1}, {-1}, {-1},
#line 160 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str199, 129},
    {-1},
#line 184 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str201, 153},
    {-1}, {-1}, {-1}, {-1},
#line 102 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str206, 71},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 135 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str212, 104},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1},
#line 117 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str223, 86},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 56 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str229, 25},
    {-1}, {-1}, {-1},
#line 98 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str233, 67},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 138 "./unictype/scripts_byname.gperf"
    {(int)(size_t)&((struct script_stringpool_t *)0)->script_stringpool_str249, 107}
  };

static const struct named_script *
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
