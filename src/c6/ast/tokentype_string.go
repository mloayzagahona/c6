// generated by stringer -type=TokenType token.go; DO NOT EDIT

package ast

import "fmt"

const _TokenType_name = "T_SPACET_COMMENT_LINET_COMMENT_BLOCKT_SEMICOLONT_COMMAT_IDENTT_URLT_MEDIAT_TRUET_FALSET_NULLT_MS_PARAM_NAMET_FUNCTION_NAMET_ID_SELECTORT_CLASS_SELECTORT_TYPE_SELECTORT_UNIVERSAL_SELECTORT_PARENT_SELECTORT_PSEUDO_SELECTORT_FUNCTIONAL_PSEUDOT_INTERPOLATION_SELECTORT_LITERAL_CONCATT_CONCATT_MS_PROGIDT_AND_SELECTORT_DESCENDANT_COMBINATORT_CHILD_COMBINATORT_ADJACENT_SIBLING_COMBINATORT_GENERAL_SIBLING_COMBINATORT_UNICODE_RANGET_IFT_ELSET_ELSE_IFT_INCLUDET_MIXINT_FUNCTIONT_GLOBALT_DEFAULTT_IMPORTANTT_OPTIONALT_FONT_FACET_ORT_ANDT_XORT_PLUST_DIVT_MULT_MINUST_MODT_BRACE_STARTT_BRACE_ENDT_LANG_CODET_BRACKET_LEFTT_ATTRIBUTE_NAMET_BRACKET_RIGHTT_EQUALT_GTT_LTT_GET_LET_ASSIGNT_ATTR_EQUALT_ATTR_TILDE_EQUALT_ATTR_HYPHEN_EQUALT_VARIABLET_IMPORTT_AT_RULET_CHARSETT_QQ_STRINGT_Q_STRINGT_UNQUOTE_STRINGT_PAREN_STARTT_PAREN_ENDT_CONSTANTT_INTEGERT_FLOATT_UNIT_NONET_UNIT_PERCENTT_UNIT_SECONDT_UNIT_MILLISECONDT_UNIT_CHT_UNIT_CMT_UNIT_EMT_UNIT_EXT_UNIT_INT_UNIT_MMT_UNIT_PCT_UNIT_PTT_UNIT_PXT_UNIT_REMT_UNIT_HZT_UNIT_KHZT_UNIT_DPIT_UNIT_DPCMT_UNIT_DPPXT_UNIT_VHT_UNIT_VWT_UNIT_VMINT_UNIT_VMAXT_UNIT_DEGT_UNIT_GRADT_UNIT_RADT_UNIT_TURNT_PROPERTY_NAME_TOKENT_PROPERTY_VALUET_HEX_COLORT_COLONT_INTERPOLATION_STARTT_INTERPOLATION_INNERT_INTERPOLATION_END"

var _TokenType_index = [...]uint16{0, 7, 21, 36, 47, 54, 61, 66, 73, 79, 86, 92, 107, 122, 135, 151, 166, 186, 203, 220, 239, 263, 279, 287, 298, 312, 335, 353, 382, 410, 425, 429, 435, 444, 453, 460, 470, 478, 487, 498, 508, 519, 523, 528, 533, 539, 544, 549, 556, 561, 574, 585, 596, 610, 626, 641, 648, 652, 656, 660, 664, 672, 684, 702, 721, 731, 739, 748, 757, 768, 778, 794, 807, 818, 828, 837, 844, 855, 869, 882, 900, 909, 918, 927, 936, 945, 954, 963, 972, 981, 991, 1000, 1010, 1020, 1031, 1042, 1051, 1060, 1071, 1082, 1092, 1103, 1113, 1124, 1145, 1161, 1172, 1179, 1200, 1221, 1240}

func (i TokenType) String() string {
	if i < 0 || i+1 >= TokenType(len(_TokenType_index)) {
		return fmt.Sprintf("TokenType(%d)", i)
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
