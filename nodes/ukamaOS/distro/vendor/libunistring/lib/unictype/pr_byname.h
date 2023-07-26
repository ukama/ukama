/* ANSI-C code produced by gperf version 3.2 */
/* Command-line: gperf -m 10 ./unictype/pr_byname.gperf  */
/* Computed positions: -k'1-2,8,14,18,$' */

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

#line 25 "./unictype/pr_byname.gperf"
struct named_property { int name; int property_index; };

#define TOTAL_KEYWORDS 206
#define MIN_WORD_LENGTH 2
#define MAX_WORD_LENGTH 34
#define MIN_HASH_VALUE 8
#define MAX_HASH_VALUE 619
/* maximum key range = 612, duplicates = 0 */

#ifdef __GNUC__
__inline
#else
#ifdef __cplusplus
inline
#endif
#endif
static unsigned int
properties_hash (register const char *str, register size_t len)
{
  static const unsigned short asso_values[] =
    {
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 620, 620, 620, 620, 620,
      620, 620, 620, 620, 620, 110, 620,  29,  17, 101,
        2,   2, 167, 230,  92,   2,  65,  62,  41, 152,
       74,   2, 104,  14,  14,  20,   5,  80,  41, 140,
      181,  68,  35, 620, 620, 620, 620, 620
    };
  register unsigned int hval = len;

  switch (hval)
    {
      default:
        hval += asso_values[(unsigned char)str[17]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 17:
      case 16:
      case 15:
      case 14:
        hval += asso_values[(unsigned char)str[13]];
#if defined __cplusplus && (__cplusplus >= 201703L || (__cplusplus >= 201103L && defined __clang_major__ && defined __clang_minor__ && __clang_major__ + (__clang_minor__ >= 9) > 3))
      [[fallthrough]];
#elif defined __GNUC__ && __GNUC__ >= 7
      __attribute__ ((__fallthrough__));
#endif
      /*FALLTHROUGH*/
      case 13:
      case 12:
      case 11:
      case 10:
      case 9:
      case 8:
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

struct properties_stringpool_t
  {
    char properties_stringpool_str8[sizeof("di")];
    char properties_stringpool_str9[sizeof("odi")];
    char properties_stringpool_str10[sizeof("ideo")];
    char properties_stringpool_str13[sizeof("idst")];
    char properties_stringpool_str16[sizeof("idstart")];
    char properties_stringpool_str19[sizeof("id_continue")];
    char properties_stringpool_str21[sizeof("deprecated")];
    char properties_stringpool_str22[sizeof("id_start")];
    char properties_stringpool_str23[sizeof("decimaldigit")];
    char properties_stringpool_str25[sizeof("idsb")];
    char properties_stringpool_str26[sizeof("sd")];
    char properties_stringpool_str27[sizeof("ids")];
    char properties_stringpool_str28[sizeof("oids")];
    char properties_stringpool_str30[sizeof("other_id_continue")];
    char properties_stringpool_str33[sizeof("other_id_start")];
    char properties_stringpool_str36[sizeof("dia")];
    char properties_stringpool_str38[sizeof("titlecase")];
    char properties_stringpool_str39[sizeof("softdotted")];
    char properties_stringpool_str40[sizeof("soft_dotted")];
    char properties_stringpool_str42[sizeof("bidiwhitespace")];
    char properties_stringpool_str44[sizeof("otheridstart")];
    char properties_stringpool_str45[sizeof("bidieuropeandigit")];
    char properties_stringpool_str46[sizeof("other_lowercase")];
    char properties_stringpool_str48[sizeof("loe")];
    char properties_stringpool_str50[sizeof("bidiembeddingoroverride")];
    char properties_stringpool_str51[sizeof("other_grapheme_extend")];
    char properties_stringpool_str55[sizeof("defaultignorablecodepoint")];
    char properties_stringpool_str58[sizeof("bidiarabicdigit")];
    char properties_stringpool_str62[sizeof("lower")];
    char properties_stringpool_str63[sizeof("olower")];
    char properties_stringpool_str66[sizeof("oalpha")];
    char properties_stringpool_str70[sizeof("ids_binary_operator")];
    char properties_stringpool_str72[sizeof("bidi_arabic_digit")];
    char properties_stringpool_str73[sizeof("ascii_hex_digit")];
    char properties_stringpool_str74[sizeof("lowercase")];
    char properties_stringpool_str76[sizeof("bidicontrol")];
    char properties_stringpool_str77[sizeof("bidi_eur_num_terminator")];
    char properties_stringpool_str78[sizeof("idsbinaryoperator")];
    char properties_stringpool_str79[sizeof("iso_control")];
    char properties_stringpool_str83[sizeof("vs")];
    char properties_stringpool_str84[sizeof("sentence_terminal")];
    char properties_stringpool_str87[sizeof("isocontrol")];
    char properties_stringpool_str89[sizeof("uideo")];
    char properties_stringpool_str90[sizeof("idcontinue")];
    char properties_stringpool_str91[sizeof("radical")];
    char properties_stringpool_str92[sizeof("bidiotherneutral")];
    char properties_stringpool_str93[sizeof("idstrinaryoperator")];
    char properties_stringpool_str96[sizeof("leftofpair")];
    char properties_stringpool_str99[sizeof("lineseparator")];
    char properties_stringpool_str102[sizeof("oupper")];
    char properties_stringpool_str103[sizeof("bidi_eur_num_separator")];
    char properties_stringpool_str104[sizeof("alpha")];
    char properties_stringpool_str107[sizeof("ci")];
    char properties_stringpool_str108[sizeof("idc")];
    char properties_stringpool_str109[sizeof("oidc")];
    char properties_stringpool_str110[sizeof("hex_digit")];
    char properties_stringpool_str111[sizeof("dep")];
    char properties_stringpool_str112[sizeof("hexdigit")];
    char properties_stringpool_str113[sizeof("othermath")];
    char properties_stringpool_str116[sizeof("diacritic")];
    char properties_stringpool_str117[sizeof("notacharacter")];
    char properties_stringpool_str119[sizeof("composite")];
    char properties_stringpool_str120[sizeof("variation_selector")];
    char properties_stringpool_str124[sizeof("joincontrol")];
    char properties_stringpool_str125[sizeof("bidic")];
    char properties_stringpool_str126[sizeof("bidi_c")];
    char properties_stringpool_str127[sizeof("dash")];
    char properties_stringpool_str129[sizeof("otheruppercase")];
    char properties_stringpool_str131[sizeof("space")];
    char properties_stringpool_str132[sizeof("decimal_digit")];
    char properties_stringpool_str133[sizeof("othergraphemeextend")];
    char properties_stringpool_str136[sizeof("bidilefttoright")];
    char properties_stringpool_str137[sizeof("cased")];
    char properties_stringpool_str138[sizeof("other_math")];
    char properties_stringpool_str139[sizeof("bidi_whitespace")];
    char properties_stringpool_str141[sizeof("zero_width")];
    char properties_stringpool_str143[sizeof("zerowidth")];
    char properties_stringpool_str146[sizeof("bidi_control")];
    char properties_stringpool_str147[sizeof("caseignorable")];
    char properties_stringpool_str148[sizeof("other_uppercase")];
    char properties_stringpool_str149[sizeof("terminal_punctuation")];
    char properties_stringpool_str155[sizeof("sentenceterminal")];
    char properties_stringpool_str157[sizeof("bidieurnumseparator")];
    char properties_stringpool_str158[sizeof("patws")];
    char properties_stringpool_str159[sizeof("pat_ws")];
    char properties_stringpool_str160[sizeof("other_default_ignorable_code_point")];
    char properties_stringpool_str163[sizeof("term")];
    char properties_stringpool_str164[sizeof("bidi_block_separator")];
    char properties_stringpool_str165[sizeof("otherlowercase")];
    char properties_stringpool_str168[sizeof("wspace")];
    char properties_stringpool_str169[sizeof("bidi_european_digit")];
    char properties_stringpool_str170[sizeof("other_alphabetic")];
    char properties_stringpool_str171[sizeof("quotationmark")];
    char properties_stringpool_str173[sizeof("joinc")];
    char properties_stringpool_str174[sizeof("join_c")];
    char properties_stringpool_str176[sizeof("non_break")];
    char properties_stringpool_str178[sizeof("bidi_hebrew_right_to_left")];
    char properties_stringpool_str179[sizeof("left_of_pair")];
    char properties_stringpool_str180[sizeof("bidiblockseparator")];
    char properties_stringpool_str184[sizeof("bidiboundaryneutral")];
    char properties_stringpool_str186[sizeof("alphabetic")];
    char properties_stringpool_str189[sizeof("line_separator")];
    char properties_stringpool_str190[sizeof("bidi_arabic_right_to_left")];
    char properties_stringpool_str191[sizeof("ext")];
    char properties_stringpool_str192[sizeof("bidihebrewrighttoleft")];
    char properties_stringpool_str193[sizeof("bidipdf")];
    char properties_stringpool_str194[sizeof("join_control")];
    char properties_stringpool_str195[sizeof("bidiarabicrighttoleft")];
    char properties_stringpool_str198[sizeof("xidcontinue")];
    char properties_stringpool_str199[sizeof("not_a_character")];
    char properties_stringpool_str201[sizeof("xidstart")];
    char properties_stringpool_str202[sizeof("xid_continue")];
    char properties_stringpool_str203[sizeof("upper")];
    char properties_stringpool_str204[sizeof("variationselector")];
    char properties_stringpool_str205[sizeof("otheridcontinue")];
    char properties_stringpool_str207[sizeof("xids")];
    char properties_stringpool_str208[sizeof("nonbreak")];
    char properties_stringpool_str210[sizeof("privateuse")];
    char properties_stringpool_str211[sizeof("xid_start")];
    char properties_stringpool_str213[sizeof("patsyn")];
    char properties_stringpool_str214[sizeof("pat_syn")];
    char properties_stringpool_str215[sizeof("uppercase")];
    char properties_stringpool_str219[sizeof("extender")];
    char properties_stringpool_str220[sizeof("ideographic")];
    char properties_stringpool_str221[sizeof("ids_trinary_operator")];
    char properties_stringpool_str222[sizeof("case_ignorable")];
    char properties_stringpool_str223[sizeof("terminalpunctuation")];
    char properties_stringpool_str225[sizeof("formatcontrol")];
    char properties_stringpool_str228[sizeof("bidi_left_to_right")];
    char properties_stringpool_str229[sizeof("otheralphabetic")];
    char properties_stringpool_str233[sizeof("qmark")];
    char properties_stringpool_str234[sizeof("quotation_mark")];
    char properties_stringpool_str235[sizeof("bidicommonseparator")];
    char properties_stringpool_str237[sizeof("bidi_common_separator")];
    char properties_stringpool_str240[sizeof("hyphen")];
    char properties_stringpool_str241[sizeof("private_use")];
    char properties_stringpool_str243[sizeof("ogrext")];
    char properties_stringpool_str244[sizeof("ogr_ext")];
    char properties_stringpool_str248[sizeof("asciihexdigit")];
    char properties_stringpool_str249[sizeof("cwt")];
    char properties_stringpool_str251[sizeof("omath")];
    char properties_stringpool_str252[sizeof("grbase")];
    char properties_stringpool_str253[sizeof("gr_base")];
    char properties_stringpool_str254[sizeof("grext")];
    char properties_stringpool_str255[sizeof("gr_ext")];
    char properties_stringpool_str260[sizeof("graphemebase")];
    char properties_stringpool_str261[sizeof("grapheme_base")];
    char properties_stringpool_str262[sizeof("numeric")];
    char properties_stringpool_str264[sizeof("graphemeextend")];
    char properties_stringpool_str273[sizeof("whitespace")];
    char properties_stringpool_str274[sizeof("punctuation")];
    char properties_stringpool_str276[sizeof("bidi_boundary_neutral")];
    char properties_stringpool_str277[sizeof("math")];
    char properties_stringpool_str278[sizeof("hex")];
    char properties_stringpool_str282[sizeof("unassigned_code_value")];
    char properties_stringpool_str284[sizeof("bidieurnumterminator")];
    char properties_stringpool_str285[sizeof("cwl")];
    char properties_stringpool_str286[sizeof("default_ignorable_code_point")];
    char properties_stringpool_str288[sizeof("xidc")];
    char properties_stringpool_str291[sizeof("bidi_other_neutral")];
    char properties_stringpool_str293[sizeof("unifiedideograph")];
    char properties_stringpool_str297[sizeof("paragraphseparator")];
    char properties_stringpool_str301[sizeof("paragraph_separator")];
    char properties_stringpool_str306[sizeof("ahex")];
    char properties_stringpool_str307[sizeof("currency_symbol")];
    char properties_stringpool_str309[sizeof("pairedpunctuation")];
    char properties_stringpool_str312[sizeof("grlink")];
    char properties_stringpool_str313[sizeof("gr_link")];
    char properties_stringpool_str314[sizeof("bidisegmentseparator")];
    char properties_stringpool_str319[sizeof("bidi_segment_separator")];
    char properties_stringpool_str320[sizeof("graphemelink")];
    char properties_stringpool_str321[sizeof("grapheme_link")];
    char properties_stringpool_str324[sizeof("cwu")];
    char properties_stringpool_str326[sizeof("logicalorderexception")];
    char properties_stringpool_str330[sizeof("bidi_non_spacing_mark")];
    char properties_stringpool_str331[sizeof("unassignedcodevalue")];
    char properties_stringpool_str335[sizeof("changes_when_titlecased")];
    char properties_stringpool_str336[sizeof("ignorable_control")];
    char properties_stringpool_str337[sizeof("grapheme_extend")];
    char properties_stringpool_str344[sizeof("ignorablecontrol")];
    char properties_stringpool_str345[sizeof("currencysymbol")];
    char properties_stringpool_str347[sizeof("patternsyntax")];
    char properties_stringpool_str349[sizeof("white_space")];
    char properties_stringpool_str361[sizeof("bidi_pdf")];
    char properties_stringpool_str362[sizeof("logical_order_exception")];
    char properties_stringpool_str366[sizeof("format_control")];
    char properties_stringpool_str383[sizeof("changes_when_lowercased")];
    char properties_stringpool_str387[sizeof("unified_ideograph")];
    char properties_stringpool_str390[sizeof("changeswhentitlecased")];
    char properties_stringpool_str396[sizeof("patternwhitespace")];
    char properties_stringpool_str397[sizeof("cwcm")];
    char properties_stringpool_str404[sizeof("bidi_embedding_or_override")];
    char properties_stringpool_str411[sizeof("bidinonspacingmark")];
    char properties_stringpool_str412[sizeof("cwcf")];
    char properties_stringpool_str416[sizeof("combining")];
    char properties_stringpool_str417[sizeof("changeswhencasefolded")];
    char properties_stringpool_str422[sizeof("changes_when_uppercased")];
    char properties_stringpool_str432[sizeof("paired_punctuation")];
    char properties_stringpool_str468[sizeof("otherdefaultignorablecodepoint")];
    char properties_stringpool_str475[sizeof("pattern_white_space")];
    char properties_stringpool_str480[sizeof("changeswhencasemapped")];
    char properties_stringpool_str489[sizeof("changeswhenuppercased")];
    char properties_stringpool_str525[sizeof("changeswhenlowercased")];
    char properties_stringpool_str581[sizeof("changes_when_casemapped")];
    char properties_stringpool_str596[sizeof("changes_when_casefolded")];
    char properties_stringpool_str619[sizeof("pattern_syntax")];
  };
static const struct properties_stringpool_t properties_stringpool_contents =
  {
    "di",
    "odi",
    "ideo",
    "idst",
    "idstart",
    "id_continue",
    "deprecated",
    "id_start",
    "decimaldigit",
    "idsb",
    "sd",
    "ids",
    "oids",
    "other_id_continue",
    "other_id_start",
    "dia",
    "titlecase",
    "softdotted",
    "soft_dotted",
    "bidiwhitespace",
    "otheridstart",
    "bidieuropeandigit",
    "other_lowercase",
    "loe",
    "bidiembeddingoroverride",
    "other_grapheme_extend",
    "defaultignorablecodepoint",
    "bidiarabicdigit",
    "lower",
    "olower",
    "oalpha",
    "ids_binary_operator",
    "bidi_arabic_digit",
    "ascii_hex_digit",
    "lowercase",
    "bidicontrol",
    "bidi_eur_num_terminator",
    "idsbinaryoperator",
    "iso_control",
    "vs",
    "sentence_terminal",
    "isocontrol",
    "uideo",
    "idcontinue",
    "radical",
    "bidiotherneutral",
    "idstrinaryoperator",
    "leftofpair",
    "lineseparator",
    "oupper",
    "bidi_eur_num_separator",
    "alpha",
    "ci",
    "idc",
    "oidc",
    "hex_digit",
    "dep",
    "hexdigit",
    "othermath",
    "diacritic",
    "notacharacter",
    "composite",
    "variation_selector",
    "joincontrol",
    "bidic",
    "bidi_c",
    "dash",
    "otheruppercase",
    "space",
    "decimal_digit",
    "othergraphemeextend",
    "bidilefttoright",
    "cased",
    "other_math",
    "bidi_whitespace",
    "zero_width",
    "zerowidth",
    "bidi_control",
    "caseignorable",
    "other_uppercase",
    "terminal_punctuation",
    "sentenceterminal",
    "bidieurnumseparator",
    "patws",
    "pat_ws",
    "other_default_ignorable_code_point",
    "term",
    "bidi_block_separator",
    "otherlowercase",
    "wspace",
    "bidi_european_digit",
    "other_alphabetic",
    "quotationmark",
    "joinc",
    "join_c",
    "non_break",
    "bidi_hebrew_right_to_left",
    "left_of_pair",
    "bidiblockseparator",
    "bidiboundaryneutral",
    "alphabetic",
    "line_separator",
    "bidi_arabic_right_to_left",
    "ext",
    "bidihebrewrighttoleft",
    "bidipdf",
    "join_control",
    "bidiarabicrighttoleft",
    "xidcontinue",
    "not_a_character",
    "xidstart",
    "xid_continue",
    "upper",
    "variationselector",
    "otheridcontinue",
    "xids",
    "nonbreak",
    "privateuse",
    "xid_start",
    "patsyn",
    "pat_syn",
    "uppercase",
    "extender",
    "ideographic",
    "ids_trinary_operator",
    "case_ignorable",
    "terminalpunctuation",
    "formatcontrol",
    "bidi_left_to_right",
    "otheralphabetic",
    "qmark",
    "quotation_mark",
    "bidicommonseparator",
    "bidi_common_separator",
    "hyphen",
    "private_use",
    "ogrext",
    "ogr_ext",
    "asciihexdigit",
    "cwt",
    "omath",
    "grbase",
    "gr_base",
    "grext",
    "gr_ext",
    "graphemebase",
    "grapheme_base",
    "numeric",
    "graphemeextend",
    "whitespace",
    "punctuation",
    "bidi_boundary_neutral",
    "math",
    "hex",
    "unassigned_code_value",
    "bidieurnumterminator",
    "cwl",
    "default_ignorable_code_point",
    "xidc",
    "bidi_other_neutral",
    "unifiedideograph",
    "paragraphseparator",
    "paragraph_separator",
    "ahex",
    "currency_symbol",
    "pairedpunctuation",
    "grlink",
    "gr_link",
    "bidisegmentseparator",
    "bidi_segment_separator",
    "graphemelink",
    "grapheme_link",
    "cwu",
    "logicalorderexception",
    "bidi_non_spacing_mark",
    "unassignedcodevalue",
    "changes_when_titlecased",
    "ignorable_control",
    "grapheme_extend",
    "ignorablecontrol",
    "currencysymbol",
    "patternsyntax",
    "white_space",
    "bidi_pdf",
    "logical_order_exception",
    "format_control",
    "changes_when_lowercased",
    "unified_ideograph",
    "changeswhentitlecased",
    "patternwhitespace",
    "cwcm",
    "bidi_embedding_or_override",
    "bidinonspacingmark",
    "cwcf",
    "combining",
    "changeswhencasefolded",
    "changes_when_uppercased",
    "paired_punctuation",
    "otherdefaultignorablecodepoint",
    "pattern_white_space",
    "changeswhencasemapped",
    "changeswhenuppercased",
    "changeswhenlowercased",
    "changes_when_casemapped",
    "changes_when_casefolded",
    "pattern_syntax"
  };
#define properties_stringpool ((const char *) &properties_stringpool_contents)

static const struct named_property properties[] =
  {
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 49 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str8, UC_PROPERTY_INDEX_DEFAULT_IGNORABLE_CODE_POINT},
#line 52 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str9, UC_PROPERTY_INDEX_OTHER_DEFAULT_IGNORABLE_CODE_POINT},
#line 187 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str10, UC_PROPERTY_INDEX_IDEOGRAPHIC},
    {-1}, {-1},
#line 197 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str13, UC_PROPERTY_INDEX_IDS_TRINARY_OPERATOR},
    {-1}, {-1},
#line 99 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str16, UC_PROPERTY_INDEX_ID_START},
    {-1}, {-1},
#line 104 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str19, UC_PROPERTY_INDEX_ID_CONTINUE},
    {-1},
#line 53 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str21, UC_PROPERTY_INDEX_DEPRECATED},
#line 98 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str22, UC_PROPERTY_INDEX_ID_START},
#line 235 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str23, UC_PROPERTY_INDEX_DECIMAL_DIGIT},
    {-1},
#line 194 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str25, UC_PROPERTY_INDEX_IDS_BINARY_OPERATOR},
#line 97 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str26, UC_PROPERTY_INDEX_SOFT_DOTTED},
#line 100 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str27, UC_PROPERTY_INDEX_ID_START},
#line 103 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str28, UC_PROPERTY_INDEX_OTHER_ID_START},
    {-1},
#line 107 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str30, UC_PROPERTY_INDEX_OTHER_ID_CONTINUE},
    {-1}, {-1},
#line 101 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str33, UC_PROPERTY_INDEX_OTHER_ID_START},
    {-1}, {-1},
#line 238 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str36, UC_PROPERTY_INDEX_DIACRITIC},
    {-1},
#line 75 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str38, UC_PROPERTY_INDEX_TITLECASE},
#line 96 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str39, UC_PROPERTY_INDEX_SOFT_DOTTED},
#line 95 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str40, UC_PROPERTY_INDEX_SOFT_DOTTED},
    {-1},
#line 169 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str42, UC_PROPERTY_INDEX_BIDI_WHITESPACE},
    {-1},
#line 102 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str44, UC_PROPERTY_INDEX_OTHER_ID_START},
#line 155 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str45, UC_PROPERTY_INDEX_BIDI_EUROPEAN_DIGIT},
#line 72 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str46, UC_PROPERTY_INDEX_OTHER_LOWERCASE},
    {-1},
#line 57 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str48, UC_PROPERTY_INDEX_LOGICAL_ORDER_EXCEPTION},
    {-1},
#line 177 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str50, UC_PROPERTY_INDEX_BIDI_EMBEDDING_OR_OVERRIDE},
#line 136 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str51, UC_PROPERTY_INDEX_OTHER_GRAPHEME_EXTEND},
    {-1}, {-1}, {-1},
#line 48 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str55, UC_PROPERTY_INDEX_DEFAULT_IGNORABLE_CODE_POINT},
    {-1}, {-1},
#line 161 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str58, UC_PROPERTY_INDEX_BIDI_ARABIC_DIGIT},
    {-1}, {-1}, {-1},
#line 71 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str62, UC_PROPERTY_INDEX_LOWERCASE},
#line 74 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str63, UC_PROPERTY_INDEX_OTHER_LOWERCASE},
    {-1}, {-1},
#line 44 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str66, UC_PROPERTY_INDEX_OTHER_ALPHABETIC},
    {-1}, {-1}, {-1},
#line 192 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str70, UC_PROPERTY_INDEX_IDS_BINARY_OPERATOR},
    {-1},
#line 160 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str72, UC_PROPERTY_INDEX_BIDI_ARABIC_DIGIT},
#line 183 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str73, UC_PROPERTY_INDEX_ASCII_HEX_DIGIT},
#line 70 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str74, UC_PROPERTY_INDEX_LOWERCASE},
    {-1},
#line 145 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str76, UC_PROPERTY_INDEX_BIDI_CONTROL},
#line 158 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str77, UC_PROPERTY_INDEX_BIDI_EUR_NUM_TERMINATOR},
#line 193 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str78, UC_PROPERTY_INDEX_IDS_BINARY_OPERATOR},
#line 203 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str79, UC_PROPERTY_INDEX_ISO_CONTROL},
    {-1}, {-1}, {-1},
#line 60 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str83, UC_PROPERTY_INDEX_VARIATION_SELECTOR},
#line 217 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str84, UC_PROPERTY_INDEX_SENTENCE_TERMINAL},
    {-1}, {-1},
#line 204 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str87, UC_PROPERTY_INDEX_ISO_CONTROL},
    {-1},
#line 190 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str89, UC_PROPERTY_INDEX_UNIFIED_IDEOGRAPH},
#line 105 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str90, UC_PROPERTY_INDEX_ID_CONTINUE},
#line 191 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str91, UC_PROPERTY_INDEX_RADICAL},
#line 179 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str92, UC_PROPERTY_INDEX_BIDI_OTHER_NEUTRAL},
#line 196 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str93, UC_PROPERTY_INDEX_IDS_TRINARY_OPERATOR},
    {-1}, {-1},
#line 231 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str96, UC_PROPERTY_INDEX_LEFT_OF_PAIR},
    {-1}, {-1},
#line 211 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str99, UC_PROPERTY_INDEX_LINE_SEPARATOR},
    {-1}, {-1},
#line 69 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str102, UC_PROPERTY_INDEX_OTHER_UPPERCASE},
#line 156 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str103, UC_PROPERTY_INDEX_BIDI_EUR_NUM_SEPARATOR},
#line 41 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str104, UC_PROPERTY_INDEX_ALPHABETIC},
    {-1}, {-1},
#line 79 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str107, UC_PROPERTY_INDEX_CASE_IGNORABLE},
#line 106 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str108, UC_PROPERTY_INDEX_ID_CONTINUE},
#line 109 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str109, UC_PROPERTY_INDEX_OTHER_ID_CONTINUE},
#line 180 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str110, UC_PROPERTY_INDEX_HEX_DIGIT},
#line 54 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str111, UC_PROPERTY_INDEX_DEPRECATED},
#line 181 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str112, UC_PROPERTY_INDEX_HEX_DIGIT},
#line 226 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str113, UC_PROPERTY_INDEX_OTHER_MATH},
    {-1}, {-1},
#line 237 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str116, UC_PROPERTY_INDEX_DIACRITIC},
#line 46 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str117, UC_PROPERTY_INDEX_NOT_A_CHARACTER},
    {-1},
#line 233 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str119, UC_PROPERTY_INDEX_COMPOSITE},
#line 58 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str120, UC_PROPERTY_INDEX_VARIATION_SELECTOR},
    {-1}, {-1}, {-1},
#line 125 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str124, UC_PROPERTY_INDEX_JOIN_CONTROL},
#line 147 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str125, UC_PROPERTY_INDEX_BIDI_CONTROL},
#line 146 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str126, UC_PROPERTY_INDEX_BIDI_CONTROL},
#line 207 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str127, UC_PROPERTY_INDEX_DASH},
    {-1},
#line 68 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str129, UC_PROPERTY_INDEX_OTHER_UPPERCASE},
    {-1},
#line 200 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str131, UC_PROPERTY_INDEX_SPACE},
#line 234 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str132, UC_PROPERTY_INDEX_DECIMAL_DIGIT},
#line 137 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str133, UC_PROPERTY_INDEX_OTHER_GRAPHEME_EXTEND},
    {-1}, {-1},
#line 149 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str136, UC_PROPERTY_INDEX_BIDI_LEFT_TO_RIGHT},
#line 76 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str137, UC_PROPERTY_INDEX_CASED},
#line 225 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str138, UC_PROPERTY_INDEX_OTHER_MATH},
#line 168 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str139, UC_PROPERTY_INDEX_BIDI_WHITESPACE},
    {-1},
#line 198 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str141, UC_PROPERTY_INDEX_ZERO_WIDTH},
    {-1},
#line 199 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str143, UC_PROPERTY_INDEX_ZERO_WIDTH},
    {-1}, {-1},
#line 144 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str146, UC_PROPERTY_INDEX_BIDI_CONTROL},
#line 78 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str147, UC_PROPERTY_INDEX_CASE_IGNORABLE},
#line 67 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str148, UC_PROPERTY_INDEX_OTHER_UPPERCASE},
#line 219 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str149, UC_PROPERTY_INDEX_TERMINAL_PUNCTUATION},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 218 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str155, UC_PROPERTY_INDEX_SENTENCE_TERMINAL},
    {-1},
#line 157 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str157, UC_PROPERTY_INDEX_BIDI_EUR_NUM_SEPARATOR},
#line 119 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str158, UC_PROPERTY_INDEX_PATTERN_WHITE_SPACE},
#line 118 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str159, UC_PROPERTY_INDEX_PATTERN_WHITE_SPACE},
#line 50 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str160, UC_PROPERTY_INDEX_OTHER_DEFAULT_IGNORABLE_CODE_POINT},
    {-1}, {-1},
#line 221 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str163, UC_PROPERTY_INDEX_TERMINAL_PUNCTUATION},
#line 164 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str164, UC_PROPERTY_INDEX_BIDI_BLOCK_SEPARATOR},
#line 73 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str165, UC_PROPERTY_INDEX_OTHER_LOWERCASE},
    {-1}, {-1},
#line 39 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str168, UC_PROPERTY_INDEX_WHITE_SPACE},
#line 154 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str169, UC_PROPERTY_INDEX_BIDI_EUROPEAN_DIGIT},
#line 42 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str170, UC_PROPERTY_INDEX_OTHER_ALPHABETIC},
#line 215 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str171, UC_PROPERTY_INDEX_QUOTATION_MARK},
    {-1},
#line 127 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str173, UC_PROPERTY_INDEX_JOIN_CONTROL},
#line 126 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str174, UC_PROPERTY_INDEX_JOIN_CONTROL},
    {-1},
#line 201 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str176, UC_PROPERTY_INDEX_NON_BREAK},
    {-1},
#line 150 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str178, UC_PROPERTY_INDEX_BIDI_HEBREW_RIGHT_TO_LEFT},
#line 230 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str179, UC_PROPERTY_INDEX_LEFT_OF_PAIR},
#line 165 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str180, UC_PROPERTY_INDEX_BIDI_BLOCK_SEPARATOR},
    {-1}, {-1}, {-1},
#line 173 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str184, UC_PROPERTY_INDEX_BIDI_BOUNDARY_NEUTRAL},
    {-1},
#line 40 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str186, UC_PROPERTY_INDEX_ALPHABETIC},
    {-1}, {-1},
#line 210 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str189, UC_PROPERTY_INDEX_LINE_SEPARATOR},
#line 152 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str190, UC_PROPERTY_INDEX_BIDI_ARABIC_RIGHT_TO_LEFT},
#line 240 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str191, UC_PROPERTY_INDEX_EXTENDER},
#line 151 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str192, UC_PROPERTY_INDEX_BIDI_HEBREW_RIGHT_TO_LEFT},
#line 175 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str193, UC_PROPERTY_INDEX_BIDI_PDF},
#line 124 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str194, UC_PROPERTY_INDEX_JOIN_CONTROL},
#line 153 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str195, UC_PROPERTY_INDEX_BIDI_ARABIC_RIGHT_TO_LEFT},
    {-1}, {-1},
#line 114 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str198, UC_PROPERTY_INDEX_XID_CONTINUE},
#line 45 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str199, UC_PROPERTY_INDEX_NOT_A_CHARACTER},
    {-1},
#line 111 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str201, UC_PROPERTY_INDEX_XID_START},
#line 113 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str202, UC_PROPERTY_INDEX_XID_CONTINUE},
#line 66 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str203, UC_PROPERTY_INDEX_UPPERCASE},
#line 59 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str204, UC_PROPERTY_INDEX_VARIATION_SELECTOR},
#line 108 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str205, UC_PROPERTY_INDEX_OTHER_ID_CONTINUE},
    {-1},
#line 112 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str207, UC_PROPERTY_INDEX_XID_START},
#line 202 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str208, UC_PROPERTY_INDEX_NON_BREAK},
    {-1},
#line 62 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str210, UC_PROPERTY_INDEX_PRIVATE_USE},
#line 110 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str211, UC_PROPERTY_INDEX_XID_START},
    {-1},
#line 123 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str213, UC_PROPERTY_INDEX_PATTERN_SYNTAX},
#line 122 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str214, UC_PROPERTY_INDEX_PATTERN_SYNTAX},
#line 65 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str215, UC_PROPERTY_INDEX_UPPERCASE},
    {-1}, {-1}, {-1},
#line 239 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str219, UC_PROPERTY_INDEX_EXTENDER},
#line 186 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str220, UC_PROPERTY_INDEX_IDEOGRAPHIC},
#line 195 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str221, UC_PROPERTY_INDEX_IDS_TRINARY_OPERATOR},
#line 77 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str222, UC_PROPERTY_INDEX_CASE_IGNORABLE},
#line 220 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str223, UC_PROPERTY_INDEX_TERMINAL_PUNCTUATION},
    {-1},
#line 206 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str225, UC_PROPERTY_INDEX_FORMAT_CONTROL},
    {-1}, {-1},
#line 148 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str228, UC_PROPERTY_INDEX_BIDI_LEFT_TO_RIGHT},
#line 43 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str229, UC_PROPERTY_INDEX_OTHER_ALPHABETIC},
    {-1}, {-1}, {-1},
#line 216 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str233, UC_PROPERTY_INDEX_QUOTATION_MARK},
#line 214 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str234, UC_PROPERTY_INDEX_QUOTATION_MARK},
#line 163 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str235, UC_PROPERTY_INDEX_BIDI_COMMON_SEPARATOR},
    {-1},
#line 162 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str237, UC_PROPERTY_INDEX_BIDI_COMMON_SEPARATOR},
    {-1}, {-1},
#line 208 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str240, UC_PROPERTY_INDEX_HYPHEN},
#line 61 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str241, UC_PROPERTY_INDEX_PRIVATE_USE},
    {-1},
#line 139 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str243, UC_PROPERTY_INDEX_OTHER_GRAPHEME_EXTEND},
#line 138 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str244, UC_PROPERTY_INDEX_OTHER_GRAPHEME_EXTEND},
    {-1}, {-1}, {-1},
#line 184 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str248, UC_PROPERTY_INDEX_ASCII_HEX_DIGIT},
#line 88 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str249, UC_PROPERTY_INDEX_CHANGES_WHEN_TITLECASED},
    {-1},
#line 227 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str251, UC_PROPERTY_INDEX_OTHER_MATH},
#line 131 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str252, UC_PROPERTY_INDEX_GRAPHEME_BASE},
#line 130 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str253, UC_PROPERTY_INDEX_GRAPHEME_BASE},
#line 135 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str254, UC_PROPERTY_INDEX_GRAPHEME_EXTEND},
#line 134 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str255, UC_PROPERTY_INDEX_GRAPHEME_EXTEND},
    {-1}, {-1}, {-1}, {-1},
#line 129 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str260, UC_PROPERTY_INDEX_GRAPHEME_BASE},
#line 128 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str261, UC_PROPERTY_INDEX_GRAPHEME_BASE},
#line 236 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str262, UC_PROPERTY_INDEX_NUMERIC},
    {-1},
#line 133 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str264, UC_PROPERTY_INDEX_GRAPHEME_EXTEND},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 38 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str273, UC_PROPERTY_INDEX_WHITE_SPACE},
#line 209 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str274, UC_PROPERTY_INDEX_PUNCTUATION},
    {-1},
#line 172 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str276, UC_PROPERTY_INDEX_BIDI_BOUNDARY_NEUTRAL},
#line 224 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str277, UC_PROPERTY_INDEX_MATH},
#line 182 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str278, UC_PROPERTY_INDEX_HEX_DIGIT},
    {-1}, {-1}, {-1},
#line 63 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str282, UC_PROPERTY_INDEX_UNASSIGNED_CODE_VALUE},
    {-1},
#line 159 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str284, UC_PROPERTY_INDEX_BIDI_EUR_NUM_TERMINATOR},
#line 82 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str285, UC_PROPERTY_INDEX_CHANGES_WHEN_LOWERCASED},
#line 47 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str286, UC_PROPERTY_INDEX_DEFAULT_IGNORABLE_CODE_POINT},
    {-1},
#line 115 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str288, UC_PROPERTY_INDEX_XID_CONTINUE},
    {-1}, {-1},
#line 178 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str291, UC_PROPERTY_INDEX_BIDI_OTHER_NEUTRAL},
    {-1},
#line 189 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str293, UC_PROPERTY_INDEX_UNIFIED_IDEOGRAPH},
    {-1}, {-1}, {-1},
#line 213 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str297, UC_PROPERTY_INDEX_PARAGRAPH_SEPARATOR},
    {-1}, {-1}, {-1},
#line 212 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str301, UC_PROPERTY_INDEX_PARAGRAPH_SEPARATOR},
    {-1}, {-1}, {-1}, {-1},
#line 185 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str306, UC_PROPERTY_INDEX_ASCII_HEX_DIGIT},
#line 222 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str307, UC_PROPERTY_INDEX_CURRENCY_SYMBOL},
    {-1},
#line 229 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str309, UC_PROPERTY_INDEX_PAIRED_PUNCTUATION},
    {-1}, {-1},
#line 143 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str312, UC_PROPERTY_INDEX_GRAPHEME_LINK},
#line 142 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str313, UC_PROPERTY_INDEX_GRAPHEME_LINK},
#line 167 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str314, UC_PROPERTY_INDEX_BIDI_SEGMENT_SEPARATOR},
    {-1}, {-1}, {-1}, {-1},
#line 166 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str319, UC_PROPERTY_INDEX_BIDI_SEGMENT_SEPARATOR},
#line 141 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str320, UC_PROPERTY_INDEX_GRAPHEME_LINK},
#line 140 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str321, UC_PROPERTY_INDEX_GRAPHEME_LINK},
    {-1}, {-1},
#line 85 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str324, UC_PROPERTY_INDEX_CHANGES_WHEN_UPPERCASED},
    {-1},
#line 56 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str326, UC_PROPERTY_INDEX_LOGICAL_ORDER_EXCEPTION},
    {-1}, {-1}, {-1},
#line 170 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str330, UC_PROPERTY_INDEX_BIDI_NON_SPACING_MARK},
#line 64 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str331, UC_PROPERTY_INDEX_UNASSIGNED_CODE_VALUE},
    {-1}, {-1}, {-1},
#line 86 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str335, UC_PROPERTY_INDEX_CHANGES_WHEN_TITLECASED},
#line 241 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str336, UC_PROPERTY_INDEX_IGNORABLE_CONTROL},
#line 132 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str337, UC_PROPERTY_INDEX_GRAPHEME_EXTEND},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 242 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str344, UC_PROPERTY_INDEX_IGNORABLE_CONTROL},
#line 223 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str345, UC_PROPERTY_INDEX_CURRENCY_SYMBOL},
    {-1},
#line 121 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str347, UC_PROPERTY_INDEX_PATTERN_SYNTAX},
    {-1},
#line 37 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str349, UC_PROPERTY_INDEX_WHITE_SPACE},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1},
#line 174 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str361, UC_PROPERTY_INDEX_BIDI_PDF},
#line 55 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str362, UC_PROPERTY_INDEX_LOGICAL_ORDER_EXCEPTION},
    {-1}, {-1}, {-1},
#line 205 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str366, UC_PROPERTY_INDEX_FORMAT_CONTROL},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 80 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str383, UC_PROPERTY_INDEX_CHANGES_WHEN_LOWERCASED},
    {-1}, {-1}, {-1},
#line 188 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str387, UC_PROPERTY_INDEX_UNIFIED_IDEOGRAPH},
    {-1}, {-1},
#line 87 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str390, UC_PROPERTY_INDEX_CHANGES_WHEN_TITLECASED},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 117 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str396, UC_PROPERTY_INDEX_PATTERN_WHITE_SPACE},
#line 94 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str397, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEMAPPED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 176 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str404, UC_PROPERTY_INDEX_BIDI_EMBEDDING_OR_OVERRIDE},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 171 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str411, UC_PROPERTY_INDEX_BIDI_NON_SPACING_MARK},
#line 91 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str412, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEFOLDED},
    {-1}, {-1}, {-1},
#line 232 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str416, UC_PROPERTY_INDEX_COMBINING},
#line 90 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str417, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEFOLDED},
    {-1}, {-1}, {-1}, {-1},
#line 83 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str422, UC_PROPERTY_INDEX_CHANGES_WHEN_UPPERCASED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 228 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str432, UC_PROPERTY_INDEX_PAIRED_PUNCTUATION},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 51 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str468, UC_PROPERTY_INDEX_OTHER_DEFAULT_IGNORABLE_CODE_POINT},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 116 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str475, UC_PROPERTY_INDEX_PATTERN_WHITE_SPACE},
    {-1}, {-1}, {-1}, {-1},
#line 93 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str480, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEMAPPED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 84 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str489, UC_PROPERTY_INDEX_CHANGES_WHEN_UPPERCASED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
#line 81 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str525, UC_PROPERTY_INDEX_CHANGES_WHEN_LOWERCASED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1},
#line 92 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str581, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEMAPPED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1},
#line 89 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str596, UC_PROPERTY_INDEX_CHANGES_WHEN_CASEFOLDED},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1}, {-1},
    {-1}, {-1}, {-1}, {-1},
#line 120 "./unictype/pr_byname.gperf"
    {(int)(size_t)&((struct properties_stringpool_t *)0)->properties_stringpool_str619, UC_PROPERTY_INDEX_PATTERN_SYNTAX}
  };

static const struct named_property *
uc_property_lookup (register const char *str, register size_t len)
{
  if (len <= MAX_WORD_LENGTH && len >= MIN_WORD_LENGTH)
    {
      register unsigned int key = properties_hash (str, len);

      if (key <= MAX_HASH_VALUE)
        {
          register int o = properties[key].name;
          if (o >= 0)
            {
              register const char *s = o + properties_stringpool;

              if (*str == *s && !strcmp (str + 1, s + 1))
                return &properties[key];
            }
        }
    }
  return 0;
}
