package utils

// List of supported regions by the Anthropics API.
var SupportedRegions = []string{
	"Albania",
	"Algeria",
	"Andorra",
	"Angola",
	"Antigua and Barbuda",
	"Argentina",
	"Armenia",
	"Australia",
	"Austria",
	"Azerbaijan",
	"Bahamas",
	"Bangladesh",
	"Barbados",
	"Belgium",
	"Belize",
	"Benin",
	"Bhutan",
	"Bolivia",
	"Botswana",
	"Brazil",
	"Brunei",
	"Bulgaria",
	"Burkina Faso",
	"Cabo Verde",
	"Canada",
	"Chile",
	"Colombia",
	"Comoros",
	"Congo, Republic of the",
	"Costa Rica",
	"Côte d'Ivoire",
	"Croatia",
	"Cyprus",
	"Czechia (Czech Republic)",
	"Denmark",
	"Djibouti",
	"Dominica",
	"Dominican Republic",
	"Ecuador",
	"El Salvador",
	"Estonia",
	"Fiji",
	"Finland",
	"France",
	"Gabon",
	"Gambia",
	"Georgia",
	"Germany",
	"Ghana",
	"Greece",
	"Grenada",
	"Guatemala",
	"Guinea",
	"Guinea-Bissau",
	"Guyana",
	"Haiti",
	"Holy See (Vatican City)",
	"Honduras",
	"Hungary",
	"Iceland",
	"India",
	"Indonesia",
	"Iraq",
	"Ireland",
	"Israel",
	"Italy",
	"Jamaica",
	"Japan",
	"Jordan",
	"Kazakhstan",
	"Kenya",
	"Kiribati",
	"Kuwait",
	"Kyrgyzstan",
	"Latvia",
	"Lebanon",
	"Lesotho",
	"Liberia",
	"Liechtenstein",
	"Lithuania",
	"Luxembourg",
	"Madagascar",
	"Malawi",
	"Malaysia",
	"Maldives",
	"Malta",
	"Marshall Islands",
	"Mauritania",
	"Mauritius",
	"Mexico",
	"Micronesia",
	"Moldova",
	"Monaco",
	"Mongolia",
	"Montenegro",
	"Morocco",
	"Mozambique",
	"Namibia",
	"Nauru",
	"Nepal",
	"Netherlands",
	"New Zealand",
	"Niger",
	"Nigeria",
	"North Macedonia",
	"Norway",
	"Oman",
	"Pakistan",
	"Palau",
	"Palestine",
	"Panama",
	"Papua New Guinea",
	"Paraguay",
	"Peru",
	"Philippines",
	"Poland",
	"Portugal",
	"Qatar",
	"Romania",
	"Rwanda",
	"Saint Kitts and Nevis",
	"Saint Lucia",
	"Saint Vincent and the Grenadines",
	"Samoa",
	"San Marino",
	"Sao Tome and Principe",
	"Senegal",
	"Serbia",
	"Seychelles",
	"Sierra Leone",
	"Singapore",
	"Slovakia",
	"Slovenia",
	"Solomon Islands",
	"South Africa",
	"South Korea",
	"Spain",
	"Sri Lanka",
	"Suriname",
	"Sweden",
	"Switzerland",
	"Taiwan",
	"Tanzania",
	"Thailand",
	"Timor-Leste, Democratic Republic of",
	"Togo",
	"Tonga",
	"Trinidad and Tobago",
	"Tunisia",
	"Turkey",
	"Tuvalu",
	"Uganda",
	"Ukraine (except Crimea, Donetsk, and Luhansk regions)",
	"United Arab Emirates",
	"United Kingdom",
	"United States of America",
	"Uruguay",
	"Vanuatu",
	"Zambia",
}

// IsRegionSupported checks if the provided region is supported by the Anthropics API.
func IsRegionSupported(region string) bool {
	for _, supported := range SupportedRegions {
		if region == supported {
			return true
		}
	}
	return false
}
