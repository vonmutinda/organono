package utils

var (
	countryCodes = map[string]string{
		"880":    "BD",
		"32":     "BE",
		"226":    "BF",
		"359":    "BG",
		"387":    "BA",
		"1246":   "BB",
		"681":    "WF",
		"1441":   "BM",
		"673":    "BN",
		"591":    "BO",
		"973":    "BH",
		"257":    "BI",
		"229":    "BJ",
		"975":    "BT",
		"1876":   "JM",
		"267":    "BW",
		"685":    "WS",
		"599":    "BQ",
		"55":     "BR",
		"1242":   "BS",
		"441534": "JE",
		"375":    "BY",
		"501":    "BZ",
		"7":      "RU",
		"250":    "RW",
		"381":    "RS",
		"670":    "TL",
		"262":    "RE",
		"993":    "TM",
		"992":    "TJ",
		"40":     "RO",
		"690":    "TK",
		"245":    "GW",
		"1671":   "GU",
		"502":    "GT",
		"":       "GS",
		"30":     "GR",
		"240":    "GQ",
		"590":    "GP",
		"81":     "JP",
		"592":    "GY",
		"441481": "GG",
		"594":    "GF",
		"995":    "GE",
		"1473":   "GD",
		"44":     "GB",
		"241":    "GA",
		"503":    "SV",
		"224":    "GN",
		"220":    "GM",
		"299":    "GL",
		"350":    "GI",
		"233":    "GH",
		"968":    "OM",
		"216":    "TN",
		"962":    "JO",
		"385":    "HR",
		"509":    "HT",
		"36":     "HU",
		"852":    "HK",
		"504":    "HN",
		"58":     "VE",
		"1787":   "PR",
		"1939":   "PR",
		"970":    "PS",
		"680":    "PW",
		"351":    "PT",
		"47":     "SJ",
		"595":    "PY",
		"964":    "IQ",
		"507":    "PA",
		"689":    "PF",
		"675":    "PG",
		"51":     "PE",
		"92":     "PK",
		"63":     "PH",
		"870":    "PN",
		"48":     "PL",
		"508":    "PM",
		"260":    "ZM",
		"372":    "EE",
		"20":     "EG",
		"27":     "ZA",
		"593":    "EC",
		"39":     "IT",
		"84":     "VN",
		"677":    "SB",
		"251":    "ET",
		"252":    "SO",
		"263":    "ZW",
		"966":    "SA",
		"34":     "ES",
		"291":    "ER",
		"382":    "ME",
		"373":    "MD",
		"261":    "MG",
		"212":    "MA",
		"377":    "MC",
		"998":    "UZ",
		"95":     "MM",
		"223":    "ML",
		"853":    "MO",
		"976":    "MN",
		"692":    "MH",
		"389":    "MK",
		"230":    "MU",
		"356":    "MT",
		"265":    "MW",
		"960":    "MV",
		"596":    "MQ",
		"1670":   "MP",
		"1664":   "MS",
		"222":    "MR",
		"441624": "IM",
		"256":    "UG",
		"255":    "TZ",
		"60":     "MY",
		"52":     "MX",
		"972":    "IL",
		"33":     "FR",
		"246":    "IO",
		"290":    "SH",
		"358":    "FI",
		"679":    "FJ",
		"500":    "FK",
		"691":    "FM",
		"298":    "FO",
		"505":    "NI",
		"31":     "NL",
		"264":    "NA",
		"678":    "VU",
		"687":    "NC",
		"227":    "NE",
		"672":    "NF",
		"234":    "NG",
		"64":     "NZ",
		"977":    "NP",
		"674":    "NR",
		"683":    "NU",
		"682":    "CK",
		"225":    "CI",
		"41":     "CH",
		"57":     "CO",
		"86":     "CN",
		"237":    "CM",
		"56":     "CL",
		"242":    "CG",
		"236":    "CF",
		"243":    "CD",
		"420":    "CZ",
		"357":    "CY",
		"506":    "CR",
		"238":    "CV",
		"53":     "CU",
		"268":    "SZ",
		"963":    "SY",
		"996":    "KG",
		"254":    "KE",
		"211":    "SS",
		"597":    "SR",
		"686":    "KI",
		"855":    "KH",
		"1869":   "KN",
		"269":    "KM",
		"239":    "ST",
		"421":    "SK",
		"82":     "KR",
		"386":    "SI",
		"850":    "KP",
		"965":    "KW",
		"221":    "SN",
		"378":    "SM",
		"232":    "SL",
		"248":    "SC",
		"1345":   "KY",
		"65":     "SG",
		"46":     "SE",
		"249":    "SD",
		"1809":   "DO",
		"1829":   "DO",
		"1767":   "DM",
		"253":    "DJ",
		"45":     "DK",
		"1284":   "VG",
		"49":     "DE",
		"967":    "YE",
		"213":    "DZ",
		"1":      "US",
		"598":    "UY",
		"961":    "LB",
		"1758":   "LC",
		"856":    "LA",
		"688":    "TV",
		"886":    "TW",
		"1868":   "TT",
		"90":     "TR",
		"94":     "LK",
		"423":    "LI",
		"371":    "LV",
		"676":    "TO",
		"370":    "LT",
		"352":    "LU",
		"231":    "LR",
		"266":    "LS",
		"66":     "TH",
		"228":    "TG",
		"235":    "TD",
		"1649":   "TC",
		"218":    "LY",
		"379":    "VA",
		"1784":   "VC",
		"971":    "AE",
		"376":    "AD",
		"1268":   "AG",
		"93":     "AF",
		"1264":   "AI",
		"1340":   "VI",
		"354":    "IS",
		"98":     "IR",
		"374":    "AM",
		"355":    "AL",
		"244":    "AO",
		"1684":   "AS",
		"54":     "AR",
		"61":     "AU",
		"43":     "AT",
		"297":    "AW",
		"91":     "IN",
		"35818":  "AX",
		"994":    "AZ",
		"353":    "IE",
		"62":     "ID",
		"380":    "UA",
		"974":    "QA",
		"258":    "MZ",
	}
)
