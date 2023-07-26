/* ANSI-C code produced by gperf version 3.2 */
/* Command-line: gperf -m 10 ./unicase/locale-languages.gperf  */
/* Computed positions: -k'1-3' */

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


#define TOTAL_KEYWORDS 260
#define MIN_WORD_LENGTH 2
#define MAX_WORD_LENGTH 3
#define MIN_HASH_VALUE 4
#define MAX_HASH_VALUE 461
/* maximum key range = 458, duplicates = 0 */

#ifdef __GNUC__
__inline
#else
#ifdef __cplusplus
inline
#endif
#endif
static unsigned int
uc_locale_language_hash (register const char *str, register size_t len)
{
  static const unsigned short asso_values[] =
    {
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462,   4,
       48,  12,  35, 124, 222, 107, 187, 191, 145, 241,
       11,  73,   0, 114,   4,  66, 213,  88,   2,  25,
       31, 209,  10, 160, 157, 154, 247, 175,  96, 462,
        0,  34,  49, 227,  52, 187, 245, 462, 207, 231,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462, 462, 462, 462, 462, 462, 462, 462,
      462, 462, 462
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[2]+2];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 2:
        hval += asso_values[(unsigned char)str[1]+17];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 1:
        hval += asso_values[(unsigned char)str[0]+3];
        break;
    }
  return hval;
}

static const char *
uc_locale_languages_lookup (register const char *str, register size_t len)
{
  struct stringpool_t
    {
      char stringpool_str4[sizeof("sq")];
      char stringpool_str5[sizeof("men")];
      char stringpool_str6[sizeof("se")];
      char stringpool_str7[sizeof("man")];
      char stringpool_str8[sizeof("sa")];
      char stringpool_str11[sizeof("sat")];
      char stringpool_str12[sizeof("mi")];
      char stringpool_str13[sizeof("min")];
      char stringpool_str14[sizeof("si")];
      char stringpool_str15[sizeof("wen")];
      char stringpool_str16[sizeof("be")];
      char stringpool_str17[sizeof("ka")];
      char stringpool_str18[sizeof("ba")];
      char stringpool_str19[sizeof("ban")];
      char stringpool_str23[sizeof("ki")];
      char stringpool_str24[sizeof("bi")];
      char stringpool_str25[sizeof("bin")];
      char stringpool_str28[sizeof("wal")];
      char stringpool_str29[sizeof("te")];
      char stringpool_str30[sizeof("bal")];
      char stringpool_str31[sizeof("ta")];
      char stringpool_str32[sizeof("tet")];
      char stringpool_str33[sizeof("mg")];
      char stringpool_str35[sizeof("sg")];
      char stringpool_str36[sizeof("mr")];
      char stringpool_str37[sizeof("ti")];
      char stringpool_str38[sizeof("sr")];
      char stringpool_str39[sizeof("ce")];
      char stringpool_str40[sizeof("or")];
      char stringpool_str41[sizeof("ca")];
      char stringpool_str42[sizeof("mad")];
      char stringpool_str44[sizeof("kg")];
      char stringpool_str45[sizeof("bg")];
      char stringpool_str47[sizeof("kr")];
      char stringpool_str48[sizeof("br")];
      char stringpool_str50[sizeof("sid")];
      char stringpool_str51[sizeof("ms")];
      char stringpool_str52[sizeof("ae")];
      char stringpool_str53[sizeof("ss")];
      char stringpool_str54[sizeof("aa")];
      char stringpool_str55[sizeof("os")];
      char stringpool_str56[sizeof("su")];
      char stringpool_str58[sizeof("tg")];
      char stringpool_str61[sizeof("tr")];
      char stringpool_str62[sizeof("ks")];
      char stringpool_str63[sizeof("bs")];
      char stringpool_str64[sizeof("ug")];
      char stringpool_str65[sizeof("ku")];
      char stringpool_str66[sizeof("kab")];
      char stringpool_str67[sizeof("ur")];
      char stringpool_str69[sizeof("tiv")];
      char stringpool_str71[sizeof("cr")];
      char stringpool_str72[sizeof("pa")];
      char stringpool_str73[sizeof("kru")];
      char stringpool_str75[sizeof("af")];
      char stringpool_str76[sizeof("ts")];
      char stringpool_str77[sizeof("pap")];
      char stringpool_str78[sizeof("pi")];
      char stringpool_str79[sizeof("la")];
      char stringpool_str84[sizeof("ar")];
      char stringpool_str85[sizeof("li")];
      char stringpool_str86[sizeof("cs")];
      char stringpool_str88[sizeof("ceb")];
      char stringpool_str89[sizeof("cu")];
      char stringpool_str90[sizeof("bem")];
      char stringpool_str91[sizeof("kam")];
      char stringpool_str92[sizeof("sd")];
      char stringpool_str97[sizeof("sas")];
      char stringpool_str98[sizeof("mo")];
      char stringpool_str99[sizeof("as")];
      char stringpool_str100[sizeof("so")];
      char stringpool_str102[sizeof("ast")];
      char stringpool_str103[sizeof("tem")];
      char stringpool_str106[sizeof("lg")];
      char stringpool_str108[sizeof("wo")];
      char stringpool_str109[sizeof("ko")];
      char stringpool_str110[sizeof("bo")];
      char stringpool_str113[sizeof("fa")];
      char stringpool_str114[sizeof("mag")];
      char stringpool_str115[sizeof("kbd")];
      char stringpool_str116[sizeof("ab")];
      char stringpool_str117[sizeof("ps")];
      char stringpool_str118[sizeof("ne")];
      char stringpool_str119[sizeof("fi")];
      char stringpool_str120[sizeof("na")];
      char stringpool_str123[sizeof("to")];
      char stringpool_str125[sizeof("nap")];
      char stringpool_str127[sizeof("lu")];
      char stringpool_str128[sizeof("de")];
      char stringpool_str130[sizeof("da")];
      char stringpool_str131[sizeof("fil")];
      char stringpool_str132[sizeof("lua")];
      char stringpool_str133[sizeof("co")];
      char stringpool_str134[sizeof("ff")];
      char stringpool_str135[sizeof("csb")];
      char stringpool_str137[sizeof("din")];
      char stringpool_str141[sizeof("lb")];
      char stringpool_str142[sizeof("ru")];
      char stringpool_str143[sizeof("fr")];
      char stringpool_str145[sizeof("sus")];
      char stringpool_str146[sizeof("pam")];
      char stringpool_str147[sizeof("ng")];
      char stringpool_str149[sizeof("ie")];
      char stringpool_str150[sizeof("nr")];
      char stringpool_str151[sizeof("ia")];
      char stringpool_str153[sizeof("tum")];
      char stringpool_str156[sizeof("ml")];
      char stringpool_str157[sizeof("ii")];
      char stringpool_str158[sizeof("sl")];
      char stringpool_str159[sizeof("mk")];
      char stringpool_str160[sizeof("za")];
      char stringpool_str161[sizeof("sk")];
      char stringpool_str162[sizeof("bej")];
      char stringpool_str163[sizeof("kaj")];
      char stringpool_str165[sizeof("zap")];
      char stringpool_str166[sizeof("oj")];
      char stringpool_str167[sizeof("kl")];
      char stringpool_str169[sizeof("yi")];
      char stringpool_str170[sizeof("kk")];
      char stringpool_str171[sizeof("lo")];
      char stringpool_str173[sizeof("kj")];
      char stringpool_str174[sizeof("bug")];
      char stringpool_str175[sizeof("kdm")];
      char stringpool_str177[sizeof("mn")];
      char stringpool_str178[sizeof("ig")];
      char stringpool_str179[sizeof("sn")];
      char stringpool_str180[sizeof("pag")];
      char stringpool_str181[sizeof("tl")];
      char stringpool_str182[sizeof("nb")];
      char stringpool_str183[sizeof("pbb")];
      char stringpool_str184[sizeof("tk")];
      char stringpool_str186[sizeof("ro")];
      char stringpool_str187[sizeof("mos")];
      char stringpool_str188[sizeof("kn")];
      char stringpool_str189[sizeof("bn")];
      char stringpool_str190[sizeof("uk")];
      char stringpool_str191[sizeof("sv")];
      char stringpool_str193[sizeof("ga")];
      char stringpool_str195[sizeof("he")];
      char stringpool_str196[sizeof("is")];
      char stringpool_str197[sizeof("ha")];
      char stringpool_str198[sizeof("mai")];
      char stringpool_str199[sizeof("iu")];
      char stringpool_str200[sizeof("kv")];
      char stringpool_str202[sizeof("tn")];
      char stringpool_str203[sizeof("hi")];
      char stringpool_str204[sizeof("nd")];
      char stringpool_str205[sizeof("fo")];
      char stringpool_str206[sizeof("fon")];
      char stringpool_str207[sizeof("ak")];
      char stringpool_str208[sizeof("zu")];
      char stringpool_str209[sizeof("my")];
      char stringpool_str210[sizeof("myn")];
      char stringpool_str211[sizeof("mh")];
      char stringpool_str212[sizeof("no")];
      char stringpool_str213[sizeof("ve")];
      char stringpool_str214[sizeof("shn")];
      char stringpool_str215[sizeof("hil")];
      char stringpool_str217[sizeof("sc")];
      char stringpool_str218[sizeof("scn")];
      char stringpool_str219[sizeof("oc")];
      char stringpool_str220[sizeof("ky")];
      char stringpool_str221[sizeof("vi")];
      char stringpool_str222[sizeof("pl")];
      char stringpool_str223[sizeof("bh")];
      char stringpool_str224[sizeof("cv")];
      char stringpool_str225[sizeof("an")];
      char stringpool_str226[sizeof("ee")];
      char stringpool_str227[sizeof("hr")];
      char stringpool_str229[sizeof("mt")];
      char stringpool_str230[sizeof("war")];
      char stringpool_str231[sizeof("st")];
      char stringpool_str234[sizeof("ty")];
      char stringpool_str235[sizeof("id")];
      char stringpool_str236[sizeof("th")];
      char stringpool_str237[sizeof("av")];
      char stringpool_str240[sizeof("raj")];
      char stringpool_str241[sizeof("gu")];
      char stringpool_str242[sizeof("luo")];
      char stringpool_str244[sizeof("cy")];
      char stringpool_str245[sizeof("hu")];
      char stringpool_str246[sizeof("ch")];
      char stringpool_str247[sizeof("ja")];
      char stringpool_str248[sizeof("mak")];
      char stringpool_str249[sizeof("sw")];
      char stringpool_str250[sizeof("ln")];
      char stringpool_str251[sizeof("sm")];
      char stringpool_str252[sizeof("srr")];
      char stringpool_str253[sizeof("om")];
      char stringpool_str254[sizeof("tt")];
      char stringpool_str255[sizeof("yo")];
      char stringpool_str257[sizeof("ay")];
      char stringpool_str258[sizeof("kw")];
      char stringpool_str259[sizeof("crh")];
      char stringpool_str260[sizeof("km")];
      char stringpool_str261[sizeof("bm")];
      char stringpool_str262[sizeof("lv")];
      char stringpool_str264[sizeof("uz")];
      char stringpool_str265[sizeof("rn")];
      char stringpool_str266[sizeof("bik")];
      char stringpool_str267[sizeof("qu")];
      char stringpool_str269[sizeof("fj")];
      char stringpool_str270[sizeof("nl")];
      char stringpool_str272[sizeof("tw")];
      char stringpool_str273[sizeof("es")];
      char stringpool_str276[sizeof("eu")];
      char stringpool_str277[sizeof("gd")];
      char stringpool_str278[sizeof("yao")];
      char stringpool_str280[sizeof("nso")];
      char stringpool_str281[sizeof("az")];
      char stringpool_str286[sizeof("gon")];
      char stringpool_str289[sizeof("ho")];
      char stringpool_str291[sizeof("nn")];
      char stringpool_str293[sizeof("nds")];
      char stringpool_str295[sizeof("pt")];
      char stringpool_str296[sizeof("jab")];
      char stringpool_str297[sizeof("am")];
      char stringpool_str298[sizeof("suk")];
      char stringpool_str300[sizeof("awa")];
      char stringpool_str302[sizeof("lt")];
      char stringpool_str303[sizeof("nv")];
      char stringpool_str304[sizeof("ik")];
      char stringpool_str307[sizeof("vo")];
      char stringpool_str308[sizeof("nah")];
      char stringpool_str309[sizeof("kmb")];
      char stringpool_str313[sizeof("dv")];
      char stringpool_str316[sizeof("fy")];
      char stringpool_str320[sizeof("eo")];
      char stringpool_str323[sizeof("ny")];
      char stringpool_str324[sizeof("nyn")];
      char stringpool_str329[sizeof("umb")];
      char stringpool_str333[sizeof("ang")];
      char stringpool_str334[sizeof("kcg")];
      char stringpool_str335[sizeof("rw")];
      char stringpool_str337[sizeof("rm")];
      char stringpool_str338[sizeof("bho")];
      char stringpool_str343[sizeof("gl")];
      char stringpool_str351[sizeof("kok")];
      char stringpool_str357[sizeof("dz")];
      char stringpool_str364[sizeof("gn")];
      char stringpool_str365[sizeof("zh")];
      char stringpool_str369[sizeof("mni")];
      char stringpool_str371[sizeof("xh")];
      char stringpool_str374[sizeof("it")];
      char stringpool_str375[sizeof("fur")];
      char stringpool_str376[sizeof("gv")];
      char stringpool_str378[sizeof("el")];
      char stringpool_str388[sizeof("ace")];
      char stringpool_str397[sizeof("nym")];
      char stringpool_str399[sizeof("en")];
      char stringpool_str400[sizeof("hy")];
      char stringpool_str414[sizeof("doi")];
      char stringpool_str416[sizeof("ilo")];
      char stringpool_str420[sizeof("ht")];
      char stringpool_str424[sizeof("hz")];
      char stringpool_str430[sizeof("jv")];
      char stringpool_str441[sizeof("hmn")];
      char stringpool_str448[sizeof("gsw")];
      char stringpool_str451[sizeof("et")];
      char stringpool_str461[sizeof("mwr")];
    };
  static const struct stringpool_t stringpool_contents =
    {
      "sq",
      "men",
      "se",
      "man",
      "sa",
      "sat",
      "mi",
      "min",
      "si",
      "wen",
      "be",
      "ka",
      "ba",
      "ban",
      "ki",
      "bi",
      "bin",
      "wal",
      "te",
      "bal",
      "ta",
      "tet",
      "mg",
      "sg",
      "mr",
      "ti",
      "sr",
      "ce",
      "or",
      "ca",
      "mad",
      "kg",
      "bg",
      "kr",
      "br",
      "sid",
      "ms",
      "ae",
      "ss",
      "aa",
      "os",
      "su",
      "tg",
      "tr",
      "ks",
      "bs",
      "ug",
      "ku",
      "kab",
      "ur",
      "tiv",
      "cr",
      "pa",
      "kru",
      "af",
      "ts",
      "pap",
      "pi",
      "la",
      "ar",
      "li",
      "cs",
      "ceb",
      "cu",
      "bem",
      "kam",
      "sd",
      "sas",
      "mo",
      "as",
      "so",
      "ast",
      "tem",
      "lg",
      "wo",
      "ko",
      "bo",
      "fa",
      "mag",
      "kbd",
      "ab",
      "ps",
      "ne",
      "fi",
      "na",
      "to",
      "nap",
      "lu",
      "de",
      "da",
      "fil",
      "lua",
      "co",
      "ff",
      "csb",
      "din",
      "lb",
      "ru",
      "fr",
      "sus",
      "pam",
      "ng",
      "ie",
      "nr",
      "ia",
      "tum",
      "ml",
      "ii",
      "sl",
      "mk",
      "za",
      "sk",
      "bej",
      "kaj",
      "zap",
      "oj",
      "kl",
      "yi",
      "kk",
      "lo",
      "kj",
      "bug",
      "kdm",
      "mn",
      "ig",
      "sn",
      "pag",
      "tl",
      "nb",
      "pbb",
      "tk",
      "ro",
      "mos",
      "kn",
      "bn",
      "uk",
      "sv",
      "ga",
      "he",
      "is",
      "ha",
      "mai",
      "iu",
      "kv",
      "tn",
      "hi",
      "nd",
      "fo",
      "fon",
      "ak",
      "zu",
      "my",
      "myn",
      "mh",
      "no",
      "ve",
      "shn",
      "hil",
      "sc",
      "scn",
      "oc",
      "ky",
      "vi",
      "pl",
      "bh",
      "cv",
      "an",
      "ee",
      "hr",
      "mt",
      "war",
      "st",
      "ty",
      "id",
      "th",
      "av",
      "raj",
      "gu",
      "luo",
      "cy",
      "hu",
      "ch",
      "ja",
      "mak",
      "sw",
      "ln",
      "sm",
      "srr",
      "om",
      "tt",
      "yo",
      "ay",
      "kw",
      "crh",
      "km",
      "bm",
      "lv",
      "uz",
      "rn",
      "bik",
      "qu",
      "fj",
      "nl",
      "tw",
      "es",
      "eu",
      "gd",
      "yao",
      "nso",
      "az",
      "gon",
      "ho",
      "nn",
      "nds",
      "pt",
      "jab",
      "am",
      "suk",
      "awa",
      "lt",
      "nv",
      "ik",
      "vo",
      "nah",
      "kmb",
      "dv",
      "fy",
      "eo",
      "ny",
      "nyn",
      "umb",
      "ang",
      "kcg",
      "rw",
      "rm",
      "bho",
      "gl",
      "kok",
      "dz",
      "gn",
      "zh",
      "mni",
      "xh",
      "it",
      "fur",
      "gv",
      "el",
      "ace",
      "nym",
      "en",
      "hy",
      "doi",
      "ilo",
      "ht",
      "hz",
      "jv",
      "hmn",
      "gsw",
      "et",
      "mwr"
    };
  #define stringpool ((const char *) &stringpool_contents)
  static const unsigned char lengthtable[] =
    {
       0,  0,  0,  0,  2,  3,  2,  3,  2,  0,  0,  3,  2,  3,
       2,  3,  2,  2,  2,  3,  0,  0,  0,  2,  2,  3,  0,  0,
       3,  2,  3,  2,  3,  2,  0,  2,  2,  2,  2,  2,  2,  2,
       3,  0,  2,  2,  0,  2,  2,  0,  3,  2,  2,  2,  2,  2,
       2,  0,  2,  0,  0,  2,  2,  2,  2,  2,  3,  2,  0,  3,
       0,  2,  2,  3,  0,  2,  2,  3,  2,  2,  0,  0,  0,  0,
       2,  2,  2,  0,  3,  2,  3,  3,  2,  0,  0,  0,  0,  3,
       2,  2,  2,  0,  3,  3,  0,  0,  2,  0,  2,  2,  2,  0,
       0,  2,  3,  3,  2,  2,  2,  2,  2,  0,  0,  2,  0,  3,
       0,  2,  2,  0,  2,  3,  3,  2,  2,  3,  0,  3,  0,  0,
       0,  2,  2,  2,  0,  3,  3,  2,  0,  2,  2,  2,  0,  3,
       0,  0,  2,  2,  2,  2,  2,  2,  3,  3,  0,  3,  2,  2,
       0,  2,  2,  2,  0,  2,  3,  3,  0,  2,  2,  2,  3,  2,
       2,  3,  2,  0,  2,  3,  2,  2,  2,  2,  0,  2,  0,  2,
       2,  2,  3,  2,  2,  0,  2,  2,  2,  2,  3,  2,  2,  2,
       3,  2,  2,  2,  3,  3,  0,  2,  3,  2,  2,  2,  2,  2,
       2,  2,  2,  2,  0,  2,  3,  2,  0,  0,  2,  2,  2,  2,
       0,  0,  3,  2,  3,  0,  2,  2,  2,  2,  3,  2,  2,  2,
       3,  2,  2,  2,  0,  2,  2,  3,  2,  2,  2,  0,  2,  2,
       3,  2,  0,  2,  2,  0,  2,  2,  0,  0,  2,  2,  3,  0,
       3,  2,  0,  0,  0,  0,  3,  0,  0,  2,  0,  2,  0,  3,
       0,  2,  3,  2,  3,  0,  3,  0,  2,  2,  2,  0,  0,  2,
       3,  3,  0,  0,  0,  2,  0,  0,  2,  0,  0,  0,  2,  0,
       0,  2,  3,  0,  0,  0,  0,  3,  0,  0,  0,  3,  3,  2,
       0,  2,  3,  0,  0,  0,  0,  2,  0,  0,  0,  0,  0,  0,
       0,  3,  0,  0,  0,  0,  0,  2,  0,  0,  0,  0,  0,  0,
       2,  2,  0,  0,  0,  3,  0,  2,  0,  0,  2,  3,  2,  0,
       2,  0,  0,  0,  0,  0,  0,  0,  0,  0,  3,  0,  0,  0,
       0,  0,  0,  0,  0,  3,  0,  2,  2,  0,  0,  0,  0,  0,
       0,  0,  0,  0,  0,  0,  0,  0,  3,  0,  3,  0,  0,  0,
       2,  0,  0,  0,  2,  0,  0,  0,  0,  0,  2,  0,  0,  0,
       0,  0,  0,  0,  0,  0,  0,  3,  0,  0,  0,  0,  0,  0,
       3,  0,  0,  2,  0,  0,  0,  0,  0,  0,  0,  0,  0,  3
    };
  static const int wordlist[] =
    {
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str4,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str5,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str6,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str7,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str8,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str11,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str12,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str13,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str14,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str15,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str16,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str17,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str18,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str19,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str23,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str24,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str25,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str28,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str29,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str30,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str31,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str32,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str33,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str35,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str36,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str37,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str38,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str39,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str40,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str41,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str42,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str44,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str45,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str47,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str48,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str50,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str51,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str52,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str53,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str54,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str55,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str56,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str58,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str61,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str62,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str63,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str64,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str65,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str66,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str67,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str69,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str71,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str72,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str73,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str75,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str76,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str77,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str78,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str79,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str84,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str85,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str86,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str88,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str89,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str90,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str91,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str92,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str97,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str98,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str99,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str100,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str102,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str103,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str106,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str108,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str109,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str110,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str113,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str114,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str115,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str116,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str117,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str118,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str119,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str120,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str123,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str125,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str127,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str128,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str130,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str131,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str132,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str133,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str134,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str135,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str137,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str141,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str142,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str143,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str145,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str146,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str147,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str149,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str150,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str151,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str153,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str156,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str157,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str158,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str159,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str160,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str161,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str162,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str163,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str165,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str166,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str167,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str169,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str170,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str171,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str173,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str174,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str175,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str177,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str178,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str179,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str180,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str181,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str182,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str183,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str184,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str186,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str187,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str188,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str189,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str190,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str191,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str193,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str195,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str196,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str197,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str198,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str199,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str200,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str202,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str203,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str204,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str205,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str206,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str207,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str208,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str209,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str210,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str211,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str212,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str213,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str214,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str215,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str217,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str218,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str219,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str220,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str221,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str222,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str223,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str224,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str225,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str226,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str227,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str229,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str230,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str231,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str234,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str235,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str236,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str237,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str240,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str241,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str242,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str244,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str245,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str246,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str247,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str248,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str249,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str250,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str251,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str252,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str253,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str254,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str255,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str257,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str258,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str259,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str260,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str261,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str262,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str264,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str265,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str266,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str267,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str269,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str270,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str272,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str273,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str276,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str277,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str278,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str280,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str281,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str286,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str289,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str291,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str293,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str295,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str296,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str297,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str298,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str300,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str302,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str303,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str304,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str307,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str308,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str309,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str313,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str316,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str320,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str323,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str324,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str329,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str333,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str334,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str335,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str337,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str338,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str343,
      -1, -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str351,
      -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str357,
      -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str364,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str365,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str369,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str371,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str374,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str375,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str376,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str378,
      -1, -1, -1, -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str388,
      -1, -1, -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str397,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str399,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str400,
      -1, -1, -1, -1, -1, -1, -1, -1, -1,
      -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str414,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str416,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str420,
      -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str424,
      -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str430,
      -1, -1, -1, -1, -1, -1, -1, -1, -1,
      -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str441,
      -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str448,
      -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str451,
      -1, -1, -1, -1, -1, -1, -1, -1, -1,
      (int)(size_t)&((struct stringpool_t *)0)->stringpool_str461
    };

  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = uc_locale_language_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        if (len == lengthtable[key])
          {
            register const char *s = wordlist[key] + stringpool;

            if (*str == *s && !memcmp (str + 1, s + 1, len - 1))
              return s;
          }
    }
  return 0;
}
#line 295 "./unicase/locale-languages.gperf"

/*
 * Local Variables:
 * coding: utf-8
 * End:
 */
